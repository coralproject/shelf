package shelf

import (
	"github.com/cayleygraph/cayley/graph/path"
	validator "gopkg.in/bluesuncorp/validator.v8"
)

//==============================================================================

// validate is used to perform model field validation.
var validate *validator.Validate

func init() {
	validate = validator.New(&validator.Config{TagName: "validate"})
}

//==============================================================================

// Relationship contains metadata about a relationship.
type Relationship struct {
	SubjectType string `bson:"subject" json:"subject" validate:"required,min=3"`
	Predicate   string `bson:"predicate" json:"predicate" validate:"required,min=3"`
	ObjectType  string `bson:"object" json:"object" validate:"required,min=3"`
	InString    string `bson:"in_string,omitempty" json:"in_string,omitempty"`
	OutString   string `bson:"out_string,omitempty" json:"out_string,omitempty"`
}

// Validate checks the Relationship value for consistency.
func (r Relationship) Validate() error {

	if err := validate.Struct(r); err != nil {
		return err
	}
	return nil
}

// View contains metadata about a view.
type View struct {
	Name          string         `bson:"name" json:"name" validate:"required,min=3"`
	Relationships []Relationship `bson:"relationships" json:"relationships"`
	Path          path.Path      `bson:"path" json:"path"`
}

// Validate checks the View value for consistency.
func (v View) Validate() error {

	if err := validate.Struct(v); err != nil {
		return err
	}
	return nil
}
