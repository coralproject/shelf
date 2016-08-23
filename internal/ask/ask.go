package ask

import (
	"errors"
	"time"

	validator "gopkg.in/bluesuncorp/validator.v8"
	"gopkg.in/mgo.v2/bson"

	"github.com/ardanlabs/kit/db"
	"github.com/ardanlabs/kit/log"
	"github.com/coralproject/shelf/internal/ask/form"
	"github.com/coralproject/shelf/internal/ask/form/gallery"
	"github.com/coralproject/shelf/internal/ask/form/submission"
)

//==============================================================================

// validate is used to perform model field validation.
var validate *validator.Validate

func init() {
	validate = validator.New(&validator.Config{TagName: "validate"})
}

//==============================================================================

// ErrInvalidID occurs when an ID is not in a valid form.
var ErrInvalidID = errors.New("ID is not in it's proper form")

//==============================================================================

// UpsertForm upserts the provided form into the MongoDB database collection and
// creates a gallery based on it.
func UpsertForm(context interface{}, db *db.DB, f *form.Form) error {
	log.Dev(context, "UpsertForm", "Started : Form[%s]", f.ID.Hex())

	var isNewForm bool

	// If there was no ID provided, we should set one. UpsertForm might optionally add
	// a form ID to ensure that we don't duplicate the FormGallery.
	if f.ID.Hex() == "" {
		isNewForm = true
	}

	if err := form.Upsert(context, db, f); err != nil {
		log.Error(context, "UpsertForm", err, "Completed")
		return err
	}

	if isNewForm {
		if _, err := gallery.Create(context, db, f.ID.Hex()); err != nil {
			log.Error(context, "UpsertForm", err, "Completed")
			return err
		}
	}

	log.Dev(context, "UpsertForm", "Completed")
	return nil
}

// CreateSubmission creates a form submission based on a given form with a set
// of answers related to it.
func CreateSubmission(context interface{}, db *db.DB, formID string, answers []submission.AnswerInput) (*submission.Submission, error) {
	log.Dev(context, "CreateSubmission", "Started : Form[%s]", formID)

	if !bson.IsObjectIdHex(formID) {
		log.Error(context, "CreateSubmission", ErrInvalidID, "Completed")
		return nil, ErrInvalidID
	}

	// validate that the answers are not invalid
	for _, answer := range answers {
		if err := answer.Validate(); err != nil {
			log.Error(context, "CreateSubmission", err, "Completed")
			return nil, err
		}
	}

	f, err := form.Retrieve(context, db, formID)
	if err != nil {
		log.Error(context, "CreateSubmission", err, "Completed")
		return nil, err
	}

	// create the new form submission
	sub := submission.Submission{
		ID:          bson.NewObjectId(),
		FormID:      bson.ObjectIdHex(formID),
		Header:      f.Header,
		Footer:      f.Footer,
		Answers:     make([]submission.Answer, 0),
		DateCreated: time.Now(),
		DateUpdated: time.Now(),
	}

	// for each answer
	for _, answer := range answers {
		var found bool

		// we must check each step of the form
		for _, step := range f.Steps {
			// and each widget
			for _, widget := range step.Widgets {
				// to see if we can find the matching widget for this answer
				if answer.WidgetID == widget.ID {
					// and push that answer into the form submission
					sub.Answers = append(sub.Answers, submission.Answer{
						WidgetID: widget.ID,
						Answer:   answer,
						Identity: widget.Identity,
						Question: widget.Title,
						Props:    widget.Props,
					})

					// mark the answer as found
					found = true

					break
				}
			}

			// so that if the answer was already found...
			if found {
				// we can break out of this step loop
				break
			}
		}
	}

	if err := submission.Create(context, db, formID, &sub); err != nil {
		log.Error(context, "CreateSubmission", err, "Completed")
		return nil, err
	}

	if _, err := form.UpdateStats(context, db, formID); err != nil {
		log.Error(context, "CreateSubmission", err, "Completed")
		return nil, err
	}

	log.Dev(context, "CreateSubmission", "Completed")
	return &sub, nil
}

// DeleteSubmission deletes a submission as well as updating a form's stats.
func DeleteSubmission(context interface{}, db *db.DB, id, formID string) error {
	log.Dev(context, "DeleteSubmission", "Started")

	if !bson.IsObjectIdHex(id) {
		log.Error(context, "DeleteSubmission", ErrInvalidID, "Completed")
		return ErrInvalidID
	}

	// FIXME: uncomment once old API has been deprecated.
	// if !bson.IsObjectIdHex(formID) {
	// 	log.Error(context, "Delete", ErrInvalidID, "Completed")
	// 	return ErrInvalidID
	// }

	if err := submission.Delete(context, db, id); err != nil {
		log.Error(context, "DeleteSubmission", err, "Completed")
		return err
	}

	// FIXME: remove once old API has been deprecated.
	if formID != "" {
		if _, err := form.UpdateStats(context, db, formID); err != nil {
			log.Error(context, "DeleteSubmission", err, "Completed")
			return err
		}
	}

	log.Dev(context, "DeleteSubmission", "Started")
	return nil
}
