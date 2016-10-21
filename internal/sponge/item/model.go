package item

import (
	"fmt"
	"time"

	validator "gopkg.in/bluesuncorp/validator.v8"
)

//==============================================================================

// validate is used to perform model field validation.
var validate *validator.Validate

func init() {
	validate = validator.New(&validator.Config{TagName: "validate"})
}

//==============================================================================

// sourceIDField is used to infer item.ID from type and data.
const sourceIDField = "id"

//==============================================================================

// Item is data, properties and behavior associated with one of a comment,
// asset, action, etc. Regardless of type (comment, asset, etc.), all Items
// are formatted with an ID, Type, and Version, and asssociated data (which
// may differ greatly between items) is encoded into the Data interface.
type Item struct {
	ID        string                 `bson:"item_id" json:"item_id" validate:"required,min=1"`
	Type      string                 `bson:"type" json:"type" validate:"required,min=2"`
	Version   int                    `bson:"version" json:"version" validate:"required,min=1"`
	Data      map[string]interface{} `bson:"data" json:"data"`
	CreatedAt time.Time              `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time              `bson:"updated_at" json:"updated_at"`
	Related   map[string]interface{} `bson:"related,omitempty" json:"related,omitempty"`
}

// Validate validates an Item value with the validator.
func (item *Item) Validate() error {
	if err := validate.Struct(item); err != nil {
		return err
	}

	return nil
}

// InferIDFromData infers an item_id from type and source id.
func (item *Item) InferIDFromData() error {

	// If a source id value is found, use it to compose the item_id for this item.
	idValue, ok := item.Data[sourceIDField]
	if !ok {
		return fmt.Errorf("Cannot Infer ID: Unable to find source id field: %s", sourceIDField)
	}

	// Set the ID via the template.
	item.ID = fmt.Sprintf("%s_%v", item.Type, idValue)

	return nil
}
