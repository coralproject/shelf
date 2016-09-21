package sfix

import (
	"encoding/json"
	"os"

	"github.com/ardanlabs/kit/db"
	"github.com/coralproject/shelf/internal/sponge"
	"github.com/coralproject/shelf/internal/sponge/item"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var path string

func init() {
	path = os.Getenv("GOPATH") + "/src/github.com/coralproject/shelf/internal/sponge/sfix/"
}

// Get loads data based from data.json.
func Get(filename string) (sponge.Data, error) {
	file, err := os.Open(path + filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var d sponge.Data
	err = json.NewDecoder(file).Decode(&d)
	if err != nil {
		return nil, err
	}

	return d, nil
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
