package query

import (
	"strings"

	"github.com/coralproject/shelf/pkg/log"
	"github.com/coralproject/shelf/pkg/srv/mongo"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// collections contains the name of the rules collection.
const collection = "query_sets"

// =============================================================================

// CreateSet is used to create Set documents in the db.
func CreateSet(context interface{}, ses *mgo.Session, qs *Set) error {
	log.Dev(context, "CreateSet", "Started : Name[%s]", qs.Name)

	f := func(c *mgo.Collection) error {
		log.Dev(context, "CreateSet", "MGO : db.%s.insert(%s)", collection, mongo.Query(qs))
		return c.Insert(&qs)
	}

	if err := mongo.ExecuteDB(context, ses, collection, f); err != nil {
		log.Error(context, "CreateSet", err, "Completed")
		return err
	}

	log.Dev(context, "CreateSet", "Completed")
	return nil
}

// =============================================================================

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

	var qsn []string
	for _, doc := range names {
		name := doc["name"].(string)
		if strings.HasPrefix(name, "test") {
			continue
		}

		qsn = append(qsn, name)
	}

	log.Dev(context, "GetSetNames", "Completed : QSN[%+v]", qsn)
	return qsn, nil
}

// GetSetByName retrieves the configuration for the specified Set.
func GetSetByName(context interface{}, ses *mgo.Session, name string) (*Set, error) {
	log.Dev(context, "GetSetByName", "Started : Name[%s]", name)

	var qs Set
	f := func(c *mgo.Collection) error {
		q := bson.M{"name": name}
		log.Dev(context, "GetSetByName", "MGO : db.%s.findOne(%s)", collection, mongo.Query(q))
		return c.Find(q).One(&qs)
	}

	if err := mongo.ExecuteDB(context, ses, collection, f); err != nil {
		log.Error(context, "GetSetByName", err, "Completed")
		return nil, err
	}

	log.Dev(context, "GetSetByName", "Completed : QS[%+v]", qs)
	return &qs, nil
}

// =============================================================================

// UpdateSet is used to update an existing Set document.
func UpdateSet(context interface{}, ses *mgo.Session, qs *Set) error {
	log.Dev(context, "UpdateSet", "Started : Name[%s]", qs.Name)

	f := func(c *mgo.Collection) error {
		q := bson.M{"name": qs.Name}

		log.Dev(context, "UpdateSet", "MGO : db.%s.upsert(%s, %s)", collection, mongo.Query(q), mongo.Query(&qs))
		_, err := c.Upsert(q, &qs)
		return err
	}

	if err := mongo.ExecuteDB(context, ses, collection, f); err != nil {
		log.Error(context, "UpdateSet", err, "Completed")
		return err
	}

	log.Dev(context, "UpdateSet", "Completed")
	return nil
}

// =============================================================================

// DeleteSet is used to remove an existing Set document.
func DeleteSet(context interface{}, ses *mgo.Session, name string) error {
	log.Dev(context, "DeleteSet", "Started : Name[%s]", name)

	qs, err := GetSetByName(context, ses, name)
	if err != nil {
		return err
	}

	f := func(c *mgo.Collection) error {
		q := bson.M{"name": qs.Name}
		log.Dev(context, "DeleteSet", "MGO : db.%s.remove(%s)", collection, mongo.Query(q))
		return c.Remove(q)
	}

	if err := mongo.ExecuteDB(context, ses, collection, f); err != nil {
		log.Error(context, "DeleteSet", err, "Completed")
		return err
	}

	log.Dev(context, "DeleteSet", "Completed")
	return nil
}
