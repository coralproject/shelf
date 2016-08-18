package ask

import (
	"time"

	"github.com/ardanlabs/kit/db"
	"github.com/ardanlabs/kit/db/mongo"
	"github.com/ardanlabs/kit/log"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// FormCollection is the mongo collection where Form documents are saved.
const FormCollection = "forms"

// FormWidget describes a specific question being asked by the Form which is
// contained within a FormStep.
type FormWidget struct {
	ID          string      `json:"id" bson:"_id"`
	Type        string      `json:"type" bson:"type"`
	Identity    bool        `json:"identity" bson:"identity"`
	Component   string      `json:"component" bson:"component"`
	Title       string      `json:"title" bson:"title"`
	Description string      `json:"description" bson:"description"`
	Wrapper     interface{} `json:"wrapper" bson:"wrapper"`
	Props       interface{} `json:"props" bson:"props"`
}

// FormStep is a collection of FormWidget's.
type FormStep struct {
	ID      string       `json:"id" bson:"_id"`
	Name    string       `json:"name" bson:"name"`
	Widgets []FormWidget `json:"widgets" bson:"widgets"`
}

// FormStats describes the statistics being recorded by a specific Form.
type FormStats struct {
	Responses int `json:"responses" bson:"responses"`
}

// Form contains the conatical representation of a Form, containing all the
// Steps, and help text relating to completing the Form.
type Form struct {
	ID             bson.ObjectId `json:"id" bson:"_id"`
	Status         string        `json:"status" bson:"status"`
	Theme          interface{}   `json:"theme" bson:"theme"`
	Settings       interface{}   `json:"settings" bson:"settings"`
	Header         interface{}   `json:"header" bson:"header"`
	Footer         interface{}   `json:"footer" bson:"footer"`
	FinishedScreen interface{}   `json:"finishedScreen" bson:"finishedScreen"`
	Steps          []FormStep    `json:"steps" bson:"steps"`
	Stats          FormStats     `json:"stats" bson:"stats"`
	CreatedBy      interface{}   `json:"created_by" bson:"created_by"`
	UpdatedBy      interface{}   `json:"updated_by" bson:"updated_by"`
	DeletedBy      interface{}   `json:"deleted_by" bson:"deleted_by"`
	DateCreated    time.Time     `json:"date_created,omitempty" bson:"date_created,omitempty"`
	DateUpdated    time.Time     `json:"date_updated,omitempty" bson:"date_updated,omitempty"`
	DateDeleted    time.Time     `json:"date_deleted,omitempty" bson:"date_deleted,omitempty"`
}

// UpsertForm upserts the provided form into the MongoDB database collection.
func UpsertForm(context interface{}, db *db.DB, form *Form) error {
	log.Dev(context, "UpsertForm", "Started")

	// TODO: validate

	// if there was no ID provided, we should set one
	if form.ID == "" {
		form.ID = bson.NewObjectId()
	}

	f := func(c *mgo.Collection) error {
		q := bson.M{"id": form.ID}
		log.Dev(context, "UpsertForm", "MGO : db.%s.upsert(%s, %s)", c.Name, mongo.Query(q), mongo.Query(form))
		_, err := c.Upsert(q, form)
		return err
	}

	if err := db.ExecuteMGO(context, FormCollection, f); err != nil {
		log.Error(context, "UpsertForm", err, "Completed")
		return err
	}

	if _, err := UpdateFormStats(context, db, form.ID.Hex()); err != nil {
		log.Error(context, "UpsertForm", err, "Completed")
		return err
	}

	log.Dev(context, "UpsertForm", "Completed")
	return nil
}

// UpdateFormStats updates the FormStats on a given Form.
func UpdateFormStats(context interface{}, db *db.DB, id string) (*FormStats, error) {
	// TODO: implement
	return nil, nil
}

// UpdateFormStatus updates the forms status and returns the updated form from
// the MongodB database collection.
func UpdateFormStatus(context interface{}, db *db.DB, id, status string) (*Form, error) {
	log.Dev(context, "UpdateFormStatus", "Started")

	if !bson.IsObjectIdHex(id) {
		log.Error(context, "UpdateFormStatus", ErrInvalidID, "Completed")
		return nil, ErrInvalidID
	}

	objectID := bson.ObjectIdHex(id)

	f := func(c *mgo.Collection) error {
		m := bson.M{}
		log.Dev(context, "UpdateFormStatus", "MGO : db.%s.update(%s, %s)", c.Name, mongo.Query(objectID), mongo.Query(m))
		return c.UpdateId(objectID, m)
	}

	if err := db.ExecuteMGO(context, FormCollection, f); err != nil {
		log.Error(context, "UpdateFormStatus", err, "Completed")
		return nil, err
	}

	form, err := RetrieveForm(context, db, id)
	if err != nil {
		log.Error(context, "UpdateFormStatus", err, "Completed")
		return nil, err
	}

	log.Dev(context, "UpdateFormStatus", "Completed")
	return form, nil
}

// RetrieveForm retrieves the form from the MongodB database collection.
func RetrieveForm(context interface{}, db *db.DB, id string) (*Form, error) {
	log.Dev(context, "RetrieveForm", "Started")

	if !bson.IsObjectIdHex(id) {
		log.Error(context, "RetrieveForm", ErrInvalidID, "Completed")
		return nil, ErrInvalidID
	}

	objectID := bson.ObjectIdHex(id)

	var form Form
	f := func(c *mgo.Collection) error {
		log.Dev(context, "RetrieveForm", "MGO : db.%s.find(%s)", c.Name, mongo.Query(objectID))
		return c.FindId(objectID).One(&form)
	}

	if err := db.ExecuteMGO(context, FormCollection, f); err != nil {
		log.Error(context, "RetrieveForm", err, "Completed")
		return nil, err
	}

	log.Dev(context, "RetrieveForm", "Completed")
	return &form, nil
}

// RetrieveManyForms retrieves a list of forms from the MongodB database
// collection.
func RetrieveManyForms(context interface{}, db *db.DB, limit, skip int) ([]Form, error) {
	log.Dev(context, "RetrieveManyForms", "Started")

	var forms = make([]Form, 0)
	f := func(c *mgo.Collection) error {
		log.Dev(context, "RetrieveManyForms", "MGO : db.%s.find().limit(%d).skip(%d)", c.Name, limit, skip)
		return c.Find(nil).Limit(limit).Skip(skip).All(&forms)
	}

	if err := db.ExecuteMGO(context, FormCollection, f); err != nil {
		log.Error(context, "RetrieveManyForms", err, "Completed")
		return nil, err
	}

	log.Dev(context, "RetrieveManyForms", "Completed")
	return forms, nil
}

// DeleteForm removes the document matching the id provided from the MongoDB
// database collection.
func DeleteForm(context interface{}, db *db.DB, id string) error {
	log.Dev(context, "DeleteForm", "Started")

	if !bson.IsObjectIdHex(id) {
		log.Error(context, "DeleteForm", ErrInvalidID, "Completed")
		return ErrInvalidID
	}

	objectID := bson.ObjectIdHex(id)

	f := func(c *mgo.Collection) error {
		q := bson.M{"_id": objectID}
		log.Dev(context, "DeleteForm", "MGO : db.%s.remove(%s)", c.Name, mongo.Query(q))
		return c.Remove(q)
	}

	if err := db.ExecuteMGO(context, FormCollection, f); err != nil {
		log.Error(context, "DeleteForm", err, "Completed")
		return err
	}

	log.Dev(context, "DeleteForm", "Completed")
	return nil
}
