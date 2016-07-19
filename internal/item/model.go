package item

import (
	"time"

	validator "gopkg.in/bluesuncorp/validator.v8"
	"gopkg.in/mgo.v2/bson"
)

const (
	DefaultVersion = 1
)

//==============================================================================

// validate is used to perform model field validation.
var validate *validator.Validate

func init() {
	validate = validator.New(&validator.Config{TagName: "validate"})
}

//==============================================================================

// Result contains Items returned via a read operation and metadata
type Result struct {
	Items []Item    `bson:"items" json:"items"`
	Date  time.Time `bson:"date" json:"date"`
}

//==============================================================================

// ItemData is what an Item can hold
//  Should be the intersection of the db and transport protocols supported
// type ItemData map[string]interface{}

//==============================================================================

// An Item is data, properties and behavior wrapped in the thinnest
//  practical wrapper: Id, Type and Version
// this will be high volume so db and json field names are truncated
type Item struct {
	Id      bson.ObjectId          `bson:"_id" json:"id"`
	Type    string                 `bson:"t" json:"t"` // ItemType.Name
	Version int                    `bson:"v" json:"v"`
	Data    map[string]interface{} `bson:"d" json:"d"`
	Rels    []Rel                  `bson:"rels,omitempty" json:"rels,omitempty"`
}

func (i *Item) Validate() error {
	if err := validate.Struct(i); err != nil {
		return err
	}

	return nil
}
