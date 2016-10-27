package handlers

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"gopkg.in/mgo.v2/bson"

	"github.com/ardanlabs/kit/web"
	"github.com/coralproject/shelf/internal/ask"
	"github.com/coralproject/shelf/internal/ask/form"
	"github.com/coralproject/shelf/internal/ask/form/submission"
	"github.com/coralproject/shelf/internal/platform/db"
)

// formHandle maintains the set of handlers for the form api.
type formHandle struct{}

// Form fronts the access to the form service functionality.
var Form formHandle

//==============================================================================

// Upsert upserts a form into the store.
// 200 Success, 400 Bad Request, 404 Not Found, 500 Internal
func (formHandle) Upsert(c *web.Context) error {
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
func (formHandle) UpdateStatus(c *web.Context) error {
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
func (formHandle) List(c *web.Context) error {
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
func (formHandle) Retrieve(c *web.Context) error {
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
func (formHandle) Delete(c *web.Context) error {
	id := c.Params["id"]

	err := form.Delete(c.SessionID, c.Ctx["DB"].(*db.DB), id)
	if err != nil {
		return err
	}

	c.Respond(nil, http.StatusOK)
	return nil
}

//==============================================================================

// AggregationKeys is a transport type that describes the json format for return value of the
// aggregate endpoint including a key lookup structure allowing consumers to easily find group keys.
type AggregationKeys struct {
	Aggregations map[string]form.Aggregation `json:"aggregations"`
}

// Aggregate performs all aggregations across a form's submissions.
// 200 Success, 400 Bad Request, 404 Not Found, 500 Internal
func (formHandle) Aggregate(c *web.Context) error {
	id := c.Params["form_id"]

	aggregations, err := form.AggregateFormSubmissions(c.SessionID, c.Ctx["DB"].(*db.DB), id)
	if err != nil {
		return err
	}

	ak := AggregationKeys{
		Aggregations: aggregations,
	}

	c.Respond(ak, http.StatusOK)

	return nil
}

// AggregateGroup performs all aggregations across a form's submissions and returns a single group.
// 200 Success, 400 Bad Request, 404 Not Found, 500 Internal
func (formHandle) AggregateGroup(c *web.Context) error {
	id := c.Params["form_id"]

	aggregations, err := form.AggregateFormSubmissions(c.SessionID, c.Ctx["DB"].(*db.DB), id)
	if err != nil {
		return err
	}

	aggregation, ok := aggregations[c.Params["group_id"]]
	if !ok {
		c.Respond(nil, http.StatusNotFound)
	}

	c.Respond(aggregation, http.StatusOK)

	return nil
}

//==============================================================================

// Form Digests: this section contains transport/bluprints for sending digests of form/question
// information to clients that do not need/won't understand full, vebose feeds.

// FormQuestionOptionDigest is the blueprint for a single multiple choice option.
type FormQuestionOptionDigest struct {
	ID    string `json:"id" bson:"id"`
	Value string `json:"value" bson:"value"`
}

// FormQuestionDigest is the bluprint for a question in a form group.
type FormQuestionDigest struct {
	ID      string                     `json:"id" bson:"id"`
	Type    string                     `json:"type" bson:"type"`
	Title   string                     `json:"title" bson:"title"`
	GroupBy bool                       `json:"group_by" bson:"group_by"`
	Options []FormQuestionOptionDigest `json:"options,omitempty" bson:"options,omitempty"`
	Order   int                        `json:"order" bson:"order"`
}

// FormDigest is the blueprint for how we send form digests to clients.
type FormDigest struct {
	Questions map[string]FormQuestionDigest `json:"questions" bson:"questions"`
}

// Digest returns a form digest.
// 200 Success, 400 Bad Request, 404 Not Found, 500 Internal
func (formHandle) Digest(c *web.Context) error {
	id := c.Params["form_id"]

	// Load the form requested.
	f, err := form.Retrieve(c.SessionID, c.Ctx["DB"].(*db.DB), id)
	if err != nil {
		return err
	}

	// Create a container for the question digests.
	questions := make(map[string]FormQuestionDigest)

	// Order is a counter to set the order questions in the form.
	order := 1

	// Loop through to form's steps/widgets to find the questions.
	for _, step := range f.Steps {
		for _, widget := range step.Widgets {

			// Unpack the question properties.
			props := widget.Props.(bson.M)

			// We are looking to only include submissions with includeInGroups or groupSubmissions.
			gs, gsok := props["groupSubmissions"]
			iig, iigok := props["includeInGroups"]

			// Skip other questions, and do it verbosely to protect against messy data.
			if (gs == nil || !gsok || gs == false) && (iig == nil || !iigok || iig == false) {
				continue
			}

			// Include options for MultipleChoice questions.
			options := []FormQuestionOptionDigest{}
			if widget.Component == "MultipleChoice" {

				// Step outside the safety of the type system...
				opts := props["options"].([]interface{})

				for _, opt := range opts {
					option := opt.(bson.M)
					fmt.Printf("\n\n%#v", option)

					// Hash the answer text for a unique key, as no actual key exists.
					hasher := md5.New()
					hasher.Write([]byte(option["title"].(string)))
					optKeyStr := hex.EncodeToString(hasher.Sum(nil))

					// Add this option to the array.
					options = append(options, FormQuestionOptionDigest{
						ID:    optKeyStr,
						Value: option["title"].(string),
					})
				}
			}

			// Add the question to the digest.
			questions[widget.ID] = FormQuestionDigest{
				ID:      widget.ID,
				Type:    widget.Component,
				Title:   widget.Title,
				GroupBy: gsok,
				Options: options,
				Order:   order,
			}

			// Increment the order counter.
			order = order + 1

		}
	}

	digest := FormDigest{
		Questions: questions,
	}

	c.Respond(digest, http.StatusOK)

	return nil
}

//==============================================================================

// Submission Groups: The submission answers marked with includeInGroups from
// the given groups, or all.

//
// 200 Success, 400 Bad Request, 404 Not Found, 500 Internal
func (formHandle) SubmissionGroup(c *web.Context) error {
	id := c.Params["form_id"]

	// unpack the search headers and create a SearchOpts for Group Submissions

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

	groups, err := form.GroupSubmissions(c.SessionID, c.Ctx["DB"].(*db.DB), id, limit, skip, opts)
	if err != nil {
		return err
	}

	groupKey, ok := c.Params["group_id"]
	if !ok {
		c.Respond(nil, http.StatusNotFound)
	}

	ta := []form.TextAggregation{}

	for group, submissions := range groups {

		if group.ID == groupKey {

			ta, err = form.TextAggregate(c.SessionID, submissions)
			if err != nil {
				return err
			}

			fmt.Printf("\n\n%#v", ta)
		}
	}

	c.Respond(ta, http.StatusOK)

	return nil
}
