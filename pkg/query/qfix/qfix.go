package qfix

import (
	"encoding/json"
	"os"

	"github.com/coralproject/xenia/pkg/query"

	"github.com/ardanlabs/kit/db"
	"github.com/ardanlabs/kit/tests"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var path string

func init() {
	path = os.Getenv("GOPATH") + "/src/github.com/coralproject/xenia/pkg/query/qfix/"
}

//==============================================================================

// Get retrieves a set document from the filesystem for testing.
func Get(fileName string) (*query.Set, error) {
	file, err := os.Open(path + fileName)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	var set query.Set
	err = json.NewDecoder(file).Decode(&set)
	if err != nil {
		return nil, err
	}

	return &set, nil
}

// Add inserts a set for testing.
func Add(db *db.DB, set *query.Set) error {
	if err := query.Upsert("", db, set); err != nil {
		return err
	}

	return nil
}

// Remove is used to clear out all the test sets from the collection.
// All test documents must start with QSTEST in their name.
func Remove(db *db.DB, pattern string) error {
	f := func(c *mgo.Collection) error {
		q := bson.M{"name": bson.RegEx{Pattern: pattern}}
		_, err := c.RemoveAll(q)
		return err
	}

	if err := db.ExecuteMGO(tests.Context, query.Collection, f); err != nil {
		return err
	}

	f = func(c *mgo.Collection) error {
		q := bson.M{"name": bson.RegEx{Pattern: pattern}}
		_, err := c.RemoveAll(q)
		return err
	}

	if err := db.ExecuteMGO(tests.Context, query.CollectionHistory, f); err != nil {
		return err
	}

	return nil
}
