package ask

import (
	"time"

	"github.com/ardanlabs/kit/db"
	"github.com/ardanlabs/kit/db/mongo"
	"github.com/ardanlabs/kit/log"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// FormGalleryCollection is the mongo collection where FormGallery documents are
// saved.
const FormGalleryCollection = "form_galleries"

// FormGalleryAnswer describes an answer from a form which has been added to a
// FormGallery.
type FormGalleryAnswer struct {
	SubmissionID    bson.ObjectId          `json:"submission_id" bson:"submission_id" validate:"required,len=24"`
	AnswerID        string                 `json:"answer_id" bson:"answer_id" validate:"required,len=24"`
	Answer          FormSubmissionAnswer   `json:"answer,omitempty" bson:"-"`
	IdentityAnswers []FormSubmissionAnswer `json:"identity_answers,omitempty" bson:"-"`
}

// FormGallery is a Form that has been moved to a shared space.
type FormGallery struct {
	ID          bson.ObjectId          `json:"id" bson:"_id" validate:"required,len=24"`
	FormID      bson.ObjectId          `json:"form_id" bson:"form_id" validate:"required,len=24"`
	Headline    string                 `json:"headline" bson:"headline"`
	Description string                 `json:"description" bson:"description"`
	Config      map[string]interface{} `json:"config" bson:"config"`
	Answers     []FormGalleryAnswer    `json:"answers" bson:"answers"`
	DateCreated time.Time              `json:"date_created,omitempty" bson:"date_created,omitempty"`
	DateUpdated time.Time              `json:"date_updated,omitempty" bson:"date_updated,omitempty"`
}

// Validate checks the FormGallery value for consistency.
func (fg *FormGallery) Validate() error {
	if err := validate.Struct(fg); err != nil {
		return err
	}

	return nil
}

// CreateFormGallery adds a form gallery based on the form id provided into the
// MongoDB database collection.
func CreateFormGallery(context interface{}, db *db.DB, formID string) (*FormGallery, error) {
	log.Dev(context, "CreateFormGallery", "Started")

	if !bson.IsObjectIdHex(formID) {
		log.Error(context, "CreateFormGallery", ErrInvalidID, "Completed")
		return nil, ErrInvalidID
	}

	gallery := FormGallery{
		ID:          bson.NewObjectId(),
		FormID:      bson.ObjectIdHex(formID),
		DateCreated: time.Now(),
		DateUpdated: time.Now(),
	}

	f := func(c *mgo.Collection) error {
		log.Dev(context, "CreateFormGallery", "MGO : db.%s.insert(%s)", c.Name, mongo.Query(gallery))
		return c.Insert(gallery)
	}

	if err := db.ExecuteMGO(context, FormGalleryCollection, f); err != nil {
		log.Error(context, "CreateFormGallery", err, "Completed")
		return nil, err
	}

	log.Dev(context, "CreateFormGallery", "Completed")
	return nil, nil
}

// RetrieveFormGallery retrieves a form gallery from the MongoDB database
// collection as well as hydrating the form gallery with form submissions.
func RetrieveFormGallery(context interface{}, db *db.DB, id string) (*FormGallery, error) {
	log.Dev(context, "RetrieveFormGallery", "Started")

	if !bson.IsObjectIdHex(id) {
		log.Error(context, "RetrieveFormGallery", ErrInvalidID, "Completed")
		return nil, ErrInvalidID
	}

	objectID := bson.ObjectIdHex(id)

	var gallery FormGallery
	f := func(c *mgo.Collection) error {
		log.Dev(context, "RetrieveFormGallery", "MGO : db.%s.find(%s)", c.Name, mongo.Query(objectID))
		return c.FindId(objectID).One(&gallery)
	}

	if err := db.ExecuteMGO(context, FormGalleryCollection, f); err != nil {
		log.Error(context, "RetrieveFormGallery", err, "Completed")
		return nil, err
	}

	if err := HydrateFormGallery(context, db, &gallery); err != nil {
		log.Error(context, "RetrieveFormGallery", err, "Completed")
		return nil, err
	}

	log.Dev(context, "RetrieveFormGallery", "Completed")
	return &gallery, nil
}

// HydrateFormGallery loads a FormGallery with form submissions from the MongoDB
// database collection.
func HydrateFormGallery(context interface{}, db *db.DB, gallery *FormGallery) error {
	log.Dev(context, "HydrateFormGallery", "Started")

	if err := gallery.Validate(); err != nil {
		log.Error(context, "HydrateFormGallery", err, "Completed")
		return err
	}

	// load all the submission id's from the answers inside the gallery
	var submissionIDs = make([]string, len(gallery.Answers))
	for i, answer := range gallery.Answers {
		// and set their hex
		submissionIDs[i] = answer.SubmissionID.Hex()
	}

	// so we can fetch all the form submissions in one request
	submissions, err := RetrieveFormSubmissions(context, db, submissionIDs)
	if err != nil {
		log.Error(context, "HydrateFormGallery", err, "Completed")
		return err
	}

	// merge the submissions into the given gallery
	MergeSubmissionsIntoGalleryAnswers(gallery, submissions)

	log.Dev(context, "HydrateFormGallery", "Completed")
	return nil
}

// MergeSubmissionsIntoGalleryAnswers associates the array of submissions onto
// matching gallery answers.
func MergeSubmissionsIntoGalleryAnswers(gallery *FormGallery, submissions []FormSubmission) {
	// walk through all their answers
	for j, answer := range gallery.Answers {

		// and for each submission
		for k, submission := range submissions {
			// if we are looking at a different submission that doesn't match the
			// answer's submission ID or the submission was to a different form that
			// the current gallery is on, then we need to skip this submission.
			if submission.ID != answer.SubmissionID || submission.FormID.Hex() != gallery.FormID.Hex() {
				continue
			}

			// we have verified that the current submission is indeed for the
			// current gallery form and matches the submission id

			// so lets walk over the current submission's answers to match to the
			// particular question/widget/answer that we want to look at
			for _, submissionAnswer := range submission.Answers {
				// and if it doesn't match
				if submissionAnswer.WidgetID != answer.AnswerID {
					// then just continue
					continue
				}

				// but as it does

				// set the answer to the current submission answer
				gallery.Answers[j].Answer = submission.Answers[k]

				// and create an empty array for the identity answers that we will
				// walk over
				gallery.Answers[j].IdentityAnswers = make([]FormSubmissionAnswer, 0)

				// specifically, walk over the the current submission's answers again
				for m, submissionAnswer := range submission.Answers {
					// to find any identity answers related to this specific answer
					if submissionAnswer.Identity {
						// and append it to the list
						gallery.Answers[j].IdentityAnswers = append(gallery.Answers[j].IdentityAnswers, submission.Answers[m])
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

// HydrateFormGalleries loads an array of form galleries with form submissions
// from the MongoDB database collection.
func HydrateFormGalleries(context interface{}, db *db.DB, galleries []FormGallery) error {
	log.Dev(context, "HydrateFormGalleries", "Started")

	for _, gallery := range galleries {
		if err := gallery.Validate(); err != nil {
			log.Error(context, "HydrateFormGalleries", err, "Completed")
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
	submissions, err := RetrieveFormSubmissions(context, db, submissionIDs)
	if err != nil {
		log.Error(context, "HydrateFormGalleries", err, "Completed")
		return err
	}

	// for each of the galleries that we're hydrating
	for i := range galleries {
		// merge the submissions into the given gallery
		MergeSubmissionsIntoGalleryAnswers(&galleries[i], submissions)
	}

	log.Dev(context, "HydrateFormGalleries", "Completed")
	return nil
}

// AddFormGalleryAnswer adds an answer to a form gallery. Duplicated answers
// are de-duplicated automatically and will not return an error.
func AddFormGalleryAnswer(context interface{}, db *db.DB, id, submissionID, answerID string) (*FormGallery, error) {
	log.Dev(context, "AddAnswerToFormGallery", "Started")

	if !bson.IsObjectIdHex(id) {
		log.Error(context, "AddAnswerToFormGallery", ErrInvalidID, "Completed")
		return nil, ErrInvalidID
	}

	if !bson.IsObjectIdHex(submissionID) {
		log.Error(context, "AddAnswerToFormGallery", ErrInvalidID, "Completed")
		return nil, ErrInvalidID
	}

	if !bson.IsObjectIdHex(answerID) {
		log.Error(context, "AddAnswerToFormGallery", ErrInvalidID, "Completed")
		return nil, ErrInvalidID
	}

	objectID := bson.ObjectIdHex(id)

	answer := FormGalleryAnswer{
		SubmissionID: bson.ObjectIdHex(submissionID),
		AnswerID:     answerID,
	}

	f := func(c *mgo.Collection) error {
		u := bson.M{
			"$addToSet": bson.M{
				"answers": answer,
			},
		}
		log.Dev(context, "AddAnswerToFormGallery", "MGO : db.%s.update(%s, %s)", c.Name, mongo.Query(objectID), mongo.Query(u))
		return c.UpdateId(objectID, u)
	}

	if err := db.ExecuteMGO(context, FormGalleryCollection, f); err != nil {
		log.Error(context, "AddAnswerToFormGallery", err, "Completed")
		return nil, err
	}

	gallery, err := RetrieveFormGallery(context, db, id)
	if err != nil {
		log.Error(context, "AddAnswerToFormGallery", err, "Completed")
		return nil, err
	}

	log.Dev(context, "AddAnswerToFormGallery", "Completed")
	return gallery, nil
}

// RemoveFormGalleryAnswer adds an answer to a form gallery. Duplicated answers
// are de-duplicated automatically and will not return an error.
func RemoveFormGalleryAnswer(context interface{}, db *db.DB, id, submissionID, answerID string) (*FormGallery, error) {
	log.Dev(context, "RemoveFormGalleryAnswer", "Started")

	if !bson.IsObjectIdHex(id) {
		log.Error(context, "RemoveFormGalleryAnswer", ErrInvalidID, "Completed")
		return nil, ErrInvalidID
	}

	if !bson.IsObjectIdHex(submissionID) {
		log.Error(context, "RemoveFormGalleryAnswer", ErrInvalidID, "Completed")
		return nil, ErrInvalidID
	}

	if !bson.IsObjectIdHex(answerID) {
		log.Error(context, "RemoveFormGalleryAnswer", ErrInvalidID, "Completed")
		return nil, ErrInvalidID
	}

	objectID := bson.ObjectIdHex(id)

	answer := FormGalleryAnswer{
		SubmissionID: bson.ObjectIdHex(submissionID),
		AnswerID:     answerID,
	}

	f := func(c *mgo.Collection) error {
		u := bson.M{
			"$pull": bson.M{
				"answers": answer,
			},
		}
		log.Dev(context, "RemoveFormGalleryAnswer", "MGO : db.%s.update(%s, %s)", c.Name, mongo.Query(objectID), mongo.Query(u))
		return c.UpdateId(objectID, u)
	}

	if err := db.ExecuteMGO(context, FormGalleryCollection, f); err != nil {
		log.Error(context, "RemoveFormGalleryAnswer", err, "Completed")
		return nil, err
	}

	gallery, err := RetrieveFormGallery(context, db, id)
	if err != nil {
		log.Error(context, "RemoveFormGalleryAnswer", err, "Completed")
		return nil, err
	}

	log.Dev(context, "RemoveFormGalleryAnswer", "Completed")
	return gallery, nil
}

// RetrieveFormGalleriesForForm retrives the form galleries for a given form
// from the MongoDB database collection.
func RetrieveFormGalleriesForForm(context interface{}, db *db.DB, formID string) ([]FormGallery, error) {
	log.Dev(context, "RetrieveFormGalleriesForForm", "Started")

	if !bson.IsObjectIdHex(formID) {
		log.Error(context, "RetrieveFormGalleriesForForm", ErrInvalidID, "Completed")
		return nil, ErrInvalidID
	}

	formObjectID := bson.ObjectIdHex(formID)

	var galleries = make([]FormGallery, 0)
	f := func(c *mgo.Collection) error {
		q := bson.M{
			"form_id": formObjectID,
		}
		log.Dev(context, "RetrieveFormGalleriesForForm", "MGO : db.%s.find(%s)", c.Name, mongo.Query(q))
		return c.Find(q).All(&galleries)
	}

	if err := db.ExecuteMGO(context, FormGalleryCollection, f); err != nil {
		log.Error(context, "RetrieveFormGalleriesForForm", err, "Completed")
		return nil, err
	}

	if err := HydrateFormGalleries(context, db, galleries); err != nil {
		log.Error(context, "RetrieveFormGalleriesForForm", err, "Completed")
		return nil, err
	}

	log.Dev(context, "RetrieveFormGalleriesForForm", "Completed")
	return galleries, nil
}

// UpdateFormGallery updates the form gallery in the MongoDB database
// collection.
func UpdateFormGallery(context interface{}, db *db.DB, id string, gallery *FormGallery) error {
	log.Dev(context, "UpdateFormGallery", "Started")

	if !bson.IsObjectIdHex(id) {
		log.Error(context, "UpdateFormGallery", ErrInvalidID, "Completed")
		return ErrInvalidID
	}

	if err := gallery.Validate(); err != nil {
		log.Error(context, "UpdateFormGallery", err, "Completed")
		return err
	}

	objectID := bson.ObjectIdHex(id)

	// update the DateUpdated timestamp.
	gallery.DateUpdated = time.Now()

	f := func(c *mgo.Collection) error {
		log.Dev(context, "UpdateFormGallery", "MGO : db.%s.update(%s, %s)", c.Name, mongo.Query(objectID), mongo.Query(gallery))
		return c.UpdateId(objectID, gallery)
	}

	if err := db.ExecuteMGO(context, FormGalleryCollection, f); err != nil {
		log.Error(context, "UpdateFormGallery", err, "Completed")
		return err
	}

	log.Dev(context, "UpdateFormGallery", "Completed")
	return nil
}
