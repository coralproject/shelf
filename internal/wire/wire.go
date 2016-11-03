// Package wire provides support for generating views.
package wire

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/ardanlabs/kit/log"
	"github.com/cayleygraph/cayley"
	"github.com/cayleygraph/cayley/graph"
	"github.com/cayleygraph/cayley/graph/path"
	"github.com/cayleygraph/cayley/quad"
	"github.com/coralproject/shelf/internal/platform/db"
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

	// Retrieve the item IDs for the view along with any related item IDs
	// that should be embedded in view items (i.e., "embeds").
	ids, embeds, err := viewIDs(v, graphPath, viewParams.ItemKey, graphDB)
	if err != nil {
		log.Error(context, "Execute", err, "Completed")
		return errResult(err), err
	}

	// Persist the items in the view, if an output Collection is provided.
	if viewParams.ResultsCollection != "" {
		if err := viewSave(context, mgoDB, v, viewParams, ids, embeds); err != nil {
			log.Error(context, "Execute", err, "Completed")
			return errResult(err), err
		}
		result := Result{
			Results: bson.M{"number_of_results": len(ids)},
		}
		return &result, nil
	}

	// Otherwise, gather the items in the view.
	items, err := viewItems(context, mgoDB, v, ids, embeds)
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

	// Loop over paths to validate the start type for each.
PathLoop:
	for _, path := range v.Paths {

		// Declare variables to track the first level relationship predicate.
		var firstRel string
		var firstDir string

		// Extract the first level relationship predicate.
		for _, segment := range path.Segments {
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
				continue PathLoop
			}
		}

		return fmt.Errorf("Start type %s does not match relationship subject types %v", v.StartType, itemTypes)
	}

	return nil
}

// viewPathToGraphPath translates the path in a view into a "path"
// utilized in graph queries.
func viewPathToGraphPath(v *view.View, key string, graphDB *cayley.Handle) (*path.Path, error) {

	// outputPath is the final tranlated graph path.
	var outputPath *path.Path

	// Loop over the paths in the view translating the metadata.
	for idx, pth := range v.Paths {

		// We create an alias prefix for tags, so we can track which
		// path a tag is in.
		alias := strconv.Itoa(idx+1) + "_"

		// Sort the view Path value.
		sort.Sort(pth.Segments)

		// graphPath will contain the entire strict graph path.
		var graphPath *path.Path

		// subPaths will contain each sub path of the full graph path,
		// as a separate graph path.
		var subPaths []path.Path

		// Loop over the path segments translating the path.
		level := 1
		for _, segment := range pth.Segments {

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
					graphPath = graphPath.Clone().Tag(alias + segment.Tag)
				}

				// Track this as a subpath.
				subPaths = append(subPaths, *graphPath.Clone())

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
				graphPath = graphPath.Clone().Tag(alias + segment.Tag)
			}

			// Add this as a subpath.
			subPaths = append(subPaths, *graphPath.Clone())

			level++
		}

		// If we are forcing a strict path, return only the resulting or
		// tagged items along the full path.
		if pth.StrictPath {
			if outputPath == nil {
				outputPath = graphPath
				continue
			}
			outputPath = outputPath.Clone().Or(graphPath)
			continue
		}

		// Otherwise add all the subpaths to the output path.
		for _, subPath := range subPaths {
			if outputPath == nil {
				addedPath := &subPath
				outputPath = addedPath.Clone()
				continue
			}
			outputPath = outputPath.Clone().Or(&subPath)
		}
	}

	return outputPath, nil
}

// embeddedRel includes information needed to embed an indication of a relationship
// into an item in a view.
type embeddedRel struct {
	itemID     string
	predicate  string
	embeddedID string
}

// embeddedRels is a slice of EmbeddedRel.
type embeddedRels []embeddedRel

// relList contains one or more related IDs.
type relList []string

// viewIDs retrieves the item IDs associated with the view.
func viewIDs(v *view.View, path *path.Path, key string, graphDB *cayley.Handle) ([]string, embeddedRels, error) {

	// Build the Cayley iterator.
	it := path.BuildIterator()
	it, _ = it.Optimize()
	defer it.Close()

	// tagOrder will allow us to look up the ordering of a tag or
	// the tag corresponding to an order on demand.
	tagOrder := make(map[string]string)

	// Extract any tags and the ordering in the View value.
	var viewTags []string
	for idx, pth := range v.Paths {
		alias := strconv.Itoa(idx+1) + "_"
		for _, segment := range pth.Segments {
			if segment.Tag != "" {
				viewTags = append(viewTags, alias+segment.Tag)
				tagOrder[alias+segment.Tag] = alias + strconv.Itoa(segment.Level)
				tagOrder[alias+strconv.Itoa(segment.Level)] = alias + segment.Tag

			}
		}
	}

	// Retrieve the end path and tagged item IDs.
	var ids []string
	var embeds embeddedRels
	for it.Next() {

		// Tag the results.
		resultTags := make(map[string]graph.Value)
		it.TagResults(resultTags)

		// Extract the tagged item IDs.
		taggedIDs := make(map[string]relList)
		for _, tag := range viewTags {
			if t, ok := resultTags[tag]; ok {

				// Append the view item ID.
				ids = append(ids, quad.NativeOf(graphDB.NameOf(t)).(string))

				// Add the tagged ID to the tagged map for embedded
				// relationship extraction.
				current, ok := taggedIDs[tag]
				if !ok {
					taggedIDs[tag] = []string{quad.NativeOf(graphDB.NameOf(t)).(string)}
					continue
				}
				updated := append(current, quad.NativeOf(graphDB.NameOf(t)).(string))
				taggedIDs[tag] = updated
			}
		}

		// Extract any IDs that need to be embedded in view items.
		embed, err := extractEmbeddedRels(v, taggedIDs, tagOrder, key)
		if err != nil {
			return ids, embeds, err
		}
		embeds = append(embeds, embed...)
	}
	if it.Err() != nil {
		return ids, embeds, it.Err()
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

	// Add root item.
	if v.ReturnRoot == true {
		ids = append(ids, key)
	}

	return ids, embeds, nil
}

// extractEmbeddedRel extracts a relationship that needs to be embedded on a view item.
func extractEmbeddedRels(v *view.View, taggedIDs map[string]relList, tagOrder map[string]string, key string) (embeddedRels, error) {

	// embeds will contain the indication of the embedded relationships.
	var embeds embeddedRels

	// Loop over the taggedIDs determining the embedding based on the ordering
	// of the view path.
	for tag, rels := range taggedIDs {

		// Get the order of the tagged ID.
		a := tagOrder[tag]
		aliasOrder := strings.Split(a, "_")
		order, err := strconv.Atoi(aliasOrder[len(aliasOrder)-1])
		if err != nil {
			return nil, err
		}
		if len(aliasOrder) <= 1 {
			return nil, fmt.Errorf("Unexpected aliased tag order")
		}

		// If the order is greater than 1, we should embed at the next tagged
		// level up from the item.
		var embedTag string
		var ok bool
		if order > 1 {

			// Get the tag of the item in which the embed should be placed.
			for order > 1 {
				alias := aliasOrder[0] + "_" + strconv.Itoa(order-1)
				embedTag, ok = tagOrder[alias]
				if ok {
					break
				}
			}
		}

		// Extract the non-alias tag.
		aliasTag := strings.Split(tag, "_")
		nonAliasTag := strings.Join(aliasTag[1:], "_")

		// If we are embedding in a non-root item, get that item ID and form
		// the EmbeddedRel value.
		if embedTag != "" {
			for tagCandidate, embedItem := range taggedIDs {
				if tagCandidate == embedTag {
					if len(embedItem) != 1 {
						return nil, fmt.Errorf("Unexpected number of IDs")
					}
					for _, rel := range rels {
						embed := embeddedRel{
							itemID:     embedItem[0],
							predicate:  nonAliasTag,
							embeddedID: rel,
						}
						embeds = append(embeds, embed)
					}
				}
			}
			continue
		}

		// Otherwise, embed in the root item.
		for _, rel := range rels {
			embed := embeddedRel{
				itemID:     key,
				predicate:  nonAliasTag,
				embeddedID: rel,
			}
			embeds = append(embeds, embed)
		}
	}

	return embeds, nil
}

// viewSave retrieve items for a view and saves those items to a new collection.
func viewSave(context interface{}, mgoDB *db.DB, v *view.View, viewParams *ViewParams, ids []string, embeds embeddedRels) error {

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

	// Group the embedded relationships by item and predicate/tag.
	embedByItem, err := groupEmbeds(embeds)
	if err != nil {
		return err
	}

	// Iterate over the view items.
	var queuedDocs int
	var result item.Item
	for results.Next(&result) {

		// Get the related IDs to embed.
		itemEmbeds, ok := embedByItem[result.ID]
		if ok {

			// Put the related IDs in the item.
			relMap := make(map[string]interface{})
			for k, v := range itemEmbeds {
				relMap[k] = v
			}
			result.Related = relMap
		}

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
func viewItems(context interface{}, db *db.DB, v *view.View, ids []string, embeds embeddedRels) ([]bson.M, error) {

	// Form the query.
	var results []item.Item
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

	// Group the embedded relationships by item and predicate/tag.
	embedByItem, err := groupEmbeds(embeds)
	if err != nil {
		return nil, err
	}

	// Embed any related item IDs in the returned items.
	var output []bson.M
	if len(embedByItem) > 0 {
		for _, result := range results {

			// Get the respective IDs to embed.
			itemEmbeds, ok := embedByItem[result.ID]
			if ok {
				relMap := make(map[string]interface{})
				for k, v := range itemEmbeds {
					relMap[k] = v
				}
				result.Related = relMap
			}

			// Convert to bson.M for output.
			itemBSON := bson.M{
				"item_id":    result.ID,
				"type":       result.Type,
				"version":    result.Version,
				"data":       result.Data,
				"created_at": result.CreatedAt,
				"updated_at": result.UpdatedAt,
				"related":    result.Related,
			}
			output = append(output, itemBSON)
		}
	}

	return output, nil
}

// predicateEmbeds includes slices of related item IDs grouped by predicate/tag.
type predicateEmbeds map[string]relList

// groupEmbeds groups embeddedRel values by item.
func groupEmbeds(embeds embeddedRels) (map[string]predicateEmbeds, error) {
	embedsOut := make(map[string]predicateEmbeds)

	// Loop over embeds to group embeds.
	for _, embed := range embeds {

		// Get the map of embededed ID for a particular item ID,
		// if it exists.  If it does not exist create the map.
		predMap, ok := embedsOut[embed.itemID]
		if !ok {
			pe := map[string]relList{
				embed.predicate: relList{embed.embeddedID},
			}
			embedsOut[embed.itemID] = pe
			continue
		}

		// Get the current IDs corresponding to this predicate/tag.
		current, ok := predMap[embed.predicate]
		if !ok {
			rl := relList{embed.embeddedID}
			embedsOut[embed.itemID][embed.predicate] = rl
			continue
		}

		// Update the current IDs.
		updated := append(current, embed.embeddedID)

		// Remove duplicates.
		found := make(map[string]bool)
		j := 0
		for i, x := range updated {
			if !found[x] {
				found[x] = true
				updated[j] = updated[i]
				j++
			}
		}
		updated = updated[:j]

		// Update the output.
		embedsOut[embed.itemID][embed.predicate] = updated
	}

	return embedsOut, nil
}
