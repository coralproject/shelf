// Package mask provides the service layer for managing masks that need
// to be applied to results before they are returned.
package mask

import (
	"errors"
	"strings"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/ardanlabs/kit/db"
	"github.com/ardanlabs/kit/db/mongo"
	"github.com/ardanlabs/kit/log"
	gc "github.com/patrickmn/go-cache"
)

// Contains the name of Mongo collections.
const (
	Collection        = "query_masks"
	CollectionHistory = "query_masks_history"
)

// Set of error variables.
var (
	ErrNotFound = errors.New("Mask Not found")
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

// Upsert is used to create or update an existing query Mask document.
func Upsert(context interface{}, db *db.DB, mask Mask) error {
	log.Dev(context, "Upsert", "Started : Mask[%+v]", mask)

	// Validate the mask that is provided.
	if err := mask.Validate(); err != nil {
		log.Error(context, "Upsert", err, "Completed")
		return err
	}

	// We need to know if this is a new query mask.
	var new bool
	if _, err := GetByName(context, db, mask.Collection, mask.Field); err != nil {
		if err != ErrNotFound {
			log.Error(context, "Upsert", err, "Completed")
			return err
		}

		new = true
	}

	// Insert or update the query mask.
	f := func(c *mgo.Collection) error {
		q := bson.M{"collection": mask.Collection, "field": mask.Field}
		log.Dev(context, "Upsert", "MGO : db.%s.upsert(%s, %s)", c.Name, mongo.Query(q), mongo.Query(mask))
		_, err := c.Upsert(q, mask)
		return err
	}

	if err := db.ExecuteMGO(context, Collection, f); err != nil {
		log.Error(context, "Upsert", err, "Completed")
		return err
	}

	// Flush the cache to invalidate everything.
	cache.Flush()

	// Add a history record if this query mask is new.
	if new {
		f = func(c *mgo.Collection) error {
			qh := bson.M{
				"collection": mask.Collection,
				"field":      mask.Field,
				"masks":      []bson.M{},
			}

			log.Dev(context, "Upsert", "MGO : db.%s.insert(%s)", c.Name, mongo.Query(qh))
			return c.Insert(qh)
		}

		if err := db.ExecuteMGO(context, CollectionHistory, f); err != nil {
			log.Error(context, "Upsert", err, "Completed")
			return err
		}
	}

	// Add this query mask to the beginning of the history.
	f = func(c *mgo.Collection) error {
		q := bson.M{"collection": mask.Collection, "field": mask.Field}
		qu := bson.M{
			"$push": bson.M{
				"masks": bson.M{
					"$each":     []Mask{mask},
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

// GetAll retrieves a list of query masks.
func GetAll(context interface{}, db *db.DB, tags []string) (map[string]Mask, error) {
	log.Dev(context, "GetAll", "Started : Tags[%v]", tags)

	key := "gms" + strings.Join(tags, "-")
	if v, found := cache.Get(key); found {
		mskMap := v.(map[string]Mask)
		log.Dev(context, "GetAll", "Completed : CACHE : Masks[%d]", len(mskMap))
		return mskMap, nil
	}

	var masks []Mask
	f := func(c *mgo.Collection) error {
		log.Dev(context, "GetAll", "MGO : db.%s.find({}).sort([\"name\"])", c.Name)
		return c.Find(nil).All(&masks)
	}

	if err := db.ExecuteMGO(context, Collection, f); err != nil {
		if err == mgo.ErrNotFound {
			err = ErrNotFound
		}

		log.Error(context, "GetAll", err, "Completed")
		return nil, err
	}

	if masks == nil {
		log.Error(context, "GetAll", ErrNotFound, "Completed")
		return nil, ErrNotFound
	}

	mskMap := make(map[string]Mask, len(masks))
	for _, msk := range masks {
		mskMap[msk.Field] = msk
	}

	cache.Set(key, mskMap, gc.DefaultExpiration)

	log.Dev(context, "GetAll", "Completed : Masks[%d]", len(mskMap))
	return mskMap, nil
}

// GetByCollection retrieves the masks for the specified collection.
func GetByCollection(context interface{}, db *db.DB, collection string) (map[string]Mask, error) {
	log.Dev(context, "GetByCollection", "Started : Collection[%s]", collection)

	key := "gbc" + collection
	if v, found := cache.Get(key); found {
		mskMap := v.(map[string]Mask)
		log.Dev(context, "GetByCollection", "Completed : CACHE : Masks[%d]", len(mskMap))
		return mskMap, nil
	}

	var masks []Mask
	f := func(c *mgo.Collection) error {
		q := bson.M{"collection": collection}
		log.Dev(context, "GetByCollection", "MGO : db.%s.findOne(%s)", c.Name, mongo.Query(q))
		return c.Find(q).All(&masks)
	}

	if err := db.ExecuteMGO(context, Collection, f); err != nil {
		if err == mgo.ErrNotFound {
			err = ErrNotFound
		}

		log.Error(context, "GetByCollection", err, "Completed")
		return nil, err
	}

	if masks == nil {
		log.Error(context, "GetByCollection", ErrNotFound, "Completed")
		return nil, ErrNotFound
	}

	mskMap := make(map[string]Mask, len(masks))
	for _, msk := range masks {
		mskMap[msk.Field] = msk
	}

	cache.Set(key, mskMap, gc.DefaultExpiration)

	log.Dev(context, "GetByCollection", "Completed : Masks[%d]", len(mskMap))
	return mskMap, nil
}

// GetByName retrieves the document for the specified query mask.
func GetByName(context interface{}, db *db.DB, collection string, field string) (Mask, error) {
	log.Dev(context, "GetByName", "Started : Collection[%s] Field[%s]", collection, field)

	key := "gbn" + collection + field
	if v, found := cache.Get(key); found {
		mask := v.(Mask)
		log.Dev(context, "GetByName", "Completed : CACHE : Mask[%+v]", mask)
		return mask, nil
	}

	var mask Mask
	f := func(c *mgo.Collection) error {
		q := bson.M{"collection": collection, "field": field}
		log.Dev(context, "GetByName", "MGO : db.%s.findOne(%s)", c.Name, mongo.Query(q))
		return c.Find(q).One(&mask)
	}

	if err := db.ExecuteMGO(context, Collection, f); err != nil {
		if err == mgo.ErrNotFound {
			err = ErrNotFound
		}

		log.Error(context, "GetByName", err, "Completed")
		return Mask{}, err
	}

	cache.Set(key, mask, gc.DefaultExpiration)

	log.Dev(context, "GetByName", "Completed : Mask[%+v]", mask)
	return mask, nil
}

// GetLastHistoryByName gets the last written query mask within the history.
func GetLastHistoryByName(context interface{}, db *db.DB, collection string, field string) (Mask, error) {
	log.Dev(context, "GetLastHistoryByName", "Started : Collection[%s] Field[%s]", collection, field)

	type rslt struct {
		Name  string `bson:"name"`
		Masks []Mask `bson:"masks"`
	}

	key := "glhbn" + collection + field
	if v, found := cache.Get(key); found {
		result := v.(rslt)
		log.Dev(context, "GetLastHistoryByName", "Completed : CACHE :  Set[%+v]", result.Masks[0])
		return result.Masks[0], nil
	}

	var result rslt

	f := func(c *mgo.Collection) error {
		q := bson.M{"collection": collection, "field": field}
		proj := bson.M{"masks": bson.M{"$slice": 1}}

		log.Dev(context, "GetLastHistoryByName", "MGO : db.%s.find(%s,%s)", c.Name, mongo.Query(q), mongo.Query(proj))
		return c.Find(q).Select(proj).One(&result)
	}

	if err := db.ExecuteMGO(context, CollectionHistory, f); err != nil {
		if err == mgo.ErrNotFound {
			err = ErrNotFound
		}

		log.Error(context, "GetLastHistoryByName", err, "Complete")
		return Mask{}, err
	}

	if result.Masks == nil {
		err := errors.New("History not found")
		log.Error(context, "GetLastHistoryByName", err, "Complete")
		return Mask{}, err
	}

	cache.Set(key, result, gc.DefaultExpiration)

	log.Dev(context, "GetLastHistoryByName", "Completed : Set[%+v]", result.Masks[0])
	return result.Masks[0], nil
}

// =============================================================================

// Delete is used to remove an existing query mask document.
func Delete(context interface{}, db *db.DB, collection string, field string) error {
	log.Dev(context, "Delete", "Started : Collection[%s] Field[%s]", collection, field)

	mask, err := GetByName(context, db, collection, field)
	if err != nil {
		return err
	}

	f := func(c *mgo.Collection) error {
		q := bson.M{"collection": mask.Collection, "field": mask.Field}
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
