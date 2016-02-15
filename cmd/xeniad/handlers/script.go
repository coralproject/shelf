package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/coralproject/xenia/pkg/script"

	"github.com/ardanlabs/kit/db"
	"github.com/ardanlabs/kit/web/app"
)

// scriptHandle maintains the set of handlers for the script api.
type scriptHandle struct{}

// Script fronts the access to the script service functionality.
var Script scriptHandle

//==============================================================================

// List returns all the existing scripts in the system.
// 200 Success, 404 Not Found, 500 Internal
func (scriptHandle) List(c *app.Context) error {
	scrs, err := script.GetScripts(c.SessionID, c.Ctx["DB"].(*db.DB), nil)
	if err != nil {
		if err == script.ErrNotFound {
			err = app.ErrNotFound
		}
		return err
	}

	c.Respond(scrs, http.StatusOK)
	return nil
}

// Retrieve returns the specified script from the system.
// 200 Success, 400 Bad Request, 404 Not Found, 500 Internal
func (scriptHandle) Retrieve(c *app.Context) error {
	scr, err := script.GetByName(c.SessionID, c.Ctx["DB"].(*db.DB), c.Params["name"])
	if err != nil {
		if err == script.ErrNotFound {
			err = app.ErrNotFound
		}
		return err
	}

	c.Respond(scr, http.StatusOK)
	return nil
}

//==============================================================================

// Upsert inserts or updates the posted Script document into the database.
// 204 SuccessNoContent, 400 Bad Request, 404 Not Found, 500 Internal
func (scriptHandle) Upsert(c *app.Context) error {
	var scr *script.Script
	if err := json.NewDecoder(c.Request.Body).Decode(&scr); err != nil {
		return err
	}

	if err := script.Upsert(c.SessionID, c.Ctx["DB"].(*db.DB), scr); err != nil {
		return err
	}

	c.Respond(nil, http.StatusNoContent)
	return nil
}

//==============================================================================

// Delete removes the specified Script from the system.
// 200 Success, 400 Bad Request, 404 Not Found, 500 Internal
func (scriptHandle) Delete(c *app.Context) error {
	if err := script.Delete(c.SessionID, c.Ctx["DB"].(*db.DB), c.Params["name"]); err != nil {
		if err == script.ErrNotFound {
			err = app.ErrNotFound
		}
		return err
	}

	c.Respond(nil, http.StatusNoContent)
	return nil
}
