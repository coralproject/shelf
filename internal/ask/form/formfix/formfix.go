package formfix

import (
	"encoding/json"
	"os"

	"github.com/ardanlabs/kit/db"
	"github.com/coralproject/shelf/internal/ask/form"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var path string

func init() {
	path = os.Getenv("GOPATH") + "/src/github.com/coralproject/shelf/internal/ask/form/formfix/"
}

// Get loads form data based on forms.json.
func Get(fixture string) ([]form.Form, error) {
	file, err := os.Open(path + fixture + ".json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var fms []form.Form
	err = json.NewDecoder(file).Decode(&fms)
	if err != nil {
		return nil, err
	}

	return fms, nil
}

// Add inserts forms for testing.
func Add(context interface{}, db *db.DB, fms []form.Form) error {
	for _, fm := range fms {
		if err := form.Upsert(context, db, &fm); err != nil {
			return err
		}
	}

	return nil
}

// Remove removes forms in Mongo that match a given pattern.
func Remove(context interface{}, db *db.DB, pattern string) error {
	f := func(c *mgo.Collection) error {
		q := bson.M{"header.title": bson.RegEx{Pattern: "^" + pattern}}
		_, err := c.RemoveAll(q)
		return err
	}

	if err := db.ExecuteMGO(context, form.Collection, f); err != nil {
		return err
	}

	return nil
}
