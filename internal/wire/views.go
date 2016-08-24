package wire

import (
	"errors"
	"fmt"
	"sort"

	"github.com/ardanlabs/kit/db"
	"github.com/cayleygraph/cayley"
	"github.com/cayleygraph/cayley/graph"
	"github.com/cayleygraph/cayley/graph/path"
	"github.com/cayleygraph/cayley/quad"
	"github.com/coralproject/shelf/internal/sponge/item"
	"github.com/coralproject/shelf/internal/wire/view"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// ErrNotFound is an error variable thrown when no results are returned from a Mongo query.
var ErrNotFound = errors.New("View items Not found")

// bufferLimit controlled the size of the batches used when upserting saved views.
const bufferLimit = 100

// viewPathToGraphPath translates the path in a view into a "path"
// utilized in graph queries.
func viewPathToGraphPath(v *view.View, key string, graphDB *cayley.Handle) (*path.Path, error) {

	// Sort the view Path value.
	sort.Sort(v.Path)

	// Loop over the path segments translating the path.
	var graphPath *path.Path
	level := 1
	for _, segment := range v.Path {

		// Check that the level is the level we expect (i.e., that the levels
		// are in order)
		if level != segment.Level {
			err := fmt.Errorf("Invalid view path level, expected %d but seeing %d", level, segment.Level)
			return graphPath, err
		}

		// Initialize the path, if we are on level 1.
		if level == 1 {

			// Add the first level relationship.
			switch segment.Direction {
			case "in":
				graphPath = cayley.StartPath(graphDB, quad.String(key)).In(quad.String(segment.Predicate))
			case "out":
				graphPath = cayley.StartPath(graphDB, quad.String(key)).Out(quad.String(segment.Predicate))
			}

			// Add the tag, if present.
			if segment.Tag != "" {
				graphPath = graphPath.Clone().Tag(segment.Tag)
			}

			level++
			continue
		}

		// Add the relationship.
		switch segment.Direction {
		case "in":
			graphPath = graphPath.Clone().In(quad.String(segment.Predicate))
		case "out":
			graphPath = graphPath.Clone().Out(quad.String(segment.Predicate))
		}

		// Add the tag, if present.
		if segment.Tag != "" {
			graphPath = graphPath.Clone().Tag(segment.Tag)
		}

		level++
	}

	return graphPath, nil
}

// viewIDs retrieves the item IDs associated with the view.
func viewIDs(v *view.View, path *path.Path, graphDB *cayley.Handle) ([]string, error) {

	// Build the Cayley iterator.
	it := path.BuildIterator()
	it, _ = it.Optimize()
	defer it.Close()

	// Extract any tags in the View value.
	var viewTags []string
	for _, segment := range v.Path {
		if segment.Tag != "" {
			viewTags = append(viewTags, segment.Tag)
		}
	}

	// Retrieve the end path and tagged item IDs.
	var ids []string
	for it.Next() {

		// Tag the results.
		resultTags := make(map[string]graph.Value)
		it.TagResults(resultTags)

		// Extract the tagged item IDs.
		for _, tag := range viewTags {
			if t, ok := resultTags[tag]; ok {
				ids = append(ids, quad.NativeOf(graphDB.NameOf(t)).(string))
			}
		}
	}
	if it.Err() != nil {
		return ids, it.Err()
	}

	// Remove duplicates.
	found := make(map[string]bool)
	j := 0
	for i, x := range ids {
		if !found[x] {
			found[x] = true
			ids[j] = ids[i]
			j++
		}
	}
	ids = ids[:j]

	return ids, nil
}

// viewSave retrieve items for a view and saves those items to a new collection.
func viewSave(context interface{}, mgoDB *db.DB, v *view.View, viewParams *ViewParams, ids []string) error {

	// Form the query.
	q := bson.M{"item_id": bson.M{"$in": ids}}
	results, err := mgoDB.BatchedQueryMGO(context, v.Collection, q)
	if err != nil {
		return err
	}

	// Set up a Bulk upsert.
	tx, err := mgoDB.BulkOperationMGO(context, viewParams.ResultsCollection)
	if err != nil {
		return err
	}

	// Iterate over the view items.
	var queuedDocs int
	var result item.Item
	for results.Next(&result) {

		// Queue the upsert of the result.
		tx.Upsert(bson.M{"item_id": result.ID}, result)
		queuedDocs++

		// If the queued documents for upsert have reached the buffer limit,
		// run the bulk upsert and re-initialize the bulk operation.
		if queuedDocs >= bufferLimit {
			if _, err := tx.Run(); err != nil {
				return err
			}
			tx, err = mgoDB.BulkOperationMGO(context, viewParams.ResultsCollection)
			if err != nil {
				return err
			}
			queuedDocs = 0
		}
	}
	if err := results.Close(); err != nil {
		return err
	}

	// Run the bulk operation for any remaining queued documents.
	if _, err := tx.Run(); err != nil {
		return err
	}

	return nil
}

// viewItems retrieves the items corresponding to the provided list of item IDs.
func viewItems(context interface{}, db *db.DB, v *view.View, ids []string) ([]bson.M, error) {

	// Form the query.
	var results []bson.M
	f := func(c *mgo.Collection) error {
		return c.Find(bson.M{"item_id": bson.M{"$in": ids}}).All(&results)
	}

	// Execute the query.
	if err := db.ExecuteMGO(context, v.Collection, f); err != nil {
		if err == mgo.ErrNotFound {
			err = ErrNotFound
		}
		return nil, err
	}

	return results, nil
}
