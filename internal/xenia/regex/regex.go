// Package regex provides the service layer for building apps using
// regex functionality.
package regex

import (
	"errors"
	"regexp"
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
	Collection        = "query_regexs"
	CollectionHistory = "query_regexs_history"
)

// Set of error variables.
var (
	ErrNotFound = errors.New("Regex Not found")
)

// =============================================================================

// c contans a cache of regex values. The cache will maintain items for one
// hour and then marked as expired. This is a very small cache so the
// gc time will be every hour.

const (
	expiration = time.Hour
	cleanup    = time.Hour
)

var cache = gc.New(expiration, cleanup)

// =============================================================================

// Upsert is used to create or update an existing Regex document.
func Upsert(context interface{}, db *db.DB, rgx Regex) error {
	log.Dev(context, "Upsert", "Started : Name[%s]", rgx.Name)

	// Validate the regex that is provided.
	if err := rgx.Validate(); err != nil {
		log.Error(context, "Upsert", err, "Completed")
		return err
	}

	// We need to know if this is a new regex.
	var new bool
	if _, err := GetByName(context, db, rgx.Name); err != nil {
		if err != ErrNotFound {
			log.Error(context, "Upsert", err, "Completed")
			return err
		}

		new = true
	}

	// Insert or update the query regex.
	f := func(c *mgo.Collection) error {
		q := bson.M{"name": rgx.Name}
		log.Dev(context, "Upsert", "MGO : db.%s.upsert(%s, %s)", c.Name, mongo.Query(q), mongo.Query(rgx))
		_, err := c.Upsert(q, rgx)
		return err
	}

	if err := db.ExecuteMGO(context, Collection, f); err != nil {
		log.Error(context, "Upsert", err, "Completed")
		return err
	}

	// Flush the cache to invalidate everything.
	cache.Flush()

	// Add a history record if this query regex is new.
	if new {
		f = func(c *mgo.Collection) error {
			qh := bson.M{
				"name":   rgx.Name,
				"regexs": []bson.M{},
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
		q := bson.M{"name": rgx.Name}
		qu := bson.M{
			"$push": bson.M{
				"regexs": bson.M{
					"$each":     []Regex{rgx},
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

// GetNames retrieves a list of query regex names.
func GetNames(context interface{}, db *db.DB) ([]string, error) {
	log.Dev(context, "GetNames", "Started")

	var rawNames []struct {
		Name string
	}

	key := "gns"
	if v, found := cache.Get(key); found {
		names := v.([]string)
		log.Dev(context, "GetNames", "Completed : CACHE : Rgxs[%d]", len(names))
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

	log.Dev(context, "GetNames", "Completed : Rgxs[%d]", len(names))
	return names, nil
}

// GetAll retrieves a list of regexs.
func GetAll(context interface{}, db *db.DB, tags []string) ([]Regex, error) {
	log.Dev(context, "GetAll", "Started : Tags[%v]", tags)

	key := "grs" + strings.Join(tags, "-")
	if v, found := cache.Get(key); found {
		rgxs := v.([]Regex)
		log.Dev(context, "GetAll", "Completed : CACHE : Rgxs[%d]", len(rgxs))
		return rgxs, nil
	}

	var rgxs []Regex
	f := func(c *mgo.Collection) error {
		log.Dev(context, "GetAll", "MGO : db.%s.find({}).sort([\"name\"])", c.Name)
		return c.Find(nil).All(&rgxs)
	}

	if err := db.ExecuteMGO(context, Collection, f); err != nil {
		if err == mgo.ErrNotFound {
			err = ErrNotFound
		}

		log.Error(context, "GetAll", err, "Completed")
		return nil, err
	}

	if rgxs == nil {
		log.Error(context, "GetAll", ErrNotFound, "Completed")
		return nil, ErrNotFound
	}

	cache.Set(key, rgxs, gc.DefaultExpiration)

	log.Dev(context, "GetAll", "Completed : Rgxs[%d]", len(rgxs))
	return rgxs, nil
}

// GetByName retrieves the document for the specified Regex.
func GetByName(context interface{}, db *db.DB, name string) (Regex, error) {
	log.Dev(context, "GetByName", "Started : Name[%s]", name)

	key := "gbn" + name
	if v, found := cache.Get(key); found {
		rgx := v.(Regex)
		log.Dev(context, "GetByName", "Completed : CACHE : Rgx[%s]", rgx.Name)
		return rgx, nil
	}

	var rgx Regex
	f := func(c *mgo.Collection) error {
		q := bson.M{"name": name}
		log.Dev(context, "GetByName", "MGO : db.%s.findOne(%s)", c.Name, mongo.Query(q))
		return c.Find(q).One(&rgx)
	}

	if err := db.ExecuteMGO(context, Collection, f); err != nil {
		if err == mgo.ErrNotFound {
			err = ErrNotFound
		}

		log.Error(context, "GetByName", err, "Completed")
		return Regex{}, err
	}

	// This call is made when the regex is required for actual use. So
	// let's compile the regex now.
	var err error
	if rgx.Compile, err = regexp.Compile(rgx.Expr); err != nil {
		return Regex{}, err
	}

	cache.Set(key, rgx, gc.DefaultExpiration)

	log.Dev(context, "GetByName", "Completed : Rgx[%s]", rgx.Name)
	return rgx, nil
}

// GetByNames retrieves the documents for the specified names.
func GetByNames(context interface{}, db *db.DB, names []string) ([]Regex, error) {
	log.Dev(context, "GetByNames", "Started : Names[%+v]", names)

	key := "gbns" + strings.Join(names, "-")
	if v, found := cache.Get(key); found {
		regexs := v.([]Regex)
		log.Dev(context, "GetByNames", "Completed : CACHE : Regexs[%+v]", regexs)
		return regexs, nil
	}

	var rgxs []Regex
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
		return c.Find(q).All(&rgxs)
	}

	if err := db.ExecuteMGO(context, Collection, f); err != nil {
		if err == mgo.ErrNotFound {
			err = ErrNotFound
		}

		log.Error(context, "GetByNames", err, "Completed")
		return nil, err
	}

	if rgxs == nil {
		log.Error(context, "GetByNames", ErrNotFound, "Completed")
		return nil, ErrNotFound
	}

	// I can't assume MongoDB will bring the results back in the order I
	// setup the query. I need the order to match on the returned slice.
	// I thought about using a map of name/value but I feel like it is overkill.

	regexs := make([]Regex, len(names))
next:
	for _, rgx := range rgxs {
		for i := range names {
			if rgx.Name == names[i] {
				regexs[i] = rgx
				continue next
			}
		}
	}

	cache.Set(key, regexs, gc.DefaultExpiration)

	log.Dev(context, "GetByNames", "Completed : Regexs[%+v]", regexs)
	return regexs, nil
}

// GetLastHistoryByName gets the last written Regex within the history.
func GetLastHistoryByName(context interface{}, db *db.DB, name string) (Regex, error) {
	log.Dev(context, "GetLastHistoryByName", "Started : Name[%s]", name)

	type rslt struct {
		Name   string  `bson:"name"`
		Regexs []Regex `bson:"regexs"`
	}

	key := "glhbn" + name
	if v, found := cache.Get(key); found {
		result := v.(rslt)
		log.Dev(context, "GetLastHistoryByName", "Completed : CACHE : Regex[%+v]", &result.Regexs[0])
		return result.Regexs[0], nil
	}

	var result rslt

	f := func(c *mgo.Collection) error {
		q := bson.M{"name": name}
		proj := bson.M{"regexs": bson.M{"$slice": 1}}

		log.Dev(context, "GetLastHistoryByName", "MGO : db.%s.find(%s,%s)", c.Name, mongo.Query(q), mongo.Query(proj))
		return c.Find(q).Select(proj).One(&result)
	}

	if err := db.ExecuteMGO(context, CollectionHistory, f); err != nil {
		if err == mgo.ErrNotFound {
			err = ErrNotFound
		}

		log.Error(context, "GetLastHistoryByName", err, "Complete")
		return Regex{}, err
	}

	if result.Regexs == nil {
		err := errors.New("History not found")
		log.Error(context, "GetLastHistoryByName", err, "Complete")
		return Regex{}, err
	}

	cache.Set(key, result, gc.DefaultExpiration)

	log.Dev(context, "GetLastHistoryByName", "Completed : Regex[%+v]", &result.Regexs[0])
	return result.Regexs[0], nil
}

// =============================================================================

// Delete is used to remove an existing Regex document.
func Delete(context interface{}, db *db.DB, name string) error {
	log.Dev(context, "Delete", "Started : Name[%s]", name)

	rgx, err := GetByName(context, db, name)
	if err != nil {
		return err
	}

	f := func(c *mgo.Collection) error {
		q := bson.M{"name": rgx.Name}
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
