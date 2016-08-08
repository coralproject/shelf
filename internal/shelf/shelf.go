package shelf

import (
	"errors"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/ardanlabs/kit/db"
	"github.com/ardanlabs/kit/log"
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
		return err
	}

	// Insert or update the default relationship manager.
	if err := upsertRelManager(context, db, rm); err != nil {
		log.Error(context, "NewRelManager", err, "Completed")
		return err
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
		return err
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
		return err
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
		return RelManager{}, err
	}

	log.Dev(context, "GetRelManager", "Completed")
	return rm, nil
}
