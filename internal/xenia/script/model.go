package script

import (
	"errors"

	"gopkg.in/bluesuncorp/validator.v8"
)

// validate is used to perform model field validation.
var validate *validator.Validate

func init() {
	validate = validator.New(&validator.Config{TagName: "validate"})
}

//==============================================================================

// Script contain pre and post commands to use per set or per query.
type Script struct {
	Name     string                   `bson:"name" json:"name" validate:"required,min=3"` // Unique name per Script document
	Commands []map[string]interface{} `bson:"commands" json:"commands"`                   // Commands to add to a query.
}

// Validate checks the query value for consistency.
func (scr Script) Validate() error {
	if err := validate.Struct(scr); err != nil {
		return err
	}

	if len(scr.Commands) == 0 {
		return errors.New("No commands exist")
	}

	return nil
}

// PrepareForInsert replaces the `$` to `_$` when found in the front of field names.
func (scr Script) PrepareForInsert() {

	// Fix the commands so it can be inserted.
	for c := range scr.Commands {
		prepareForInsert(scr.Commands[c])
	}
}

// PrepareForUse replaces the `_$` to `$` when found in the front of field names.
func (scr Script) PrepareForUse() {

	// Fix the commands so it can be inserted.
	for c := range scr.Commands {
		prepareForUse(scr.Commands[c])
	}
}
