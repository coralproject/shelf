package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ardanlabs/kit/cfg"
	"github.com/ardanlabs/kit/web"
	"github.com/coralproject/shelf/internal/sponge/item"
	"github.com/coralproject/shelf/internal/talk"
)

// actionHandle maintains the set of handlers for the action api.
type actionHandle struct{}

// Action fronts the access to the action service functionality.
var Action actionHandle

var (
	// ErrParameterNotFound is action not found error.
	ErrParameterNotFound = errors.New("Parameter not found")
)

//==============================================================================

// Create an action (flag, like, etc) on an Item.
// 200 Success, 404 Not Found, 500 Internal
func (actionHandle) Create(c *web.Context) error {

	spongedURL := cfg.MustURL("SPONGED_URL").String()

	// Get the user's item_id.
	var user item.Item
	if err := json.NewDecoder(c.Request.Body).Decode(&user); err != nil {
		return err
	}

	// Action to be applied to the target with targetID.
	action, ok := c.Params["action"]
	if !ok {
		return ErrParameterNotFound
	}

	// Target where the action is on.
	targetID, ok := c.Params["item_key"]
	if !ok {
		return ErrParameterNotFound
	}

	// Add the action to the target
	target, err := talk.AddAction(spongedURL, user, action, targetID)
	if err != nil {
		return err
	}

	c.Respond(target, http.StatusOK)
	return nil
}

// Remove an action (flag, like, etc) from an Item.
// 200 Success, 404 Not Found, 500 Internal
func (actionHandle) Remove(c *web.Context) error {
	spongedURL := cfg.MustURL("SPONGED_URL").String()

	// Get the user's item_id.
	var user item.Item
	if err := json.NewDecoder(c.Request.Body).Decode(&user); err != nil {
		return err
	}

	// Action to be applied to the target with targetID.
	action, ok := c.Params["action"]
	if !ok {
		return ErrParameterNotFound
	}

	// Target where the action is on.
	targetID, ok := c.Params["item_key"]
	if !ok {
		return ErrParameterNotFound
	}

	// Add the action to the target
	target, err := talk.RemoveAction(spongedURL, user, action, targetID)
	if err != nil {
		return err
	}

	c.Respond(target, http.StatusOK)
	return nil
}
