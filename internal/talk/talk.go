package talk

// Talk compiles all the functionality related with the Talk product.
// It is used by Corald to massage the data before sending to all other services it uses.

import (
	"errors"

	"github.com/coralproject/shelf/internal/sponge/item"
)

var (
	// ErrActionNotFound is action not found error.
	ErrActionNotFound = errors.New("Action not found")

	//ErrActionNotAllowed is an error when the action is not allowed.
	ErrActionNotAllowed = errors.New("Action not allowed")

	//ErrTypeNotExpected comes when the type asserted was not expected
	ErrTypeNotExpected = errors.New("Type not expected")
)

//==============================================================================

// AddAction an action (flag, like, etc) on an Item.
func AddAction(spongedURL string, user item.Item, action string, targetID string) (item.Item, error) {

	// The action must be in an allowed actions list.
	// To Do: Is this user allowed to do this action on this target?
	if err := userAllowed(user, action); err != nil {
		return item.Item{}, err
	}

	// Users may only add on action of a certain type to one item.

	// Get the Target by targetID.
	target, err := getItemByID(spongedURL, targetID)
	if err != nil {
		return item.Item{}, err
	}

	// Get the actions that the target already has.
	var actions []interface{}
	if target.Data[action] != nil {
		actions = target.Data[action].([]interface{})
	} else {
		target.Data[action] = make([]interface{}, 0)
	}

	// Store the action in the target's action array.
	// Is the user already store for the action?
	found := false
	for _, u := range actions {
		if u == user.ID {
			found = true
			break
		}
	}
	// If the user did not add that action before, then add the action for the user to the target.
	if !found {
		target.Data[action] = append(actions, user.ID)
	}

	// Update the target with the new actions's slice.
	err = upsertItem(spongedURL, target)
	if err != nil {
		return item.Item{}, err
	}

	return target, nil
}

// RemoveAction an action (flag, like, etc) from an Item.
func RemoveAction(spongedURL string, user item.Item, action string, targetID string) (item.Item, error) {

	// The action must be in an allowed actions list.
	// To Do: Is this user allowed to do this action on this target?
	if err := userAllowed(user, action); err != nil {
		return item.Item{}, err
	}

	// Users may only add on action of a certain type to one item.

	// Get the Target by targetID.
	target, err := getItemByID(spongedURL, targetID)
	if err != nil {
		return item.Item{}, err
	}

	// Get the actions that the target already has.
	var actions []interface{}
	if target.Data[action] != nil {
		actions = target.Data[action].([]interface{})
	} else {
		target.Data[action] = make([]interface{}, 0)
	}

	// Remove the action in the target's action array.
	for i, u := range actions {
		if u == user.ID {
			target.Data[action] = append(actions[:i], actions[i+1:]...)
			break
		}
	}

	// Update the target with the new actions's slice.
	err = upsertItem(spongedURL, target)
	if err != nil {
		return item.Item{}, err
	}

	return target, nil
}

// UserAllowed check if the action is allowed or not.
// Returns an error if the user is not allowed to execute that action.
func userAllowed(user item.Item, action string) error {

	// To Do: this will be moved out of the code
	allowedActions := []string{"flagged_by", "liked_by"}

	for _, a := range allowedActions {
		if a == action {
			return nil
		}
	}

	return ErrActionNotAllowed
}
