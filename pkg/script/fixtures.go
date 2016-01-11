package script

import (
	"github.com/ardanlabs/kit/db"
	"github.com/ardanlabs/kit/tests"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// RemoveTestData is used to clear out all the test data from the collection.
// All test documents must start with QSTEST in their name.
func RemoveTestData(db *db.DB) error {
	f := func(c *mgo.Collection) error {
		q := bson.M{"name": bson.RegEx{Pattern: "STEST"}}
		_, err := c.RemoveAll(q)
		return err
	}

	if err := db.ExecuteMGO(tests.Context, Collection, f); err != nil {
		return err
	}

	f = func(c *mgo.Collection) error {
		q := bson.M{"name": bson.RegEx{Pattern: "STEST"}}
		_, err := c.RemoveAll(q)
		return err
	}

	if err := db.ExecuteMGO(tests.Context, CollectionHistory, f); err != nil {
		return err
	}

	return nil
}
