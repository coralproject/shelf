package item

import (
	"errors"
	"fmt"
	"strings"

	"github.com/ardanlabs/kit/db"
	"github.com/ardanlabs/kit/db/mongo"
	"github.com/ardanlabs/kit/log"

	"gopkg.in/mgo.v2"
)

// =============================================================================

// EnsureIndexes perform index create commands against Mongo for the indexes
// required by each item type
func EnsureTypeIndexes(context interface{}, db *db.DB, types map[string]Type) error {
	log.Dev(context, "Item.EnsureTypeIndexes", "Started")

	var errStr string
	var indexes []mgo.Index

	// create a list of indexes to create
	//  note, this list may be incomplete or too aggressive all
	//  indexes should reference the funcs / cases when they are needed

	// The type + _id index (type limits first)
	//  Used to find one item by type
	//  ? - does introducing type make this more efficient than the
	//      default _id index?
	indexes = append(indexes, mgo.Index{
		Key:        []string{"t", "_id"},
		Unique:     true,
		DropDups:   false,
		Background: true,
		Sparse:     true,
	})

	// for each type
	for _, t := range types {

		// Index the "source identity" index: type + foreign id
		// This index powers all the initial rel lookups
		if t.IdField != "" {
			indexes = append(indexes, mgo.Index{
				Key:        []string{"t", "d." + t.IdField},
				Unique:     true,
				DropDups:   false,
				Background: true,
				Sparse:     true,
			})
		}

		// Index the rel subdocs for "reverse relationship" queries
		//   Aka: all items that have a relationship with a given item
		//
		// Mongo only allows a single multi-key index
		//   Will this index be engaged in a subset of the rel fields is queried?
		indexes = append(indexes, mgo.Index{
			Key:        []string{"rels"},
			Unique:     false,
			DropDups:   false,
			Background: true,
			Sparse:     true,
		})

	}

	// create the indexes we figure we need
	for _, mgoIdx := range indexes {

		f := func(c *mgo.Collection) error {

			log.Dev(context, "EnsureTypeIndexes", "MGO : db.%s.ensureindex(%s)", c.Name, mongo.Query(mgoIdx))
			if err := c.EnsureIndex(mgoIdx); err != nil {
				log.Error(context, "EnsureIndexes", err, "Ensuring Index")
				errStr += fmt.Sprintf("[%s:%s] ", strings.Join(mgoIdx.Key, ","), err.Error())
			}

			return nil
		}

		if err := db.ExecuteMGO(context, Collection, f); err != nil {
			log.Error(context, "EnsureIndexes", err, "Completed")
			return err
		}
	}

	if errStr != "" {
		return errors.New(errStr)
	}

	log.Dev(context, "EnsureIndexes", "Completed")
	return nil
}
