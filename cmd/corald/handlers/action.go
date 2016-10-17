package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ardanlabs/kit/cfg"
	"github.com/ardanlabs/kit/web"
	"github.com/coralproject/shelf/internal/platform/db"
	"github.com/coralproject/shelf/internal/sponge/item"
)

// actionHandle maintains the set of handlers for the action api.
type actionHandle struct{}

// Action fronts the access to the action service functionality.
var Action actionHandle

var (
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
func (actionHandle) Create(c *web.Context) error {

	spongedURL := cfg.MustURL("SPONGED_URL").String()

	// Users may only add on action of a certain type to one item.
	// Get the user's item_id.
	var user item.Item
	if err := json.NewDecoder(c.Request.Body).Decode(&user); err != nil {
		return err
	}

	// The action must be in an allowed actions list.
	action := c.Params["action"]
	if err := actionAllowed(action); err != nil {
		return err
	}

	// To Do: Is this user allowed to do this action on this target?
	if err := userAllowed(user, action); err != nil {
		return err
	}

	db := c.Ctx["DB"].(*db.DB)

	// Get the Target by item_key.
	itemKey := c.Params["item_key"]
	target, err := item.GetByID(c, db, itemKey)
	if err != nil {
		return err
	}

	// This is a new action and we need to create an empty []string.
	if target.Data[action] == nil {
		target.Data[action] = make([]interface{}, 0)
	}
	// Convert interface{} to []string
	s, ok := target.Data[action].([]interface{})
	if !ok {
		return ErrTypeNotExpected
	}

	// Store the action in the target's action array.
	// Is the user already store for the action?
	found := false
	for _, u := range s {
		if u == user.ID {
			found = true
			break
		}
	}
	// If the user did not add that action before, then add the action for the user to the target.
	if !found {
		target.Data[action] = append(s, user.ID)
	}

	// Upsert the target with the new actions
	url := spongedURL + "/v1/item"
	body, err := json.Marshal(target)
	if err != nil {
		return err
	}

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	c.Respond(target, http.StatusOK)

	return nil
}

// Remove an action (flag, like, etc) from an Item.
// 200 Success, 404 Not Found, 500 Internal
func (actionHandle) Remove(c *web.Context) error {

	spongedURL := cfg.MustURL("SPONGED_URL").String()

	// Users may only add on action of a certain type to one item.
	// Get the user's item_id.
	var user item.Item
	if err := json.NewDecoder(c.Request.Body).Decode(&user); err != nil {
		return err
	}

	// The action must be in an allowed actions list.
	action := c.Params["action"]
	if err := actionAllowed(action); err != nil {
		return err
	}

	// To Do: Is this user allowed to do this action on this target?
	if err := userAllowed(user, action); err != nil {
		return err
	}

	db := c.Ctx["DB"].(*db.DB)

	// Get the Target by item_key.
	itemKey := c.Params["item_key"]
	target, err := item.GetByID(c, db, itemKey)
	if err != nil {
		return err
	}

	// This is a new action and we need to create an empty []string.
	if target.Data[action] == nil {
		target.Data[action] = make([]string, 0)
	}
	// Convert interface{} to []string
	s, ok := target.Data[action].([]interface{})
	if !ok {
		return ErrTypeNotExpected
	}

	// Remove the action from the target's action array.
	for i, u := range s {
		if u == user.ID {
			// Remove the action
			target.Data[action] = append(s[:i], s[i+1:]...)
			break
		}
	}

	// Upsert the target with the new actions
	url := spongedURL + "/v1/item"
	body, err := json.Marshal(target)
	if err != nil {
		return err
	}

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	c.Respond(target, http.StatusOK)

	return nil
}

// ActionAllowed check if the action is allowed or not.
// Returns an error if it is not allowed.
func actionAllowed(action string) error {

	// To Do: this will be moved out of the code
	allowedActions := []string{"flagged_by", "liked_by"}

	for _, a := range allowedActions {
		if a == action {
			return nil
		}
	}

	return ErrActionNotAllowed
}

// UserAllowed check if the action is allowed or not.
// Returns an error if the user is not allowed to execute that action.
func userAllowed(user item.Item, action string) error {
	return nil
}
