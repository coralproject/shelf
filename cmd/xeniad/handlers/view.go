// Package handlers contains the handler logic for processing requests.
package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ardanlabs/kit/web"
	"github.com/coralproject/shelf/internal/platform/db"
	"github.com/coralproject/shelf/internal/wire/view"
)

// viewHandle maintains the set of handlers for the view api.
type viewHandle struct{}

// View fronts the access to the view service functionality.
var View viewHandle

//==============================================================================

// List returns all the existing views in the system.
// 200 Success, 404 Not Found, 500 Internal
func (viewHandle) List(c *web.Context) error {
	views, err := view.GetAll(c.SessionID, c.Ctx["DB"].(*db.DB))
	if err != nil {
		if err == view.ErrNotFound {
			err = web.ErrNotFound
		}
		return err
	}

	c.Respond(views, http.StatusOK)
	return nil
}

// Retrieve returns the specified View from the system.
// 200 Success, 400 Bad Request, 404 Not Found, 500 Internal
func (viewHandle) Retrieve(c *web.Context) error {
	v, err := view.GetByName(c.SessionID, c.Ctx["DB"].(*db.DB), c.Params["name"])
	if err != nil {
		if err == view.ErrNotFound {
			err = web.ErrNotFound
		}
		return err
	}

	c.Respond(v, http.StatusOK)
	return nil
}

//==============================================================================

// Upsert inserts or updates the posted View document into the database.
// 204 SuccessNoContent, 400 Bad Request, 404 Not Found, 500 Internal
func (viewHandle) Upsert(c *web.Context) error {
	var v view.View
	if err := json.NewDecoder(c.Request.Body).Decode(&v); err != nil {
		return err
	}

	if err := view.Upsert(c.SessionID, c.Ctx["DB"].(*db.DB), &v); err != nil {
		return err
	}

	c.Respond(nil, http.StatusNoContent)
	return nil
}

//==============================================================================

// Delete removes the specified View from the system.
// 200 Success, 400 Bad Request, 404 Not Found, 500 Internal
func (viewHandle) Delete(c *web.Context) error {
	if err := view.Delete(c.SessionID, c.Ctx["DB"].(*db.DB), c.Params["name"]); err != nil {
		if err == view.ErrNotFound {
			err = web.ErrNotFound
		}
		return err
	}

	c.Respond(nil, http.StatusNoContent)
	return nil
}
