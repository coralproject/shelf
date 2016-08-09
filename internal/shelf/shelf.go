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

// NewRelsAndViews creates new relationships and views, based on input JSON.
func NewRelsAndViews(context interface{}, db *db.DB, rv RelsAndViews) error {
	log.Dev(context, "NewRelsAndViews", "Started")

	// Validate the RelsAndViews value.
	if err := rv.Validate(); err != nil {
		log.Error(context, "NewRelsAndViews", err, "Completed")
		return err
	}

	// Insert or update the relationships and views.
	if err := upsertRelsAndViews(context, db, rv); err != nil {
		log.Error(context, "NewRelsAndViews", err, "Completed")
		return err
	}

	log.Dev(context, "NewRelsAndViews", "Completed")
	return nil
}

// upsertRelsAndViews upserts relationships and views into Mongo.
func upsertRelsAndViews(context interface{}, db *db.DB, rv RelsAndViews) error {

	// Upsert the relationships.
	for _, rel := range rv.Relationships {
		if _, err := AddRelationship(context, db, rel); err != nil {
			return err
		}
	}

	// Upsert the views.
	for _, view := range rv.Views {
		if _, err := AddView(context, db, view); err != nil {
			return err
		}
	}

	return nil
}

// ClearRelsAndViews clears current relationships and views from Mongo.
func ClearRelsAndViews(context interface{}, db *db.DB) error {
	log.Dev(context, "ClearRelsAndViews", "Started")

	// Clear relationships.
	f := func(c *mgo.Collection) error {
		log.Dev(context, "Remove", "MGO : db.%s.remove({})", c.Name)
		_, err := c.RemoveAll(nil)
		return err
	}
	if err := db.ExecuteMGO(context, RelCollection, f); err != nil {
		log.Error(context, "ClearRelsAndViews", err, "Completed")
		return err
	}

	// Clear views.
	if err := db.ExecuteMGO(context, ViewCollection, f); err != nil {
		log.Error(context, "ClearRelsAndViews", err, "Completed")
		return err
	}

	log.Dev(context, "ClearRelsAndViews", "Completed")
	return nil
}

// GetRelsAndViews retrieves the current relationships and views from Mongo.
func GetRelsAndViews(context interface{}, db *db.DB) (RelsAndViews, error) {
	log.Dev(context, "GetRelsAndViews", "Started")

	// Get the relationships and views from Mongo.
	var rels []Relationship
	var views []View
	relFunc := func(c *mgo.Collection) error {
		log.Dev(context, "Find", "MGO : db.%s.find()", c.Name)
		return c.Find(nil).All(&rels)
	}
	viewFunc := func(c *mgo.Collection) error {
		log.Dev(context, "Find", "MGO : db.%s.find()", c.Name)
		return c.Find(nil).All(&views)
	}

	if err := db.ExecuteMGO(context, RelCollection, relFunc); err != nil {
		if err == mgo.ErrNotFound {
			err = ErrNotFound
		}
		log.Error(context, "GetRelsAndViews", err, "Completed")
		return RelsAndViews{}, err
	}
	if err := db.ExecuteMGO(context, ViewCollection, viewFunc); err != nil {
		if err == mgo.ErrNotFound {
			err = ErrNotFound
		}
		log.Error(context, "GetRelsAndViews", err, "Completed")
		return RelsAndViews{}, err
	}

	// Form a RelsAndViews value.
	rv := RelsAndViews{
		Relationships: rels,
		Views:         views,
	}

	log.Dev(context, "GetRelsAndViews", "Completed")
	return rv, nil
}
