// Package handlers contains the handler logic for processing requests.
package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ardanlabs/kit/db"
	"github.com/ardanlabs/kit/web/app"
	"github.com/coralproject/shelf/internal/sponge"
	"github.com/coralproject/shelf/internal/sponge/item"
)

// dataHandle maintains the set of handlers for the data api.
type dataHandle struct{}

// Data fronts the access to the data service functionality.
var Data dataHandle

//==============================================================================

// Upsert receives POSTed data, itemizes it then Upserts it via the item service
// 204 SuccessNoContent, 400 Bad Request, 404 Not Found, 500 Internal
func (dataHandle) Upsert(c *app.Context) error {

	var da sponge.Data
	if err := json.NewDecoder(c.Request.Body).Decode(&da); err != nil {
		return err
	}

	ty := c.Params["type"]

	// Pass the version from the type config for now. TODO: Possibly add version
	// option in a separate endpoint in the future.
	var it item.Item
	it, err := sponge.Itemize(c.SessionID, c.Ctx["DB"].(*db.DB), ty, 1, da)
	if err != nil {
		return err
	}

	if err := item.Upsert(c.SessionID, c.Ctx["DB"].(*db.DB), &it); err != nil {
		return err
	}

	c.Respond(it, http.StatusNoContent)
	return nil
}
