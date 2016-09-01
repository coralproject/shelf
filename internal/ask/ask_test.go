package ask_test

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/ardanlabs/kit/cfg"
	"github.com/ardanlabs/kit/db"
	"github.com/ardanlabs/kit/db/mongo"
	"github.com/ardanlabs/kit/tests"
	"github.com/coralproject/shelf/internal/ask"
	"github.com/coralproject/shelf/internal/ask/form"
	"github.com/coralproject/shelf/internal/ask/form/formfix"
	"github.com/coralproject/shelf/internal/ask/form/gallery"
	"github.com/coralproject/shelf/internal/ask/form/gallery/galleryfix"
	"github.com/coralproject/shelf/internal/ask/form/submission"
	"github.com/coralproject/shelf/internal/ask/form/submission/submissionfix"
)

// prefix is what we are looking to delete after the test.
const prefix = "ASKTEST"

func TestMain(m *testing.M) {
	// Initialize the configuration and logging systems. Plus anything
	// else the web app layer needs.
	tests.Init("XENIA")

	// Initialize MongoDB using the `tests.TestSession` as the name of the
	// master session.
	cfg := mongo.Config{
		Host:     cfg.MustString("MONGO_HOST"),
		AuthDB:   cfg.MustString("MONGO_AUTHDB"),
		DB:       cfg.MustString("MONGO_DB"),
		User:     cfg.MustString("MONGO_USER"),
		Password: cfg.MustString("MONGO_PASS"),
	}
	tests.InitMongo(cfg)

	os.Exit(m.Run())
}

func setup(t *testing.T) *db.DB {
	tests.ResetLog()

	db, err := db.NewMGO(tests.Context, tests.TestSession)
	if err != nil {
		t.Fatalf("Should be able to get a Mongo session : %v", err)
	}

	return db
}

func teardown(t *testing.T, db *db.DB) {
	if err := formfix.Remove(tests.Context, db, prefix); err != nil {
		t.Fatalf("%s\tShould be able to remove the form fixtures from the database : %v", tests.Failed, err)
	}
	t.Logf("%s\tShould be able to remove the form fixtures from the database.", tests.Success)

	if err := galleryfix.Remove(tests.Context, db, prefix); err != nil {
		t.Fatalf("%s\tShould be able to remove the gallery fixtures from the database : %v", tests.Failed, err)
	}
	t.Logf("%s\tShould be able to remove the gallery fixtures from the database.", tests.Success)

	if err := submissionfix.Remove(tests.Context, db, prefix); err != nil {
		t.Fatalf("%s\tShould be able to remove the submission fixtures from the database : %v", tests.Failed, err)
	}
	t.Logf("%s\tShould be able to remove the submission fixtures from the database.", tests.Success)

	db.CloseMGO(tests.Context)
	tests.DisplayLog()
}

func Test_UpsertForm(t *testing.T) {
	db := setup(t)
	defer teardown(t, db)

	t.Log("Given the need to upsert a form.")
	{

		//----------------------------------------------------------------------
		// Get the form fixture.

		fms, err := formfix.Get("ask_form")
		if err != nil {
			t.Fatalf("%s\tShould be able to get the form fixture : %v", tests.Failed, err)
		}
		t.Logf("%s\tShould be able to get the form fixture", tests.Success)

		//----------------------------------------------------------------------
		// Select a specific form.

		fm := fms[0]

		//----------------------------------------------------------------------
		// Update it's ID to a new one to ensure we aren't updating.

		fm.ID = bson.ObjectId("")

		t.Log("\tWhen starting from an empty forms collection")
		{
			//----------------------------------------------------------------------
			// Upsert the form.
			if err := ask.UpsertForm(tests.Context, db, &fm); err != nil {
				t.Fatalf("\t%s\tShould be able to upsert the form : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to upsert the form", tests.Success)

			if fm.ID.Hex() == "" {
				t.Fatalf("\t%s\tShould be able to update the ID when upserted as a new record : ID was not updated", tests.Failed)
			}
			t.Logf("\t%s\tShould be able to update the ID when upserted as a new record", tests.Success)

			//----------------------------------------------------------------------
			// Retrieve the form to ensure it was created.

			if _, err := form.Retrieve(tests.Context, db, fm.ID.Hex()); err != nil {
				t.Fatalf("\t%s\tShould be able to retrieve the form : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to retrieve the form", tests.Success)

			//----------------------------------------------------------------------
			// Retrieve the gallery to ensure it was created.

			gs, err := gallery.List(tests.Context, db, fm.ID.Hex())
			if err != nil {
				t.Fatalf("\t%s\tShould be able to retrieve the gallery : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to retrieve the gallery", tests.Success)

			//----------------------------------------------------------------------
			// Cleanup the galleries created.

			defer func(gs []gallery.Gallery) {
				for _, g := range gs {
					if err := gallery.Delete(tests.Context, db, g.ID.Hex()); err != nil {
						t.Fatalf("\t%s\tShould be able to remove the created galleries : %v", tests.Failed, err)
					}
				}
				t.Logf("\t%s\tShould be able to remove the created galleries.", tests.Success)
			}(gs)

		}

		t.Log("\tWhen starting from an non-empty forms collection")
		{
			//----------------------------------------------------------------------
			// Update the form.

			newFooter := bson.M{"key": "value"}

			fm.Footer = newFooter

			//----------------------------------------------------------------------
			// Upsert the form.
			if err := ask.UpsertForm(tests.Context, db, &fm); err != nil {
				t.Fatalf("\t%s\tShould be able to upsert the form : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to upsert the form", tests.Success)

			if fm.ID.Hex() == "" {
				t.Fatalf("\t%s\tShould be able to update the ID when upserted as a new record : ID was not updated", tests.Failed)
			}
			t.Logf("\t%s\tShould be able to update the ID when upserted as a new record", tests.Success)

			//----------------------------------------------------------------------
			// Retrieve the form to ensure it was created.

			rf, err := form.Retrieve(tests.Context, db, fm.ID.Hex())
			if err != nil {
				t.Fatalf("\t%s\tShould be able to retrieve the form : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to retrieve the form", tests.Success)

			rFooter, ok := rf.Footer.(bson.M)
			if !ok {
				t.Fatalf("\t%s\tShould have a bson document in the footer value : Does not", tests.Failed)
			}
			t.Logf("\t%s\tShould have a bson document in the footer value", tests.Success)

			value, ok := rFooter["key"]
			if !ok {
				t.Fatalf("\t%s\tShould have a bson key in the footer value : Does not", tests.Failed)
			}
			t.Logf("\t%s\tShould have a bson key in the footer value", tests.Success)

			if value != "value" {
				t.Fatalf("\t%s\tShould have expected value : Expected \"%s\", got \"%v\"", tests.Failed, "value", value)
			}
			t.Logf("\t%s\tShould have expected value", tests.Success)
		}
	}
}

func Test_CreateDeleteSubmission(t *testing.T) {
	db := setup(t)
	defer teardown(t, db)

	// CreateSubmission(context interface{}, db *db.DB, formID string, answers []submission.AnswerInput) (*submission.Submission, error)

	t.Log("Given the need to add a submission.")
	{

		//----------------------------------------------------------------------
		// Get the form fixture.

		fms, err := formfix.Get("ask_form")
		if err != nil {
			t.Fatalf("%s\tShould be able to get the form fixture : %v", tests.Failed, err)
		}
		t.Logf("%s\tShould be able to get the form fixture", tests.Success)

		if err := formfix.Add(tests.Context, db, fms); err != nil {
			t.Fatalf("%s\tShould be able to add the form fixture : %v", tests.Failed, err)
		}
		t.Logf("%s\tShould be able to add the form fixture", tests.Success)

		fm := fms[0]

		t.Log("\tWhen starting from an empty submission collection")
		{

			var answers []submission.AnswerInput

			// Create the answers based on the form layout.

			answer := time.Now().Unix()

			for _, step := range fm.Steps {
				for _, widget := range step.Widgets {
					answers = append(answers, submission.AnswerInput{
						WidgetID: widget.ID,
						Answer:   answer,
					})
				}
			}

			// Create the submission.

			sub, err := ask.CreateSubmission(tests.Context, db, fm.ID.Hex(), answers)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to create a submission : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to create a submission.", tests.Success)

			// Ensure that the answers match.

			matchSubmissionsAndAnswers(t, sub, fm, answers)

			// Get the submission from the database.

			rsub, err := submission.Retrieve(tests.Context, db, sub.ID.Hex())
			if err != nil {
				t.Fatalf("\t%s\tShould be able to retrieve a created submission : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to retrieve a created submission.", tests.Success)

			// Ensure that their answers match.

			matchSubmissionsAndAnswers(t, rsub, fm, answers)

			// Ensure that the form's stats were updated.

			rfm, err := form.Retrieve(tests.Context, db, fm.ID.Hex())
			if err != nil {
				t.Fatalf("\t%s\tShould be able to retrieve a form : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to retrieve a form.", tests.Success)

			if rfm.Stats.Responses != 1 {
				t.Fatalf("\t%s\tShould be able to update the stats on a form : Expected %d, got %d", tests.Failed, 1, rfm.Stats.Responses)
			}
			t.Logf("\t%s\tShould be able to update the stats on a form", tests.Success)

			// Delete the submission.

			if err := ask.DeleteSubmission(tests.Context, db, sub.ID.Hex(), fm.ID.Hex()); err != nil {
				t.Fatalf("\t%s\tShould be able to delete a submission : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to delete a submission.", tests.Success)

			// Ensure that it is deleted.

			if _, err := submission.Retrieve(tests.Context, db, sub.ID.Hex()); err == nil {
				t.Fatalf("\t%s\tShould return not found when trying to retrieve a deleted submission : No error", tests.Failed)
			} else if err != mgo.ErrNotFound {
				t.Fatalf("\t%s\tShould return not found when trying to retrieve a deleted submission : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould return not found when trying to retrieve a deleted submission.", tests.Success)

			// Ensure that the form's stats were updated.

			rfm, err = form.Retrieve(tests.Context, db, fm.ID.Hex())
			if err != nil {
				t.Fatalf("\t%s\tShould be able to retrieve a form : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to retrieve a form.", tests.Success)

			if rfm.Stats.Responses != 0 {
				t.Fatalf("\t%s\tShould be able to update the stats on a form : Expected %d, got %d", tests.Failed, 0, rfm.Stats.Responses)
			}
			t.Logf("\t%s\tShould be able to update the stats on a form", tests.Success)
		}
	}
}

func matchSubmissionsAndAnswers(t *testing.T, sub *submission.Submission, fm form.Form, answers []submission.AnswerInput) {
	// Match that the questions matched.
	for _, subAnswer := range sub.Answers {
		var found bool

		for _, step := range fm.Steps {
			for _, widget := range step.Widgets {
				if widget.ID != subAnswer.WidgetID {
					continue
				}

				found = true

				if !jsonEqual(subAnswer.Question, widget.Title) {
					t.Fatalf("\t%s\tShould have matching answers : Expected \"%v\", got \"%v\"", tests.Failed, widget.Title, subAnswer.Question)
				}

				if !jsonEqual(subAnswer.Identity, widget.Identity) {
					t.Fatalf("\t%s\tShould have matching answers : Expected \"%v\", got \"%v\"", tests.Failed, widget.Identity, subAnswer.Identity)
				}

				if !jsonEqual(subAnswer.Props, widget.Props) {
					t.Fatalf("\t%s\tShould have matching answers : Expected \"%v\", got \"%v\"", tests.Failed, widget.Props, subAnswer.Props)
				}

				for _, answer := range answers {
					if answer.WidgetID != widget.ID {
						continue
					}

					if !jsonEqual(answer.Answer, subAnswer.Answer) {
						t.Fatalf("\t%s\tShould have matching answers : Expected \"%v\", got \"%v\"", tests.Failed, answer.Answer, subAnswer.Answer)
					}
				}

				// If we got to this point, it means that we did match on the submission
				// and the widget.
				break
			}

			if found {

				// If we got to this point, it means that we did find a match on the
				// submission on the widget and it's not necessary to match on a
				// different step.
				break
			}
		}

		if !found {
			t.Fatalf("\t%s\tShould have matching answers : Could not find the widget which is related to answer %s", tests.Failed, subAnswer.WidgetID)
		}
	}
	t.Logf("\t%s\tShould have matching answers", tests.Success)
}

func jsonEqual(a, b interface{}) bool {
	ab, err := json.Marshal(a)
	if err != nil {
		return false
	}

	bb, err := json.Marshal(b)
	if err != nil {
		return false
	}

	return strings.Compare(string(ab), string(bb)) == 0
}
