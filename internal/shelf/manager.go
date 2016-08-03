package shelf

import (
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/ardanlabs/kit/db"
	"github.com/pkg/errors"
)

// Collection is the MongoDB collection housing metadata about relationships and views.
const Collection = "relationship_manager"

// NewRelManager creates a new relationship manager, either with defaults
// or based on a provided JSON config.
func NewRelManager(context interface{}, db *db.DB, rm RelManager) error {

	// Validate the relationship manager.
	if err := rm.Validate(); err != nil {
		return errors.Wrap(err, "Could not validate the provided relationship manager")
	}

	// Insert or update the default relationship manager.
	if err := upsertRelManager(context, db, rm); err != nil {
		return errors.Wrap(err, "Could not upsert default relationship manager")
	}

	return nil
}

// upsertRelManager upserts a relationship manager into Mongo.
func upsertRelManager(context interface{}, db *db.DB, rm RelManager) error {
	f := func(c *mgo.Collection) error {
		q := bson.M{"id": 1}
		_, err := c.Upsert(q, &rm)
		return err
	}
	if err := db.ExecuteMGO(context, Collection, f); err != nil {
		return errors.Wrap(err, "Could not execute Mongo upsert statement")
	}
	return nil
}

// ClearRelManager clears a current relationship manager from Mongo.
func ClearRelManager(context interface{}, db *db.DB) error {
	f := func(c *mgo.Collection) error {
		q := bson.M{"id": 1}
		err := c.Remove(q)
		return err
	}
	if err := db.ExecuteMGO(context, Collection, f); err != nil {
		return errors.Wrap(err, "Could not execute Mongo remove statement")
	}
	return nil
}
