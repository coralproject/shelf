package view

import (
	"fmt"

	validator "gopkg.in/bluesuncorp/validator.v8"
)

//==============================================================================

// validate is used to perform model field validation.
var validate *validator.Validate

func init() {
	validate = validator.New(&validator.Config{TagName: "validate"})
}

//==============================================================================

// PathSegment contains metadata about a segment of a path,
// which path partially defines a View.
type PathSegment struct {
	Level     int    `bson:"level" json:"level" validate:"required,min=1"`
	Direction string `bson:"direction" json:"direction" validate:"required,min=2"`
	Predicate string `bson:"predicate" json:"predicate" validate:"required,min=1"`
	Tag       string `bson:"tag,omitempty" json:"tag,omitempty"`
}

// Validate checks the PathSegment value for consistency.
func (ps *PathSegment) Validate() error {
	if err := validate.Struct(ps); err != nil {
		return err
	}
	return nil
}

// View contains metadata about a view.
type View struct {
	Name      string        `bson:"name" json:"name" validate:"required,min=3"`
	StartType string        `bson:"start_type" json:"start_type" validate:"required,min=3"`
	Path      []PathSegment `bson:"path" json:"path" validate:"required,min=1"`
}

// Validate checks the View value for consistency.
func (v *View) Validate() error {

	// Validate the View value.
	if err := validate.Struct(v); err != nil {
		return err
	}

	// Validate each of the PathSegment values in the View.
	for _, segment := range v.Path {

		// Validate the PathSegment using the validator.
		if err := segment.Validate(); err != nil {
			return err
		}

		// Ensure that the Direction value is either "in" or "out."
		switch segment.Direction {
		case "in", "out":
			continue
		default:
			return fmt.Errorf("Path segment includes undefined direction")
		}
	}
	return nil
}
