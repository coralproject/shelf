package script

import (
	"errors"
	"strings"

	"github.com/ardanlabs/kit/db"
	"github.com/ardanlabs/kit/db/mongo"
	"github.com/ardanlabs/kit/log"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Contains the name of Mongo collections.
const (
	Collection         = "scripts"
	CollectionHistory  = "scripts_history"
	CollectionExecTest = "test_scripts"
)

// =============================================================================

// Upsert is used to create or update an existing Script document.
func Upsert(context interface{}, db *db.DB, scr *Script) error {
	log.Dev(context, "Upsert", "Started : Name[%s]", scr.Name)

	// Validate the set that is provided.
	if err := scr.Validate(); err != nil {
		log.Error(context, "Upsert", err, "Completed")
		return err
	}

	// We need to know if this is a new Set.
	var new bool
	if _, err := GetByName(context, db, scr.Name); err != nil {
		if err != mgo.ErrNotFound {
			log.Error(context, "Upsert", err, "Completed")
			return err
		}

		new = true
	}

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
					"$each":     []*Script{scr},
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

	var names []bson.M
	f := func(c *mgo.Collection) error {
		s := bson.M{"name": 1}
		log.Dev(context, "GetNames", "MGO : db.%s.find({}, %s).sort([\"name\"])", mongo.Query(s))
		return c.Find(nil).Select(s).Sort("name").All(&names)
	}

	if err := db.ExecuteMGO(context, Collection, f); err != nil {
		log.Error(context, "GetNames", err, "Completed")
		return nil, err
	}

	var sets []string
	for _, doc := range names {
		name := doc["name"].(string)
		if strings.HasPrefix(name, "test") {
			continue
		}

		sets = append(sets, name)
	}

	log.Dev(context, "GetNames", "Completed : Sets[%+v]", sets)
	return sets, nil
}

// GetByName retrieves the document for the specified name.
func GetByName(context interface{}, db *db.DB, name string) (*Script, error) {
	log.Dev(context, "GetByName", "Started : Name[%s]", name)

	var scr Script
	f := func(c *mgo.Collection) error {
		q := bson.M{"name": name}
		log.Dev(context, "GetByName", "MGO : db.%s.findOne(%s)", c.Name, mongo.Query(q))
		return c.Find(q).One(&scr)
	}

	if err := db.ExecuteMGO(context, Collection, f); err != nil {
		log.Error(context, "GetByName", err, "Completed")
		return nil, err
	}

	log.Dev(context, "GetByName", "Completed : Script[%+v]", &scr)
	return &scr, nil
}

// GetByNames retrieves the documents for the specified names.
func GetByNames(context interface{}, db *db.DB, names []string) ([]Script, error) {
	log.Dev(context, "GetByNames", "Started : Names[%+v]", names)

	var scrs []Script
	f := func(c *mgo.Collection) error {
		qn := make([]bson.M, len(names))
		for i, name := range names {
			if name != "" {
				qn[i] = bson.M{"name": name}
			}
		}

		q := bson.M{"$or": qn}
		log.Dev(context, "GetByNames", "MGO : db.%s.find(%s)", c.Name, mongo.Query(q))
		return c.Find(q).All(&scrs)
	}

	if err := db.ExecuteMGO(context, Collection, f); err != nil {
		log.Error(context, "GetByNames", err, "Completed")
		return nil, err
	}

	// I can't assume MongoDB will bring it back in the order I
	// setup the query. I want the order to match on the returned slice.
	// I thought about using a map but I feel like it is overkill.

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

	log.Dev(context, "GetByNames", "Completed : Scripts[%+v]", scripts)
	return scripts, nil
}

// GetLastHistoryByName gets the last written Set within the query_history
// collection and returns the last one else returns a non-nil error if it fails.
func GetLastHistoryByName(context interface{}, db *db.DB, name string) (*Script, error) {
	log.Dev(context, "GetLastHistoryByName", "Started : Name[%s]", name)

	var result struct {
		Name    string   `bson:"name"`
		Scripts []Script `bson:"scripts"`
	}

	f := func(c *mgo.Collection) error {
		q := bson.M{"name": name}
		proj := bson.M{"scripts": bson.M{"$slice": 1}}

		log.Dev(context, "GetLastHistoryByName", "MGO : db.%s.find(%s,%s)", c.Name, mongo.Query(q), mongo.Query(proj))
		return c.Find(q).Select(proj).One(&result)
	}

	err := db.ExecuteMGO(context, CollectionHistory, f)
	if err != nil {
		log.Error(context, "GetLastHistoryByName", err, "Complete")
		return nil, err
	}

	if result.Scripts == nil {
		err := errors.New("History not found")
		log.Error(context, "GetLastHistoryByName", err, "Complete")
		return nil, err
	}

	log.Dev(context, "GetLastHistoryByName", "Completed : Scripts[%+v]", &result.Scripts[0])
	return &result.Scripts[0], nil
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

	log.Dev(context, "Delete", "Completed")
	return nil
}
