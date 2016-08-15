package item

import (
	"errors"

	"github.com/ardanlabs/kit/db"
	"github.com/ardanlabs/kit/db/mongo"
	"github.com/ardanlabs/kit/log"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Collection is the Mongo collection containing item values.
const Collection = "items"

// ErrNotFound is an error variable thrown when no results are returned from a Mongo query.
var ErrNotFound = errors.New("Set Not found")

// Upsert upserts an item to the items collections.
func Upsert(context interface{}, db *db.DB, item *Item) error {
	log.Dev(context, "Upsert", "Started : ID[%s]", item.ID)

	// Validate the item.
	if err := item.Validate(); err != nil {
		log.Error(context, "Upsert", err, "Completed")
		return err
	}

	// Upsert the item.
	f := func(c *mgo.Collection) error {
		q := bson.M{"item_id": item.ID}
		log.Dev(context, "Upsert", "MGO : db.%s.upsert(%s, %s)", c.Name, mongo.Query(q), mongo.Query(item))
		_, err := c.Upsert(q, item)
		return err
	}
	if err := db.ExecuteMGO(context, Collection, f); err != nil {
		log.Error(context, "Upsert", err, "Completed")
		return err
	}

	log.Dev(context, "Upsert", "Completed")
	return nil
}

// GetByIDs retrieves items by ID from Mongo.
func GetByIDs(context interface{}, db *db.DB, ids []string) ([]Item, error) {
	log.Dev(context, "GetByIDs", "Started : IDs%v", ids)

	// Get the items from Mongo.
	var items []Item
	f := func(c *mgo.Collection) error {
		q := bson.M{"item_id": bson.M{"$in": ids}}
		log.Dev(context, "Find", "MGO : db.%s.find(%s)", c.Name, mongo.Query(q))
		return c.Find(q).All(&items)
	}
	if err := db.ExecuteMGO(context, Collection, f); err != nil {
		if err == mgo.ErrNotFound {
			err = ErrNotFound
		}
		log.Error(context, "GetByIDs", err, "Completed")
		return items, err
	}

	log.Dev(context, "GetByIDs", "Completed")
	return items, nil
}

// Delete removes an item from from Mongo.
func Delete(context interface{}, db *db.DB, id string) error {
	log.Dev(context, "Delete", "Started : ID[%s]", id)

	// Remove the item.
	f := func(c *mgo.Collection) error {
		q := bson.M{"item_id": id}
		log.Dev(context, "Remove", "MGO : db.%s.remove(%s)", c.Name, mongo.Query(q))
		return c.Remove(q)
	}
	if err := db.ExecuteMGO(context, Collection, f); err != nil {
		log.Error(context, "Delete", err, "Completed")
		return err
	}

	log.Dev(context, "Delete", "Completed")
	return nil
}
