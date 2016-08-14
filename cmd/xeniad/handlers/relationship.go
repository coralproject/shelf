// Package handlers contains the handler logic for processing requests.
package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ardanlabs/kit/db"
	"github.com/ardanlabs/kit/web/app"
	"github.com/coralproject/xenia/internal/wire/relationship"
)

// relationshipHandle maintains the set of handlers for the relationship api.
type relationshipHandle struct{}

// Relationship fronts the access to the relationship service functionality.
var Relationship relationshipHandle

//==============================================================================

// List returns all the existing relationships in the system.
// 200 Success, 404 Not Found, 500 Internal
func (relationshipHandle) List(c *app.Context) error {
	rels, err := relationship.GetAll(c.SessionID, c.Ctx["DB"].(*db.DB))
	if err != nil {
		if err == relationship.ErrNotFound {
			err = app.ErrNotFound
		}
		return err
	}

	c.Respond(rels, http.StatusOK)
	return nil
}

// Retrieve returns the specified Relationship from the system.
// 200 Success, 400 Bad Request, 404 Not Found, 500 Internal
func (relationshipHandle) Retrieve(c *app.Context) error {
	rel, err := relationship.GetByPredicate(c.SessionID, c.Ctx["DB"].(*db.DB), c.Params["predicate"])
	if err != nil {
		if err == relationship.ErrNotFound {
			err = app.ErrNotFound
		}
		return err
	}

	c.Respond(rel, http.StatusOK)
	return nil
}

//==============================================================================

// Upsert inserts or updates the posted Relationship document into the database.
// 204 SuccessNoContent, 400 Bad Request, 404 Not Found, 500 Internal
func (relationshipHandle) Upsert(c *app.Context) error {
	var rel relationship.Relationship
	if err := json.NewDecoder(c.Request.Body).Decode(&rel); err != nil {
		return err
	}

	if err := relationship.Upsert(c.SessionID, c.Ctx["DB"].(*db.DB), &rel); err != nil {
		return err
	}

	c.Respond(nil, http.StatusNoContent)
	return nil
}

//==============================================================================

// Delete removes the specified Relationship from the system.
// 200 Success, 400 Bad Request, 404 Not Found, 500 Internal
func (relationshipHandle) Delete(c *app.Context) error {
	if err := relationship.Delete(c.SessionID, c.Ctx["DB"].(*db.DB), c.Params["predicate"]); err != nil {
		if err == relationship.ErrNotFound {
			err = app.ErrNotFound
		}
		return err
	}

	c.Respond(nil, http.StatusNoContent)
	return nil
}
