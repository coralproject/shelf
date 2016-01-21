// Package regex provides the service layer for building apps using
// regex functionality.
package regex

import (
	"errors"

	"github.com/ardanlabs/kit/db"
	"github.com/ardanlabs/kit/db/mongo"
	"github.com/ardanlabs/kit/log"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Contains the name of Mongo collections.
const (
	Collection        = "query_regexs"
	CollectionHistory = "query_regexs_history"
)

// Set of error variables.
var (
	ErrNotFound = errors.New("Set Not found")
)

// =============================================================================

// Upsert is used to create or update an existing Regex document.
func Upsert(context interface{}, db *db.DB, rgx *Regex) error {
	log.Dev(context, "Upsert", "Started : Name[%s]", rgx.Name)

	// Validate the regex that is provided.
	if err := rgx.Validate(); err != nil {
		log.Error(context, "Upsert", err, "Completed")
		return err
	}

	// We need to know if this is a new regex.
	var new bool
	if _, err := GetByName(context, db, rgx.Name); err != nil {
		if err != ErrNotFound {
			log.Error(context, "Upsert", err, "Completed")
			return err
		}

		new = true
	}

	// Insert or update the query regex.
	f := func(c *mgo.Collection) error {
		q := bson.M{"name": rgx.Name}
		log.Dev(context, "Upsert", "MGO : db.%s.upsert(%s, %s)", c.Name, mongo.Query(q), mongo.Query(rgx))
		_, err := c.Upsert(q, rgx)
		return err
	}

	if err := db.ExecuteMGO(context, Collection, f); err != nil {
		log.Error(context, "Upsert", err, "Completed")
		return err
	}

	// Add a history record if this query regex is new.
	if new {
		f = func(c *mgo.Collection) error {
			qh := bson.M{
				"name":   rgx.Name,
				"regexs": []bson.M{},
			}

			log.Dev(context, "Upsert", "MGO : db.%s.insert(%s)", c.Name, mongo.Query(qh))
			return c.Insert(qh)
		}

		if err := db.ExecuteMGO(context, CollectionHistory, f); err != nil {
			log.Error(context, "Upsert", err, "Completed")
			return err
		}
	}

	// Add this query set to the beginning of the history.
	f = func(c *mgo.Collection) error {
		q := bson.M{"name": rgx.Name}
		qu := bson.M{
			"$push": bson.M{
				"regexs": bson.M{
					"$each":     []*Regex{rgx},
					"$position": 0,
				},
			},
		}

		log.Dev(context, "Upsert", "MGO : db.%s.update(%s, %s)", c.Name, mongo.Query(q), mongo.Query(qu))
		_, err := c.Upsert(q, qu)
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

// GetNames retrieves a list of query regex names.
func GetNames(context interface{}, db *db.DB) ([]string, error) {
	log.Dev(context, "GetNames", "Started")

	var rawNames []struct {
		Name string
	}

	f := func(c *mgo.Collection) error {
		s := bson.M{"name": 1}
		log.Dev(context, "GetNames", "MGO : db.%s.find({}, %s).sort([\"name\"])", c.Name, mongo.Query(s))
		return c.Find(nil).Select(s).Sort("name").All(&rawNames)
	}

	if err := db.ExecuteMGO(context, Collection, f); err != nil {
		if err == mgo.ErrNotFound {
			err = ErrNotFound
		}

		log.Error(context, "GetNames", err, "Completed")
		return nil, err
	}

	names := make([]string, len(rawNames))
	for i := range rawNames {
		names[i] = rawNames[i].Name
	}

	log.Dev(context, "GetNames", "Completed : Sets[%d]", len(names))
	return names, nil
}

// GetRegexs retrieves a list of regexs.
func GetRegexs(context interface{}, db *db.DB, tags []string) ([]Regex, error) {
	log.Dev(context, "GetSets", "Started : Tags[%v]", tags)

	var rgxs []Regex
	f := func(c *mgo.Collection) error {
		log.Dev(context, "GetSets", "MGO : db.%s.find({}).sort([\"name\"])", c.Name)
		return c.Find(nil).All(&rgxs)
	}

	if err := db.ExecuteMGO(context, Collection, f); err != nil {
		if err == mgo.ErrNotFound {
			err = ErrNotFound
		}

		log.Error(context, "GetSets", err, "Completed")
		return nil, err
	}

	log.Dev(context, "GetSets", "Completed : Sets[%d]", len(rgxs))
	return rgxs, nil
}

// GetByName retrieves the document for the specified Regex.
func GetByName(context interface{}, db *db.DB, name string) (*Regex, error) {
	log.Dev(context, "GetByName", "Started : Name[%s]", name)

	var rgx Regex
	f := func(c *mgo.Collection) error {
		q := bson.M{"name": name}
		log.Dev(context, "GetByName", "MGO : db.%s.findOne(%s)", c.Name, mongo.Query(q))
		return c.Find(q).One(&rgx)
	}

	if err := db.ExecuteMGO(context, Collection, f); err != nil {
		if err == mgo.ErrNotFound {
			err = ErrNotFound
		}

		log.Error(context, "GetByName", err, "Completed")
		return nil, err
	}

	log.Dev(context, "GetByName", "Completed : Set[%+v]", &rgx)
	return &rgx, nil
}

// GetByNames retrieves the documents for the specified names.
func GetByNames(context interface{}, db *db.DB, names []string) ([]Regex, error) {
	log.Dev(context, "GetByNames", "Started : Names[%+v]", names)

	var rgxs []Regex
	f := func(c *mgo.Collection) error {

		// Build a list of documents to find by name.
		qn := make([]bson.M, len(names))
		for i, name := range names {
			if name != "" {
				qn[i] = bson.M{"name": name}
			}
		}

		// Place that list in an $or operation.
		q := bson.M{"$or": qn}

		log.Dev(context, "GetByNames", "MGO : db.%s.find(%s)", c.Name, mongo.Query(q))
		return c.Find(q).All(&rgxs)
	}

	if err := db.ExecuteMGO(context, Collection, f); err != nil {
		if err == mgo.ErrNotFound {
			err = ErrNotFound
		}

		log.Error(context, "GetByNames", err, "Completed")
		return nil, err
	}

	// I can't assume MongoDB will bring the results back in the order I
	// setup the query. I need the order to match on the returned slice.
	// I thought about using a map of name/value but I feel like it is overkill.

	regexs := make([]Regex, len(names))
next:
	for _, rgx := range rgxs {
		for i := range names {
			if rgx.Name == names[i] {
				regexs[i] = rgx
				continue next
			}
		}
	}

	log.Dev(context, "GetByNames", "Completed : Regexs[%+v]", regexs)
	return regexs, nil
}

// GetLastHistoryByName gets the last written Regex within the history.
func GetLastHistoryByName(context interface{}, db *db.DB, name string) (*Regex, error) {
	log.Dev(context, "GetLastHistoryByName", "Started : Name[%s]", name)

	var result struct {
		Name   string  `bson:"name"`
		Regexs []Regex `bson:"regexs"`
	}

	f := func(c *mgo.Collection) error {
		q := bson.M{"name": name}
		proj := bson.M{"regexs": bson.M{"$slice": 1}}

		log.Dev(context, "GetLastHistoryByName", "MGO : db.%s.find(%s,%s)", c.Name, mongo.Query(q), mongo.Query(proj))
		return c.Find(q).Select(proj).One(&result)
	}

	if err := db.ExecuteMGO(context, CollectionHistory, f); err != nil {
		if err == mgo.ErrNotFound {
			err = ErrNotFound
		}

		log.Error(context, "GetLastHistoryByName", err, "Complete")
		return nil, err
	}

	if result.Regexs == nil {
		err := errors.New("History not found")
		log.Error(context, "GetLastHistoryByName", err, "Complete")
		return nil, err
	}

	log.Dev(context, "GetLastHistoryByName", "Completed : Regex[%+v]", &result.Regexs[0])
	return &result.Regexs[0], nil
}

// =============================================================================

// Delete is used to remove an existing Regex document.
func Delete(context interface{}, db *db.DB, name string) error {
	log.Dev(context, "Delete", "Started : Name[%s]", name)

	rgx, err := GetByName(context, db, name)
	if err != nil {
		return err
	}

	f := func(c *mgo.Collection) error {
		q := bson.M{"name": rgx.Name}
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
