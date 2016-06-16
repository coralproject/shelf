package regex

import (
	"regexp"

	"gopkg.in/bluesuncorp/validator.v8"
)

//==============================================================================

// validate is used to perform model field validation.
var validate *validator.Validate

func init() {
	validate = validator.New(&validator.Config{TagName: "validate"})
}

//==============================================================================

// Regex contains a single regular expresion bound to a name.
type Regex struct {
	Name string `bson:"name" json:"name" validate:"required,min=3"`
	Expr string `bson:"expr" json:"expr" validate:"required,min=3"`

	Compile *regexp.Regexp
}

// Validate checks the regex value for consistency and that it compiles.
func (r Regex) Validate() error {
	if err := validate.Struct(r); err != nil {
		return err
	}

	if _, err := regexp.Compile(r.Expr); err != nil {
		return err
	}

	return nil
}
