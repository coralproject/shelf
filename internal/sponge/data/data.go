// Package data handles raw data, defines types and turns them into items.
package data

import (
	"errors"

	"github.com/ardanlabs/kit/db"
	"github.com/coralproject/shelf/internal/sponge/item"
	"gopkg.in/mgo.v2/bson"
)

//==============================================================================

// Data contains an arbitrary dataset to be converted into an item
type Data map[string]interface{}

// ErrTypeNotFound is an error variable thrown when the type to convert into is not recognized.
var ErrTypeNotFound = errors.New("Type Not found")

//==============================================================================

// Itemize takes a data type and unstructured data packet and returns that data in item form.
// If the data corresponds to an item already in the data store, the returned item will
// have the existing item's _ids consistent with an update operation.
func Itemize(context interface{}, db *db.DB, t string, v int, d Data) (item.Item, error) {
	i := item.Item{}

	// validate and set type
	if isRegistered(t) == false {
		// TODO, how can we use
		return i, errors.New("Type not found: " + t)
	}
	i.Type = t

	// This data may correspond to an item already present. Check the _source id_ to see
	// if there's a source key to look for in the item store.

	// Get this data's IdField value for this type.
	idValue := d[Types[t].IDField]

	// If a source id value is found, look to see if this item already exists.
	if idValue != nil {

		// Create a query referencing the source_id in data and type.
		dbIDField := "data." + Types[t].IDField
		q := bson.M{"type": t, dbIDField: idValue}

		// Look up items with the souce id and type
		dbItem, err := item.GetOneByQuery(context, db, q)
		if err != nil {
			return item.Item{}, err
		}

		// If we found an item, assign the existing id.
		if dbItem != nil {
			i.ID = dbItem.ID
		}
	}

	i.Version = v
	i.Data = d

	return i, nil

}
