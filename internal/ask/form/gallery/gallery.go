package gallery

import (
	"errors"
	"time"

	"github.com/ardanlabs/kit/db"
	"github.com/ardanlabs/kit/db/mongo"
	"github.com/ardanlabs/kit/log"
	"github.com/coralproject/shelf/internal/ask/form/submission"
	validator "gopkg.in/bluesuncorp/validator.v8"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
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

// GalleryCollection is the mongo collection where Gallery documents are
// saved.
const GalleryCollection = "form_galleries"

// Answer describes an answer from a form which has been added to a
// Gallery.
type Answer struct {
	SubmissionID    bson.ObjectId       `json:"submission_id" bson:"submission_id" validate:"required"`
	AnswerID        string              `json:"answer_id" bson:"answer_id" validate:"required,len=24"`
	Answer          submission.Answer   `json:"answer,omitempty" bson:"-"`
	IdentityAnswers []submission.Answer `json:"identity_answers,omitempty" bson:"-"`
}

// Gallery is a Form that has been moved to a shared space.
type Gallery struct {
	ID          bson.ObjectId          `json:"id" bson:"_id" validate:"required"`
	FormID      bson.ObjectId          `json:"form_id" bson:"form_id" validate:"required"`
	Headline    string                 `json:"headline" bson:"headline"`
	Description string                 `json:"description" bson:"description"`
	Config      map[string]interface{} `json:"config" bson:"config"`
	Answers     []Answer               `json:"answers" bson:"answers"`
	DateCreated time.Time              `json:"date_created,omitempty" bson:"date_created,omitempty"`
	DateUpdated time.Time              `json:"date_updated,omitempty" bson:"date_updated,omitempty"`
}

// Validate checks the Gallery value for consistency.
func (fg *Gallery) Validate() error {
	if err := validate.Struct(fg); err != nil {
		return err
	}

	return nil
}

// Create adds a form gallery based on the form id provided into the
// MongoDB database collection.
func Create(context interface{}, db *db.DB, formID string) (*Gallery, error) {
	log.Dev(context, "Create", "Started")

	if !bson.IsObjectIdHex(formID) {
		log.Error(context, "Create", ErrInvalidID, "Completed")
		return nil, ErrInvalidID
	}

	gallery := Gallery{
		ID:          bson.NewObjectId(),
		FormID:      bson.ObjectIdHex(formID),
		DateCreated: time.Now(),
		DateUpdated: time.Now(),
	}

	f := func(c *mgo.Collection) error {
		log.Dev(context, "Create", "MGO : db.%s.insert(%s)", c.Name, mongo.Query(gallery))
		return c.Insert(gallery)
	}

	if err := db.ExecuteMGO(context, GalleryCollection, f); err != nil {
		log.Error(context, "Create", err, "Completed")
		return nil, err
	}

	log.Dev(context, "Create", "Completed")
	return nil, nil
}

// Retrieve retrieves a form gallery from the MongoDB database
// collection as well as hydrating the form gallery with form submissions.
func Retrieve(context interface{}, db *db.DB, id string) (*Gallery, error) {
	log.Dev(context, "Retrieve", "Started")

	if !bson.IsObjectIdHex(id) {
		log.Error(context, "Retrieve", ErrInvalidID, "Completed")
		return nil, ErrInvalidID
	}

	objectID := bson.ObjectIdHex(id)

	var gallery Gallery
	f := func(c *mgo.Collection) error {
		log.Dev(context, "Retrieve", "MGO : db.%s.find(%s)", c.Name, mongo.Query(objectID))
		return c.FindId(objectID).One(&gallery)
	}

	if err := db.ExecuteMGO(context, GalleryCollection, f); err != nil {
		log.Error(context, "Retrieve", err, "Completed")
		return nil, err
	}

	if err := hydrate(context, db, &gallery); err != nil {
		log.Error(context, "Retrieve", err, "Completed")
		return nil, err
	}

	log.Dev(context, "Retrieve", "Completed")
	return &gallery, nil
}

// hydrate loads a Gallery with form submissions from the MongoDB
// database collection.
func hydrate(context interface{}, db *db.DB, gallery *Gallery) error {
	log.Dev(context, "hydrate", "Started")

	if err := gallery.Validate(); err != nil {
		log.Error(context, "hydrate", err, "Completed")
		return err
	}

	// load all the submission id's from the answers inside the gallery
	var submissionIDs = make([]string, len(gallery.Answers))
	for i, answer := range gallery.Answers {
		// and set their hex
		submissionIDs[i] = answer.SubmissionID.Hex()
	}

	// so we can fetch all the form submissions in one request
	submissions, err := submission.List(context, db, submissionIDs)
	if err != nil {
		log.Error(context, "hydrate", err, "Completed")
		return err
	}

	// merge the submissions into the given gallery
	mergeSubmissionsIntoGalleryAnswers(gallery, submissions)

	log.Dev(context, "hydrate", "Completed")
	return nil
}

// mergeSubmissionsIntoGalleryAnswers associates the array of submissions onto
// matching gallery answers.
func mergeSubmissionsIntoGalleryAnswers(gallery *Gallery, submissions []submission.Submission) {
	// walk through all their answers
	for j, answer := range gallery.Answers {

		// and for each submission
		for k, sub := range submissions {
			// if we are looking at a different submission that doesn't match the
			// answer's submission ID or the submission was to a different form that
			// the current gallery is on, then we need to skip this submission.
			if sub.ID != answer.SubmissionID || sub.FormID.Hex() != gallery.FormID.Hex() {
				continue
			}

			// we have verified that the current submission is indeed for the
			// current gallery form and matches the submission id

			// so lets walk over the current submission's answers to match to the
			// particular question/widget/answer that we want to look at
			for _, submissionAnswer := range sub.Answers {
				// and if it doesn't match
				if submissionAnswer.WidgetID != answer.AnswerID {
					// then just continue
					continue
				}

				// but as it does

				// set the answer to the current submission answer
				gallery.Answers[j].Answer = sub.Answers[k]

				// and create an empty array for the identity answers that we will
				// walk over
				gallery.Answers[j].IdentityAnswers = make([]submission.Answer, 0)

				// specifically, walk over the the current submission's answers again
				for m, submissionAnswer := range sub.Answers {
					// to find any identity answers related to this specific answer
					if submissionAnswer.Identity {
						// and append it to the list
						gallery.Answers[j].IdentityAnswers = append(gallery.Answers[j].IdentityAnswers, sub.Answers[m])
					}
				}

				// because we found a match for the specific answer/widget, we can't
				// possibly have another duplicate, so stop looping over the the
				// current submissions answers
				break
			}

			// and because we found a match for the submission id and the form id
			// there can't possibly be another match, so stop looping over the
			// submissions
			break
		}
	}
}

// hydrateMany loads an array of form galleries with form submissions
// from the MongoDB database collection.
func hydrateMany(context interface{}, db *db.DB, galleries []Gallery) error {
	log.Dev(context, "hydrateMany", "Started")

	for _, gallery := range galleries {
		if err := gallery.Validate(); err != nil {
			log.Error(context, "hydrateMany", err, "Completed")
			return err
		}
	}

	// load all the submission id's from the answers inside the gallery
	var submissionIDs []string

	for _, gallery := range galleries {
		for i, answer := range gallery.Answers {
			// and set their hex
			submissionIDs[i] = answer.SubmissionID.Hex()
		}
	}

	// so we can fetch all the form submissions in one request
	submissions, err := submission.List(context, db, submissionIDs)
	if err != nil {
		log.Error(context, "hydrateMany", err, "Completed")
		return err
	}

	// for each of the galleries that we're hydrating
	for i := range galleries {
		// merge the submissions into the given gallery
		mergeSubmissionsIntoGalleryAnswers(&galleries[i], submissions)
	}

	log.Dev(context, "hydrateMany", "Completed")
	return nil
}

// AddAnswer adds an answer to a form gallery. Duplicated answers
// are de-duplicated automatically and will not return an error.
func AddAnswer(context interface{}, db *db.DB, id, submissionID, answerID string) (*Gallery, error) {
	log.Dev(context, "AddAnswer", "Started")

	if !bson.IsObjectIdHex(id) {
		log.Error(context, "AddAnswer", ErrInvalidID, "Completed")
		return nil, ErrInvalidID
	}

	if !bson.IsObjectIdHex(submissionID) {
		log.Error(context, "AddAnswer", ErrInvalidID, "Completed")
		return nil, ErrInvalidID
	}

	if !bson.IsObjectIdHex(answerID) {
		log.Error(context, "AddAnswer", ErrInvalidID, "Completed")
		return nil, ErrInvalidID
	}

	objectID := bson.ObjectIdHex(id)

	answer := Answer{
		SubmissionID: bson.ObjectIdHex(submissionID),
		AnswerID:     answerID,
	}

	f := func(c *mgo.Collection) error {
		u := bson.M{
			"$addToSet": bson.M{
				"answers": answer,
			},
		}
		log.Dev(context, "AddAnswer", "MGO : db.%s.update(%s, %s)", c.Name, mongo.Query(objectID), mongo.Query(u))
		return c.UpdateId(objectID, u)
	}

	if err := db.ExecuteMGO(context, GalleryCollection, f); err != nil {
		log.Error(context, "AddAnswer", err, "Completed")
		return nil, err
	}

	gallery, err := Retrieve(context, db, id)
	if err != nil {
		log.Error(context, "AddAnswer", err, "Completed")
		return nil, err
	}

	log.Dev(context, "AddAnswer", "Completed")
	return gallery, nil
}

// RemoveAnswer adds an answer to a form gallery. Duplicated answers
// are de-duplicated automatically and will not return an error.
func RemoveAnswer(context interface{}, db *db.DB, id, submissionID, answerID string) (*Gallery, error) {
	log.Dev(context, "RemoveAnswer", "Started")

	if !bson.IsObjectIdHex(id) {
		log.Error(context, "RemoveAnswer", ErrInvalidID, "Completed")
		return nil, ErrInvalidID
	}

	if !bson.IsObjectIdHex(submissionID) {
		log.Error(context, "RemoveAnswer", ErrInvalidID, "Completed")
		return nil, ErrInvalidID
	}

	if !bson.IsObjectIdHex(answerID) {
		log.Error(context, "RemoveAnswer", ErrInvalidID, "Completed")
		return nil, ErrInvalidID
	}

	objectID := bson.ObjectIdHex(id)

	answer := Answer{
		SubmissionID: bson.ObjectIdHex(submissionID),
		AnswerID:     answerID,
	}

	f := func(c *mgo.Collection) error {
		u := bson.M{
			"$pull": bson.M{
				"answers": answer,
			},
		}
		log.Dev(context, "RemoveAnswer", "MGO : db.%s.update(%s, %s)", c.Name, mongo.Query(objectID), mongo.Query(u))
		return c.UpdateId(objectID, u)
	}

	if err := db.ExecuteMGO(context, GalleryCollection, f); err != nil {
		log.Error(context, "RemoveAnswer", err, "Completed")
		return nil, err
	}

	gallery, err := Retrieve(context, db, id)
	if err != nil {
		log.Error(context, "RemoveAnswer", err, "Completed")
		return nil, err
	}

	log.Dev(context, "RemoveAnswer", "Completed")
	return gallery, nil
}

// List retrives the form galleries for a given form from the MongoDB database
// collection.
func List(context interface{}, db *db.DB, formID string) ([]Gallery, error) {
	log.Dev(context, "List", "Started")

	if !bson.IsObjectIdHex(formID) {
		log.Error(context, "List", ErrInvalidID, "Completed")
		return nil, ErrInvalidID
	}

	formObjectID := bson.ObjectIdHex(formID)

	var galleries = make([]Gallery, 0)
	f := func(c *mgo.Collection) error {
		q := bson.M{
			"form_id": formObjectID,
		}
		log.Dev(context, "List", "MGO : db.%s.find(%s)", c.Name, mongo.Query(q))
		return c.Find(q).All(&galleries)
	}

	if err := db.ExecuteMGO(context, GalleryCollection, f); err != nil {
		log.Error(context, "List", err, "Completed")
		return nil, err
	}

	if err := hydrateMany(context, db, galleries); err != nil {
		log.Error(context, "List", err, "Completed")
		return nil, err
	}

	log.Dev(context, "List", "Completed")
	return galleries, nil
}

// Update updates the form gallery in the MongoDB database
// collection.
func Update(context interface{}, db *db.DB, id string, gallery *Gallery) error {
	log.Dev(context, "Update", "Started")

	if !bson.IsObjectIdHex(id) {
		log.Error(context, "Update", ErrInvalidID, "Completed")
		return ErrInvalidID
	}

	if err := gallery.Validate(); err != nil {
		log.Error(context, "Update", err, "Completed")
		return err
	}

	objectID := bson.ObjectIdHex(id)

	// update the DateUpdated timestamp.
	gallery.DateUpdated = time.Now()

	f := func(c *mgo.Collection) error {
		log.Dev(context, "Update", "MGO : db.%s.update(%s, %s)", c.Name, mongo.Query(objectID), mongo.Query(gallery))
		return c.UpdateId(objectID, gallery)
	}

	if err := db.ExecuteMGO(context, GalleryCollection, f); err != nil {
		log.Error(context, "Update", err, "Completed")
		return err
	}

	log.Dev(context, "Update", "Completed")
	return nil
}
