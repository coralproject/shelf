package query

import (
	"strings"

	"github.com/coralproject/shelf/pkg/db"
	"github.com/coralproject/shelf/pkg/db/mongo"
	"github.com/coralproject/shelf/pkg/log"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// collection contains the name of the query_sets collection.
const collection = "query_sets"
const collectionHistory = "query_set_history"

// =============================================================================

// UpsertSet is used to create or update an existing Set document.
func UpsertSet(context interface{}, db *db.DB, qs *Set) error {
	log.Dev(context, "UpsertSet", "Started : Name[%s]", qs.Name)

	// TODO: do we need this here, or should we create a custom upsert for this?
	// We need to keep track of history of upserts, ensuring that the last
	// we attained was safe before we update.
	// cqs, err := GetSetByName(context, db, qs.Name)
	// if err == nil {
	// 	err2 := UpsertHistory(context, db, cqs)
	// 	if err2 != nil {
	// 		return err2
	// 	}
	// }

	f := func(c *mgo.Collection) error {
		q := bson.M{"name": qs.Name}
		log.Dev(context, "UpsertSet", "MGO : db.%s.upsert(%s, %s)", collection, mongo.Query(q), mongo.Query(&qs))
		_, err := c.Upsert(q, &qs)
		return err
	}

	if err := db.ExecuteMGO(context, collection, f); err != nil {
		log.Error(context, "UpsertSet", err, "Completed")
		return err
	}

	log.Dev(context, "UpsertSet", "Completed")
	return nil
}

// =============================================================================

// GetSetNames retrieves a list of rule names.
func GetSetNames(context interface{}, db *db.DB) ([]string, error) {
	log.Dev(context, "GetSetNames", "Started")

	var names []bson.M
	f := func(c *mgo.Collection) error {
		q := bson.M{"name": 1}
		log.Dev(context, "GetSetNames", "MGO : db.%s.find({}, %s).sort([\"name\"])", collection, mongo.Query(q))
		return c.Find(nil).Select(q).Sort("name").All(&names)
	}

	if err := db.ExecuteMGO(context, collection, f); err != nil {
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
func GetSetByName(context interface{}, db *db.DB, name string) (*Set, error) {
	log.Dev(context, "GetSetByName", "Started : Name[%s]", name)

	var qs Set
	f := func(c *mgo.Collection) error {
		q := bson.M{"name": name}
		log.Dev(context, "GetSetByName", "MGO : db.%s.findOne(%s)", collection, mongo.Query(q))
		return c.Find(q).One(&qs)
	}

	if err := db.ExecuteMGO(context, collection, f); err != nil {
		log.Error(context, "GetSetByName", err, "Completed")
		return nil, err
	}

	log.Dev(context, "GetSetByName", "Completed : QS[%+v]", qs)
	return &qs, nil
}

// =============================================================================

// DeleteSet is used to remove an existing Set document.
func DeleteSet(context interface{}, db *db.DB, name string) error {
	log.Dev(context, "DeleteSet", "Started : Name[%s]", name)

	qs, err := GetSetByName(context, db, name)
	if err != nil {
		return err
	}

	f := func(c *mgo.Collection) error {
		q := bson.M{"name": qs.Name}
		log.Dev(context, "DeleteSet", "MGO : db.%s.remove(%s)", collection, mongo.Query(q))
		return c.Remove(q)
	}

	if err := db.ExecuteMGO(context, collection, f); err != nil {
		log.Error(context, "DeleteSet", err, "Completed")
		return err
	}

	log.Dev(context, "DeleteSet", "Completed")
	return nil
}

// =============================================================================

// GetLastSetHistoryByName gets the last written Set within the query_history
// collection and returns the last one else returns a non-nil error if it fails.
func GetLastSetHistoryByName(context interface{}, db *db.DB, name string) (*Set, error) {
	log.Dev(context, "GetLastSetHistoryByName", "Started : Name[%s]", name)

	var qs Set

	f := func(c *mgo.Collection) error {
		q := bson.M{"name": name}
		log.Dev(context, "GetLastSetHistoryByName", "MGO : db.%s.find(%s).count()", collection, mongo.Query(q))
		total, err := c.Find(q).Count()
		if err != nil {
			return err
		}

		beforeLast := total - 1
		log.Dev(context, "GetLastSetHistoryByName", "MGO : db.%s.find(%s).skip(%d).one()", collection, mongo.Query(q), mongo.Query(beforeLast))
		return c.Find(q).Skip(beforeLast).One(&qs)
	}

	err := db.ExecuteMGO(context, collectionHistory, f)
	if err != nil {
		log.Error(context, "GetLastSetHistoryByName", err, "Complete : Name[%s]", name)
		return nil, err
	}

	log.Dev(context, "GetLastSetHistoryByName", "Complete : Name[%s]", name)
	return &qs, nil
}

// =============================================================================

// UpsertHistory adds the last query record in the query collection into the
// list of query collection for that name.
// Providing a corruption or bad save mitigation tactic.
// Returns an error if the query was not found in the db records or if the
// update was not successful.
func UpsertHistory(context interface{}, db *db.DB, qs *Set) error {
	log.Dev(context, "upsertHistory", "Started : Name[%s]", qs.Name)

	f := func(c *mgo.Collection) error {
		q := bson.M{"name": qs.Name}
		qu := bson.M{
			"$push": bson.M{
				"history": qs,
			},
		}

		log.Dev(context, "upsertHistory", "MGO : db.%s.upsert(%s, %s)", collection, mongo.Query(q), mongo.Query(qu))
		_, err := c.Upsert(q, qu)
		return err
	}

	err := db.ExecuteMGO(context, collectionHistory, f)
	if err != nil {
		log.Error(context, "upsertHistory", err, "Complete : Name[%s]", qs.Name)
		return err
	}

	log.Dev(context, "upsertHistory", "Complete : Name[%s]", qs.Name)
	return nil
}
