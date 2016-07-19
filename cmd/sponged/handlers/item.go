// Package handlers contains the handler logic for processing requests.
package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ardanlabs/kit/db"
	"github.com/ardanlabs/kit/web/app"

	"github.com/coralproject/xenia/internal/item"
)

type itemHandle struct{}

// Version fronts the access to the ver service functionality.
var Item itemHandle

//==============================================================================

const (
	unrecognizedTypeError = "Unrecognized Type"
)

//==============================================================================

// Types returns all the item type and relationships information
//  currently registered in this system
func (itemHandle) Types(c *app.Context) error {
	c.Respond(item.Types, http.StatusOK)
	return nil
}

// Upsert receives a type and a data payload, validates it
//  and passes it on to item service code for upserting. TODO: auth
//  and other request specific actions will be handled here.
func (itemHandle) Upsert(c *app.Context) error {

	// get the type from the query string
	tn := c.Params["type"]

	// validate type
	_, ok := item.Types[tn]
	if !ok {
		c.Respond(unrecognizedTypeError, http.StatusNotFound)
		return nil
	}

	// get the data from the body
	data := make(map[string]interface{})
	if err := json.NewDecoder(c.Request.Body).Decode(&data); err != nil {
		c.Respond(err, http.StatusInternalServerError)
		return nil
	}

	// create an item from it
	i, err := item.Create(c, c.Ctx["DB"].(*db.DB), tn, 0, data)
	if err != nil {
		c.Respond(err, http.StatusInternalServerError)
		return nil
	}

	// upsert the item
	err = item.Upsert(c, c.Ctx["DB"].(*db.DB), &i)
	if err != nil {
		c.Respond(err, http.StatusInternalServerError)
		return nil
	}

	c.Respond(i, http.StatusOK)
	return nil
}
