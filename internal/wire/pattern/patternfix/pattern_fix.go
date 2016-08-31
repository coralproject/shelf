package patternfix

import (
	"encoding/json"
	"os"

	"github.com/ardanlabs/kit/db"
	"github.com/coralproject/shelf/internal/wire/pattern"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var path string

func init() {
	path = os.Getenv("GOPATH") + "/src/github.com/coralproject/shelf/internal/wire/pattern/patternfix/"
}

// Get loads pattern data based on patterns.json.
func Get() ([]pattern.Pattern, error) {
	file, err := os.Open(path + "patterns.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var patterns []pattern.Pattern
	err = json.NewDecoder(file).Decode(&patterns)
	if err != nil {
		return nil, err
	}

	return patterns, nil
}

// Add inserts patterns for testing.
func Add(context interface{}, db *db.DB, patterns []pattern.Pattern) error {
	for _, pat := range patterns {
		if err := pattern.Upsert(context, db, &pat); err != nil {
			return err
		}
	}

	return nil
}

// Remove removes patterns in Mongo that match a given prefix.
func Remove(context interface{}, db *db.DB, prefix string) error {
	f := func(c *mgo.Collection) error {
		q := bson.M{"type": bson.RegEx{Pattern: prefix}}
		_, err := c.RemoveAll(q)
		return err
	}

	if err := db.ExecuteMGO(context, pattern.Collection, f); err != nil {
		return err
	}

	return nil
}
