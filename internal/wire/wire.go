// Package wire provides support for generating views.
package wire

import (
	"gopkg.in/mgo.v2/bson"

	"github.com/ardanlabs/kit/db"
	"github.com/ardanlabs/kit/log"
	"github.com/cayleygraph/cayley"
	"github.com/coralproject/shelf/internal/wire/view"
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
	ViewName string
	ItemKey  string
}

//==============================================================================

// Generate generates the specified view.
func Generate(context interface{}, mgoDB *db.DB, graphDB *cayley.Handle, viewParams *ViewParams) (*Result, error) {
	log.Dev(context, "Generate", "Started : Name[%s]", viewParams.ViewName)

	// Get the view.
	v, err := view.GetByName(context, mgoDB, viewParams.ViewName)
	if err != nil {
		log.Error(context, "Generate", err, "Completed")
		return errResult(err), err
	}

	// Validate the start type.
	if err := verifyStartType(context, mgoDB, v); err != nil {
		log.Error(context, "Generate", err, "Completed")
		return errResult(err), err
	}

	// Translate the view path into a graph query path.
	graphPath, err := viewPathToGraphPath(v, viewParams, graphDB)
	if err != nil {
		log.Error(context, "Generate", err, "Completed")
		return errResult(err), err
	}

	// Retrieve the item IDs for the view.
	ids, err := viewIDs(v, graphPath, graphDB)
	if err != nil {
		log.Error(context, "Generate", err, "Completed")
		return errResult(err), err
	}

	// Save the items out to a collection if enabled.
	if v.CollectionOutPrefix != "" {

		// Output to the collection.
		if err := viewSave(context, mgoDB, v, viewParams, ids); err != nil {
			log.Error(context, "Generate", err, "Completed")
			return errResult(err), err
		}

		// Output the results.
		result := Result{
			Results: bson.M{"number_of_results": len(ids)},
		}
		return &result, nil
	}

	// Get the view items if the view is not persisted.
	items, err := viewItems(context, mgoDB, v, ids)
	if err != nil {
		log.Error(context, "Generate", err, "Completed")
		return errResult(err), err
	}

	// Form the result.
	result := Result{
		Results: items,
	}

	log.Dev(context, "Generate", "Completed")
	return &result, nil
}
