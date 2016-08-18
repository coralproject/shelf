package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ardanlabs/kit/db"
	"github.com/ardanlabs/kit/web/app"
	"github.com/coralproject/shelf/internal/ask"
)

// formSubmissionHandle maintains the set of handlers for the form submission api.
type formSubmissionHandle struct{}

// FormSubmission fronts the access to the form service functionality.
var FormSubmission formSubmissionHandle

// Create creates a new FormSubmission based on the payload of replies and the
// formID that is being submitted.
// 200 Success, 400 Bad Request, 404 Not Found, 500 Internal
func (formSubmissionHandle) Create(c *app.Context) error {
	var payload struct {
		Answers []ask.FormSubmissionAnswerInput `json:"replies"`
	}
	if err := json.NewDecoder(c.Request.Body).Decode(&payload); err != nil {
		return err
	}

	formID := c.Params["form_id"]

	submission, err := ask.CreateFormSubmission(c, c.Ctx["DB"].(*db.DB), formID, payload.Answers)
	if err != nil {
		return err
	}

	c.Respond(submission, http.StatusOK)

	return nil
}

// UpdateStatus updates the status of a FormSubmission based on the route
// params.
// 200 Success, 400 Bad Request, 404 Not Found, 500 Internal
func (formSubmissionHandle) UpdateStatus(c *app.Context) error {
	id := c.Params["id"]
	status := c.Params["status"]

	submission, err := ask.UpdateFormStatus(c, c.Ctx["DB"].(*db.DB), id, status)
	if err != nil {
		return err
	}

	c.Respond(submission, http.StatusOK)

	return nil
}

// UpdateAnswer updates an answer based on the payload submitted to the
// endpoint.
// 200 Success, 400 Bad Request, 404 Not Found, 500 Internal
func (formSubmissionHandle) UpdateAnswer(c *app.Context) error {
	var editedAnswer interface{}
	if err := json.NewDecoder(c.Request.Body).Decode(&editedAnswer); err != nil {
		return err
	}

	id := c.Params["id"]
	answerID := c.Params["answer_id"]

	submission, err := ask.UpdateFormSubmissionAnswer(c, c.Ctx["DB"].(*db.DB), id, answerID, editedAnswer)
	if err != nil {
		return err
	}

	c.Respond(submission, http.StatusOK)

	return nil
}

// Search retrieves a set of FormSubmission's based on the search params
// provided in the query string.
// 200 Success, 400 Bad Request, 404 Not Found, 500 Internal
func (formSubmissionHandle) Search(c *app.Context) error {
	formID := c.Params["form_id"]

	limit, err := strconv.Atoi(c.Request.URL.Query().Get("limit"))
	if err != nil {
		limit = 0
	}

	skip, err := strconv.Atoi(c.Request.URL.Query().Get("skip"))
	if err != nil {
		skip = 0
	}

	opts := ask.FormSubmissionSearchOpts{
		Query:    c.Request.URL.Query().Get("search"),
		FilterBy: c.Request.URL.Query().Get("filterby"),
	}

	if c.Request.URL.Query().Get("orderby") == "dsc" {
		opts.DscOrder = true
	}

	results, err := ask.SearchFormSubmissions(c, c.Ctx["DB"].(*db.DB), formID, limit, skip, opts)
	if err != nil {
		return err
	}

	c.Respond(results, http.StatusOK)

	return nil
}

// Retrieves a given FormSubmission based on the route params.
// 200 Success, 400 Bad Request, 404 Not Found, 500 Internal
func (formSubmissionHandle) Retrieve(c *app.Context) error {
	id := c.Params["id"]

	submission, err := ask.RetrieveFormSubmission(c, c.Ctx["DB"].(*db.DB), id)
	if err != nil {
		return err
	}

	c.Respond(submission, http.StatusOK)

	return nil
}

// Removes a FormSubmission based on the route params.
// 200 Success, 400 Bad Request, 404 Not Found, 500 Internal
func (formSubmissionHandle) Delete(c *app.Context) error {
	id := c.Params["id"]

	// FIXME: pull out ", ok" check once old API has been deprecated.
	formID, ok := c.Params["form_id"]
	if !ok {
		// currently will fall back on no form ID in the event that there is a
		// deleted form submission, but it should not once the old API route has
		// been deprecated.
		formID = ""
	}

	err := ask.DeleteFormSubmission(c, c.Ctx["DB"].(*db.DB), id, formID)
	if err != nil {
		return err
	}

	c.Respond(nil, http.StatusOK)

	return nil
}

// AddFlag adds a new flag to a given FormSubmission based on the provided route
// params.
// 200 Success, 400 Bad Request, 404 Not Found, 500 Internal
func (formSubmissionHandle) AddFlag(c *app.Context) error {
	id := c.Params["id"]
	flag := c.Params["flag"]

	submission, err := ask.AddFlagToFormSubmission(c, c.Ctx["DB"].(*db.DB), id, flag)
	if err != nil {
		return err
	}

	c.Respond(submission, http.StatusOK)

	return nil
}

// RemoveFlag removes a given flag from a FormSubmission based on the provided
// route params.
// 200 Success, 400 Bad Request, 404 Not Found, 500 Internal
func (formSubmissionHandle) RemoveFlag(c *app.Context) error {
	id := c.Params["id"]
	flag := c.Params["flag"]

	submission, err := ask.RemoveFlagFromFormSubmission(c, c.Ctx["DB"].(*db.DB), id, flag)
	if err != nil {
		return err
	}

	c.Respond(submission, http.StatusOK)

	return nil
}
