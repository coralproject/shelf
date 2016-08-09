package shelf

import (
	"fmt"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/ardanlabs/kit/db"
	"github.com/ardanlabs/kit/log"
)

// AddView adds a view to the current views.
func AddView(context interface{}, db *db.DB, view View) (string, error) {
	log.Dev(context, "AddView", "Started")

	// Get the current relationships and views.
	rm, err := GetRelsAndViews(context, db)
	if err != nil {
		log.Error(context, "AddView", err, "Completed")
		return view.ID, err
	}

	// Make sure the given view name does not exist already.
	var names []string
	for _, prevView := range rm.Views {
		names = append(names, prevView.Name)
	}
	if stringContains(names, view.Name) {
		log.Error(context, "AddView", err, "Completed")
		return view.ID, fmt.Errorf("View name already exists")
	}

	// Make sure that the relationships referenced in the view exist.
	var existingRels []string
	for _, existingRel := range rm.Relationships {
		existingRels = append(existingRels, existingRel.ID)
	}
	for _, segment := range view.Path {
		if !stringContains(existingRels, segment.RelationshipID) {
			log.Error(context, "AddView", err, "Completed")
			return view.ID, fmt.Errorf("Referenced relationship %s does not exist", segment.RelationshipID)
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
		err := c.Remove(q)
		return err
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
		err := c.Update(q, &view)
		return err
	}
	if err := db.ExecuteMGO(context, ViewCollection, f); err != nil {
		log.Error(context, "UpdateView", err, "Completed")
		return err
	}

	log.Dev(context, "UpdateView", "Completed")
	return nil
}
