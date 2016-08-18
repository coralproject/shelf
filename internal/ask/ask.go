package ask

import (
	"errors"

	validator "gopkg.in/bluesuncorp/validator.v8"

	"github.com/ardanlabs/kit/db"
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
