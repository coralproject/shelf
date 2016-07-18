// Package handlers contains the handler logic for processing requests.
package handlers

import (
	"net/http"

	"github.com/ardanlabs/kit/web/app"

	"github.com/coralproject/xenia/internal/item"
)

type itemHandle struct{}

// Version fronts the access to the ver service functionality.
var Item itemHandle

//==============================================================================

// Types returns all the item type and relationships information
//  currently registered in this system
func (itemHandle) Types(c *app.Context) error {
	c.Respond(item.Types, http.StatusOK)
	return nil
}
