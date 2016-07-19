package item

import (
	"errors"
	"time"

	"github.com/ardanlabs/kit/db"
	"github.com/ardanlabs/kit/db/mongo"
	"github.com/ardanlabs/kit/log"

	gc "github.com/patrickmn/go-cache"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Contains the name of Mongo collections.
const (
	Collection        = "coral_items"
	CollectionHistory = "coral_items_history"
)

// Set of error variables.
var (
	ErrNotFound = errors.New("Item Not found")
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

// create an item out of its type, version and data or die trying
//  if the item already exists, adopt it's Id but do not presume to update
func Create(context interface{}, db *db.DB, t string, v int, d map[string]interface{}) (Item, error) {

	i := Item{}

	// Check to see if this item already exists in the db
	//  this is done by checking it's "foreign identity":
	//  Checking for its type and it's Type.Id fields

	// get this data's IdField value for this type
	idValue := getDatumByKey(Types[t].IdField, d)

	// if there is not a value, generate a new id and continue
	if idValue == nil {
		i.Id = bson.NewObjectId()
	} else { // if there is an id, look it up

		// create the field path for the id field in the d subdoc
		dbIdField := "d." + Types[t].IdField

		// build a query
		var q = bson.M{"t": t, dbIdField: idValue}

		// get one by query, assuming data consistency
		dbItem, err := GetOneByQuery(context, db, q)
		if err != nil {
			return Item{}, err
		}

		// if we found an item, assign the id to the new item
		if dbItem != nil {
			i.Id = dbItem.Id
		} else { // otherwise, new id it is
			i.Id = bson.NewObjectId()
		}

	}

	// validate and set type
	if isRegistered(t) == false {
		return i, errors.New("Type not recognized: " + t)
	}
	i.Type = t

	// set default version if zero value
	if v == 0 {
		v = DefaultVersion
	}
	i.Version = v

	// set the data into the item
	i.Data = d

	// get the relationships for this item
	rels, err := GetRels(context, db, &i)
	if err != nil {
		return Item{}, err // todo, clean up empty type return
	}
	i.Rels = *rels

	return i, nil
}

// Items are trasparently created or updated depending on thier existence
func Upsert(context interface{}, db *db.DB, item *Item) error {

	// validate our item
	if err := item.Validate(); err != nil {
		log.Error(context, "Upsert", err, "Completed")
		return err
	}

	// We need to know if this is a new set.
	var new bool
	if _, err := GetById(context, db, item.Id); err != nil {
		if err != ErrNotFound {
			log.Error(context, "Upsert", err, "Completed")
			return err
		}

		new = true
	}

	// Insert or update the item.
	f := func(c *mgo.Collection) error {
		q := bson.M{"_id": item.Id}
		log.Dev(context, "Upsert", "MGO : db.%s.upsert(%s, %s)", c.Name, mongo.Query(q), mongo.Query(item))
		_, err := c.Upsert(q, item)
		return err
	}

	if err := db.ExecuteMGO(context, Collection, f); err != nil {
		log.Error(context, "Upsert", err, "Completed")
		return err
	}

	if new {
		// historical code
	}

	// if the item isn't new it may be in various caches
	//   flush the whole cache
	if !new {
		cache.Flush()
	}

	return nil
}

// GetById retrieves an item by its id, given the id as a string.
func GetByIdString(context interface{}, db *db.DB, i string) (*Item, error) {

	// can we make this into a valid bson ObjectId?
	id := bson.ObjectIdHex(i)
	//	if err != nil {
	//		return nil, err
	//	}

	// if so, use the traditional GetById
	return GetById(context, db, id)
}

// GetById retrieves an item by its id.
func GetById(context interface{}, db *db.DB, id bson.ObjectId) (*Item, error) {
	log.Dev(context, "GetById", "Started : Id[%s]", id.Hex())

	var item Item

	// check if the item is in the cache
	key := "item-" + id.Hex()
	if v, found := cache.Get(key); found {
		item := v.(Item)
		log.Dev(context, "GetById", "Completed : CACHE : Item[%+v]", &item)
		return &item, nil
	}

	// query the database for the item
	f := func(c *mgo.Collection) error {
		q := bson.M{"_id": id}
		log.Dev(context, "GetById", "MGO : db.%s.findOne(%s)", c.Name, mongo.Query(q))
		return c.Find(q).One(&item)
	}

	if err := db.ExecuteMGO(context, Collection, f); err != nil {
		if err == mgo.ErrNotFound {
			err = ErrNotFound
		}

		log.Error(context, "GetById", err, "Completed")
		return nil, err
	}

	// set the cache: TODO, caching based on type params
	cache.Set(key, item, gc.DefaultExpiration)

	log.Dev(context, "GetById", "Completed : Item[%+v]", &item)
	return &item, nil
}

// GetById retrieves items by an array of ids
func GetByIds(context interface{}, db *db.DB, ids []bson.ObjectId) (*[]Item, error) {
	log.Dev(context, "GetByIds", "Started : Looking for %s ids", len(ids))

	var items []Item

	// query the database for the item
	f := func(c *mgo.Collection) error {
		q := bson.M{"_id": bson.M{"$in": ids}}
		log.Dev(context, "GetByIds", "MGO : ", c.Name, mongo.Query(q))
		return c.Find(q).All(&items)
	}

	if err := db.ExecuteMGO(context, Collection, f); err != nil {
		if err == mgo.ErrNotFound {
			err = ErrNotFound
		}

		log.Error(context, "GetByIds", err, "Completed")
		return nil, err
	}

	log.Dev(context, "GetByIds", "Completed : Found %+v items", len(items))
	return &items, nil
}

// GetByQuery accepts a bson.M query and runs it against the item collection
//  caution should be used to only query against indexed fields
func GetByQuery(context interface{}, db *db.DB, q bson.M) (*[]Item, error) {
	log.Dev(context, "GetByQuery", "Started : Looking for %#v", q)

	var items []Item

	// query the database for the item
	f := func(c *mgo.Collection) error {
		log.Dev(context, "GetByQuery", "MGO : %#v", q)
		return c.Find(q).All(&items)
	}

	if err := db.ExecuteMGO(context, Collection, f); err != nil {
		if err == mgo.ErrNotFound {
			err = ErrNotFound
		}

		log.Error(context, "GetByQuery", err, "Completed")
		return nil, err
	}

	log.Dev(context, "GetByQuery", "Completed : Found %+v items", len(items))
	return &items, nil

}

// GetOneByQuery accepts a bson.M query and runs it against the item collection
//  returning the first record found
func GetOneByQuery(context interface{}, db *db.DB, q bson.M) (*Item, error) {
	log.Dev(context, "GetByQuery", "Started : Looking for %#v", q)

	var item Item

	// query the database for the item
	f := func(c *mgo.Collection) error {
		log.Dev(context, "GetByQuery", "MGO : %#v", q)
		return c.Find(q).One(&item)
	}

	if err := db.ExecuteMGO(context, Collection, f); err != nil {
		// it's ok to return an empty query, not all other errors are thrown
		if err != mgo.ErrNotFound {
			return nil, err
		}

		// return nil nil for no results found
		log.Dev(context, "GetByQuery", "Completed : No Items found")
		return nil, nil
	}

	log.Dev(context, "GetByQuery", "Completed : Found %+v", item.Id)
	return &item, nil

}
