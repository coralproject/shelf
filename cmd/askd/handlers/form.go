package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ardanlabs/kit/db"
	"github.com/ardanlabs/kit/web/app"
	"github.com/coralproject/shelf/internal/ask"
	"github.com/coralproject/shelf/internal/ask/form"
)

// formHandle maintains the set of handlers for the form api.
type formHandle struct{}

// Form fronts the access to the form service functionality.
var Form formHandle

//==============================================================================

// Upsert upserts a form into the store.
// 200 Success, 400 Bad Request, 404 Not Found, 500 Internal
func (formHandle) Upsert(c *app.Context) error {
	var form form.Form
	if err := json.NewDecoder(c.Request.Body).Decode(&form); err != nil {
		return err
	}

	// perform the upsert operation
	err := ask.UpsertForm(c.SessionID, c.Ctx["DB"].(*db.DB), &form)
	if err != nil {
		return err
	}

	c.Respond(form, http.StatusOK)
	return nil
}

// UpdateStatus updates the status of a form in the store.
// 200 Success, 400 Bad Request, 404 Not Found, 500 Internal
func (formHandle) UpdateStatus(c *app.Context) error {
	id := c.Params["id"]
	status := c.Params["status"]

	f, err := form.UpdateStatus(c.SessionID, c.Ctx["DB"].(*db.DB), id, status)
	if err != nil {
		return err
	}

	c.Respond(f, http.StatusOK)
	return nil
}

// List retrieves a list of forms based on the query parameters from the
// store.
// 200 Success, 400 Bad Request, 404 Not Found, 500 Internal
func (formHandle) List(c *app.Context) error {
	limit, err := strconv.Atoi(c.Request.URL.Query().Get("limit"))
	if err != nil {
		limit = 0
	}

	skip, err := strconv.Atoi(c.Request.URL.Query().Get("skip"))
	if err != nil {
		skip = 0
	}

	forms, err := form.List(c.SessionID, c.Ctx["DB"].(*db.DB), limit, skip)
	if err != nil {
		return err
	}

	c.Respond(forms, http.StatusOK)
	return nil
}

// Retrieve retrieves a single form from the store.
// 200 Success, 400 Bad Request, 404 Not Found, 500 Internal
func (formHandle) Retrieve(c *app.Context) error {
	id := c.Params["id"]

	f, err := form.Retrieve(c.SessionID, c.Ctx["DB"].(*db.DB), id)
	if err != nil {
		return err
	}

	c.Respond(f, http.StatusOK)
	return nil
}

// Delete removes a single form from the store.
// 200 Success, 400 Bad Request, 404 Not Found, 500 Internal
func (formHandle) Delete(c *app.Context) error {
	id := c.Params["id"]

	err := form.Delete(c.SessionID, c.Ctx["DB"].(*db.DB), id)
	if err != nil {
		return err
	}

	c.Respond(nil, http.StatusOK)
	return nil
}
