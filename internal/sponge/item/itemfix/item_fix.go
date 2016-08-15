package itemfix

import (
	"encoding/json"
	"os"

	"github.com/ardanlabs/kit/db"
	"github.com/coralproject/shelf/internal/sponge/item"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var path string

func init() {
	path = os.Getenv("GOPATH") + "/src/github.com/coralproject/shelf/internal/sponge/item/itemfix/"
}

// Get loads item data based on item.json.
func Get() ([]item.Item, error) {
	file, err := os.Open(path + "item.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var items []item.Item
	err = json.NewDecoder(file).Decode(&items)
	if err != nil {
		return nil, err
	}

	return items, nil
}

// Add inserts items for testing.
func Add(context interface{}, db *db.DB, items []item.Item) error {
	for _, it := range items {
		if err := item.Upsert(context, db, &it); err != nil {
			return err
		}
	}

	return nil
}

// Remove removes items in Mongo that match a given pattern.
func Remove(context interface{}, db *db.DB, pattern string) error {
	f := func(c *mgo.Collection) error {
		q := bson.M{"item_id": bson.RegEx{Pattern: pattern}}
		_, err := c.RemoveAll(q)
		return err
	}

	if err := db.ExecuteMGO(context, item.Collection, f); err != nil {
		return err
	}

	return nil
}
