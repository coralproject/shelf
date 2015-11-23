package query

import (
	"strings"

	"github.com/coralproject/shelf/pkg/log"
	"github.com/coralproject/shelf/pkg/srv/mongo"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// collections contains the name of the rules collection.
const collection = "rules"

// GetSetNames retrieves a list of rule names.
func GetSetNames(context interface{}, ses *mgo.Session) ([]string, error) {
	log.Dev(context, "GetSetNames", "Started")

	var names []bson.M
	f := func(c *mgo.Collection) error {
		q := bson.M{"name": 1}
		log.Dev(context, "GetSetNames", "MGO : db.%s.find({}, %s).sort([\"name\"])", collection, mongo.Query(q))
		return c.Find(nil).Select(q).Sort("name").All(&names)
	}

	if err := mongo.ExecuteDB(context, ses, collection, f); err != nil {
		log.Error(context, "GetSetNames", err, "Completed")
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

	log.Dev(context, "GetSetNames", "Completed : RSN[%+v]", rsn)
	return rsn, nil
}

// GetSet retrieves the configuration for the specified Set.
func GetSet(context interface{}, ses *mgo.Session, name string) (*Set, error) {
	log.Dev(context, "GetSet", "Started : Name[%s]", name)

	var rs Set
	f := func(c *mgo.Collection) error {
		q := bson.M{"name": name}
		log.Dev(context, "GetSet", "MGO : db.%s.findOne(%s)", collection, mongo.Query(q))
		return c.Find(q).One(&rs)
	}

	if err := mongo.ExecuteDB(context, ses, collection, f); err != nil {
		log.Error(context, "GetSet", err, "Completed")
		return nil, err
	}

	log.Dev(context, "GetSet", "Completed : RS[%+v]", rs)
	return &rs, nil
}

// Create is used to create Set document/record in the db.
func Create(context interface{}, ses *mgo.Session, rs *Set) error {
	log.Dev(context, "Create", "Started : Name[%s]", rs.Name)

	f := func(c *mgo.Collection) error {
		log.Dev(context, "Create", "MGO :\n\ndb.%s.Insert(%s)\n", collection, mongo.Query(rs))
		return c.Insert(rs)
	}

	if err := mongo.ExecuteDB(context, ses, collection, f); err != nil {
		log.Error(context, "Create", err, "Completed")
		return err
	}

	log.Dev(context, "Create", "Completed")
	return nil
}

// UpdateSet is used to create or update existing Set documents.
func UpdateSet(context interface{}, ses *mgo.Session, rs *Set) error {
	log.Dev(context, "UpdateSet", "Started : Name[%s]", rs.Name)

	f := func(c *mgo.Collection) error {
		q := bson.M{"name": rs.Name}

		log.Dev(context, "UpdateSet", "MGO :\n\ndb.%s.upsert(%s, %s)\n", collection, mongo.Query(q), mongo.Query(rs))
		_, err := c.Upsert(q, rs)
		return err
	}

	if err := mongo.ExecuteDB(context, ses, collection, f); err != nil {
		log.Error(context, "UpdateSet", err, "Completed")
		return err
	}

	log.Dev(context, "UpdateSet", "Completed")
	return nil
}

// DeleteSet is used to remove an existing Set documents.
func DeleteSet(context interface{}, ses *mgo.Session, name string) (*Set, error) {
	log.Dev(context, "RemoveSet", "Started : Name[%s]", name)

	var rs Set
	f := func(c *mgo.Collection) error {
		q := bson.M{"name": name}
		log.Dev(context, "DeleteSet", "MGO :\n\ndb.%s.remove(%s)\n", collection, mongo.Query(q))
		return c.Remove(q)
	}

	if err := mongo.ExecuteDB(context, ses, collection, f); err != nil {
		log.Error(context, "DeleteSet", err, "Completed")
		return nil, err
	}

	log.Dev(context, "DeleteSet", "Completed")
	return &rs, nil
}
