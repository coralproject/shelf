package query

// import (
// 	"errors"
// 	"strings"

// 	"github.com/ardanlabs/kit/db"
// 	"github.com/ardanlabs/kit/db/mongo"
// 	"github.com/ardanlabs/kit/log"

// 	"gopkg.in/mgo.v2"
// 	"gopkg.in/mgo.v2/bson"
// )

// // =============================================================================

// // UpsertScript is used to create or update an existing Script document.
// func UpsertScript(context interface{}, db *db.DB, scr *Script) error {
// 	log.Dev(context, "UpsertSscript", "Started : Name[%s]", scr.Name)

// 	// Validate the set that is provided.
// 	if err := scr.Validate(); err != nil {
// 		log.Error(context, "UpsertScript", err, "Completed")
// 		return err
// 	}

// 	// We need to know if this is a new script.
// 	var new bool
// 	if _, err := GetScriptByName(context, db, scr.Name); err != nil {
// 		if err != mgo.ErrNotFound {
// 			log.Error(context, "UpsertScript", err, "Completed")
// 			return err
// 		}

// 		new = true
// 	}

// 	// Insert or update the script.
// 	f := func(c *mgo.Collection) error {
// 		q := bson.M{"name": set.Name}
// 		log.Dev(context, "UpsertScript", "MGO : db.%s.upsert(%s, %s)", c.Name, mongo.Query(q), mongo.Query(src))
// 		_, err := c.Upsert(q, set)
// 		return err
// 	}

// 	if err := db.ExecuteMGO(context, Collection, f); err != nil {
// 		log.Error(context, "UpsertScript", err, "Completed")
// 		return err
// 	}

// 	// Add a history record if this script set is new.
// 	if new {
// 		f = func(c *mgo.Collection) error {
// 			sh := bson.M{
// 				"name": scr.Name,
// 				"scripts": []bson.M{},
// 			}

// 			log.Dev(context, "UpsertScript", "MGO : db.%s.insert(%s)", c.Name, mongo.Query(sh))
// 			return c.Insert(sh)
// 		}

// 		if err := db.ExecuteMGO(context, CollectionHistory, f); err != nil {
// 			log.Error(context, "UpsertScript", err, "Completed")
// 			return err
// 		}
// 	}

// 	// Add this script to the beginning of the history.
// 	f = func(c *mgo.Collection) error {
// 		q := bson.M{"name": scr.Name}
// 		su := bson.M{
// 			"$push": bson.M{
// 				"sets": bson.M{
// 					"$each":     []*Script{scr},
// 					"$position": 0,
// 				},
// 			},
// 		}

// 		log.Dev(context, "UpsertScript", "MGO : db.%s.update(%s, %s)", c.Name, mongo.Query(q), mongo.Query(su))
// 		_, err := c.Upsert(q, qu)
// 		return err
// 	}

// 	if err := db.ExecuteMGO(context, CollectionHistory, f); err != nil {
// 		log.Error(context, "UpsertScript", err, "Completed")
// 		return err
// 	}

// 	log.Dev(context, "UpsertScript", "Completed")
// 	return nil
// }

// // =============================================================================

// // GetScriptNames retrieves a list of script names.
// func GetScriptNames(context interface{}, db *db.DB) ([]string, error) {
// 	log.Dev(context, "GetScriptNames", "Started")

// 	var names []bson.M
// 	f := func(c *mgo.Collection) error {
// 		q := bson.M{"name": 1}
// 		log.Dev(context, "GetScriptNames", "MGO : db.%s.find({}, %s).sort([\"name\"])", c.Name, mongo.Query(q))
// 		return c.Find(nil).Select(q).Sort("name").All(&names)
// 	}

// 	if err := db.ExecuteMGO(context, Collection, f); err != nil {
// 		log.Error(context, "GetScriptNames", err, "Completed")
// 		return nil, err
// 	}

// 	var sets []string
// 	for _, doc := range names {
// 		name := doc["name"].(string)
// 		if strings.HasPrefix(name, "test") {
// 			continue
// 		}

// 		sets = append(sets, name)
// 	}

// 	log.Dev(context, "GetScriptNames", "Completed : Sets[%+v]", sets)
// 	return sets, nil
// }

// // GetSetByName retrieves the configuration for the specified Set.
// func GetSetByName(context interface{}, db *db.DB, name string) (*Set, error) {
// 	log.Dev(context, "GetSetByName", "Started : Name[%s]", name)

// 	var set Set
// 	f := func(c *mgo.Collection) error {
// 		q := bson.M{"name": name}
// 		log.Dev(context, "GetSetByName", "MGO : db.%s.findOne(%s)", c.Name, mongo.Query(q))
// 		return c.Find(q).One(&set)
// 	}

// 	if err := db.ExecuteMGO(context, Collection, f); err != nil {
// 		log.Error(context, "GetSetByName", err, "Completed")
// 		return nil, err
// 	}

// 	log.Dev(context, "GetSetByName", "Completed : Set[%+v]", &set)
// 	return &set, nil
// }

// // GetLastSetHistoryByName gets the last written Set within the query_history
// // collection and returns the last one else returns a non-nil error if it fails.
// func GetLastSetHistoryByName(context interface{}, db *db.DB, name string) (*Set, error) {
// 	log.Dev(context, "GetLastSetHistoryByName", "Started : Name[%s]", name)

// 	var result struct {
// 		Name string `bson:"name"`
// 		Sets []Set  `bson:"sets"`
// 	}

// 	f := func(c *mgo.Collection) error {
// 		q := bson.M{"name": name}
// 		proj := bson.M{"sets": bson.M{"$slice": 1}}

// 		log.Dev(context, "GetLastSetHistoryByName", "MGO : db.%s.find(%s,%s)", c.Name, mongo.Query(q), mongo.Query(proj))
// 		return c.Find(q).Select(proj).One(&result)
// 	}

// 	err := db.ExecuteMGO(context, CollectionHistory, f)
// 	if err != nil {
// 		log.Error(context, "GetLastSetHistoryByName", err, "Complete")
// 		return nil, err
// 	}

// 	if result.Sets == nil {
// 		err := errors.New("History not found")
// 		log.Error(context, "GetLastSetHistoryByName", err, "Complete")
// 		return nil, err
// 	}

// 	log.Dev(context, "GetLastSetHistoryByName", "Completed : QS[%+v]", &result.Sets[0])
// 	return &result.Sets[0], nil
// }
