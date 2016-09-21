// Package sponge handles all data flowing into the Coral Platform.
package sponge

import (
	"fmt"

	"github.com/ardanlabs/kit/db"
	"github.com/coralproject/shelf/internal/sponge/item"
)

//==============================================================================

// defaultSourceIDField the source id field to "id" by convention.
const defaultSourceIDField = "id"

// Data contains an arbitrary dataset to be converted into an item
type Data map[string]interface{}

//==============================================================================

// Itemize takes a data type and unstructured data packet and returns that data in item form.
// If the data corresponds to an item already in the data store, the returned item will
// have the existing item's _ids consistent with an update operation.
func Itemize(context interface{}, db *db.DB, typ string, ver int, dat Data) (item.Item, error) {

	// Create an new item.
	it := item.Item{}
	it.Version = ver
	it.Data = dat
	it.Type = typ

	// Does the source_id field exist in the data?
	idValue, ok := dat[defaultSourceIDField]

	// If a source id value is found, use it to compose the item_id for this item.
	if ok {

		itemID, err := makeItemID(typ, idValue)
		if err != nil {
			return item.Item{}, err
		}

		it.ID = itemID
	}

	// If no source id provided, leave it.ID unset to be assigned upon Insert
	return it, nil

}

// makeItemID takes a type and source_id and composes a unique item_id for that pair
func makeItemID(typ string, sourceID interface{}) (string, error) {

	id := typ + "_" + fmt.Sprintf("%v", sourceID)

	return id, nil
}
