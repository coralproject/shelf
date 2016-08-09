package shelf

import (
	"errors"

	mgo "gopkg.in/mgo.v2"

	"github.com/ardanlabs/kit/db"
	"github.com/ardanlabs/kit/log"
)

const (
	// ViewCollection is the Mongo collection containing view metadata.
	ViewCollection = "views"
	// RelCollection is the Mongo collection containing relationship metadata.
	RelCollection = "relationships"
)

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

	// Insert or update the relationships and views.
	if err := upsertRelManager(context, db, rm); err != nil {
		log.Error(context, "NewRelManager", err, "Completed")
		return err
	}

	log.Dev(context, "NewRelManager", "Completed")
	return nil
}

// upsertRelManager upserts a relationship manager into Mongo.
func upsertRelManager(context interface{}, db *db.DB, rm RelManager) error {

	// Upsert the relationships.
	for _, rel := range rm.Relationships {
		if _, err := AddRelationship(context, db, rel); err != nil {
			return err
		}
	}

	// Upsert the views.
	for _, view := range rm.Views {
		if _, err := AddView(context, db, view); err != nil {
			return err
		}
	}

	return nil
}

// ClearRelManager clears a current relationship manager from Mongo.
func ClearRelManager(context interface{}, db *db.DB) error {
	log.Dev(context, "ClearRelManager", "Started")

	// Clear relationships.
	f := func(c *mgo.Collection) error {
		_, err := c.RemoveAll(nil)
		return err
	}
	if err := db.ExecuteMGO(context, RelCollection, f); err != nil {
		log.Error(context, "ClearRelManager", err, "Completed")
		return err
	}

	// Clear views.
	if err := db.ExecuteMGO(context, ViewCollection, f); err != nil {
		log.Error(context, "ClearRelManager", err, "Completed")
		return err
	}

	log.Dev(context, "ClearRelManager", "Completed")
	return nil
}

// GetRelManager retrieves the current relationship manager from Mongo.
func GetRelManager(context interface{}, db *db.DB) (RelManager, error) {
	log.Dev(context, "GetRelManager", "Started")

	var rels []Relationship
	var views []View
	relFunc := func(c *mgo.Collection) error {
		return c.Find(nil).All(&rels)
	}
	viewFunc := func(c *mgo.Collection) error {
		return c.Find(nil).All(&views)
	}

	if err := db.ExecuteMGO(context, RelCollection, relFunc); err != nil {
		if err == mgo.ErrNotFound {
			err = ErrNotFound
		}
		log.Error(context, "GetRelManager", err, "Completed")
		return RelManager{}, err
	}
	if err := db.ExecuteMGO(context, ViewCollection, viewFunc); err != nil {
		if err == mgo.ErrNotFound {
			err = ErrNotFound
		}
		log.Error(context, "GetRelManager", err, "Completed")
		return RelManager{}, err
	}

	rm := RelManager{
		Relationships: rels,
		Views:         views,
	}

	log.Dev(context, "GetRelManager", "Completed")
	return rm, nil
}
