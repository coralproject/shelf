package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ardanlabs/kit/db"
	"github.com/ardanlabs/kit/web/app"
	"github.com/coralproject/shelf/internal/ask/form/gallery"
)

// formGalleryHandle maintains the set of handlers for the form gallery api.
type formGalleryHandle struct{}

// FormGallery fronts the access to the form service functionality.
var FormGallery formGalleryHandle

// AddAnswer adds an answer to a form gallery in the store.
// 200 Success, 400 Bad Request, 404 Not Found, 500 Internal
func (formGalleryHandle) AddAnswer(c *app.Context) error {
	id := c.Params["id"]
	submissionID := c.Params["submission_id"]
	answerID := c.Params["answer_id"]

	gallery, err := gallery.AddAnswer(c.SessionID, c.Ctx["DB"].(*db.DB), id, submissionID, answerID)
	if err != nil {
		return err
	}

	c.Respond(gallery, http.StatusOK)
	return nil
}

// AddAnswer removes an answer from a form gallery in the store.
// 200 Success, 400 Bad Request, 404 Not Found, 500 Internal
func (formGalleryHandle) RemoveAnswer(c *app.Context) error {
	id := c.Params["id"]
	submissionID := c.Params["submission_id"]
	answerID := c.Params["answer_id"]

	gallery, err := gallery.RemoveAnswer(c.SessionID, c.Ctx["DB"].(*db.DB), id, submissionID, answerID)
	if err != nil {
		return err
	}

	c.Respond(gallery, http.StatusOK)
	return nil
}

// RetrieveForForm retrieves a collection of galleries based on a specific form
// id.
// 200 Success, 400 Bad Request, 404 Not Found, 500 Internal
func (formGalleryHandle) RetrieveForForm(c *app.Context) error {
	formID := c.Params["form_id"]

	galleries, err := gallery.List(c.SessionID, c.Ctx["DB"].(*db.DB), formID)
	if err != nil {
		return err
	}

	c.Respond(galleries, http.StatusOK)
	return nil
}

// Retrieve retrieves a FormGallery based on it's id.
// 200 Success, 400 Bad Request, 404 Not Found, 500 Internal
func (formGalleryHandle) Retrieve(c *app.Context) error {
	id := c.Params["id"]

	gallery, err := gallery.Retrieve(c.SessionID, c.Ctx["DB"].(*db.DB), id)
	if err != nil {
		return err
	}

	c.Respond(gallery, http.StatusOK)
	return nil
}

// Update updates a FormGallery based on it's id and it's provided payload.
// 200 Success, 400 Bad Request, 404 Not Found, 500 Internal
func (formGalleryHandle) Update(c *app.Context) error {
	var g gallery.Gallery
	if err := json.NewDecoder(c.Request.Body).Decode(&g); err != nil {
		return err
	}

	id := c.Params["id"]

	err := gallery.Update(c.SessionID, c.Ctx["DB"].(*db.DB), id, &g)
	if err != nil {
		return err
	}

	c.Respond(g, http.StatusOK)
	return nil
}
