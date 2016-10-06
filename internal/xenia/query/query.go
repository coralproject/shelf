// Package query provides the service layer for building apps using
// query functionality.
package query

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/ardanlabs/kit/log"
	"github.com/coralproject/shelf/internal/platform/db"
	"github.com/coralproject/shelf/internal/platform/db/mongo"
	gc "github.com/patrickmn/go-cache"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Contains the name of Mongo collections.
const (
	Collection        = "query_sets"
	CollectionHistory = "query_sets_history"
)

// Set of error variables.
var (
	ErrNotFound = errors.New("Set Not found")
)

// =============================================================================

// c contans a cache of set values. The cache will maintain items for one
// second and then marked as expired. This is a very small cache so the
// gc time will be every hour.

const (
	expiration = time.Second
	cleanup    = time.Hour
)

var cache = gc.New(expiration, cleanup)

// =============================================================================

// EnsureIndexes perform index create commands against Mongo for the indexes
// specied in each query for the set. It will attempt to ensure all indexes
// regardless if one fails. Then reports all failures.
func EnsureIndexes(context interface{}, db *db.DB, set *Set) error {
	log.Dev(context, "EnsureIndexes", "Started : Name[%s]", set.Name)

	var errStr string

	for _, q := range set.Queries {
		if len(q.Indexes) == 0 {
			continue
		}

		f := func(c *mgo.Collection) error {
			for _, idx := range q.Indexes {
				mgoIdx := mgo.Index{
					Key:        idx.Key,
					Unique:     idx.Unique,
					DropDups:   idx.DropDups,
					Background: idx.Background,
					Sparse:     idx.Sparse,
				}

				log.Dev(context, "EnsureIndexes", "MGO : db.%s.ensureindex(%s)", c.Name, mongo.Query(mgoIdx))
				if err := c.EnsureIndex(mgoIdx); err != nil {
					log.Error(context, "EnsureIndexes", err, "Ensuring Index")
					errStr += fmt.Sprintf("[%s:%s] ", strings.Join(idx.Key, ","), err.Error())
				}
			}

			return nil
		}

		if err := db.ExecuteMGO(context, q.Collection, f); err != nil {
			log.Error(context, "EnsureIndexes", err, "Completed")
			return err
		}
	}

	if errStr != "" {
		return errors.New(errStr)
	}

	log.Dev(context, "EnsureIndexes", "Completed")
	return nil
}

// =============================================================================

// Upsert is used to create or update an existing Set document.
func Upsert(context interface{}, db *db.DB, set *Set) error {
	log.Dev(context, "Upsert", "Started : Name[%s]", set.Name)

	// Validate the set that is provided.
	if err := set.Validate(); err != nil {
		log.Error(context, "Upsert", err, "Completed")
		return err
	}

	// We need to know if this is a new set.
	var new bool
	if _, err := GetByName(context, db, set.Name); err != nil {
		if err != ErrNotFound {
			log.Error(context, "Upsert", err, "Completed")
			return err
		}

		new = true
	}

	// Fix the set so it can be inserted.
	set.PrepareForInsert()
	defer set.PrepareForUse()

	// Insert or update the query set.
	f := func(c *mgo.Collection) error {
		q := bson.M{"name": set.Name}
		log.Dev(context, "Upsert", "MGO : db.%s.upsert(%s, %s)", c.Name, mongo.Query(q), mongo.Query(set))
		_, err := c.Upsert(q, set)
		return err
	}

	if err := db.ExecuteMGO(context, Collection, f); err != nil {
		log.Error(context, "Upsert", err, "Completed")
		return err
	}

	// Flush the cache to invalidate everything.
	cache.Flush()

	// Add a history record if this query set is new.
	if new {
		f = func(c *mgo.Collection) error {
			qh := bson.M{
				"name": set.Name,
				"sets": []bson.M{},
			}

			log.Dev(context, "Upsert", "MGO : db.%s.insert(%s)", c.Name, mongo.Query(qh))
			return c.Insert(qh)
		}

		if err := db.ExecuteMGO(context, CollectionHistory, f); err != nil {
			log.Error(context, "Upsert", err, "Completed")
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

		log.Dev(context, "Upsert", "MGO : db.%s.update(%s, %s)", c.Name, mongo.Query(q), mongo.Query(qu))
		_, err := c.Upsert(q, qu)
		return err
	}

	if err := db.ExecuteMGO(context, CollectionHistory, f); err != nil {
		log.Error(context, "Upsert", err, "Completed")
		return err
	}

	log.Dev(context, "Upsert", "Completed")
	return nil
}

// =============================================================================

// GetNames retrieves a list of query names.
func GetNames(context interface{}, db *db.DB) ([]string, error) {
	log.Dev(context, "GetNames", "Started")

	var rawNames []struct {
		Name string
	}

	key := "gns"
	if v, found := cache.Get(key); found {
		names := v.([]string)
		log.Dev(context, "GetNames", "Completed : CACHE : Sets[%d]", len(names))
		return names, nil
	}

	f := func(c *mgo.Collection) error {
		s := bson.M{"name": 1}
		log.Dev(context, "GetNames", "MGO : db.%s.find({}, %s).sort([\"name\"])", c.Name, mongo.Query(s))
		return c.Find(nil).Select(s).Sort("name").All(&rawNames)
	}

	if err := db.ExecuteMGO(context, Collection, f); err != nil {
		if err == mgo.ErrNotFound {
			err = ErrNotFound
		}

		log.Error(context, "GetNames", err, "Completed")
		return nil, err
	}

	if rawNames == nil {
		log.Error(context, "GetNames", ErrNotFound, "Completed")
		return nil, ErrNotFound
	}

	names := make([]string, len(rawNames))
	for i := range rawNames {
		names[i] = rawNames[i].Name
	}

	cache.Set(key, names, gc.DefaultExpiration)

	log.Dev(context, "GetNames", "Completed : Sets[%d]", len(names))
	return names, nil
}

// GetAll retrieves a list of sets.
func GetAll(context interface{}, db *db.DB, tags []string) ([]Set, error) {
	log.Dev(context, "GetAll", "Started : Tags[%v]", tags)

	key := "gss" + strings.Join(tags, "-")
	if v, found := cache.Get(key); found {
		sets := v.([]Set)
		log.Dev(context, "GetAll", "Completed : CACHE : Sets[%d]", len(sets))
		return sets, nil
	}

	var sets []Set
	f := func(c *mgo.Collection) error {
		log.Dev(context, "GetAll", "MGO : db.%s.find({}).sort([\"name\"])", c.Name)
		return c.Find(nil).All(&sets)
	}

	if err := db.ExecuteMGO(context, Collection, f); err != nil {
		if err == mgo.ErrNotFound {
			err = ErrNotFound
		}

		log.Error(context, "GetAll", err, "Completed")
		return nil, err
	}

	if sets == nil {
		log.Error(context, "GetAll", ErrNotFound, "Completed")
		return nil, ErrNotFound
	}

	// Fix the sets so they can be used for processing.
	for i := range sets {
		sets[i].PrepareForUse()
	}

	cache.Set(key, sets, gc.DefaultExpiration)

	log.Dev(context, "GetAll", "Completed : Sets[%d]", len(sets))
	return sets, nil
}

// GetByName retrieves the document for the specified Set.
func GetByName(context interface{}, db *db.DB, name string) (*Set, error) {
	log.Dev(context, "GetByName", "Started : Name[%s]", name)

	key := "gbn" + name
	if v, found := cache.Get(key); found {
		set := v.(Set)
		log.Dev(context, "GetByName", "Completed : CACHE : Set[%+v]", &set)
		return &set, nil
	}

	var set Set
	f := func(c *mgo.Collection) error {
		q := bson.M{"name": name}
		log.Dev(context, "GetByName", "MGO : db.%s.findOne(%s)", c.Name, mongo.Query(q))
		return c.Find(q).One(&set)
	}

	if err := db.ExecuteMGO(context, Collection, f); err != nil {
		if err == mgo.ErrNotFound {
			err = ErrNotFound
		}

		log.Error(context, "GetByName", err, "Completed")
		return nil, err
	}

	// Fix the set so it can be used for processing.
	set.PrepareForUse()

	cache.Set(key, set, gc.DefaultExpiration)

	log.Dev(context, "GetByName", "Completed : Set[%+v]", &set)
	return &set, nil
}

// GetLastHistoryByName gets the last written Set within the history.
func GetLastHistoryByName(context interface{}, db *db.DB, name string) (*Set, error) {
	log.Dev(context, "GetLastHistoryByName", "Started : Name[%s]", name)

	type rslt struct {
		Name string `bson:"name"`
		Sets []Set  `bson:"sets"`
	}

	key := "glhbn" + name
	if v, found := cache.Get(key); found {
		result := v.(rslt)
		log.Dev(context, "GetLastHistoryByName", "Completed : CACHE :  Set[%+v]", &result.Sets[0])
		return &result.Sets[0], nil
	}

	var result rslt

	f := func(c *mgo.Collection) error {
		q := bson.M{"name": name}
		proj := bson.M{"sets": bson.M{"$slice": 1}}

		log.Dev(context, "GetLastHistoryByName", "MGO : db.%s.find(%s,%s)", c.Name, mongo.Query(q), mongo.Query(proj))
		return c.Find(q).Select(proj).One(&result)
	}

	if err := db.ExecuteMGO(context, CollectionHistory, f); err != nil {
		if err == mgo.ErrNotFound {
			err = ErrNotFound
		}

		log.Error(context, "GetLastHistoryByName", err, "Complete")
		return nil, err
	}

	if result.Sets == nil {
		err := errors.New("History not found")
		log.Error(context, "GetLastHistoryByName", err, "Complete")
		return nil, err
	}

	// Fix the set so it can be used for processing.
	result.Sets[0].PrepareForUse()

	cache.Set(key, result, gc.DefaultExpiration)

	log.Dev(context, "GetLastHistoryByName", "Completed : Set[%+v]", &result.Sets[0])
	return &result.Sets[0], nil
}

// =============================================================================

// Delete is used to remove an existing Set document.
func Delete(context interface{}, db *db.DB, name string) error {
	log.Dev(context, "Delete", "Started : Name[%s]", name)

	set, err := GetByName(context, db, name)
	if err != nil {
		return err
	}

	f := func(c *mgo.Collection) error {
		q := bson.M{"name": set.Name}
		log.Dev(context, "Delete", "MGO : db.%s.remove(%s)", c.Name, mongo.Query(q))
		return c.Remove(q)
	}

	if err := db.ExecuteMGO(context, Collection, f); err != nil {
		log.Error(context, "Delete", err, "Completed")
		return err
	}

	cache.Flush()

	log.Dev(context, "Delete", "Completed")
	return nil
}
