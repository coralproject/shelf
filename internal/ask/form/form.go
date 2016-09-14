package form

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

// Collection is the mongo collection where Form documents are saved.
const Collection = "forms"

// Widget describes a specific question being asked by the Form which is
// contained within a Step.
type Widget struct {
	ID          string      `json:"id" bson:"_id"`
	Type        string      `json:"type" bson:"type"`
	Identity    bool        `json:"identity" bson:"identity"`
	Component   string      `json:"component" bson:"component"`
	Title       string      `json:"title" bson:"title"`
	Description string      `json:"description" bson:"description"`
	Wrapper     interface{} `json:"wrapper" bson:"wrapper"`
	Props       interface{} `json:"props" bson:"props"`
}

// Step is a collection of Widget's.
type Step struct {
	ID      string   `json:"id" bson:"_id"`
	Name    string   `json:"name" bson:"name"`
	Widgets []Widget `json:"widgets" bson:"widgets"`
}

// Stats describes the statistics being recorded by a specific Form.
type Stats struct {
	Responses int `json:"responses" bson:"responses"`
}

//==============================================================================

// Form contains the conatical representation of a Form, containing all the
// Steps, and help text relating to completing the Form.
type Form struct {
	ID             bson.ObjectId          `json:"id" bson:"_id" validate:"required"`
	Status         string                 `json:"status" bson:"status"`
	Theme          interface{}            `json:"theme" bson:"theme"`
	Settings       map[string]interface{} `json:"settings" bson:"settings"`
	Header         interface{}            `json:"header" bson:"header"`
	Footer         interface{}            `json:"footer" bson:"footer"`
	FinishedScreen interface{}            `json:"finishedScreen" bson:"finishedScreen"`
	Steps          []Step                 `json:"steps" bson:"steps"`
	Stats          Stats                  `json:"stats" bson:"stats"`
	CreatedBy      interface{}            `json:"created_by" bson:"created_by"`
	UpdatedBy      interface{}            `json:"updated_by" bson:"updated_by"`
	DeletedBy      interface{}            `json:"deleted_by" bson:"deleted_by"`
	DateCreated    time.Time              `json:"date_created,omitempty" bson:"date_created,omitempty"`
	DateUpdated    time.Time              `json:"date_updated,omitempty" bson:"date_updated,omitempty"`
	DateDeleted    time.Time              `json:"date_deleted,omitempty" bson:"date_deleted,omitempty"`
}

// Validate checks the Form value for consistency.
func (f *Form) Validate() error {
	if err := validate.Struct(f); err != nil {
		return err
	}

	return nil
}

//==============================================================================

// Upsert upserts the provided form into the MongoDB database collection.
func Upsert(context interface{}, db *db.DB, form *Form) error {
	log.Dev(context, "Upsert", "Started")

	var isNewForm bool

	// If there is no ID probided, we should set one as this is an Upsert
	// operation. It is also important to remember if this was a new form or not
	// because we need to update the stats if this wasn't a new form.
	if form.ID.Hex() == "" {
		form.ID = bson.NewObjectId()
		isNewForm = true
		form.DateCreated = time.Now()
	}

	form.DateUpdated = time.Now()

	if err := form.Validate(); err != nil {
		log.Error(context, "Upsert", err, "Completed")
		return err
	}

	f := func(c *mgo.Collection) error {
		log.Dev(context, "Upsert", "MGO : db.%s.upsert(%s, %s)", c.Name, mongo.Query(form.ID.Hex()), mongo.Query(form))
		_, err := c.UpsertId(form.ID, form)
		return err
	}

	if err := db.ExecuteMGO(context, Collection, f); err != nil {
		log.Error(context, "Upsert", err, "Completed")
		return err
	}

	// New forms don't have any stats so don't bother updating it.
	if !isNewForm {
		if _, err := UpdateStats(context, db, form.ID.Hex()); err != nil {
			log.Error(context, "Upsert", err, "Completed")
			return err
		}
	}

	log.Dev(context, "Upsert", "Completed")
	return nil
}

// UpdateStats updates the Stats on a given Form.
func UpdateStats(context interface{}, db *db.DB, id string) (*Stats, error) {
	log.Dev(context, "UpdateStats", "Started : Form[%s]", id)

	if !bson.IsObjectIdHex(id) {
		log.Error(context, "UpdateStats", ErrInvalidID, "Completed")
		return nil, ErrInvalidID
	}

	objectID := bson.ObjectIdHex(id)

	count, err := submission.Count(context, db, id)
	if err != nil {
		log.Error(context, "UpdateStats", ErrInvalidID, "Completed")
		return nil, err
	}

	stats := Stats{
		Responses: count,
	}
	f := func(c *mgo.Collection) error {
		u := bson.M{
			"$set": bson.M{
				"stats":        stats,
				"date_updated": time.Now(),
			},
		}
		log.Dev(context, "UpdateStats", "MGO : db.%s.update(%s, %s)", c.Name, mongo.Query(objectID), mongo.Query(u))
		return c.UpdateId(objectID, u)
	}

	if err := db.ExecuteMGO(context, Collection, f); err != nil {
		log.Error(context, "UpdateStats", err, "Completed")
		return nil, err
	}

	log.Dev(context, "UpdateStats", "Completed")
	return &stats, nil
}

// UpdateStatus updates the forms status and returns the updated form from
// the MongodB database collection.
func UpdateStatus(context interface{}, db *db.DB, id, status string) (*Form, error) {
	log.Dev(context, "UpdateStatus", "Started : Form[%s]", id)

	if !bson.IsObjectIdHex(id) {
		log.Error(context, "UpdateStatus", ErrInvalidID, "Completed")
		return nil, ErrInvalidID
	}

	objectID := bson.ObjectIdHex(id)

	f := func(c *mgo.Collection) error {
		u := bson.M{
			"$set": bson.M{
				"status":       status,
				"date_updated": time.Now(),
			},
		}
		log.Dev(context, "UpdateStatus", "MGO : db.%s.update(%s, %s)", c.Name, mongo.Query(objectID), mongo.Query(u))
		return c.UpdateId(objectID, u)
	}

	if err := db.ExecuteMGO(context, Collection, f); err != nil {
		log.Error(context, "UpdateStatus", err, "Completed")
		return nil, err
	}

	form, err := Retrieve(context, db, id)
	if err != nil {
		log.Error(context, "UpdateStatus", err, "Completed")
		return nil, err
	}

	log.Dev(context, "UpdateStatus", "Completed")
	return form, nil
}

// Retrieve retrieves the form from the MongodB database collection.
func Retrieve(context interface{}, db *db.DB, id string) (*Form, error) {
	log.Dev(context, "Retrieve", "Started : Form[%s]", id)

	if !bson.IsObjectIdHex(id) {
		log.Error(context, "Retrieve", ErrInvalidID, "Completed")
		return nil, ErrInvalidID
	}

	objectID := bson.ObjectIdHex(id)

	var form Form
	f := func(c *mgo.Collection) error {
		log.Dev(context, "Retrieve", "MGO : db.%s.find(%s)", c.Name, mongo.Query(objectID))
		return c.FindId(objectID).One(&form)
	}

	if err := db.ExecuteMGO(context, Collection, f); err != nil {
		log.Error(context, "Retrieve", err, "Completed")
		return nil, err
	}

	log.Dev(context, "Retrieve", "Completed")
	return &form, nil
}

// List retrieves a list of forms from the MongodB database collection.
func List(context interface{}, db *db.DB, limit, skip int) ([]Form, error) {
	log.Dev(context, "List", "Started")

	var forms = make([]Form, 0)
	f := func(c *mgo.Collection) error {
		log.Dev(context, "List", "MGO : db.%s.find().limit(%d).skip(%d)", c.Name, limit, skip)
		return c.Find(nil).Limit(limit).Skip(skip).All(&forms)
	}

	if err := db.ExecuteMGO(context, Collection, f); err != nil {
		log.Error(context, "List", err, "Completed")
		return nil, err
	}

	log.Dev(context, "List", "Completed")
	return forms, nil
}

// Delete removes the document matching the id provided from the MongoDB
// database collection.
func Delete(context interface{}, db *db.DB, id string) error {
	log.Dev(context, "Delete", "Started : Form[%s]", id)

	if !bson.IsObjectIdHex(id) {
		log.Error(context, "Delete", ErrInvalidID, "Completed")
		return ErrInvalidID
	}

	objectID := bson.ObjectIdHex(id)

	f := func(c *mgo.Collection) error {
		log.Dev(context, "Delete", "MGO : db.%s.remove(%s)", c.Name, mongo.Query(objectID))
		return c.RemoveId(objectID)
	}

	if err := db.ExecuteMGO(context, Collection, f); err != nil {
		log.Error(context, "Delete", err, "Completed")
		return err
	}

	log.Dev(context, "Delete", "Completed")
	return nil
}
