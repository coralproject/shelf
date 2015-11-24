// Package query provides CRUD API methods for handling database operations for query records.
package query

import (
	"strings"

	"github.com/coralproject/shelf/pkg/log"
	"github.com/coralproject/shelf/pkg/srv/mongo"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// collections contains the name of the rules collection.
const collection = "queries"

// ==============================================================

// Create is used to create Set document/record in the db.
func Create(context interface{}, ses *mgo.Session, rs Set) error {
	log.Dev(context, "Create", "Started : Name[%s]", rs.Name)

	f := func(c *mgo.Collection) error {
		log.Dev(context, "Create", "MGO :\n\ndb.%s.Insert(%s)\n", collection, mongo.Query(rs))
		rs.ID = bson.NewObjectId()
		return c.Insert(&rs)
	}

	if err := mongo.ExecuteDB(context, ses, collection, f); err != nil {
		log.Error(context, "Create", err, "Completed")
		return err
	}

	log.Dev(context, "Create", "Completed")
	return nil
}

// ==============================================================

// GetSetNames retrieves a list of rule names.
func GetSetNames(context interface{}, ses *mgo.Session) ([]string, error) {
	log.Dev(context, "GetNames", "Started")

	var names []bson.M
	f := func(c *mgo.Collection) error {
		q := bson.M{"name": 1}
		log.Dev(context, "GetNames", "MGO : db.%s.find({}, %s).sort([\"name\"])", collection, mongo.Query(q))
		return c.Find(nil).Select(q).Sort("name").All(&names)
	}

	if err := mongo.ExecuteDB(context, ses, collection, f); err != nil {
		log.Error(context, "GetNames", err, "Completed")
		return nil, err
	}

	var rsn []string
	for _, doc := range names {
		name := doc["name"].(string)
		if strings.HasPrefix(name, "test") {
			continue
		}

		rsn = append(rsn, name)
	}

	log.Dev(context, "GetNames", "Completed : RSN[%+v]", rsn)
	return rsn, nil
}

// ==============================================================

// GetSetByName retrieves the configuration for the specified Set.
func GetSetByName(context interface{}, ses *mgo.Session, name string) (*Set, error) {
	log.Dev(context, "Get", "Started : Name[%s]", name)

	var rs Set
	f := func(c *mgo.Collection) error {
		q := bson.M{"name": name}
		log.Dev(context, "Get", "MGO : db.%s.findOne(%s)", collection, mongo.Query(q))
		return c.Find(q).One(&rs)
	}

	if err := mongo.ExecuteDB(context, ses, collection, f); err != nil {
		log.Error(context, "Get", err, "Completed")
		return nil, err
	}

	log.Dev(context, "Get", "Completed : RS[%+v]", rs)
	return &rs, nil
}

// ==============================================================

// Update is used to update existing Set documents.
func Update(context interface{}, ses *mgo.Session, rs Set) error {
	log.Dev(context, "Update", "Started : Name[%s]", rs.Name)

	getid := func(c *mgo.Collection) error {
		q := bson.M{"name": rs.Name}
		qs := bson.M{"_id": 1}

		log.Dev(context, "Update", "MGO :\n\ndb.%s.find(%s).select(%s)\n", collection, mongo.Query(q), mongo.Query(qs))
		id := make(map[string]bson.ObjectId)
		if err := c.Find(q).Select(qs).One(&id); err != nil {
			return err
		}

		rs.ID = id["_id"]
		return nil
	}

	if err := mongo.ExecuteDB(context, ses, collection, getid); err != nil {
		rs.ID = bson.NewObjectId()
	}

	f := func(c *mgo.Collection) error {
		q := bson.M{"name": rs.Name}

		log.Dev(context, "Update", "MGO :\n\ndb.%s.upsert(%s, %s)\n", collection, mongo.Query(q), mongo.Query(rs))
		_, err := c.Upsert(q, &rs)
		return err
	}

	if err := mongo.ExecuteDB(context, ses, collection, f); err != nil {
		log.Error(context, "Update", err, "Completed")
		return err
	}

	log.Dev(context, "Update", "Completed")
	return nil
}

// ==============================================================

// Delete is used to remove an existing Set documents.
func Delete(context interface{}, ses *mgo.Session, name string) error {
	log.Dev(context, "Delete", "Started : Name[%s]", name)

	f := func(c *mgo.Collection) error {
		q := bson.M{"name": name}
		log.Dev(context, "Delete", "MGO :\n\ndb.%s.remove(%s)\n", collection, mongo.Query(q))
		return c.Remove(q)
	}

	if err := mongo.ExecuteDB(context, ses, collection, f); err != nil {
		log.Error(context, "Delete", err, "Completed")
		return err
	}

	log.Dev(context, "Delete", "Completed")
	return nil
}

// ==============================================================
