package handlers

import (
	"errors"
	"net/http"

	"github.com/ardanlabs/kit/web"
	"github.com/coralproject/shelf/cmd/corald/service"
)

// actionHandle maintains the set of handlers for the action api.
type actionHandle struct{}

// Action fronts the access to the action service functionality.
var Action actionHandle

var (
	// ErrParameterNotFound is action not found error.
	ErrParameterNotFound = errors.New("Parameter not found")

	// ErrActionNotFound is action not found error.
	ErrActionNotFound = errors.New("Action not found")

	//ErrActionNotAllowed is an error when the action is not allowed.
	ErrActionNotAllowed = errors.New("Action not allowed")

	//ErrTypeNotExpected comes when the type asserted was not expected
	ErrTypeNotExpected = errors.New("Type not expected")
)

//==============================================================================

// Create an action (flag, like, etc) on an Item.
// 200 Success, 404 Not Found, 500 Internal
func (actionHandle) Add(c *web.Context) error {

	// Action to be applied to the target with targetID id.
	action := c.Params["action"]

	// Target where the action is on.
	targetID := c.Params["item_key"]

	// Item that is performing the action.
	itemID := c.Params["user_key"]

	// Add the action by itm to the target targetID
	err := addAction(c, itemID, action, targetID)
	if err != nil {
		return err
	}

	c.Respond(nil, http.StatusOK)
	return nil
}

// Remove an action (flag, like, etc) from an Item.
// 200 Success, 404 Not Found, 500 Internal
func (actionHandle) Remove(c *web.Context) error {

	// Action to be applied to the target with targetID.
	action := c.Params["action"]

	// Target where the action is on.
	targetID := c.Params["item_key"]

	// Item that is performing the action.
	itemID := c.Params["user_key"]

	// Add the action to the target
	err := removeAction(c, itemID, action, targetID)
	if err != nil {
		return err
	}

	c.Respond(nil, http.StatusOK)
	return nil
}

//==============================================================================

// addAction ads an action (flag, like, etc) perform by an user UserID on a Target.
func addAction(c *web.Context, userID string, action string, targetID string) error {

	// Get the Target's data by targetID.
	target, err := service.GetItemByID(c, targetID)
	if err != nil {
		return err
	}

	// Get the actions that the target already has.
	// If it has no action 'action' then create the field to store the new one.
	var actions []interface{}
	actions, ok := target.Data[action].([]interface{})
	if !ok {
		target.Data[action] = make([]interface{}, 0)
	}

	// Store the action in the target's action array.
	var found bool
	for _, usrID := range actions {
		if usrID.(string) == userID {
			found = true
			break
		}
	}

	// If we did not find the user in the actions
	// If the user did not add that action before, then add the action for the user to the target.
	if !found {
		target.Data[action] = append(actions, userID)
	}

	// Update the target with the new actions's slice.
	if err = service.UpsertItem(c, target); err != nil {
		return err
	}

	return nil
}

// removeAction removes an action (flag, like, etc) from an Item.
func removeAction(c *web.Context, userID string, action string, targetID string) error {

	// Get the Target by targetID.
	target, err := service.GetItemByID(c, targetID)
	if err != nil {
		return err
	}

	// Get the actions that the target already has.
	var actions []string
	actions, ok := target.Data[action].([]string)
	if !ok {
		return ErrActionNotFound
	}

	// Remove the action in the target's action array.
	for i, u := range actions {
		if u == userID {
			target.Data[action] = append(actions[:i], actions[i+1:]...)
			break
		}
	}

	// Update the target with the new actions's slice.
	if err = service.UpsertItem(c, target); err != nil {
		return err
	}

	return nil
}
