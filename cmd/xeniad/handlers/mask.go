package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/coralproject/xenia/pkg/mask"
	"github.com/coralproject/xenia/pkg/regex"

	"github.com/ardanlabs/kit/db"
	"github.com/ardanlabs/kit/web/app"
)

// maskHandle maintains the set of handlers for the mask api.
type maskHandle struct{}

// Mask fronts the access to the mask service functionality.
var Mask maskHandle

//==============================================================================

// List returns all the existing mask in the system.
// 200 Success, 404 Not Found, 500 Internal
func (maskHandle) List(c *app.Context) error {
	masks, err := mask.GetAll(c.SessionID, c.Ctx["DB"].(*db.DB), nil)
	if err != nil {
		if err == regex.ErrNotFound {
			err = app.ErrNotFound
		}
		return err
	}

	c.Respond(masks, http.StatusOK)
	return nil
}

// Retrieve returns the specified mask from the system.
// 200 Success, 400 Bad Request, 404 Not Found, 500 Internal
func (maskHandle) Retrieve(c *app.Context) error {
	collection := c.Params["collection"]
	field := c.Params["field"]

	if collection == "" {
		collection = "*"
	}

	if field == "" {
		masks, err := mask.GetByCollection(c.SessionID, c.Ctx["DB"].(*db.DB), collection)
		if err != nil {
			if err == regex.ErrNotFound {
				err = app.ErrNotFound
			}
			return err
		}

		c.Respond(masks, http.StatusOK)
		return nil
	}

	msk, err := mask.GetByName(c.SessionID, c.Ctx["DB"].(*db.DB), collection, field)
	if err != nil {
		if err == regex.ErrNotFound {
			err = app.ErrNotFound
		}
		return err
	}

	c.Respond(msk, http.StatusOK)
	return nil
}

//==============================================================================

// Upsert inserts or updates the posted mask document into the database.
// 204 SuccessNoContent, 400 Bad Request, 404 Not Found, 500 Internal
func (maskHandle) Upsert(c *app.Context) error {
	var msk mask.Mask
	if err := json.NewDecoder(c.Request.Body).Decode(&msk); err != nil {
		return err
	}

	if err := mask.Upsert(c.SessionID, c.Ctx["DB"].(*db.DB), msk); err != nil {
		return err
	}

	c.Respond(nil, http.StatusNoContent)
	return nil
}

//==============================================================================

// Delete removes the specified mask from the system.
// 200 Success, 400 Bad Request, 404 Not Found, 500 Internal
func (maskHandle) Delete(c *app.Context) error {
	if err := mask.Delete(c.SessionID, c.Ctx["DB"].(*db.DB), c.Params["collection"], c.Params["field"]); err != nil {
		if err == regex.ErrNotFound {
			err = app.ErrNotFound
		}
		return err
	}

	c.Respond(nil, http.StatusNoContent)
	return nil
}
