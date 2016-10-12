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

// Validate checks the pathsegment value for consistency.
func (ps *PathSegment) Validate() error {
	if err := validate.Struct(ps); err != nil {
		return err
	}
	return nil
}

// PathSegments is a slice of PathSegment values.
type PathSegments []PathSegment

// Len is required to sort a slice of PathSegment.
func (slice PathSegments) Len() int {
	return len(slice)
}

// Less is required to sort a slice of PathSegment.
func (slice PathSegments) Less(i, j int) bool {
	return slice[i].Level < slice[j].Level
}

// Swap is required to sort a slice of PathSegment.
func (slice PathSegments) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

// Path includes information defining one or multiple graph paths,
// along with a boolean choice for whether or not the path is a strict graph path.
type Path struct {
	StrictPath bool         `bson:"strict_path,omitempty" json:"strict_path,omitempty"`
	Segments   PathSegments `bson:"path_segments" json:"path_segments" validate:"required,min=1"`
}

// Validate checks the pathsegment value for consistency.
func (path *Path) Validate() error {
	if err := validate.Struct(path); err != nil {
		return err
	}
	return nil
}

// View contains metadata about a view.
type View struct {
	Name       string `bson:"name" json:"name" validate:"required,min=3"`
	Collection string `bson:"collection" json:"collection" validate:"required,min=2"`
	StartType  string `bson:"start_type" json:"start_type" validate:"required,min=3"`
	ReturnRoot bool   `bson:"return_root,omitempty" json:"return_root,omitempty"`
	Paths      []Path `bson:"paths" json:"paths" validate:"required,min=1"`
}

// Validate checks the View value for consistency.
func (v *View) Validate() error {

	// Validate the View value.
	if err := validate.Struct(v); err != nil {
		return err
	}

	// Validate each Path in the Paths.
	for _, path := range v.Paths {

		// Validate the Path using the validator.
		if err := path.Validate(); err != nil {
			return err
		}

		// Validate each of the PathSegment values in the Path.
		for _, segment := range path.Segments {

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
	}

	return nil
}
