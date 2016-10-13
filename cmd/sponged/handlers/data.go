// Package handlers contains the handler logic for processing requests.
package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ardanlabs/kit/web"
	"github.com/coralproject/shelf/internal/platform/db"
	"github.com/coralproject/shelf/internal/sponge"
	"github.com/coralproject/shelf/internal/sponge/item"
)

// dataHandle maintains the set of handlers for the data api, which is responsible
// for all requests sending/requesting unstructured data not yet in item form.
type dataHandle struct{}

// Data fronts the access to the data service functionality.
var Data dataHandle

//==============================================================================

// defaultVersion is set to item.Version when no version is provided.
const defaultVersion = 1

//==============================================================================

// Import receives POSTed data, itemizes it then imports it via the item API.
// 204 SuccessNoContent, 400 Bad Request, 404 Not Found, 500 Internal.
func (dataHandle) Import(c *web.Context) error {

	// Unmarshall the data packet from the Request Body.
	var dat map[string]interface{}
	if err := json.NewDecoder(c.Request.Body).Decode(&dat); err != nil {
		return err
	}

	// Create a new item with known Type, Version and Data.
	itm := item.Item{
		Type:    c.Params["type"],
		Version: defaultVersion,
		Data:    dat,
	}

	// Item.ID must be inferred from the source_id in the data.
	if err := itm.InferIDFromData(); err != nil {
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

	// Respond with no content success.
	c.Respond(itm, http.StatusOK)
	return nil
}
