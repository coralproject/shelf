package query

import (
	"strings"

	"github.com/coralproject/shelf/pkg/db/mongo"
	"github.com/coralproject/shelf/pkg/log"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// collections contains the name of the rules collection.
const collection = "rules"

// GetQuerySetNames retrieves a list of rule names.
func GetQuerySetNames(context interface{}, ses *mgo.Session) ([]string, error) {
	log.Dev(context, "GetQuerySetNames", "Started")

	var names []bson.M
	f := func(c *mgo.Collection) error {
		q := bson.M{"name": 1}
		log.Dev(context, "GetQuerySetNames", "MGO : db.%s.find({}, %s).sort([\"name\"])", collection, mongo.Query(q))
		return c.Find(nil).Select(q).Sort("name").All(&names)
	}

	if err := mongo.ExecuteDB(context, ses, collection, f); err != nil {
		log.Error(context, "GetQuerySetNames", err, "Completed")
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

	log.Dev(context, "GetQuerySetNames", "Completed : RSN[%+v]", rsn)
	return rsn, nil
}

// GetQuerySet retrieves the configuration for the specified QuerySet.
func GetQuerySet(context interface{}, ses *mgo.Session, name string) (*QuerySet, error) {
	log.Dev(context, "GetQuerySet", "Started : Name[%s]", name)

	var rs QuerySet
	f := func(c *mgo.Collection) error {
		q := bson.M{"name": name}
		log.Dev(context, "GetQuerySet", "MGO : db.%s.findOne(%s)", collection, mongo.Query(q))
		return c.Find(q).One(&rs)
	}

	if err := mongo.ExecuteDB(context, ses, collection, f); err != nil {
		log.Error(context, "GetQuerySet", err, "Completed")
		return nil, err
	}

	log.Dev(context, "GetQuerySet", "Completed : RS[%+v]", rs)
	return &rs, nil
}

// UpdateQuerySet is used to create or update existing QuerySet documents.
func UpdateQuerySet(context interface{}, ses *mgo.Session, rs *QuerySet) error {
	log.Dev(context, "UpdateQuerySet", "Started : Name[%s]", rs.Name)

	f := func(c *mgo.Collection) error {
		q := bson.M{"name": rs.Name}

		log.Dev(context, "UpdateQuerySet", "MGO :\n\ndb.%s.upsert(%s, %s)\n", collection, mongo.Query(q), mongo.Query(rs))
		_, err := c.Upsert(q, rs)
		return err
	}

	if err := mongo.ExecuteDB(context, ses, collection, f); err != nil {
		log.Error(context, "UpdateQuerySet", err, "Completed")
		return err
	}

	log.Dev(context, "UpdateQuerySet", "Completed")
	return nil
}

// RemoveQuerySet is used to remove an existing QuerySet documents.
func RemoveQuerySet(context interface{}, ses *mgo.Session, name string) (*QuerySet, error) {
	log.Dev(context, "RemoveQuerySet", "Started : Name[%s]", name)

	var rs QuerySet
	f := func(c *mgo.Collection) error {
		q := bson.M{"name": name}

		log.Dev(context, "RemoveQuerySet", "MGO :\n\ndb.%s.remove(%s)\n", collection, mongo.Query(q))
		return c.Remove(q)
	}

	if err := mongo.ExecuteDB(context, ses, collection, f); err != nil {
		log.Error(context, "RemoveQuerySet", err, "Completed")
		return nil, err
	}

	log.Dev(context, "RemoveQuerySet", "Completed")
	return &rs, nil
}
