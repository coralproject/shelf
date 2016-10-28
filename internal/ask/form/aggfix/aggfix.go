package aggfix

import (
	"encoding/json"
	"os"

	"github.com/coralproject/shelf/internal/ask/form"
	"github.com/coralproject/shelf/internal/ask/form/submission"
	"github.com/coralproject/shelf/internal/platform/db"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var path string

func init() {
	path = os.Getenv("GOPATH") + "/src/github.com/coralproject/shelf/internal/ask/form/aggfix/"
}

// Get loads form and submission data based on forms.json.
func Get() (*form.Form, []submission.Submission, error) {
	file, err := os.Open(path + "forms.json")
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	var fm form.Form
	err = json.NewDecoder(file).Decode(&fm)
	if err != nil {
		return nil, nil, err
	}

	file, err = os.Open(path + "form_submissions.json")
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	var subs []submission.Submission
	err = json.NewDecoder(file).Decode(&subs)
	if err != nil {
		return nil, nil, err
	}

	return &fm, subs, nil
}

// Add inserts forms for testing.
func Add(context interface{}, db *db.DB, fm *form.Form, subs []submission.Submission) error {
	if err := form.Upsert(context, db, fm); err != nil {
		return err
	}

	for _, sub := range subs {
		if err := submission.Create(context, db, fm.ID.Hex(), &sub); err != nil {
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

	if err := db.ExecuteMGO(context, submission.Collection, f); err != nil {
		return err
	}

	return nil
}
