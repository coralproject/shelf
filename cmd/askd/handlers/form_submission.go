package handlers

import (
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"net/url"
	"strconv"

	"github.com/ardanlabs/kit/db"
	"github.com/ardanlabs/kit/web/app"
	"github.com/coralproject/shelf/internal/ask"
	"github.com/coralproject/shelf/internal/ask/form"
	"github.com/coralproject/shelf/internal/ask/form/submission"
)

// ErrInvalidCaptcha is returned when a captcha is required for a form but it
// is not valid on the request.
var ErrInvalidCaptcha = errors.New("captcha invalid")

// formSubmissionHandle maintains the set of handlers for the form submission api.
type formSubmissionHandle struct{}

// FormSubmission fronts the access to the form service functionality.
var FormSubmission formSubmissionHandle

// Create creates a new FormSubmission based on the payload of replies and the
// formID that is being submitted.
// 200 Success, 400 Bad Request, 404 Not Found, 500 Internal
func (formSubmissionHandle) Create(c *app.Context) error {
	var payload struct {
		Recaptcha string
		Answers   []submission.AnswerInput `json:"replies"`
	}
	if err := json.NewDecoder(c.Request.Body).Decode(&payload); err != nil {
		return err
	}

	formID := c.Params["form_id"]

	{
		// We should check to see if the form has a recaptcha property.
		f, err := form.Retrieve(c.SessionID, c.Ctx["DB"].(*db.DB), formID)
		if err != nil {
			return err
		}

		// If the recaptcha is enabled on the form, then we should check that the
		// response contains the data we need and if it's valid.
		if f.Settings["recaptcha"].(bool) {
			if len(payload.Recaptcha) <= 0 {
				return ErrInvalidCaptcha
			}

			body := url.Values{
				"secret":   []string{c.Ctx["recaptcha"].(string)},
				"response": []string{payload.Recaptcha},
			}

			ip, _, err := net.SplitHostPort(c.Request.RemoteAddr)
			if err == nil {
				body["remoteip"] = []string{ip}
			}

			resp, err := http.PostForm("https://www.google.com/recaptcha/api/siteverify", body)
			if err != nil {
				return err
			}
			defer resp.Body.Close()

			var rr struct {
				Success bool `json:"success"`
			}
			if err := json.NewDecoder(resp.Body).Decode(&rr); err != nil {
				return err
			}

			if !rr.Success {
				return ErrInvalidCaptcha
			}
		}
	}

	s, err := ask.CreateSubmission(c.SessionID, c.Ctx["DB"].(*db.DB), formID, payload.Answers)
	if err != nil {
		return err
	}

	c.Respond(s, http.StatusOK)

	return nil
}

// UpdateStatus updates the status of a FormSubmission based on the route
// params.
// 200 Success, 400 Bad Request, 404 Not Found, 500 Internal
func (formSubmissionHandle) UpdateStatus(c *app.Context) error {
	id := c.Params["id"]
	status := c.Params["status"]

	s, err := submission.UpdateStatus(c.SessionID, c.Ctx["DB"].(*db.DB), id, status)
	if err != nil {
		return err
	}

	c.Respond(s, http.StatusOK)

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

	s, err := submission.UpdateAnswer(c.SessionID, c.Ctx["DB"].(*db.DB), id, answerID, editedAnswer)
	if err != nil {
		return err
	}

	c.Respond(s, http.StatusOK)

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

	opts := submission.SearchOpts{
		Query:    c.Request.URL.Query().Get("search"),
		FilterBy: c.Request.URL.Query().Get("filterby"),
	}

	if c.Request.URL.Query().Get("orderby") == "dsc" {
		opts.DscOrder = true
	}

	results, err := submission.Search(c.SessionID, c.Ctx["DB"].(*db.DB), formID, limit, skip, opts)
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

	s, err := submission.Retrieve(c.SessionID, c.Ctx["DB"].(*db.DB), id)
	if err != nil {
		return err
	}

	c.Respond(s, http.StatusOK)

	return nil
}

// Removes a FormSubmission based on the route params.
// 200 Success, 400 Bad Request, 404 Not Found, 500 Internal
func (formSubmissionHandle) Delete(c *app.Context) error {
	id := c.Params["id"]

	// FIXME: pull out ", ok" check once old API has been deprecated.
	formID, ok := c.Params["form_id"]
	if !ok {

		// If in the event that the url does not contain the form id, we should get
		// it from the database by looking up the requested submission.

		sub, err := submission.Retrieve(c.SessionID, c.Ctx["DB"].(*db.DB), id)
		if err != nil {
			return err
		}

		formID = sub.FormID.Hex()
	}

	err := ask.DeleteSubmission(c.SessionID, c.Ctx["DB"].(*db.DB), id, formID)
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

	s, err := submission.AddFlag(c.SessionID, c.Ctx["DB"].(*db.DB), id, flag)
	if err != nil {
		return err
	}

	c.Respond(s, http.StatusOK)

	return nil
}

// RemoveFlag removes a given flag from a FormSubmission based on the provided
// route params.
// 200 Success, 400 Bad Request, 404 Not Found, 500 Internal
func (formSubmissionHandle) RemoveFlag(c *app.Context) error {
	id := c.Params["id"]
	flag := c.Params["flag"]

	s, err := submission.RemoveFlag(c.SessionID, c.Ctx["DB"].(*db.DB), id, flag)
	if err != nil {
		return err
	}

	c.Respond(s, http.StatusOK)

	return nil
}
