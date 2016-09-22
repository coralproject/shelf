// Package handlers contains the handler logic for processing requests.
package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ardanlabs/kit/db"
	"github.com/ardanlabs/kit/web/app"
	"github.com/coralproject/shelf/internal/sponge/item"
)

// dataHandle maintains the set of handlers for the data api, which is responsible
// for all requests sending/requesting unstructured data not yet in item form.
type dataHandle struct{}

// Data fronts the access to the data service functionality.
var Data dataHandle

//==============================================================================

// defaultVersion is set to item.Version when no version is provided
const defaultVersion = 1

//==============================================================================

// Upsert receives POSTed data, itemizes it then Upserts it via the item service
// 204 SuccessNoContent, 400 Bad Request, 404 Not Found, 500 Internal
func (dataHandle) Upsert(c *app.Context) error {

	// Unmarshall the data packet from the Request Body.
	var dat map[string]interface{}
	if err := json.NewDecoder(c.Request.Body).Decode(&dat); err != nil {
		return err
	}

	// Create a new item with known Type, Version and Data.
	it := item.Item{
		Type:    c.Params["type"],
		Version: defaultVersion,
		Data:    dat,
	}

	// Item.ID must be inferred from the source_id in the data.
	if err := it.InferIDFromData(); err != nil {
		return err
	}

	// Upsert the item.
	if err := item.Upsert(c.SessionID, c.Ctx["DB"].(*db.DB), &it); err != nil {
		return err
	}

	// Respond with no content success.
	c.Respond(nil, http.StatusNoContent)
	return nil
}
