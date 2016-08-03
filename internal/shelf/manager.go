package shelf

import (
	"encoding/json"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/ardanlabs/kit/db"
	"github.com/coralproject/xenia/internal/shelf/sfix"
	"github.com/pkg/errors"
)

// Collection is the MongoDB collection housing metadata about relationships and views.
const Collection = "relationship_manager"

// NewRelManager creates a new relationship manager, either with defaults
// or based on a provided JSON config.
func NewRelManager(context interface{}, db *db.DB, rm RelManager) error {

	// Check if the relationship manager provided is empty.  If so,
	// create a default relationship manager.
	if rm.Relationships == nil || rm.Views == nil {
		if err := defaultRelManager(context, db); err != nil {
			return errors.Wrap(err, "Could not create default relationship manager")
		}
	}

	return nil
}

// defaultRelManager creates a default relationship manager.
func defaultRelManager(context interface{}, db *db.DB) error {

	// Import the default relationship manager.
	raw, err := sfix.LoadDefaultRelManager()
	if err != nil {
		return errors.Wrap(err, "Could not get default relationship manager data")
	}
	var rm RelManager
	if err := json.Unmarshal(raw, &rm); err != nil {
		return errors.Wrap(err, "Could not unmarshal default relationship manager data")
	}

	// Validate the default relationship manager.
	if err := rm.Validate(); err != nil {
		return errors.Wrap(err, "Could not validate the default relationship manager")
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
