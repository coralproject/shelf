package ifix

import (
	"encoding/json"
	"os"

	"github.com/coralproject/xenia/internal/item"

	"github.com/ardanlabs/kit/db"
	"github.com/ardanlabs/kit/tests"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var path string

func init() {
	path = os.Getenv("GOPATH") + "/src/github.com/coralproject/xenia/internal/item/ifix/"
}

//==============================================================================

// Get retrieves some ItemData for testing
func Get(fileName string) (*[]map[string]interface{}, error) {
	file, err := os.Open(path + fileName)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	var data []map[string]interface{}
	err = json.NewDecoder(file).Decode(&data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func RegisterTypes(fileName string) error {

	file, err := os.Open(path + fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	// grab the item type fixture file
	var types []item.Type
	err = json.NewDecoder(file).Decode(&types)
	if err != nil {
		return err
	}

	// retister the item types
	for _, t := range types {
		item.RegisterType(t)
	}

	return nil
}

// InserFromDataFile reads a file of data and attempts to
//  create items of a provided type and insert them
func InsertItemsFromDataFile(context interface{}, db *db.DB, fileName string, t string) error {

	file, err := os.Open(path + fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	// grab the item type fixture file
	var data []map[string]interface{}
	err = json.NewDecoder(file).Decode(&data)

	if err != nil {
		return err
	}

	for _, d := range data {
		i, err := item.Create(context, db, t, 1, d)
		if err != nil {
			return err
		}
		err = item.Upsert(tests.Context, db, &i)
		if err != nil {
			return err
		}

	}

	return nil
}

// Add inserts a set for testing.
func Add(db *db.DB, i *item.Item) error {

	if err := item.Upsert("", db, i); err != nil {
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

	if err := db.ExecuteMGO(tests.Context, item.Collection, f); err != nil {
		return err
	}

	f = func(c *mgo.Collection) error {
		q := bson.M{"name": bson.RegEx{Pattern: pattern}}
		_, err := c.RemoveAll(q)
		return err
	}

	if err := db.ExecuteMGO(tests.Context, item.CollectionHistory, f); err != nil {
		return err
	}

	return nil
}
