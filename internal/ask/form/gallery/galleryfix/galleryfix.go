package galleryfix

import (
	"encoding/json"
	"os"

	"github.com/ardanlabs/kit/db"
	"github.com/coralproject/shelf/internal/ask/form/gallery"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var path string

func init() {
	path = os.Getenv("GOPATH") + "/src/github.com/coralproject/shelf/internal/ask/form/gallery/galleryfix/"
}

// Get loads gallery data.
func Get(fixture string) ([]gallery.Gallery, error) {
	file, err := os.Open(path + fixture + ".json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var gs []gallery.Gallery
	err = json.NewDecoder(file).Decode(&gs)
	if err != nil {
		return nil, err
	}

	return gs, nil
}

// Add inserts gallerys for testing.
func Add(context interface{}, db *db.DB, gs []gallery.Gallery) error {
	for i := range gs {
		// The gallery.Create function will add/update fields so we need to pass
		// the correct reference.
		if err := gallery.Create(context, db, &gs[i]); err != nil {
			return err
		}
	}

	return nil
}

// Remove removes gallerys in Mongo that match a given pattern.
func Remove(context interface{}, db *db.DB, pattern string) error {
	f := func(c *mgo.Collection) error {
		q := bson.M{"description": bson.RegEx{Pattern: "^" + pattern}}
		_, err := c.RemoveAll(q)
		return err
	}

	if err := db.ExecuteMGO(context, gallery.Collection, f); err != nil {
		return err
	}

	return nil
}
