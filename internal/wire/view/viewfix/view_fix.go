package viewfix

import (
	"encoding/json"
	"os"

	"github.com/ardanlabs/kit/db"
	"github.com/coralproject/xenia/internal/wire/view"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var path string

func init() {
	path = os.Getenv("GOPATH") + "/src/github.com/coralproject/xenia/internal/wire/view/viewfix/"
}

// Get loads view data based on view.json.
func Get() ([]view.View, error) {
	file, err := os.Open(path + "view.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var views []view.View
	err = json.NewDecoder(file).Decode(&views)
	if err != nil {
		return nil, err
	}

	return views, nil
}

// Add inserts views for testing.
func Add(context interface{}, db *db.DB, views []view.View) error {
	for _, v := range views {
		if err := view.Upsert(context, db, &v); err != nil {
			return err
		}
	}

	return nil
}

// Remove removes views in Mongo that match a given pattern.
func Remove(context interface{}, db *db.DB, pattern string) error {
	f := func(c *mgo.Collection) error {
		q := bson.M{"name": bson.RegEx{Pattern: pattern}}
		_, err := c.RemoveAll(q)
		return err
	}

	if err := db.ExecuteMGO(context, view.Collection, f); err != nil {
		return err
	}

	return nil
}
