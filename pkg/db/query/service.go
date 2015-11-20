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

// GetRuleSetNames retrieves a list of rule names.
func GetRuleSetNames(context interface{}, ses *mgo.Session) ([]string, error) {
	log.Dev(context, "GetRuleSetNames", "Started")

	var names []bson.M
	f := func(c *mgo.Collection) error {
		q := bson.M{"name": 1}
		log.Dev(context, "GetRuleSetNames", "MGO : db.%s.find({}, %s).sort([\"name\"])", collection, mongo.Query(q))
		return c.Find(nil).Select(q).Sort("name").All(&names)
	}

	if err := mongo.ExecuteDB(context, ses, collection, f); err != nil {
		log.Error(context, "GetRuleSetNames", err, "Completed")
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

	log.Dev(context, "GetRuleSetNames", "Completed : RSN[%+v]", rsn)
	return rsn, nil
}

// GetRuleSet retrieves the configuration for the specified RuleSet.
func GetRuleSet(context interface{}, ses *mgo.Session, name string) (*RuleSet, error) {
	log.Dev(context, "GetRuleSet", "Started : Name[%s]", name)

	var rs RuleSet
	f := func(c *mgo.Collection) error {
		q := bson.M{"name": name}
		log.Dev(context, "GetRuleSet", "MGO : db.%s.findOne(%s)", collection, mongo.Query(q))
		return c.Find(q).One(&rs)
	}

	if err := mongo.ExecuteDB(context, ses, collection, f); err != nil {
		log.Error(context, "GetRuleSet", err, "Completed")
		return nil, err
	}

	log.Dev(context, "GetRuleSet", "Completed : RS[%+v]", rs)
	return &rs, nil
}

// UpdateRuleSet is used to create or update existing RuleSet documents.
func UpdateRuleSet(context interface{}, ses *mgo.Session, rs *RuleSet) error {
	log.Dev(context, "UpdateRuleSet", "Started : Name[%s]", rs.Name)

	f := func(c *mgo.Collection) error {
		q := bson.M{"name": rs.Name}

		log.Dev(context, "UpdateRuleSet", "MGO :\n\ndb.%s.upsert(%s, %s)\n", collection, mongo.Query(q), mongo.Query(rs))
		_, err := c.Upsert(q, rs)
		return err
	}

	if err := mongo.ExecuteDB(context, ses, collection, f); err != nil {
		log.Error(context, "UpdateRuleSet", err, "Completed")
		return err
	}

	log.Dev(context, "UpdateRuleSet", "Completed")
	return nil
}

// RemoveRuleSet is used to remove an existing RuleSet documents.
func RemoveRuleSet(context interface{}, ses *mgo.Session, name string) (*RuleSet, error) {
	log.Dev(context, "RemoveRuleSet", "Started : Name[%s]", name)

	var rs RuleSet
	f := func(c *mgo.Collection) error {
		q := bson.M{"name": name}

		log.Dev(context, "RemoveRuleSet", "MGO :\n\ndb.%s.remove(%s)\n", collection, mongo.Query(q))
		return c.Remove(q)
	}

	if err := mongo.ExecuteDB(context, ses, collection, f); err != nil {
		log.Error(context, "RemoveRuleSet", err, "Completed")
		return nil, err
	}

	log.Dev(context, "RemoveRuleSet", "Completed")
	return &rs, nil
}
