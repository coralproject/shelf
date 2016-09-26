// Package handlers contains the handler logic for processing requests.
package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/ardanlabs/kit/db"
	"github.com/ardanlabs/kit/web/app"
	"github.com/cayleygraph/cayley"
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

	// See if the item already exists.
	if it.ID != "" {
		items, err := item.GetByIDs(c.SessionID, c.Ctx["DB"].(*db.DB), []string{it.ID})
		if err != nil {
			if err != item.ErrNotFound {
				return err
			}
		}

		// If we got more than one item, there is a problem.
		if len(items) > 1 {
			return fmt.Errorf("Expected one item, retrieved %d items.", len(items))
		}

		// Determine if the existing item is identical to the provided item.
		if len(items) == 1 {

			// If the item is identical, we don't have to do anything.
			if reflect.DeepEqual(items[0], it) {
				return nil
			}

			// If the item is not identical, remove the stale relationships by
			// preparing an item map.
			itMap := map[string]interface{}{
				"item_id": items[0].ID,
				"type":    items[0].Type,
				"version": items[0].Version,
				"data":    items[0].Data,
			}

			// Remove the corresponding relationships from the graph.
			if err := wire.RemoveFromGraph(c.SessionID, c.Ctx["DB"].(*db.DB), c.Ctx["Graph"].(*cayley.Handle), itMap); err != nil {
				return err
			}
		}
	}

	// Add the item to the items collection.
	if err := item.Upsert(c.SessionID, c.Ctx["DB"].(*db.DB), &it); err != nil {
		return err
	}

	// Prepare the generic item data map.
	itMap := map[string]interface{}{
		"item_id": it.ID,
		"type":    it.Type,
		"version": it.Version,
		"data":    it.Data,
	}

	// Infer relationships and add them to the graph.
	if err := wire.AddToGraph(c.SessionID, c.Ctx["DB"].(*db.DB), c.Ctx["Graph"].(*cayley.Handle), itMap); err != nil {
		return err
	}

	c.Respond(nil, http.StatusNoContent)
	return nil
}

//==============================================================================

// Delete removes the specified Item from the system.
// 200 Success, 400 Bad Request, 404 Not Found, 500 Internal
func (itemHandle) Delete(c *app.Context) error {

	// Get the item from the items collection.
	items, err := item.GetByIDs(c.SessionID, c.Ctx["DB"].(*db.DB), []string{c.Params["id"]})
	if err != nil {
		if err == item.ErrNotFound {
			err = app.ErrNotFound
		}
		return err
	}

	// Check to make sure we have the single item.
	if len(items) != 1 {
		return fmt.Errorf("Item ID corresponds to an unexpected number, %d, of items", len(items))
	}

	// Delete the item.
	if err := item.Delete(c.SessionID, c.Ctx["DB"].(*db.DB), c.Params["id"]); err != nil {
		if err == item.ErrNotFound {
			err = app.ErrNotFound
		}
		return err
	}

	// Prepare the item map data.
	itMap := map[string]interface{}{
		"item_id": items[0].ID,
		"type":    items[0].Type,
		"version": items[0].Version,
		"data":    items[0].Data,
	}

	// Remove the corresponding relationships from the graph.
	if err := wire.RemoveFromGraph(c.SessionID, c.Ctx["DB"].(*db.DB), c.Ctx["Graph"].(*cayley.Handle), itMap); err != nil {
		return err
	}

	c.Respond(nil, http.StatusNoContent)
	return nil
}
