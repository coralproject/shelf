// Package wire provides support for generating views.
package wire

import (
	"time"

	"github.com/ardanlabs/kit/db"
	"github.com/ardanlabs/kit/log"
	"github.com/cayleygraph/cayley"
	"github.com/coralproject/shelf/internal/wire/view"
	"gopkg.in/mgo.v2/bson"
)

// Result represents what a user will receive after generating a view.
type Result struct {
	Name          string    `json:"name"`
	Updated       time.Time `json:"last_updated,omitempty"`
	Results       int       `json:"number_of_results"`
	CollectionOut string    `json:"collection_out,omitempty"`
	CollectionIn  string    `json:"collection_in"`
	Items         []bson.M  `json:"items,omitempty"`
	Error         string    `json:"error,omitempty"`
}

// errResult returns a Result value with an error message.
func errResult(v *view.View, err error) *Result {
	result := Result{
		Name:  v.Name,
		Error: err.Error(),
	}

	return &result
}

// ViewParams represents how the View will be generated and persisted.
type ViewParams struct {
	StartID       string
	StartType     string
	Persist       bool
	CollectionOut string
	CollectionIn  string
	GraphHandle   *cayley.Handle
}

//==============================================================================

// Generate generates the specified view.
func Generate(context interface{}, db *db.DB, v *view.View, viewParams *ViewParams) *Result {
	log.Dev(context, "Generate", "Started : Name[%s]", v.Name)

	// Validate the view that is provided.
	if err := v.Validate(); err != nil {
		log.Error(context, "Generate", err, "Completed")
		return errResult(v, err)
	}

	// Validate the start type.
	if err := verifyStartType(context, db, v, viewParams); err != nil {
		log.Error(context, "Generate", err, "Completed")
		return errResult(v, err)
	}

	// Translate the view path into a graph query path.
	graphPath, err := viewPathToGraphPath(v, viewParams)
	if err != nil {
		log.Error(context, "Generate", err, "Completed")
		return errResult(v, err)
	}

	// Retrieve the item IDs for the view.
	ids, err := viewIDs(v, graphPath, viewParams)
	if err != nil {
		log.Error(context, "Generate", err, "Completed")
		return errResult(v, err)
	}

	// Save the items out to a collection if enabled.
	if viewParams.Persist {

		// Output to the collection.
		if err := viewSave(context, db, viewParams, ids); err != nil {
			log.Error(context, "Generate", err, "Completed")
			return errResult(v, err)
		}

		// Output the results.
		result := Result{
			Name:          v.Name,
			Updated:       time.Now(),
			Results:       len(ids),
			CollectionIn:  viewParams.CollectionIn,
			CollectionOut: viewParams.CollectionOut,
		}
		return &result
	}

	// Get the view items if it is not persisted.
	items, err := viewItems(context, db, viewParams, ids)
	if err != nil {
		log.Error(context, "Generate", err, "Completed")
		return errResult(v, err)
	}

	// Form the result.
	result := Result{
		Name:         v.Name,
		Updated:      time.Now(),
		Results:      len(ids),
		CollectionIn: viewParams.CollectionIn,
		Items:        items,
	}
	return &result
}
