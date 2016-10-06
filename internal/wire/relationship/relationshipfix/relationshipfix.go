package relationshipfix

import (
	"encoding/json"
	"os"

	"github.com/coralproject/shelf/internal/platform/db"
	"github.com/coralproject/shelf/internal/wire/relationship"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var path string

func init() {
	path = os.Getenv("GOPATH") + "/src/github.com/coralproject/shelf/internal/wire/relationship/relationshipfix/"
}

// Get loads relationship data based on relationships.json.
func Get() ([]relationship.Relationship, error) {
	file, err := os.Open(path + "relationship.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var rels []relationship.Relationship
	err = json.NewDecoder(file).Decode(&rels)
	if err != nil {
		return nil, err
	}

	return rels, nil
}

// Add inserts relationships for testing.
func Add(context interface{}, db *db.DB, rels []relationship.Relationship) error {
	for _, rel := range rels {
		if err := relationship.Upsert(context, db, &rel); err != nil {
			return err
		}
	}

	return nil
}

// Remove removes relationships in Mongo that match a given pattern.
func Remove(context interface{}, db *db.DB, pattern string) error {
	f := func(c *mgo.Collection) error {
		q := bson.M{"predicate": bson.RegEx{Pattern: pattern}}
		_, err := c.RemoveAll(q)
		return err
	}

	if err := db.ExecuteMGO(context, relationship.Collection, f); err != nil {
		return err
	}

	return nil
}
