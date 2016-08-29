// Package wire provides support for generating views.
package wire

import (
	"errors"
	"fmt"
	"sort"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/ardanlabs/kit/db"
	"github.com/ardanlabs/kit/log"
	"github.com/cayleygraph/cayley"
	"github.com/cayleygraph/cayley/graph"
	"github.com/cayleygraph/cayley/graph/path"
	"github.com/cayleygraph/cayley/quad"
	"github.com/coralproject/shelf/internal/sponge/item"
	"github.com/coralproject/shelf/internal/wire/relationship"
	"github.com/coralproject/shelf/internal/wire/view"
)

var (
	// ErrNotFound is an error variable thrown when no results are returned from a Mongo query.
	ErrNotFound = errors.New("View items Not found")

	// bufferLimit controls the size of the batches used when upserting saved views.
	bufferLimit = 100
)

const (
	inString  = "in"
	outString = "out"
)

// Result represents what a user will receive after generating a view.
type Result struct {
	Results interface{} `json:"results"`
}

// errResult returns a Result value with an error message.
func errResult(err error) *Result {
	result := Result{
		Results: bson.M{"error": err.Error()},
	}

	return &result
}

// ViewParams represents how the View will be generated and persisted.
type ViewParams struct {
	ViewName          string `json:"view_name"`
	ItemKey           string `json:"item_key"`
	ResultsCollection string `json:"results_collection"`
	BufferLimit       int    `json:"buffer_limit"`
}

//==============================================================================

// Execute executes a graph query to generate the specified view.
func Execute(context interface{}, mgoDB *db.DB, graphDB *cayley.Handle, viewParams *ViewParams) (*Result, error) {
	log.Dev(context, "Execute", "Started : Name[%s]", viewParams.ViewName)

	// Get the view.
	v, err := view.GetByName(context, mgoDB, viewParams.ViewName)
	if err != nil {
		log.Error(context, "Execute", err, "Completed")
		return errResult(err), err
	}

	// Validate the start type.
	if err := validateStartType(context, mgoDB, v); err != nil {
		log.Error(context, "Execute", err, "Completed")
		return errResult(err), err
	}

	// Translate the view path into a graph query path.
	graphPath, err := viewPathToGraphPath(v, viewParams.ItemKey, graphDB)
	if err != nil {
		log.Error(context, "Execute", err, "Completed")
		return errResult(err), err
	}

	// Retrieve the item IDs for the view.
	ids, err := viewIDs(v, graphPath, graphDB)
	if err != nil {
		log.Error(context, "Execute", err, "Completed")
		return errResult(err), err
	}

	// Persist the items in the view, if an output Collection is provided.
	if viewParams.ResultsCollection != "" {
		if err := viewSave(context, mgoDB, v, viewParams, ids); err != nil {
			log.Error(context, "Execute", err, "Completed")
			return errResult(err), err
		}
		result := Result{
			Results: bson.M{"number_of_results": len(ids)},
		}
		return &result, nil
	}

	// Otherwise, gather the items in the view.
	items, err := viewItems(context, mgoDB, v, ids)
	if err != nil {
		log.Error(context, "Execute", err, "Completed")
		return errResult(err), err
	}
	result := Result{
		Results: items,
	}

	log.Dev(context, "Execute", "Completed")
	return &result, nil
}

//==============================================================================

// validateStartType verifies the start type of a view path.
func validateStartType(context interface{}, db *db.DB, v *view.View) error {

	// Extract the first level relationship predicate.
	var firstRel string
	var firstDir string
	for _, segment := range v.Path {
		if segment.Level == 1 {
			firstRel = segment.Predicate
			firstDir = segment.Direction
		}
	}

	// Get the relationship metadata.
	rel, err := relationship.GetByPredicate(context, db, firstRel)
	if err != nil {
		return err
	}

	// Get the relevant item types based on the direction of the
	// first relationship in the path.
	var itemTypes []string
	switch firstDir {
	case outString:
		itemTypes = rel.SubjectTypes
	case inString:
		itemTypes = rel.ObjectTypes
	}

	// Validate the starting type provided in the view.
	for _, itemType := range itemTypes {
		if itemType == v.StartType {
			return nil
		}
	}

	return fmt.Errorf("Start type %s does not match relationship subject types %v", v.StartType, itemTypes)
}

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
			case inString:
				graphPath = cayley.StartPath(graphDB, quad.String(key)).In(quad.String(segment.Predicate))
			case outString:
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
		case inString:
			graphPath = graphPath.Clone().In(quad.String(segment.Predicate))
		case outString:
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

	// Determine the buffer limit that will be used for saving this view.
	if viewParams.BufferLimit != 0 {
		bufferLimit = viewParams.BufferLimit
	}

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
