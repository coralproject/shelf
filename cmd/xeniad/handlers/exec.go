// Package handlers contains the handler logic for processing requests.
package handlers

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/ardanlabs/kit/db"
	"github.com/ardanlabs/kit/web/app"
	"github.com/cayleygraph/cayley"
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
func (execHandle) Name(c *app.Context) error {
	set, err := query.GetByName(c.SessionID, c.Ctx["DB"].(*db.DB), c.Params["name"])
	if err != nil {
		if err == query.ErrNotFound {
			err = app.ErrNotFound
		}
		return err
	}

	var vars map[string]string

	return execute(c, set, vars)
}

// NameOnView runs the specified Set on a view and returns results.
// 200 Success, 400 Bad Request, 404 Not Found, 500 Internal
func (execHandle) NameOnView(c *app.Context) error {
	set, err := query.GetByName(c.SessionID, c.Ctx["DB"].(*db.DB), c.Params["name"])
	if err != nil {
		if err == query.ErrNotFound {
			err = app.ErrNotFound
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
func (execHandle) Custom(c *app.Context) error {
	var set *query.Set
	if err := json.NewDecoder(c.Request.Body).Decode(&set); err != nil {
		return err
	}

	var vars map[string]string

	return execute(c, set, vars)
}

//==============================================================================

// execute takes a context and Set and executes the set returning
// any possible response.
func execute(c *app.Context, set *query.Set, vars map[string]string) error {

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
	var result *query.Result

	// If the query is a query on a view, provide the graph handle.
	// Otherwise just provide the db.DB value.
	for _, q := range set.Queries {
		if q.Collection == "view" {
			result = xenia.Exec(c.SessionID, c.Ctx["DB"].(*db.DB), c.Ctx["Graph"].(*cayley.Handle), set, vars)
			c.Respond(result, http.StatusOK)
			return nil
		}
	}
	result = xenia.Exec(c.SessionID, c.Ctx["DB"].(*db.DB), nil, set, vars)

	c.Respond(result, http.StatusOK)
	return nil
}
