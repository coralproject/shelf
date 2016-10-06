// Package handlers contains the handler logic for processing requests.
package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/ardanlabs/kit/web"
	"github.com/coralproject/shelf/internal/platform/db"
	"github.com/coralproject/shelf/internal/sponge"
	"github.com/coralproject/shelf/internal/sponge/item"
)

// itemHandle maintains the set of handlers for theitem api.
type itemHandle struct{}

// Item fronts the access to the item service functionality.
var Item itemHandle

//==============================================================================

// Retrieve returns the items, specified by IDs, from the system.
// 200 Success, 400 Bad Request, 404 Not Found, 500 Internal
func (itemHandle) Retrieve(c *web.Context) error {
	var items []item.Item
	ids := strings.Split(c.Params["id"], ",")
	items, err := item.GetByIDs(c.SessionID, c.Ctx["DB"].(*db.DB), ids)
	if err != nil {
		if err == item.ErrNotFound {
			err = web.ErrNotFound
		}
		return err
	}

	c.Respond(items, http.StatusOK)
	return nil
}

//==============================================================================

// Import inserts or updates the posted Item document into the items collection
// and adds/removes any necessary quads to/from the relationship graph.
// 204 SuccessNoContent, 400 Bad Request, 404 Not Found, 500 Internal
func (itemHandle) Import(c *web.Context) error {

	// Decode the item.
	var itm item.Item
	if err := json.NewDecoder(c.Request.Body).Decode(&itm); err != nil {
		return err
	}

	db := c.Ctx["DB"].(*db.DB)

	graphHandle, err := db.GraphHandle(c.SessionID)
	if err != nil {
		return err
	}

	// Upsert the item into the items collection and add/remove necessary
	// quads to/from the graph.
	if err := sponge.Import(c.SessionID, db, graphHandle, &itm); err != nil {
		return err
	}

	c.Respond(itm, http.StatusOK)
	return nil
}

//==============================================================================

// Remove removes the specified Item from the items collection and removes any
// relevant quads from the graph database.
// 200 Success, 400 Bad Request, 404 Not Found, 500 Internal
func (itemHandle) Remove(c *web.Context) error {
	db := c.Ctx["DB"].(*db.DB)

	graphHandle, err := db.GraphHandle(c.SessionID)
	if err != nil {
		return err
	}

	if err := sponge.Remove(c.SessionID, db, graphHandle, c.Params["id"]); err != nil {
		return err
	}

	c.Respond(nil, http.StatusNoContent)
	return nil
}
