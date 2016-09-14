// Package handlers contains the handler logic for processing requests.
package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ardanlabs/kit/db"
	"github.com/ardanlabs/kit/web/app"
	"github.com/coralproject/shelf/internal/wire/pattern"
)

// patternHandle maintains the set of handlers for the pattern api.
type patternHandle struct{}

// Pattern fronts the access to the pattern service functionality.
var Pattern patternHandle

//==============================================================================

// List returns all the existing patterns in the system.
// 200 Success, 404 Not Found, 500 Internal
func (patternHandle) List(c *app.Context) error {
	ps, err := pattern.GetAll(c.SessionID, c.Ctx["DB"].(*db.DB))
	if err != nil {
		if err == pattern.ErrNotFound {
			err = app.ErrNotFound
		}
		return err
	}

	c.Respond(ps, http.StatusOK)
	return nil
}

// Retrieve returns the specified Pattern from the system.
// 200 Success, 400 Bad Request, 404 Not Found, 500 Internal
func (patternHandle) Retrieve(c *app.Context) error {
	p, err := pattern.GetByType(c.SessionID, c.Ctx["DB"].(*db.DB), c.Params["type"])
	if err != nil {
		if err == pattern.ErrNotFound {
			err = app.ErrNotFound
		}
		return err
	}

	c.Respond(p, http.StatusOK)
	return nil
}

//==============================================================================

// Upsert inserts or updates the posted Pattern document into the database.
// 204 SuccessNoContent, 400 Bad Request, 404 Not Found, 500 Internal
func (patternHandle) Upsert(c *app.Context) error {
	var p pattern.Pattern
	if err := json.NewDecoder(c.Request.Body).Decode(&p); err != nil {
		return err
	}

	if err := pattern.Upsert(c.SessionID, c.Ctx["DB"].(*db.DB), &p); err != nil {
		return err
	}

	c.Respond(nil, http.StatusNoContent)
	return nil
}

//==============================================================================

// Delete removes the specified Pattern from the system.
// 200 Success, 400 Bad Request, 404 Not Found, 500 Internal
func (patternHandle) Delete(c *app.Context) error {
	if err := pattern.Delete(c.SessionID, c.Ctx["DB"].(*db.DB), c.Params["type"]); err != nil {
		if err == pattern.ErrNotFound {
			err = app.ErrNotFound
		}
		return err
	}

	c.Respond(nil, http.StatusNoContent)
	return nil
}
