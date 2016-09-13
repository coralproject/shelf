package pattern

import validator "gopkg.in/bluesuncorp/validator.v8"

//==============================================================================

// validate is used to perform model field validation.
var validate *validator.Validate

func init() {
	validate = validator.New(&validator.Config{TagName: "validate"})
}

//==============================================================================

// Inference includes information used to infer a particular relationship
// within an item.
type Inference struct {
	RelIDField string `bson:"related_ID_field" json:"related_ID_field" validate:"required,min=2"`
	Predicate  string `bson:"predicate" json:"predicate" validate:"required,min=2"`
	Direction  string `bson:"direction" json:"direction" validate:"required,min=2"`
	Required   bool   `bson:"required" json:"required"`
}

// Validate checks the Inference value for consistency.
func (inf *Inference) Validate() error {
	if err := validate.Struct(inf); err != nil {
		return err
	}
	return nil
}

// Pattern includes information used to infer relationships given an
// item of an certain type.
type Pattern struct {
	Type       string      `bson:"type" json:"type" validate:"required,min=2"`
	Inferences []Inference `bson:"inferences" json:"inferences" validate:"required,min=1"`
}

// Validate checks the Pattern value for consistency.
func (p *Pattern) Validate() error {
	if err := validate.Struct(p); err != nil {
		return err
	}

	for _, infer := range p.Inferences {
		if err := validate.Struct(infer); err != nil {
			return err
		}
	}
	return nil
}
