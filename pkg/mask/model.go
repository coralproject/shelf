package mask

import (
	"fmt"

	"gopkg.in/bluesuncorp/validator.v8"
)

// Set of query types we expect to receive.
const (
	MaskRemove = "remove" // Field is removed.
	MaskAll    = "all"    // Everything is masked.
	MaskEmail  = "email"  // Email based masking.
	MaskRight  = "right"  // Mask everything except last n characters. Default 4.
	MaskLeft   = "left"   // Mask everything except first n characters. Default 4.
)

//==============================================================================

// validate is used to perform model field validation.
var validate *validator.Validate

func init() {
	validate = validator.New(&validator.Config{TagName: "validate"})
}

//==============================================================================

// Mask contains information about what needs to be masked.
type Mask struct {
	Collection string `bson:"collection" json:"collection" validate:"required,min=3"`
	Field      string `bson:"field" json:"field" validate:"required,min=3"`
	Type       string `bson:"type" json:"type" validate:"required,min=3"`
}

// Validate checks the set value for consistency.
func (m Mask) Validate() error {
	if err := validate.Struct(m); err != nil {
		return err
	}

	switch m.Type[0:3] {
	case MaskAll, MaskRemove[0:3], MaskEmail[0:3], MaskRight[0:3], MaskLeft[0:3]:
		return nil
	default:
		return fmt.Errorf("Invalid mask type %s", m.Type)
	}
}
