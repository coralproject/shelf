package relationship

import (
	"errors"

	"github.com/ardanlabs/kit/log"
	"github.com/coralproject/shelf/internal/platform/db"
	"github.com/coralproject/shelf/internal/platform/db/mongo"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Collection is the Mongo collection containing relationship metadata.
const Collection = "relationships"

// ErrNotFound is an error variable thrown when no results are returned from a Mongo query.
var ErrNotFound = errors.New("Relationship Not found")

// Upsert upserts a relationship to the collection of currently utilized relationships.
func Upsert(context interface{}, db *db.DB, rel *Relationship) error {
	log.Dev(context, "Upsert", "Started : Predicate[%s]", rel.Predicate)

	// Validate the relationship.
	if err := rel.Validate(); err != nil {
		log.Error(context, "Upsert", err, "Completed")
		return err
	}

	// Upsert the relationship.
	f := func(c *mgo.Collection) error {
		q := bson.M{"predicate": rel.Predicate}
		log.Dev(context, "Upsert", "MGO : db.%s.upsert(%s, %s)", c.Name, mongo.Query(q), mongo.Query(rel))
		_, err := c.Upsert(q, rel)
		return err
	}
	if err := db.ExecuteMGO(context, Collection, f); err != nil {
		log.Error(context, "Upsert", err, "Completed")
		return err
	}

	log.Dev(context, "Upsert", "Completed")
	return nil
}

// GetAll retrieves the current relationships from Mongo.
func GetAll(context interface{}, db *db.DB) ([]Relationship, error) {
	log.Dev(context, "GetAll", "Started")

	// Get the relationships from Mongo.
	var rels []Relationship
	f := func(c *mgo.Collection) error {
		log.Dev(context, "Find", "MGO : db.%s.find()", c.Name)
		return c.Find(nil).All(&rels)
	}
	if err := db.ExecuteMGO(context, Collection, f); err != nil {
		if err == mgo.ErrNotFound {
			err = ErrNotFound
		}
		log.Error(context, "GetAll", err, "Completed")
		return nil, err
	}

	log.Dev(context, "GetAll", "Completed")
	return rels, nil
}

// GetByPredicate retrieves a relationship by predicate from Mongo.
func GetByPredicate(context interface{}, db *db.DB, predicate string) (*Relationship, error) {
	log.Dev(context, "GetByPredicate", "Started : Predicate[%s]", predicate)

	// Get the relationship from Mongo.
	var rel Relationship
	f := func(c *mgo.Collection) error {
		q := bson.M{"predicate": predicate}
		log.Dev(context, "Find", "MGO : db.%s.find(%s)", c.Name, mongo.Query(q))
		return c.Find(q).One(&rel)
	}
	if err := db.ExecuteMGO(context, Collection, f); err != nil {
		if err == mgo.ErrNotFound {
			err = ErrNotFound
		}
		log.Error(context, "GetByPredicate", err, "Completed")
		return &rel, err
	}

	log.Dev(context, "GetByPredicate", "Completed")
	return &rel, nil
}

// Delete removes a relationship from from Mongo.
func Delete(context interface{}, db *db.DB, predicate string) error {
	log.Dev(context, "Delete", "Started : Predicate[%s]", predicate)

	// Remove the relationship.
	f := func(c *mgo.Collection) error {
		q := bson.M{"predicate": predicate}
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
