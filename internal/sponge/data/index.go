package data

import (
	"errors"
	"fmt"
	"strings"

	"github.com/ardanlabs/kit/db"
	"github.com/ardanlabs/kit/db/mongo"
	"github.com/ardanlabs/kit/log"
	"github.com/coralproject/sponge/pkg/item"
	mgo "gopkg.in/mgo.v2"
)

// =============================================================================

// EnsureTypeIndexes perform index create commands against Mongo for the indexes
// required by each item type.
func EnsureTypeIndexes(context interface{}, db *db.DB, types map[string]Type) error {
	log.Dev(context, "Item.EnsureTypeIndexes", "Started")

	var errStr string
	var indexes []mgo.Index

	// The type + _id index increase efficiency for queries looking for items of a type
	// for example {type: "comment", _id: "xzy"}.
	indexes = append(indexes, mgo.Index{
		Key:        []string{"type", "_id"},
		Unique:     true,
		DropDups:   false,
		Background: true,
		Sparse:     true,
	})

	// Index the "source identity" index: type + foreign id allowing for lookups
	// based on source ids.
	for _, t := range types {

		if t.IDField != "" {
			indexes = append(indexes, mgo.Index{
				Key:        []string{"t", "d." + t.IDField},
				Unique:     true,
				DropDups:   false,
				Background: true,
				Sparse:     true,
			})
		}
	}

	// Create the indexes. We can blindly ensure all indexes every time. Ensuring an existing
	// operation is a noop.
	for _, mgoIdx := range indexes {

		f := func(c *mgo.Collection) error {

			log.Dev(context, "EnsureTypeIndexes", "MGO : db.%s.ensureindex(%s)", c.Name, mongo.Query(mgoIdx))
			if err := c.EnsureIndex(mgoIdx); err != nil {
				log.Error(context, "EnsureIndexes", err, "Ensuring Index")
				errStr += fmt.Sprintf("[%s:%s] ", strings.Join(mgoIdx.Key, ","), err.Error())
			}

			return nil
		}

		if err := db.ExecuteMGO(context, item.Collection, f); err != nil {
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
