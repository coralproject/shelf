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
	SubmissionID    bson.ObjectId          `json:"submission_id" bson:"submission_id"`
	AnswerID        string                 `json:"answer_id" bson:"answer_id"`
	Answer          FormSubmissionAnswer   `json:"answer,omitempty" bson:"-"`
	IdentityAnswers []FormSubmissionAnswer `json:"identity_answers,omitempty" bson:"-"`
}

// FormGallery is a Form that has been moved to a shared space.
type FormGallery struct {
	ID          bson.ObjectId          `json:"id" bson:"_id"`
	FormID      bson.ObjectId          `json:"form_id" bson:"form_id"`
	Headline    string                 `json:"headline" bson:"headline"`
	Description string                 `json:"description" bson:"description"`
	Config      map[string]interface{} `json:"config" bson:"config"`
	Answers     []FormGalleryAnswer    `json:"answers" bson:"answers"`
	DateCreated time.Time              `json:"date_created,omitempty" bson:"date_created,omitempty"`
	DateUpdated time.Time              `json:"date_updated,omitempty" bson:"date_updated,omitempty"`
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
		log.Dev(context, "RetrieveFormGallery", "MGO : db.%s.find(%s)", c.Name, mongo.Query(objectID.Hex()))
		return c.FindId(objectID).One(&gallery)
	}

	if err := db.ExecuteMGO(context, FormGalleryCollection, f); err != nil {
		log.Error(context, "RetrieveFormGallery", err, "Completed")
		return nil, err
	}

	// TODO: hydrate the form galleries?

	log.Dev(context, "RetrieveFormGallery", "Completed")
	return &gallery, nil
}

// HydrateFormGallery loads a FormGallery with form submissions from the MongoDB
// database collection.
func HydrateFormGallery(context interface{}, db *db.DB, gallery *FormGallery) error {
	// TODO: implement
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
		log.Dev(context, "AddAnswerToFormGallery", "MGO : db.%s.update(%s, %s)", c.Name, mongo.Query(objectID.Hex()), mongo.Query(u))
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
		log.Dev(context, "RemoveFormGalleryAnswer", "MGO : db.%s.update(%s, %s)", c.Name, mongo.Query(objectID.Hex()), mongo.Query(u))
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

	// TODO: hydrate the form galleries?

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

	objectID := bson.ObjectIdHex(id)

	// update the DateUpdated timestamp.
	gallery.DateUpdated = time.Now()

	f := func(c *mgo.Collection) error {
		log.Dev(context, "UpdateFormGallery", "MGO : db.%s.update(%s, %s)", c.Name, mongo.Query(objectID.Hex()), mongo.Query(gallery))
		return c.UpdateId(objectID, gallery)
	}

	if err := db.ExecuteMGO(context, FormGalleryCollection, f); err != nil {
		log.Error(context, "UpdateFormGallery", err, "Completed")
		return err
	}

	log.Dev(context, "UpdateFormGallery", "Completed")
	return nil
}
