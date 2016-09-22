package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ardanlabs/kit/web/app"
	"github.com/coralproject/shelf/internal/coral"
)

// SpongeHost holds temporarily the sponge API url
const SpongeHost = "http://127.0.0.1:3001"

// XeniaHost holds temporarily the xenia API url
const XeniaHost = "http://127.0.0.1:4000"

// itemHandle maintains the set of handlers for the form api.
type itemHandle struct{}

// Item fronts the access to the comment service functionality.
var Item itemHandle

//==============================================================================

// Retrieve runs the query_set in the Xenia service.
// 200 Success, 404 Not Found, 500 Internal
func (itemHandle) Retrieve(c *app.Context) error {

	URL := XeniaHost + "/1.0/exec/" + c.Params["query_set"]

	respond, err := coral.DoRequest(c, "GET", URL, c.Request.Body)
	if err != nil {
		return err
	}

	var results map[string]interface{}

	if err := json.Unmarshal(respond.Body, &results); err != nil {
		return err
	}

	c.Respond(results, http.StatusOK)

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
