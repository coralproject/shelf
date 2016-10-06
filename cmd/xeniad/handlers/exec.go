// Package handlers contains the handler logic for processing requests.
package handlers

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/ardanlabs/kit/web"
	"github.com/coralproject/shelf/internal/platform/db"
	"github.com/coralproject/shelf/internal/xenia"
	"github.com/coralproject/shelf/internal/xenia/query"
)

// execHandle maintains the set of handlers for the exec api.
type execHandle struct{}

// Exec fronts the access to the exec service functionality.
var Exec execHandle

//==============================================================================

// Name runs the specified Set and returns results.
// 200 Success, 400 Bad Request, 404 Not Found, 500 Internal
func (execHandle) Name(c *web.Context) error {
	set, err := query.GetByName(c.SessionID, c.Ctx["DB"].(*db.DB), c.Params["name"])
	if err != nil {
		if err == query.ErrNotFound {
			err = web.ErrNotFound
		}
		return err
	}

	var vars map[string]string

	return execute(c, set, vars)
}

// NameOnView runs the specified Set on a view and returns results.
// 200 Success, 400 Bad Request, 404 Not Found, 500 Internal
func (execHandle) NameOnView(c *web.Context) error {
	set, err := query.GetByName(c.SessionID, c.Ctx["DB"].(*db.DB), c.Params["name"])
	if err != nil {
		if err == query.ErrNotFound {
			err = web.ErrNotFound
		}
		return err
	}

	vars := map[string]string{
		"view": c.Params["view"],
		"item": c.Params["item"],
	}

	return execute(c, set, vars)
}

// Custom runs the provided Set and return results.
// 200 Success, 400 Bad Request, 404 Not Found, 500 Internal
func (execHandle) Custom(c *web.Context) error {
	var set *query.Set
	if err := json.NewDecoder(c.Request.Body).Decode(&set); err != nil {
		return err
	}

	var vars map[string]string

	return execute(c, set, vars)
}

// CustomOnView runs the provided Set on a view and return results.
// 200 Success, 400 Bad Request, 404 Not Found, 500 Internal
func (execHandle) CustomOnView(c *web.Context) error {
	var set *query.Set
	if err := json.NewDecoder(c.Request.Body).Decode(&set); err != nil {
		return err
	}

	vars := map[string]string{
		"view": c.Params["view"],
		"item": c.Params["item"],
	}

	return execute(c, set, vars)
}

//==============================================================================

// execute takes a context and Set and executes the set returning
// any possible response.
func execute(c *web.Context, set *query.Set, vars map[string]string) error {

	// Parse the vars in the query string.
	if c.Request.URL.RawQuery != "" {
		if m, err := url.ParseQuery(c.Request.URL.RawQuery); err == nil {
			if vars == nil {
				vars = make(map[string]string)
			}
			for k, v := range m {
				vars[k] = v[0]
			}
		}
	}

	// Get the result.
	result := xenia.Exec(c.SessionID, c.Ctx["DB"].(*db.DB), set, vars)

	c.Respond(result, http.StatusOK)
	return nil
}
