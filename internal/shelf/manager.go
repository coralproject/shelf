package shelf

import (
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/ardanlabs/kit/db"
	"github.com/ardanlabs/kit/log"
	"github.com/pkg/errors"
)

// Collection is the MongoDB collection housing metadata about relationships and views.
const Collection = "relationship_manager"

// Set of error variables.
var (
	ErrNotFound = errors.New("Set Not found")
)

// NewRelManager creates a new relationship manager, either with defaults
// or based on a provided JSON config.
func NewRelManager(context interface{}, db *db.DB, rm RelManager) error {
	log.Dev(context, "NewRelManager", "Started")

	// Validate the relationship manager.
	if err := rm.Validate(); err != nil {
		log.Error(context, "NewRelManager", err, "Completed")
		return errors.Wrap(err, "Could not validate the provided relationship manager")
	}

	// Insert or update the default relationship manager.
	if err := upsertRelManager(context, db, rm); err != nil {
		log.Error(context, "NewRelManager", err, "Completed")
		return errors.Wrap(err, "Could not upsert default relationship manager")
	}

	log.Dev(context, "NewRelManager", "Completed")
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
	log.Dev(context, "ClearRelManager", "Started")
	f := func(c *mgo.Collection) error {
		q := bson.M{"id": 1}
		err := c.Remove(q)
		return err
	}
	if err := db.ExecuteMGO(context, Collection, f); err != nil {
		log.Error(context, "ClearRelManager", err, "Completed")
		return errors.Wrap(err, "Could not execute Mongo remove statement")
	}
	log.Dev(context, "ClearRelManager", "Completed")
	return nil
}

// GetRelManager retrieves the current relationship manager from Mongo.
func GetRelManager(context interface{}, db *db.DB) (RelManager, error) {
	log.Dev(context, "GetRelManager", "Started")

	var rm RelManager
	f := func(c *mgo.Collection) error {
		return c.Find(bson.M{"id": 1}).One(&rm)
	}

	if err := db.ExecuteMGO(context, Collection, f); err != nil {
		if err == mgo.ErrNotFound {
			err = ErrNotFound
		}
		log.Error(context, "GetRelManager", err, "Completed")
		return RelManager{}, errors.Wrap(err, "Could not retrieve relationship manager")
	}

	log.Dev(context, "GetRelManager", "Completed")
	return rm, nil
}
