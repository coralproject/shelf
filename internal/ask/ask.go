package ask

import (
	"errors"

	validator "gopkg.in/bluesuncorp/validator.v8"
	mgo "gopkg.in/mgo.v2"

	"github.com/ardanlabs/kit/db"
	"github.com/ardanlabs/kit/db/mongo"
	"github.com/ardanlabs/kit/log"
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

// EnsureIndexes perform index create commands against Mongo for the indexes
// needed for the ask package to run.
func EnsureIndexes(context interface{}, db *db.DB) error {
	log.Dev(context, "EnsureIndexes", "Started")

	f := func(c *mgo.Collection) error {
		index := mgo.Index{
			Key:        []string{"$text:replies.answer"},
			Unique:     false,
			DropDups:   false,
			Background: false,
			Sparse:     true,
			Name:       "replies.answer.text",
		}
		log.Dev(context, "EnsureIndexes", "MGO : db.%s.ensureIndex(%s)", c.Name, mongo.Query(index))
		return c.EnsureIndex(index)
	}

	if err := db.ExecuteMGO(context, FormSubmissionsCollection, f); err != nil {
		log.Error(context, "EnsureIndexes", err, "Completed")
		return err
	}

	log.Dev(context, "EnsureIndexes", "Completed")
	return nil
}

// Upsert upserts the provided form into the MongoDB database collection.
func Upsert(context interface{}, db *db.DB, form *Form) error {
	log.Dev(context, "Upsert", "Started")

	var isNewForm bool

	// if there was no ID provided, we should set one
	if form.ID == "" {
		isNewForm = true
	}

	err := UpsertForm(context, db, form)
	if err != nil {
		log.Error(context, "Upsert", err, "Completed")
		return err
	}

	if isNewForm {
		_, err := CreateFormGallery(context, db, form.ID.Hex())
		if err != nil {
			log.Error(context, "Upsert", err, "Completed")
			return err
		}
	}

	log.Dev(context, "Upsert", "Completed")
	return nil
}
