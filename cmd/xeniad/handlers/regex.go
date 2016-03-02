package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/coralproject/xenia/pkg/regex"

	"github.com/ardanlabs/kit/db"
	"github.com/ardanlabs/kit/web/app"
)

// regexHandle maintains the set of handlers for the regex api.
type regexHandle struct{}

// Regex fronts the access to the regex service functionality.
var Regex regexHandle

//==============================================================================

// List returns all the existing regex in the system.
// 200 Success, 404 Not Found, 500 Internal
func (regexHandle) List(c *app.Context) error {
	rgxs, err := regex.GetAll(c.SessionID, c.Ctx["DB"].(*db.DB), nil)
	if err != nil {
		if err == regex.ErrNotFound {
			err = app.ErrNotFound
		}
		return err
	}

	c.Respond(rgxs, http.StatusOK)
	return nil
}

// Retrieve returns the specified regex from the system.
// 200 Success, 400 Bad Request, 404 Not Found, 500 Internal
func (regexHandle) Retrieve(c *app.Context) error {
	rgx, err := regex.GetByName(c.SessionID, c.Ctx["DB"].(*db.DB), c.Params["name"])
	if err != nil {
		if err == regex.ErrNotFound {
			err = app.ErrNotFound
		}
		return err
	}

	c.Respond(rgx, http.StatusOK)
	return nil
}

//==============================================================================

// Upsert inserts or updates the posted Regex document into the database.
// 204 SuccessNoContent, 400 Bad Request, 404 Not Found, 500 Internal
func (regexHandle) Upsert(c *app.Context) error {
	var rgx regex.Regex
	if err := json.NewDecoder(c.Request.Body).Decode(&rgx); err != nil {
		return err
	}

	if err := regex.Upsert(c.SessionID, c.Ctx["DB"].(*db.DB), rgx); err != nil {
		return err
	}

	c.Respond(nil, http.StatusNoContent)
	return nil
}

//==============================================================================

// Delete removes the specified Regex from the system.
// 200 Success, 400 Bad Request, 404 Not Found, 500 Internal
func (regexHandle) Delete(c *app.Context) error {
	if err := regex.Delete(c.SessionID, c.Ctx["DB"].(*db.DB), c.Params["name"]); err != nil {
		if err == regex.ErrNotFound {
			err = app.ErrNotFound
		}
		return err
	}

	c.Respond(nil, http.StatusNoContent)
	return nil
}
