// Package wire provides support for generating views.
package wire

import (
	"gopkg.in/mgo.v2/bson"

	"github.com/ardanlabs/kit/db"
	"github.com/ardanlabs/kit/db/mongo"
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
	ViewName          string
	ItemKey           string
	ResultsCollection string
}

//==============================================================================

// Execute executes a graph query to generate the specified view.
func Execute(context interface{}, mgoDB *db.DB, mgoCfg mongo.Config, graphDB *cayley.Handle, viewParams *ViewParams) (*Result, error) {
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
		if err := viewSave(context, mgoCfg, v, viewParams, ids); err != nil {
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
