package shelf

import validator "gopkg.in/bluesuncorp/validator.v8"

//==============================================================================

// validate is used to perform model field validation.
var validate *validator.Validate

func init() {
	validate = validator.New(&validator.Config{TagName: "validate"})
}

//==============================================================================

// RelManager contains metadata about what relationships and views are currenlty
// being utilized in the system.
type RelManager struct {
	Relationships []Relationship `bson:"relationships" json:"relationships" validate:"required,min=1"`
	Views         []View         `bson:"views" json:"views" validate:"required,min=1"`
}

// Validate checks the RelManager value for consistency.
func (rm RelManger) Validate() error {

	if err := validate.Struct(rm); err != nil {
		return err
	}

	for _, rel := range rm.Relationships {
		if err := rel.Validate(); err != nil {
			return err
		}
	}
	for _, view := range rm.Views {
		if err := view.Validate(); err != nil {
			return err
		}
	}
	return nil
}

// Relationship contains metadata about a relationship.
// Note, predicate should be unique.
type Relationship struct {
	ID          string `bson:"id" json:"id" validate:"required, min=1"`
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
	ID        string        `bson:"id" json:"id" validate:"required,min=1"`
	Name      string        `bson:"name" json:"name" validate:"required,min=3"`
	StartType string        `bson:"start_type" json:"start_type" validate:"required,min=3"`
	Path      []PathSegment `bson:"path" json:"path" validate:"required,min=1"`
}

// PathSegment contains metadata about a segment of a path,
// which path partially defines a View.
type PathSegment struct {
	Direction      string `bson:"direction" json:"direction" validate:"required,min=2"`
	RelationshipID string `bson:"relationship_id" json:"relationship_id" validate:"required, min=1"`
	Tag            string `bson:"tag,omitempty" json:"tag,omitempty"`
}

// Validate checks the View value for consistency.
func (v View) Validate() error {

	if err := validate.Struct(v); err != nil {
		return err
	}

	for _, segment := range v.PathSegment {
		if err := segment.Validate(); err != nil {
			return err
		}
	}
	return nil
}
