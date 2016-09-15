// Package handlers contains the handler logic for processing requests.
package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/ardanlabs/kit/db"
	"github.com/ardanlabs/kit/web/app"
	"github.com/coralproject/shelf/internal/sponge/item"
	"github.com/coralproject/shelf/internal/wire"
)

// itemHandle maintains the set of handlers for theitem api.
type itemHandle struct{}

// Item fronts the access to the item service functionality.
var Item itemHandle

//==============================================================================

// Retrieve returns the items, specified by IDs, from the system.
// 200 Success, 400 Bad Request, 404 Not Found, 500 Internal
func (itemHandle) Retrieve(c *app.Context) error {
	var items []item.Item
	ids := strings.Split(c.Params["id"], ",")
	items, err := item.GetByIDs(c.SessionID, c.Ctx["DB"].(*db.DB), ids)
	if err != nil {
		if err == item.ErrNotFound {
			err = app.ErrNotFound
		}
		return err
	}

	c.Respond(items, http.StatusOK)
	return nil
}

//==============================================================================

// Upsert inserts or updates the posted Item document into the database.
// 204 SuccessNoContent, 400 Bad Request, 404 Not Found, 500 Internal
func (itemHandle) Upsert(c *app.Context) error {

	// Decode the item.
	var it item.Item
	if err := json.NewDecoder(c.Request.Body).Decode(&it); err != nil {
		return err
	}

	// Add the item to the items collection.
	if err := item.Upsert(c.SessionID, c.Ctx["DB"].(*db.DB), &it); err != nil {
		return err
	}

	// Prepare generic item data.
	itMap, ok := item.(map[string]interface{})
	if !ok {
		return fmt.Errorf("Could not convert item data to a map.")
	}

	// Infer relationships and add them to the graph.
	if err := wire.AddToGraph(c.SessionID, c.Ctx["DB"], c.Ctx["Graph"], itMap); err != nil {
		return err
	}

	c.Respond(nil, http.StatusNoContent)
	return nil
}

//==============================================================================

// Delete removes the specified Item from the system.
// 200 Success, 400 Bad Request, 404 Not Found, 500 Internal
func (itemHandle) Delete(c *app.Context) error {
	if err := item.Delete(c.SessionID, c.Ctx["DB"].(*db.DB), c.Params["id"]); err != nil {
		if err == item.ErrNotFound {
			err = app.ErrNotFound
		}
		return err
	}

	c.Respond(nil, http.StatusNoContent)
	return nil
}
