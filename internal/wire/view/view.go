package view

import (
	"errors"

	"github.com/ardanlabs/kit/log"
	"github.com/coralproject/shelf/internal/platform/db"
	"github.com/coralproject/shelf/internal/platform/db/mongo"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Collection is the Mongo collection containing view metadata.
const Collection = "views"

// ErrNotFound is an error variable thrown when no results are returned from a Mongo query.
var ErrNotFound = errors.New("View Not found")

// Upsert upserts a view to the collection of currently utilized views.
func Upsert(context interface{}, db *db.DB, view *View) error {
	log.Dev(context, "Upsert", "Started : Name[%s]", view.Name)

	// Validate the view.
	if err := view.Validate(); err != nil {
		log.Error(context, "Upsert", err, "Completed")
		return err
	}

	// Upsert the view.
	f := func(c *mgo.Collection) error {
		q := bson.M{"name": view.Name}
		log.Dev(context, "Upsert", "MGO : db.%s.upsert(%s, %s)", c.Name, mongo.Query(q), mongo.Query(view))
		_, err := c.Upsert(q, view)
		return err
	}
	if err := db.ExecuteMGO(context, Collection, f); err != nil {
		log.Error(context, "Upsert", err, "Completed")
		return err
	}

	log.Dev(context, "Upsert", "Completed")
	return nil
}

// GetAll retrieves the current views from Mongo.
func GetAll(context interface{}, db *db.DB) ([]View, error) {
	log.Dev(context, "GetAll", "Started")

	// Get the views from Mongo.
	var views []View
	f := func(c *mgo.Collection) error {
		log.Dev(context, "Find", "MGO : db.%s.find()", c.Name)
		return c.Find(nil).All(&views)
	}
	if err := db.ExecuteMGO(context, Collection, f); err != nil {
		if err == mgo.ErrNotFound {
			err = ErrNotFound
		}
		log.Error(context, "GetAll", err, "Completed")
		return nil, err
	}

	log.Dev(context, "GetAll", "Completed")
	return views, nil
}

// GetByName retrieves a view by name from Mongo.
func GetByName(context interface{}, db *db.DB, name string) (*View, error) {
	log.Dev(context, "GetByName", "Started : Name[%s]", name)

	// Get the view from Mongo.
	var view View
	f := func(c *mgo.Collection) error {
		q := bson.M{"name": name}
		log.Dev(context, "Find", "MGO : db.%s.find(%s)", c.Name, mongo.Query(q))
		return c.Find(q).One(&view)
	}
	if err := db.ExecuteMGO(context, Collection, f); err != nil {
		if err == mgo.ErrNotFound {
			err = ErrNotFound
		}
		log.Error(context, "GetByName", err, "Completed")
		return &view, err
	}

	log.Dev(context, "GetByName", "Completed")
	return &view, nil
}

// Delete removes a view from from Mongo.
func Delete(context interface{}, db *db.DB, name string) error {
	log.Dev(context, "Delete", "Started : Name[%s]", name)

	// Remove the view.
	f := func(c *mgo.Collection) error {
		q := bson.M{"name": name}
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
