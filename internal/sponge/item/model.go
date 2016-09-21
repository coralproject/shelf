package item

import validator "gopkg.in/bluesuncorp/validator.v8"

//==============================================================================

// validate is used to perform model field validation.
var validate *validator.Validate

func init() {
	validate = validator.New(&validator.Config{TagName: "validate"})
}

//==============================================================================

// Item is data, properties and behavior associated with one of a comment,
// asset, action, etc. Regardless of type (comment, asset, etc.), all Items
// are formatted with an ID, Type, and Version, and asssociated data (which
// may differ greatly between items) is encoded into the Data interface.
type Item struct {
	ID      string      `bson:"item_id" json:"item_id" validate:"required,min=1"`
	Type    string      `bson:"type" json:"type" validate:"required,min=2"`
	Version int         `bson:"version" json:"version" validate:"required,min=1"`
	Data    interface{} `bson:"data" json:"data"`
}

// Validate validates an Item value with the validator.
func (item *Item) Validate() error {
	if err := validate.Struct(item); err != nil {
		return err
	}

	return nil
}
