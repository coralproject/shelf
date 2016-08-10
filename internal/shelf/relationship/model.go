package relationship

import validator "gopkg.in/bluesuncorp/validator.v8"

//==============================================================================

// validate is used to perform model field validation.
var validate *validator.Validate

func init() {
	validate = validator.New(&validator.Config{TagName: "validate"})
}

//==============================================================================

// Relationship contains metadata about a relationship.
// Note, predicate should be unique.
type Relationship struct {
	SubjectTypes []string `bson:"subject_types" json:"subject_types" validate:"required,min=1"`
	Predicate    string   `bson:"predicate" json:"predicate" validate:"required,min=2"`
	ObjectTypes  []string `bson:"object_types" json:"object_types" validate:"required,min=1"`
	InString     string   `bson:"in_string,omitempty" json:"in_string,omitempty"`
	OutString    string   `bson:"out_string,omitempty" json:"out_string,omitempty"`
}

// Validate checks the Relationship value for consistency.
func (r *Relationship) Validate() error {
	if err := validate.Struct(r); err != nil {
		return err
	}
	return nil
}
