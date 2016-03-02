package mfix

import (
	"encoding/json"
	"os"

	"github.com/coralproject/xenia/pkg/mask"

	"github.com/ardanlabs/kit/db"
	"github.com/ardanlabs/kit/tests"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var path string

func init() {
	path = os.Getenv("GOPATH") + "/src/github.com/coralproject/xenia/pkg/mask/mfix/"
}

//==============================================================================

// Get retrieves a slice of mask documents from the filesystem for testing.
func Get(fileName string) ([]mask.Mask, error) {
	file, err := os.Open(path + fileName)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	var masks []mask.Mask
	err = json.NewDecoder(file).Decode(&masks)
	if err != nil {
		return nil, err
	}

	return masks, nil
}

// Add inserts a mask for testing.
func Add(db *db.DB, msk mask.Mask) error {
	f := func(c *mgo.Collection) error {
		q := bson.M{"collection": msk.Collection, "field": msk.Field}
		_, err := c.Upsert(q, msk)
		return err
	}

	if err := db.ExecuteMGO(tests.Context, mask.Collection, f); err != nil {
		return err
	}

	return nil
}

// Remove is used to clear out all the test masks from the collection.
func Remove(db *db.DB, collection string) error {
	f := func(c *mgo.Collection) error {
		q := bson.M{"$or": []bson.M{bson.M{"collection": collection}, bson.M{"collection": "*", "field": "test"}}}
		_, err := c.RemoveAll(q)
		return err
	}

	if err := db.ExecuteMGO(tests.Context, mask.Collection, f); err != nil {
		return err
	}

	f = func(c *mgo.Collection) error {
		q := bson.M{"$or": []bson.M{bson.M{"collection": collection}, bson.M{"collection": "*", "field": "test"}}}
		_, err := c.RemoveAll(q)
		return err
	}

	if err := db.ExecuteMGO(tests.Context, mask.CollectionHistory, f); err != nil {
		return err
	}

	return nil
}
