package rfix

import (
	"encoding/json"
	"os"

	"github.com/coralproject/xenia/pkg/regex"

	"github.com/ardanlabs/kit/db"
	"github.com/ardanlabs/kit/tests"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var path string

func init() {
	path = os.Getenv("GOPATH") + "/src/github.com/coralproject/xenia/pkg/regex/rfix/"
}

//==============================================================================

// Get retrieves a regex document from the filesystem for testing.
func Get(fileName string) (*regex.Regex, error) {
	file, err := os.Open(path + fileName)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	var rgx regex.Regex
	err = json.NewDecoder(file).Decode(&rgx)
	if err != nil {
		return nil, err
	}

	return &rgx, nil
}

// Add inserts a regex for testing.
func Add(db *db.DB, rgx *regex.Regex) error {
	f := func(c *mgo.Collection) error {
		q := bson.M{"name": rgx.Name}
		_, err := c.Upsert(q, rgx)
		return err
	}

	if err := db.ExecuteMGO(tests.Context, regex.Collection, f); err != nil {
		return err
	}

	return nil
}

// Remove is used to clear out all the test regexs from the collection.
// All test documents must start with QSTEST in their name.
func Remove(db *db.DB, pattern string) error {
	f := func(c *mgo.Collection) error {
		q := bson.M{"name": bson.RegEx{Pattern: pattern}}
		_, err := c.RemoveAll(q)
		return err
	}

	if err := db.ExecuteMGO(tests.Context, regex.Collection, f); err != nil {
		return err
	}

	f = func(c *mgo.Collection) error {
		q := bson.M{"name": bson.RegEx{Pattern: pattern}}
		_, err := c.RemoveAll(q)
		return err
	}

	if err := db.ExecuteMGO(tests.Context, regex.CollectionHistory, f); err != nil {
		return err
	}

	return nil
}
