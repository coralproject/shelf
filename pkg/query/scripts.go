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

// 	// Insert or update the query set.
// 	f := func(c *mgo.Collection) error {
// 		q := bson.M{"name": set.Name}
// 		log.Dev(context, "UpsertScript", "MGO : db.%s.upsert(%s, %s)", c.Name, mongo.Query(q), mongo.Query(src))
// 		_, err := c.Upsert(q, set)
// 		return err
// 	}

// 	if err := db.ExecuteMGO(context, Collection, f); err != nil {
// 		log.Error(context, "UpsertSet", err, "Completed")
// 		return err
// 	}

// 	// Add a history record if this query set is new.
// 	if new {
// 		f = func(c *mgo.Collection) error {
// 			qh := bson.M{
// 				"name": set.Name,
// 				"sets": []bson.M{},
// 			}

// 			log.Dev(context, "UpsertSet", "MGO : db.%s.insert(%s)", c.Name, mongo.Query(qh))
// 			return c.Insert(qh)
// 		}

// 		if err := db.ExecuteMGO(context, CollectionHistory, f); err != nil {
// 			log.Error(context, "UpsertSet", err, "Completed")
// 			return err
// 		}
// 	}

// 	// Add this query set to the beginning of the history.
// 	f = func(c *mgo.Collection) error {
// 		q := bson.M{"name": set.Name}
// 		qu := bson.M{
// 			"$push": bson.M{
// 				"sets": bson.M{
// 					"$each":     []*Set{set},
// 					"$position": 0,
// 				},
// 			},
// 		}

// 		log.Dev(context, "UpsertSet", "MGO : db.%s.update(%s, %s)", c.Name, mongo.Query(q), mongo.Query(qu))
// 		_, err := c.Upsert(q, qu)
// 		return err
// 	}

// 	if err := db.ExecuteMGO(context, CollectionHistory, f); err != nil {
// 		log.Error(context, "UpsertSet", err, "Completed")
// 		return err
// 	}

// 	log.Dev(context, "UpsertSet", "Completed")
// 	return nil
// }
