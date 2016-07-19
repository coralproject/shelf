package item

import (
	"errors"
	"strconv"

	"gopkg.in/mgo.v2/bson"

	"github.com/ardanlabs/kit/db"
)

//==============================================================================

var (
	ErrRelsTypesNotFound = errors.New("Could not retrieve RelTypes for item")
)

//==============================================================================

// Rel holds an item's relationship to another item
type Rel struct {
	Name string `bson:"n" json:"n"`   // Name of relationship
	Type string `bson:"t" json:"t"`   // Item Type of target
	Id   string `bson:"id" json:"id"` // Id of target (consider storing as native bson.Id?)
}

func typeof(v interface{}) string {
	switch t := v.(type) {
	case int:
		return "int"
	case int64:
		return "int64"
	case float64:
		return "float64"
	case string:
		return "string"
	//... etc
	default:
		_ = t
		return "unknown"
	}
}

//==============================================================================

func getDatumByKey(k string, d interface{}) interface{} {

	// get the data as a map for searching
	m := d.(map[string]interface{})

	return m[k]

}

func GetRelsByIdString(context interface{}, db *db.DB, idString string) (*[]Rel, error) {
	// can we make this into a valid bson ObjectId?
	id := bson.ObjectIdHex(idString)
	//	if err != nil {
	//		return nil, err
	//	}

	// if so, use the traditional GetById to find the item
	i, err := GetById(context, db, id)
	if err != nil {
		return nil, err
	}

	// and look up it's rels
	return GetRels(context, db, i)
}

// GetRels looks up an item's relationships and returns them
func GetRels(context interface{}, db *db.DB, i *Item) (*[]Rel, error) {

	var rels []Rel

	// get the rel types for this item's type
	rts := Types[i.Type].Rels

	// for each reltype
	for _, rt := range rts {

		// find the foreign key value in the item data
		fkv := getDatumByKey(rt.Field, i.Data)

		// if there is not value, skip this rel
		if fkv == nil {
			continue
		}

		// create the field path for the foreign key field
		fkf := "d." + Types[i.Type].IdField

		// try string or int keys force keys to strings
		//   todo, better handle keys, all should be bson?

		// try with the default value
		var q = bson.M{"t": rt.Type, fkf: fkv}
		items, err := GetByQuery(context, db, q)

		if len(*items) == 0 && typeof(fkv) == "int" {

			// try with a string value
			fkvs := strconv.Itoa(fkv.(int))
			var q = bson.M{"t": rt.Type, fkf: fkvs}
			items, err = GetByQuery(context, db, q)

		}

		if err != nil {
			// how should we handle not being able to look up related items?
			//  we probably don't want to prevent the insert without measures to recover the item
			//  although we would be causing data inconsistencies?  maybe a flag that
			//  relations need to be re-queried?
		}

		// for each item
		for _, i := range *items {

			// create the relationship
			r := Rel{
				Name: rt.Name,
				Type: i.Type,
				Id:   i.Id.Hex(),
			}

			// add it to the list
			rels = append(rels, r)
		}

	}

	return &rels, nil
}
