package wirefix

import (
	"encoding/json"
	"os"

	"github.com/ardanlabs/kit/db"
	"github.com/coralproject/shelf/internal/wire/pattern"
	"github.com/coralproject/shelf/internal/wire/relationship"
	"github.com/coralproject/shelf/internal/wire/view"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var path string

func init() {
	path = os.Getenv("GOPATH") + "/src/github.com/coralproject/shelf/internal/wire/wirefix/"
}

// Get loads relationship, view, and pattern data based on *.json files.
func Get() ([]relationship.Relationship, []view.View, []pattern.Pattern, error) {
	file, err := os.Open(path + "relationship.json")
	if err != nil {
		return nil, nil, nil, err
	}

	var rels []relationship.Relationship
	err = json.NewDecoder(file).Decode(&rels)
	if err != nil {
		return nil, nil, nil, err
	}
	file.Close()

	file, err = os.Open(path + "view.json")
	if err != nil {
		return nil, nil, nil, err
	}

	var views []view.View
	err = json.NewDecoder(file).Decode(&views)
	if err != nil {
		return nil, nil, nil, err
	}
	file.Close()

	file, err = os.Open(path + "pattern.json")
	if err != nil {
		return nil, nil, nil, err
	}

	var patterns []pattern.Pattern
	err = json.NewDecoder(file).Decode(&patterns)
	if err != nil {
		return nil, nil, nil, err
	}
	file.Close()

	return rels, views, patterns, nil
}

// Add inserts relationships, views, and patterns for testing.
func Add(context interface{}, db *db.DB, rels []relationship.Relationship, views []view.View, patterns []pattern.Pattern) error {
	for _, rel := range rels {
		if err := relationship.Upsert(context, db, &rel); err != nil {
			return err
		}
	}

	for _, pat := range patterns {
		if err := pattern.Upsert(context, db, &pat); err != nil {
			return err
		}
	}

	for _, vw := range views {
		if err := view.Upsert(context, db, &vw); err != nil {
			return err
		}
	}

	return nil
}

// Remove removes relationships, views, and patterns in Mongo that match a given pattern.
func Remove(context interface{}, db *db.DB, prefix string) error {
	f := func(c *mgo.Collection) error {
		q := bson.M{"predicate": bson.RegEx{Pattern: prefix}}
		_, err := c.RemoveAll(q)
		return err
	}

	if err := db.ExecuteMGO(context, relationship.Collection, f); err != nil {
		return err
	}

	f = func(c *mgo.Collection) error {
		q := bson.M{"name": bson.RegEx{Pattern: prefix}}
		_, err := c.RemoveAll(q)
		return err
	}

	if err := db.ExecuteMGO(context, view.Collection, f); err != nil {
		return err
	}

	f = func(c *mgo.Collection) error {
		q := bson.M{"type": bson.RegEx{Pattern: prefix}}
		_, err := c.RemoveAll(q)
		return err
	}

	if err := db.ExecuteMGO(context, pattern.Collection, f); err != nil {
		return err
	}

	return nil
}
