package shelf

import (
	"fmt"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/ardanlabs/kit/db"
	"github.com/ardanlabs/kit/db/mongo"
	"github.com/ardanlabs/kit/log"
)

// AddView adds a view to the current views.
func AddView(context interface{}, db *db.DB, view View) (string, error) {
	log.Dev(context, "AddView", "Started")

	// Get the current relationships and views.
	rv, err := GetRelsAndViews(context, db)
	if err != nil {
		log.Error(context, "AddView", err, "Completed")
		return view.ID, err
	}

	// Make sure the given view name does not exist already.
	var names []string
	for _, prevView := range rv.Views {
		names = append(names, prevView.Name)
	}
	if stringContains(names, view.Name) {
		err = fmt.Errorf("View name already exists")
		log.Error(context, "AddView", err, "Completed")
		return view.ID, err
	}

	// Make sure that the relationships referenced in the view exist.
	var existingRels []string
	for _, existingRel := range rv.Relationships {
		existingRels = append(existingRels, existingRel.ID)
	}
	for _, segment := range view.Path {
		if !stringContains(existingRels, segment.RelationshipID) {
			err = fmt.Errorf("Referenced relationship %s does not exist", segment.RelationshipID)
			log.Error(context, "AddView", err, "Completed")
			return view.ID, err
		}
	}

	// Assign a relationship ID, if necessary.
	if view.ID == "" {
		viewID, err := newUUID()
		if err != nil {
			log.Error(context, "AddView", err, "Completed")
			return view.ID, err
		}
		view.ID = viewID
	}

	// Upsert the view.
	f := func(c *mgo.Collection) error {
		q := bson.M{"id": view.ID}
		log.Dev(context, "Upsert", "MGO : db.%s.upsert(%s, %s)", c.Name, mongo.Query(q), mongo.Query(view))
		_, err := c.Upsert(q, &view)
		return err
	}
	if err := db.ExecuteMGO(context, ViewCollection, f); err != nil {
		log.Error(context, "AddView", err, "Completed")
		return view.ID, err
	}

	log.Dev(context, "AddView", "Completed")
	return view.ID, nil
}

// RemoveView removes a view from the current views.
func RemoveView(context interface{}, db *db.DB, viewID string) error {
	log.Dev(context, "RemoveView", "Started")

	f := func(c *mgo.Collection) error {
		q := bson.M{"id": viewID}
		log.Dev(context, "Remove", "MGO : db.%s.remove(%s)", c.Name, mongo.Query(q))
		return c.Remove(q)
	}
	if err := db.ExecuteMGO(context, ViewCollection, f); err != nil {
		log.Error(context, "RemoveView", err, "Completed")
		return err
	}

	log.Dev(context, "RemoveView", "Completed")
	return nil
}

// UpdateView updates a view in the current views.
func UpdateView(context interface{}, db *db.DB, view View) error {
	log.Dev(context, "UpdateView", "Started")

	// Validate the view.
	if err := view.Validate(); err != nil {
		log.Error(context, "UpdateView", err, "Completed")
		return err
	}

	// Update the view.
	f := func(c *mgo.Collection) error {
		q := bson.M{"id": view.ID}
		log.Dev(context, "Update", "MGO : db.%s.update(%s, %s)", c.Name, mongo.Query(q), mongo.Query(view))
		return c.Update(q, &view)
	}
	if err := db.ExecuteMGO(context, ViewCollection, f); err != nil {
		log.Error(context, "UpdateView", err, "Completed")
		return err
	}

	log.Dev(context, "UpdateView", "Completed")
	return nil
}
