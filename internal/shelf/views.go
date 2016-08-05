package shelf

import (
	"fmt"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/ardanlabs/kit/db"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

// AddView adds a view to the relationship manager.
func AddView(context interface{}, db *db.DB, view View) (string, error) {

	// Get the current relationship manager.
	rm, err := GetRelManager(context, db)
	if err != nil {
		return "", errors.Wrap(err, "Could not get the current relationship manager")
	}

	// Make sure the given view name does not exist already.
	var names []string
	for _, prevView := range rm.Views {
		names = append(names, prevView.Name)
	}
	if stringContains(names, view.Name) {
		return "", fmt.Errorf("View name already exists")
	}

	// Make sure that the relationships referenced in the view exist.
	var existingRels []string
	for _, existingRel := range rm.Relationships {
		existingRels = append(existingRels, existingRel.ID)
	}
	for _, segment := range view.Path {
		if !stringContains(existingRels, segment.RelationshipID) {
			return "", fmt.Errorf("Referenced relationship %s does not exist", segment.RelationshipID)
		}
	}

	// Assign a relationship ID, and add the relationship to the relationship manager.
	view.ID = uuid.NewV4().String()
	rm.Views = append(rm.Views, view)

	// Update the relationship manager.
	if err := NewRelManager(context, db, rm); err != nil {
		return "", errors.Wrap(err, "Could not update the relationship manager")
	}

	return view.ID, nil
}

// RemoveView removes a view from the relationship manager.
func RemoveView(context interface{}, db *db.DB, viewID string) error {
	f := func(c *mgo.Collection) error {
		q := bson.M{"id": 1}
		err := c.Update(q, bson.M{"$pull": bson.M{"views": bson.M{"id": viewID}}})
		return err
	}
	if err := db.ExecuteMGO(context, Collection, f); err != nil {
		return errors.Wrap(err, "Could not execute Mongo update statement")
	}

	return nil
}

// UpdateView updates a view in the relationship manager.
func UpdateView(context interface{}, db *db.DB, view View) error {

	// Validate the view.
	if err := view.Validate(); err != nil {
		return errors.Wrap(err, "Could not validate the provided view")
	}

	// Update the view.
	f := func(c *mgo.Collection) error {
		q := bson.M{"id": 1, "views.id": view.ID}
		err := c.Update(q, bson.M{"$set": bson.M{"views.$": &view}})
		return err
	}
	if err := db.ExecuteMGO(context, Collection, f); err != nil {
		return errors.Wrap(err, "Could not execute Mongo update statement")
	}

	return nil
}
