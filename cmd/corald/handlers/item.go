package handlers

import (
	"net/http"

	"github.com/ardanlabs/kit/web/app"
	"github.com/coralproject/shelf/internal/coral"
)

// SpongeHost holds temporarily the sponge API url
const SpongeHost = "http://127.0.0.1:3001"

// itemHandle maintains the set of handlers for the form api.
type itemHandle struct{}

// Item fronts the access to the comment service functionality.
var Item itemHandle

//==============================================================================

// Retrieve returns the item by item_id in the system.
// 200 Success, 404 Not Found, 500 Internal
func (itemHandle) Retrieve(c *app.Context) error {

	URL := SpongeHost + "/1.0/item/" + c.Params["item_id"]

	result, err := coral.DoRequest(c, "POST", URL, c.Request.Body)
	if err != nil {
		return err
	}

	c.Respond(result, http.StatusOK)

	return nil
}

//==============================================================================

// Upsert insert or update an item in the system.
// 200 Success, 404 Not Found, 500 Internal
func (itemHandle) Upsert(c *app.Context) error {

	URL := SpongeHost + "/1.0/item"

	result, err := coral.DoRequest(c, "PUT", URL, c.Request.Body)
	if err != nil {
		return err
	}

	c.Respond(result, http.StatusOK)

	return nil
}
