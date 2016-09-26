package pattern

import (
	"errors"

	"github.com/ardanlabs/kit/db"
	"github.com/ardanlabs/kit/db/mongo"
	"github.com/ardanlabs/kit/log"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Collection is the Mongo collection containing pattern metadata.
const Collection = "patterns"

// ErrNotFound is an error variable thrown when no results are returned from a Mongo query.
var ErrNotFound = errors.New("Set Not found")

// Upsert upserts a pattern to the collection of currently utilized patterns.
func Upsert(context interface{}, db *db.DB, pattern *Pattern) error {
	log.Dev(context, "Upsert", "Started : Type[%s]", pattern.Type)

	// Validate the pattern.
	if err := pattern.Validate(); err != nil {
		log.Error(context, "Upsert", err, "Completed")
		return err
	}

	// Upsert the pattern.
	f := func(c *mgo.Collection) error {
		q := bson.M{"type": pattern.Type}
		log.Dev(context, "Upsert", "MGO : db.%s.upsert(%s, %s)", c.Name, mongo.Query(q), mongo.Query(pattern))
		_, err := c.Upsert(q, pattern)
		return err
	}
	if err := db.ExecuteMGO(context, Collection, f); err != nil {
		log.Error(context, "Upsert", err, "Completed")
		return err
	}

	log.Dev(context, "Upsert", "Completed")
	return nil
}

// GetAll retrieves the current patterns from Mongo.
func GetAll(context interface{}, db *db.DB) ([]Pattern, error) {
	log.Dev(context, "GetAll", "Started")

	// Get the relationships from Mongo.
	var patterns []Pattern
	f := func(c *mgo.Collection) error {
		log.Dev(context, "Find", "MGO : db.%s.find()", c.Name)
		return c.Find(nil).All(&patterns)
	}
	if err := db.ExecuteMGO(context, Collection, f); err != nil {
		if err == mgo.ErrNotFound {
			err = ErrNotFound
		}
		log.Error(context, "GetAll", err, "Completed")
		return nil, err
	}

	log.Dev(context, "GetAll", "Completed")
	return patterns, nil
}

// GetByType retrieves a pattern by type from Mongo.
func GetByType(context interface{}, db *db.DB, itemType string) (*Pattern, error) {
	log.Dev(context, "GetByType", "Started : Type[%s]", itemType)

	// Get the pattern from Mongo.
	var pattern Pattern
	f := func(c *mgo.Collection) error {
		q := bson.M{"type": itemType}
		log.Dev(context, "Find", "MGO : db.%s.find(%s)", c.Name, mongo.Query(q))
		return c.Find(q).One(&pattern)
	}
	if err := db.ExecuteMGO(context, Collection, f); err != nil {
		if err == mgo.ErrNotFound {
			err = ErrNotFound
		}
		log.Error(context, "GetByType", err, "Completed")
		return &pattern, err
	}

	log.Dev(context, "GetByType", "Completed")
	return &pattern, nil
}

// Delete removes a pattern from from Mongo.
func Delete(context interface{}, db *db.DB, itemType string) error {
	log.Dev(context, "Delete", "Started : Type[%s]", itemType)

	// Remove the relationship.
	f := func(c *mgo.Collection) error {
		q := bson.M{"type": itemType}
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
