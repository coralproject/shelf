package sfix

import (
	"encoding/json"
	"os"

	"github.com/ardanlabs/kit/db"
	"github.com/ardanlabs/kit/tests"
	"github.com/coralproject/shelf/internal/xenia/script"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var path string

func init() {
	path = os.Getenv("GOPATH") + "/src/github.com/coralproject/shelf/internal/xenia/script/sfix/"
}

//==============================================================================

// Get retrieves a set document from the filesystem for testing.
func Get(fileName string) (script.Script, error) {
	file, err := os.Open(path + fileName)
	if err != nil {
		return script.Script{}, err
	}

	defer file.Close()

	var scr script.Script
	err = json.NewDecoder(file).Decode(&scr)
	if err != nil {
		return script.Script{}, err
	}

	return scr, nil
}

// Add inserts a script for testing.
func Add(db *db.DB, scr script.Script) error {
	if err := script.Upsert("", db, scr); err != nil {
		return err
	}

	return nil
}

// Remove is used to clear out all the test sets from the collection.
// All test documents must start with STEST in their name.
func Remove(db *db.DB, pattern string) error {
	f := func(c *mgo.Collection) error {
		q := bson.M{"name": bson.RegEx{Pattern: pattern}}
		_, err := c.RemoveAll(q)
		return err
	}

	if err := db.ExecuteMGO(tests.Context, script.Collection, f); err != nil {
		return err
	}

	f = func(c *mgo.Collection) error {
		q := bson.M{"name": bson.RegEx{Pattern: pattern}}
		_, err := c.RemoveAll(q)
		return err
	}

	if err := db.ExecuteMGO(tests.Context, script.CollectionHistory, f); err != nil {
		return err
	}

	return nil
}
