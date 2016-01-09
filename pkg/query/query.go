// Package query provides the service layer for building apps using
// query functionality.
package query

import (
	"errors"
	"strings"

	"github.com/ardanlabs/kit/db"
	"github.com/ardanlabs/kit/db/mongo"
	"github.com/ardanlabs/kit/log"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Contains the name of Mongo collections.
const (
	Collection         = "query_sets"
	CollectionHistory  = "query_sets_history"
	CollectionExecTest = "test_query"
)

// sets maintains the set of set related API calls.
type sets struct{}

// Sets fronts the access to the set functionality.
var Sets sets

// =============================================================================

// Upsert is used to create or update an existing Set document.
func (sets) Upsert(context interface{}, db *db.DB, set *Set) error {
	log.Dev(context, "sets.Upsert", "Started : Name[%s]", set.Name)

	// Validate the set that is provided.
	if err := set.Validate(); err != nil {
		log.Error(context, "sets.Upsert", err, "Completed")
		return err
	}

	// We need to know if this is a new set.
	var new bool
	if _, err := Sets.GetByName(context, db, set.Name); err != nil {
		if err != mgo.ErrNotFound {
			log.Error(context, "sets.Upsert", err, "Completed")
			return err
		}

		new = true
	}

	// Insert or update the query set.
	f := func(c *mgo.Collection) error {
		q := bson.M{"name": set.Name}
		log.Dev(context, "sets.Upsert", "MGO : db.%s.upsert(%s, %s)", c.Name, mongo.Query(q), mongo.Query(set))
		_, err := c.Upsert(q, set)
		return err
	}

	if err := db.ExecuteMGO(context, Collection, f); err != nil {
		log.Error(context, "sets.Upsert", err, "Completed")
		return err
	}

	// Add a history record if this query set is new.
	if new {
		f = func(c *mgo.Collection) error {
			qh := bson.M{
				"name": set.Name,
				"sets": []bson.M{},
			}

			log.Dev(context, "sets.Upsert", "MGO : db.%s.insert(%s)", c.Name, mongo.Query(qh))
			return c.Insert(qh)
		}

		if err := db.ExecuteMGO(context, CollectionHistory, f); err != nil {
			log.Error(context, "sets.Upsert", err, "Completed")
			return err
		}
	}

	// Add this query set to the beginning of the history.
	f = func(c *mgo.Collection) error {
		q := bson.M{"name": set.Name}
		qu := bson.M{
			"$push": bson.M{
				"sets": bson.M{
					"$each":     []*Set{set},
					"$position": 0,
				},
			},
		}

		log.Dev(context, "sets.Upsert", "MGO : db.%s.update(%s, %s)", c.Name, mongo.Query(q), mongo.Query(qu))
		_, err := c.Upsert(q, qu)
		return err
	}

	if err := db.ExecuteMGO(context, CollectionHistory, f); err != nil {
		log.Error(context, "sets.Upsert", err, "Completed")
		return err
	}

	log.Dev(context, "sets.Upsert", "Completed")
	return nil
}

// =============================================================================

// GetNames retrieves a list of query names.
func (sets) GetNames(context interface{}, db *db.DB) ([]string, error) {
	log.Dev(context, "sets.GetNames", "Started")

	var names []bson.M
	f := func(c *mgo.Collection) error {
		q := bson.M{"name": 1}
		log.Dev(context, "sets.GetNames", "MGO : db.%s.find({}, %s).sort([\"name\"])", c.Name, mongo.Query(q))
		return c.Find(nil).Select(q).Sort("name").All(&names)
	}

	if err := db.ExecuteMGO(context, Collection, f); err != nil {
		log.Error(context, "sets.GetNames", err, "Completed")
		return nil, err
	}

	var sets []string
	for _, doc := range names {
		name := doc["name"].(string)
		if strings.HasPrefix(name, "test") {
			continue
		}

		sets = append(sets, name)
	}

	log.Dev(context, "sets.GetNames", "Completed : Sets[%+v]", sets)
	return sets, nil
}

// GetByName retrieves the configuration for the specified Set.
func (sets) GetByName(context interface{}, db *db.DB, name string) (*Set, error) {
	log.Dev(context, "sets.GetByName", "Started : Name[%s]", name)

	var set Set
	f := func(c *mgo.Collection) error {
		q := bson.M{"name": name}
		log.Dev(context, "sets.GetByName", "MGO : db.%s.findOne(%s)", c.Name, mongo.Query(q))
		return c.Find(q).One(&set)
	}

	if err := db.ExecuteMGO(context, Collection, f); err != nil {
		log.Error(context, "sets.GetByName", err, "Completed")
		return nil, err
	}

	log.Dev(context, "sets.GetByName", "Completed : Set[%+v]", &set)
	return &set, nil
}

// GetLastHistoryByName gets the last written Set within the query_history
// collection and returns the last one else returns a non-nil error if it fails.
func (sets) GetLastHistoryByName(context interface{}, db *db.DB, name string) (*Set, error) {
	log.Dev(context, "sets.GetLastHistoryByName", "Started : Name[%s]", name)

	var result struct {
		Name string `bson:"name"`
		Sets []Set  `bson:"sets"`
	}

	f := func(c *mgo.Collection) error {
		q := bson.M{"name": name}
		proj := bson.M{"sets": bson.M{"$slice": 1}}

		log.Dev(context, "sets.GetLastHistoryByName", "MGO : db.%s.find(%s,%s)", c.Name, mongo.Query(q), mongo.Query(proj))
		return c.Find(q).Select(proj).One(&result)
	}

	err := db.ExecuteMGO(context, CollectionHistory, f)
	if err != nil {
		log.Error(context, "sets.GetLastHistoryByName", err, "Complete")
		return nil, err
	}

	if result.Sets == nil {
		err := errors.New("History not found")
		log.Error(context, "sets.GetLastHistoryByName", err, "Complete")
		return nil, err
	}

	log.Dev(context, "sets.GetLastHistoryByName", "Completed : QS[%+v]", &result.Sets[0])
	return &result.Sets[0], nil
}

// =============================================================================

// Delete is used to remove an existing Set document.
func (sets) Delete(context interface{}, db *db.DB, name string) error {
	log.Dev(context, "sets.Delete", "Started : Name[%s]", name)

	set, err := Sets.GetByName(context, db, name)
	if err != nil {
		return err
	}

	f := func(c *mgo.Collection) error {
		q := bson.M{"name": set.Name}
		log.Dev(context, "sets.Delete", "MGO : db.%s.remove(%s)", c.Name, mongo.Query(q))
		return c.Remove(q)
	}

	if err := db.ExecuteMGO(context, Collection, f); err != nil {
		log.Error(context, "sets.Delete", err, "Completed")
		return err
	}

	log.Dev(context, "sets.Delete", "Completed")
	return nil
}
