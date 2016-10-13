package script

import (
	"errors"
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
	Collection        = "query_scripts"
	CollectionHistory = "query_scripts_history"
)

// Set of error variables.
var (
	ErrNotFound = errors.New("Script Not found")
)

// =============================================================================

// c contans a cache of script values. The cache will maintain items for one
// hour and then marked as expired. This is a very small cache so the
// gc time will be every hour.

const (
	expiration = time.Hour
	cleanup    = time.Hour
)

var cache = gc.New(expiration, cleanup)

// =============================================================================

// Upsert is used to create or update an existing Script document.
func Upsert(context interface{}, db *db.DB, scr Script) error {
	log.Dev(context, "Upsert", "Started : Name[%s]", scr.Name)

	// Validate the set that is provided.
	if err := scr.Validate(); err != nil {
		log.Error(context, "Upsert", err, "Completed")
		return err
	}

	// We need to know if this is a new Set.
	var new bool
	if _, err := GetByName(context, db, scr.Name); err != nil {
		if err != ErrNotFound {
			log.Error(context, "Upsert", err, "Completed")
			return err
		}

		new = true
	}

	// Fix the set so it can be inserted.
	scr.PrepareForInsert()
	defer scr.PrepareForUse()

	// Insert or update the Set.
	f := func(c *mgo.Collection) error {
		q := bson.M{"name": scr.Name}
		log.Dev(context, "Upsert", "MGO : db.%s.upsert(%s, %s)", c.Name, mongo.Query(q), mongo.Query(scr))
		_, err := c.Upsert(q, scr)
		return err
	}

	if err := db.ExecuteMGO(context, Collection, f); err != nil {
		log.Error(context, "Upsert", err, "Completed")
		return err
	}

	// Flush the cache to invalidate everything.
	cache.Flush()

	// Add a history record if this script set is new.
	if new {
		f = func(c *mgo.Collection) error {
			sh := bson.M{
				"name":    scr.Name,
				"scripts": []bson.M{},
			}

			log.Dev(context, "Upsert", "MGO : db.%s.insert(%s)", c.Name, mongo.Query(sh))
			return c.Insert(sh)
		}

		if err := db.ExecuteMGO(context, CollectionHistory, f); err != nil {
			log.Error(context, "Upsert", err, "Completed")
			return err
		}
	}

	// Add this script to the beginning of the history.
	f = func(c *mgo.Collection) error {
		q := bson.M{"name": scr.Name}
		su := bson.M{
			"$push": bson.M{
				"scripts": bson.M{
					"$each":     []Script{scr},
					"$position": 0,
				},
			},
		}

		log.Dev(context, "Upsert", "MGO : db.%s.update(%s, %s)", c.Name, mongo.Query(q), mongo.Query(su))
		_, err := c.Upsert(q, su)
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

// GetNames retrieves a list of script names.
func GetNames(context interface{}, db *db.DB) ([]string, error) {
	log.Dev(context, "GetNames", "Started")

	var rawNames []struct {
		Name string
	}

	key := "gns"
	if v, found := cache.Get(key); found {
		names := v.([]string)
		log.Dev(context, "GetNames", "Completed : CACHE : Scripts[%d]", len(names))
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

	log.Dev(context, "GetNames", "Completed : Scripts[%d]", len(names))
	return names, nil
}

// GetAll retrieves a list of scripts.
func GetAll(context interface{}, db *db.DB, tags []string) ([]Script, error) {
	log.Dev(context, "GetAll", "Started : Tags[%v]", tags)

	key := "gss" + strings.Join(tags, "-")
	if v, found := cache.Get(key); found {
		scrs := v.([]Script)
		log.Dev(context, "GetAll", "Completed : CACHE : Scripts[%d]", len(scrs))
		return scrs, nil
	}

	var scrs []Script
	f := func(c *mgo.Collection) error {
		log.Dev(context, "GetAll", "MGO : db.%s.find({}).sort([\"name\"])", c.Name)
		return c.Find(nil).All(&scrs)
	}

	if err := db.ExecuteMGO(context, Collection, f); err != nil {
		if err == mgo.ErrNotFound {
			err = ErrNotFound
		}

		log.Error(context, "GetAll", err, "Completed")
		return nil, err
	}

	if scrs == nil {
		log.Error(context, "GetAll", ErrNotFound, "Completed")
		return nil, ErrNotFound
	}

	// Fix the scripts so they can be used for processing.
	for i := range scrs {
		scrs[i].PrepareForUse()
	}

	cache.Set(key, scrs, gc.DefaultExpiration)

	log.Dev(context, "GetAll", "Completed : Scripts[%d]", len(scrs))
	return scrs, nil
}

// GetByName retrieves the document for the specified name.
func GetByName(context interface{}, db *db.DB, name string) (Script, error) {
	log.Dev(context, "GetByName", "Started : Name[%s]", name)

	key := "gbn" + name
	if v, found := cache.Get(key); found {
		scr := v.(Script)
		log.Dev(context, "GetByName", "Completed : CACHE : Script[%+v]", &scr)
		return scr, nil
	}

	var scr Script
	f := func(c *mgo.Collection) error {
		q := bson.M{"name": name}
		log.Dev(context, "GetByName", "MGO : db.%s.findOne(%s)", c.Name, mongo.Query(q))
		return c.Find(q).One(&scr)
	}

	if err := db.ExecuteMGO(context, Collection, f); err != nil {
		if err == mgo.ErrNotFound {
			err = ErrNotFound
		}

		log.Error(context, "GetByName", err, "Completed")
		return Script{}, err
	}

	// Fix the script so it can be used for processing.
	scr.PrepareForUse()

	cache.Set(key, scr, gc.DefaultExpiration)

	log.Dev(context, "GetByName", "Completed : Script[%+v]", &scr)
	return scr, nil
}

// GetByNames retrieves the documents for the specified names.
func GetByNames(context interface{}, db *db.DB, names []string) ([]Script, error) {
	log.Dev(context, "GetByNames", "Started : Names[%+v]", names)

	key := "gbns" + strings.Join(names, "-")
	if v, found := cache.Get(key); found {
		scripts := v.([]Script)
		log.Dev(context, "GetByNames", "Completed : CACHE : Scripts[%+v]", scripts)
		return scripts, nil
	}

	var scrs []Script
	f := func(c *mgo.Collection) error {

		// Build a list of documents to find by name.
		qn := make([]bson.M, len(names))
		for i, name := range names {
			if name != "" {
				qn[i] = bson.M{"name": name}
			}
		}

		// Place that list in an $or operation.
		q := bson.M{"$or": qn}

		log.Dev(context, "GetByNames", "MGO : db.%s.find(%s)", c.Name, mongo.Query(q))
		return c.Find(q).All(&scrs)
	}

	if err := db.ExecuteMGO(context, Collection, f); err != nil {
		if err == mgo.ErrNotFound {
			err = ErrNotFound
		}

		log.Error(context, "GetByNames", err, "Completed")
		return nil, err
	}

	if scrs == nil {
		log.Error(context, "GetByNames", ErrNotFound, "Completed")
		return nil, ErrNotFound
	}

	// I can't assume MongoDB will bring the results back in the order I
	// setup the query. I need the order to match on the returned slice.
	// I thought about using a map of name/value but I feel like it is overkill.

	scripts := make([]Script, len(names))
next:
	for _, scr := range scrs {
		for i := range names {
			if scr.Name == names[i] {
				scripts[i] = scr
				continue next
			}
		}
	}

	// Fix the scripts so they can be used for processing.
	for i := range scripts {
		scripts[i].PrepareForUse()
	}

	cache.Set(key, scripts, gc.DefaultExpiration)

	log.Dev(context, "GetByNames", "Completed : Scripts[%+v]", scripts)
	return scripts, nil
}

// GetLastHistoryByName gets the last written Script within the history.
func GetLastHistoryByName(context interface{}, db *db.DB, name string) (Script, error) {
	log.Dev(context, "GetLastHistoryByName", "Started : Name[%s]", name)

	type rslt struct {
		Name    string   `bson:"name"`
		Scripts []Script `bson:"scripts"`
	}

	key := "glhbn" + name
	if v, found := cache.Get(key); found {
		result := v.(rslt)
		log.Dev(context, "GetLastHistoryByName", "Completed : CACHE : Script[%+v]", &result.Scripts[0])
		return result.Scripts[0], nil
	}

	var result rslt

	f := func(c *mgo.Collection) error {
		q := bson.M{"name": name}
		proj := bson.M{"scripts": bson.M{"$slice": 1}}

		log.Dev(context, "GetLastHistoryByName", "MGO : db.%s.find(%s,%s)", c.Name, mongo.Query(q), mongo.Query(proj))
		return c.Find(q).Select(proj).One(&result)
	}

	if err := db.ExecuteMGO(context, CollectionHistory, f); err != nil {
		if err == mgo.ErrNotFound {
			err = ErrNotFound
		}

		log.Error(context, "GetLastHistoryByName", err, "Complete")
		return Script{}, err
	}

	if result.Scripts == nil {
		err := errors.New("History not found")
		log.Error(context, "GetLastHistoryByName", err, "Complete")
		return Script{}, err
	}

	// Fix the script so it can be used for processing.
	result.Scripts[0].PrepareForUse()

	cache.Set(key, result, gc.DefaultExpiration)

	log.Dev(context, "GetLastHistoryByName", "Completed : Script[%+v]", &result.Scripts[0])
	return result.Scripts[0], nil
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
