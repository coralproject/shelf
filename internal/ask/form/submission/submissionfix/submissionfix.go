package submissionfix

import (
	"encoding/json"
	"os"

	"github.com/ardanlabs/kit/db"
	"github.com/coralproject/shelf/internal/ask/form/submission"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var path string

func init() {
	path = os.Getenv("GOPATH") + "/src/github.com/coralproject/shelf/internal/ask/form/submission/submissionfix/"
}

// Get retrieves a submission document from the filesystem for testing.
func GetMany(fileName string) ([]submission.Submission, error) {
	file, err := os.Open(path + fileName)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	var subs []submission.Submission
	err = json.NewDecoder(file).Decode(&subs)
	if err != nil {
		return nil, err
	}

	return subs, nil
}

// Add inserts submissions to the DB for testing.
func Add(context interface{}, db *db.DB, subs []submission.Submission) error {
	for _, sub := range subs {
		if err := submission.Create(context, db, sub.FormID.Hex(), &sub); err != nil {
			return err
		}
	}

	return nil
}

// Remove removes forms in Mongo that match a given pattern.
func Remove(context interface{}, db *db.DB, prefix string) error {
	f := func(c *mgo.Collection) error {
		q := bson.M{"header.title": bson.RegEx{Pattern: "^" + prefix}}
		_, err := c.RemoveAll(q)
		return err
	}

	if err := db.ExecuteMGO(context, submission.Collection, f); err != nil {
		return err
	}

	return nil
}
