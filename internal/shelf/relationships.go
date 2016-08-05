package shelf

import (
	"fmt"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/ardanlabs/kit/db"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

// AddRelationship adds a relationship to the relationship manager.
func AddRelationship(context interface{}, db *db.DB, rel Relationship) (string, error) {

	// Get the current relationship manager.
	rm, err := GetRelManager(context, db)
	if err != nil {
		return "", errors.Wrap(err, "Could not get the current relationship manager")
	}

	// Make sure the given predicate does not exist already.
	var predicates []string
	for _, prevRel := range rm.Relationships {
		predicates = append(predicates, prevRel.Predicate)
	}
	if stringContains(predicates, rel.Predicate) {
		return "", fmt.Errorf("Predicate already exists")
	}

	// Assign a relationship ID, and add the relationship to the relationship manager.
	rel.ID = uuid.NewV4().String()
	rm.Relationships = append(rm.Relationships, rel)

	// Update the relationship manager.
	if err := NewRelManager(context, db, rm); err != nil {
		return "", errors.Wrap(err, "Could not update the relationship manager")
	}

	return rel.ID, nil
}

// RemoveRelationship removes a relationship from the relationship manager.
func RemoveRelationship(context interface{}, db *db.DB, relID string) error {

	// Get the current relationship manager.
	rm, err := GetRelManager(context, db)
	if err != nil {
		return errors.Wrap(err, "Could not get the current relationship manager")
	}

	// Make sure the given ID is not used in an active view.
	var relIDs []string
	for _, view := range rm.Views {
		for _, segment := range view.Path {
			relIDs = append(relIDs, segment.RelationshipID)
		}
	}
	if stringContains(relIDs, relID) {
		return fmt.Errorf("Active view is utilizing relationship %s", relID)
	}

	// Remove the relationship.
	f := func(c *mgo.Collection) error {
		q := bson.M{"id": 1}
		err := c.Update(q, bson.M{"$pull": bson.M{"relationships": bson.M{"id": relID}}})
		return err
	}
	if err := db.ExecuteMGO(context, Collection, f); err != nil {
		return errors.Wrap(err, "Could not execute Mongo update statement")
	}

	return nil
}

// UpdateRelationship updates a relationship in the relationship manager.
func UpdateRelationship(context interface{}, db *db.DB, rel Relationship) error {

	// Validate the relationship.
	if err := rel.Validate(); err != nil {
		return errors.Wrap(err, "Could not validate the provided relationship")
	}

	// Remove the relationship.
	f := func(c *mgo.Collection) error {
		q := bson.M{"id": 1, "relationships.id": rel.ID}
		err := c.Update(q, bson.M{"$set": bson.M{"relationships.$": &rel}})
		return err
	}
	if err := db.ExecuteMGO(context, Collection, f); err != nil {
		return errors.Wrap(err, "Could not execute Mongo update statement")
	}

	return nil
}
