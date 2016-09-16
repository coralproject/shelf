package gallery_test

import (
	"os"
	"testing"

	"gopkg.in/mgo.v2/bson"

	"github.com/ardanlabs/kit/cfg"
	"github.com/ardanlabs/kit/db"
	"github.com/ardanlabs/kit/db/mongo"
	"github.com/ardanlabs/kit/tests"
	"github.com/coralproject/shelf/internal/ask/form/gallery"
	"github.com/coralproject/shelf/internal/ask/form/gallery/galleryfix"
	"github.com/coralproject/shelf/internal/ask/form/submission"
	"github.com/coralproject/shelf/internal/ask/form/submission/submissionfix"
)

// prefix is what we are looking to delete after the test.
const prefix = "FGTEST"

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

func setup(t *testing.T, fixture string) ([]gallery.Gallery, *db.DB) {
	tests.ResetLog()

	gs, err := galleryfix.Get(fixture)
	if err != nil {
		t.Fatalf("%s\tShould be able retrieve gallery fixture : %s", tests.Failed, err)
	}
	t.Logf("%s\tShould be able retrieve gallery fixture.", tests.Success)

	db, err := db.NewMGO(tests.Context, tests.TestSession)
	if err != nil {
		t.Fatalf("Should be able to get a Mongo session : %v", err)
	}

	return gs, db
}

func teardown(t *testing.T, db *db.DB) {
	if err := galleryfix.Remove(tests.Context, db, prefix); err != nil {
		t.Fatalf("%s\tShould be able to remove the galleries : %v", tests.Failed, err)
	}
	t.Logf("%s\tShould be able to remove the galleries.", tests.Success)

	db.CloseMGO(tests.Context)
	tests.DisplayLog()
}

func Test_CreateDelete(t *testing.T) {
	gs, db := setup(t, "gallery")
	defer teardown(t, db)

	t.Log("Given the need to create and delete galleries.")
	{
		t.Log("\tWhen starting from an empty galleries collection")
		{
			//----------------------------------------------------------------------
			// Starting with a single gallery.
			g := gs[0]

			//----------------------------------------------------------------------
			// Create the gallery.

			if err := gallery.Create(tests.Context, db, &g); err != nil {
				t.Fatalf("\t%s\tShould be able to upsert a gallery : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to upsert a gallery.", tests.Success)

			//----------------------------------------------------------------------
			// Get the gallery.

			rg, err := gallery.Retrieve(tests.Context, db, g.ID.Hex())
			if err != nil {
				t.Fatalf("\t%s\tShould be able to get the gallery by id : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to get the gallery by id.", tests.Success)

			//----------------------------------------------------------------------
			// Check that we got the gallery we expected.

			if rg.ID.Hex() != g.ID.Hex() {
				t.Fatalf("\t%s\tShould be able to get back the same gallery.", tests.Failed)
			}
			t.Logf("\t%s\tShould be able to get back the same gallery.", tests.Success)

			//----------------------------------------------------------------------
			// Delete the gallery.

			if err := gallery.Delete(tests.Context, db, g.ID.Hex()); err != nil {
				t.Fatalf("\t%s\tShould be able to delete the gallery : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to delete the gallery.", tests.Success)

			//----------------------------------------------------------------------
			// Get the gallery.

			_, err = gallery.Retrieve(tests.Context, db, g.ID.Hex())
			if err == nil {
				t.Fatalf("\t%s\tShould generate an error when getting a gallery with the deleted id : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould generate an error when getting a gallery with the deleted id.", tests.Success)
		}
	}
}

func Test_Answers(t *testing.T) {
	gs, db := setup(t, "gallery")
	defer teardown(t, db)

	t.Log("Given the need to add and remove answers from galleries.")
	{
		t.Log("\tWhen starting from an empty galleries collection but saturated submissions collection")
		{

			//----------------------------------------------------------------------
			// Starting with a single gallery.
			g := gs[0]

			//----------------------------------------------------------------------
			// Create the gallery.

			if err := gallery.Create(tests.Context, db, &g); err != nil {
				t.Fatalf("\t%s\tShould be able to upsert a gallery : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to upsert a gallery.", tests.Success)

			//----------------------------------------------------------------------
			// We need to get/fill in real form submissions.

			subs, err := submissionfix.GetMany("gallery_submissions.json")
			if err != nil {
				t.Fatalf("Should be able to fetch submission fixtures : %v", err)
			}

			// Set the form ID on these new submissions.
			for i := range subs {
				// We need to assign a new submission ID to ensure that we don't collide
				// with existing submissions.
				subs[i].ID = bson.NewObjectId()
				subs[i].FormID = g.FormID
			}

			// Add all the submission fixtures.
			if err := submissionfix.Add(tests.Context, db, subs); err != nil {
				t.Fatalf("Should be able to add submission fixtures : %v", err)
			}

			// Remove the new fixtures after the tests are completed.
			defer func() {
				for _, sub := range subs {
					if err := submission.Delete(tests.Context, db, sub.ID.Hex()); err != nil {
						t.Fatalf("%s\tShould be able to remove submission fixtures : %v", tests.Failed, err)
					}
				}
				t.Logf("%s\tShould be able to remove submission fixtures.", tests.Success)
			}()

			//----------------------------------------------------------------------
			// Add an answer to the gallery.

			crg, err := gallery.AddAnswer(tests.Context, db, g.ID.Hex(), subs[0].ID.Hex(), subs[0].Answers[0].WidgetID)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to add an answer to a gallery : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to add an answer to a gallery.", tests.Success)

			//----------------------------------------------------------------------
			// Check that the returned gallery has the correct answers.

			if len(crg.Answers) != 1 {
				t.Fatalf("\t%s\tShould have at least one answer on the returned gallery : Expected 1, got %d", tests.Failed, len(crg.Answers))
			}
			t.Logf("\t%s\tShould have at least one answer on the returned gallery.", tests.Success)

			matchAnswers(t, crg.Answers[0], subs[0], subs[0].Answers[0])

			//----------------------------------------------------------------------
			// Check that the retrieved gallery has the correct answers.

			rg, err := gallery.Retrieve(tests.Context, db, g.ID.Hex())
			if err != nil {
				t.Fatalf("\t%s\tShould be able to retrieve a gallery : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to retrieve a gallery.", tests.Success)

			if len(rg.Answers) != 1 {
				t.Fatalf("\t%s\tShould have at least one answer on the returned gallery : Expected 1, got %d", tests.Failed, len(rg.Answers))
			}
			t.Logf("\t%s\tShould have at least one answer on the returned gallery.", tests.Success)

			matchAnswers(t, rg.Answers[0], subs[0], subs[0].Answers[0])

			//----------------------------------------------------------------------
			// Remove the answer.

			drg, err := gallery.RemoveAnswer(tests.Context, db, g.ID.Hex(), subs[0].ID.Hex(), subs[0].Answers[0].WidgetID)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to remove an answer from a gallery : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to remove an answer from a gallery.", tests.Success)

			//----------------------------------------------------------------------
			// Check that the returned gallery has the correct answers.

			if len(drg.Answers) != 0 {
				t.Fatalf("\t%s\tShould have at no answers on the returned gallery : Expected 0, got %d", tests.Failed, len(drg.Answers))
			}
			t.Logf("\t%s\tShould have at no answers on the returned gallery.", tests.Success)

			//----------------------------------------------------------------------
			// Check that the retrieved gallery has the correct answers.

			rg, err = gallery.Retrieve(tests.Context, db, g.ID.Hex())
			if err != nil {
				t.Fatalf("\t%s\tShould be able to retrieve a gallery : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to retrieve a gallery.", tests.Success)

			if len(rg.Answers) != 0 {
				t.Fatalf("\t%s\tShould have at no answers on the returned gallery : Expected 0, got %d", tests.Failed, len(drg.Answers))
			}
			t.Logf("\t%s\tShould have at no answers on the returned gallery.", tests.Success)
		}
	}
}

func matchAnswers(t *testing.T, ga gallery.Answer, sub submission.Submission, sa submission.Answer) {
	if ga.SubmissionID.Hex() != sub.ID.Hex() {
		t.Fatalf("\t%s\tShould match the submission ID : Expected %s, got %s", tests.Failed, sub.ID.Hex(), ga.SubmissionID.Hex())
	}
	t.Logf("\t%s\tShould match the submission ID", tests.Success)

	if ga.AnswerID != sa.WidgetID {
		t.Fatalf("\t%s\tShould match the widget ID : Expected %s, got %s", tests.Failed, sa.WidgetID, ga.AnswerID)
	}
	t.Logf("\t%s\tShould match the widget ID", tests.Success)

	if mongo.Query(ga.Answer.Answer) != mongo.Query(sa.Answer) {
		t.Fatalf("\t%s\tShould match the answer : Expected %s, got %s", tests.Failed, mongo.Query(sa.Answer), mongo.Query(ga.Answer.Answer))
	}
	t.Logf("\t%s\tShould match the answer.", tests.Success)
}

func Test_List(t *testing.T) {
	gs, db := setup(t, "gallery_list")
	defer teardown(t, db)

	t.Log("Given the need to list galleries.")
	{
		t.Log("\tWhen starting from an empty galleries collection.")
		{
			lgs, err := gallery.List(tests.Context, db, gs[0].FormID.Hex())
			if err != nil {
				t.Fatalf("\t%s\tShould be able to list no galleries : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to list no galleries.", tests.Success)

			if len(lgs) != 0 {
				t.Fatalf("\t%s\tShould be able to list the correct amount of galleries: Expected 0, found %d", tests.Failed, len(lgs))
			}
			t.Logf("\t%s\tShould be able to list the correct amount of galleries.", tests.Success)

			if err := galleryfix.Add(tests.Context, db, gs); err != nil {
				t.Fatalf("\t%s\tShould be able to load gallery fixtures : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to load gallery fixtures.", tests.Success)

			lgs, err = gallery.List(tests.Context, db, gs[0].FormID.Hex())
			if err != nil {
				t.Fatalf("\t%s\tShould be able to load gallery fixtures : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to load gallery fixtures.", tests.Success)

			if len(lgs) != len(gs) {
				t.Fatalf("\t%s\tShould be able to list the correct amount of galleries: Expected %d, found %d", tests.Failed, len(gs), len(lgs))
			}
			t.Logf("\t%s\tShould be able to list the correct amount of galleries.", tests.Success)

			for _, g := range lgs {
				if g.FormID.Hex() != gs[0].FormID.Hex() {
					t.Fatalf("\t%s\tShould have the correct form id : Expected %s, got %s", tests.Failed, gs[0].FormID.Hex(), g.FormID.Hex())
				}
			}
			t.Logf("\t%s\tShould have the correct form id.", tests.Success)

			matches := 0
			for _, fg := range gs {
				for _, lg := range lgs {
					if lg.ID.Hex() == fg.ID.Hex() {
						matches++
					}
				}
			}

			if matches != len(lgs) {
				t.Fatalf("\t%s\tShould contain all the fixtures in the listed contents : Not all fixtures found", tests.Failed)
			}
			t.Logf("\t%s\tShould contain all the fixtures in the listed contents.", tests.Success)
		}
	}
}

func Test_Update(t *testing.T) {
	gs, db := setup(t, "gallery")
	defer teardown(t, db)

	t.Log("Given the need to list galleries.")
	{
		t.Log("\tWhen starting from an empty galleries collection.")
		{
			//----------------------------------------------------------------------
			// Starting with a single gallery.
			g := gs[0]

			//----------------------------------------------------------------------
			// Create the gallery.

			if err := gallery.Create(tests.Context, db, &g); err != nil {
				t.Fatalf("\t%s\tShould be able to upsert a gallery : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to upsert a gallery.", tests.Success)

			//----------------------------------------------------------------------
			// Update the gallery.

			newHeadline := "my new headline"

			g.Headline = newHeadline

			if err := gallery.Update(tests.Context, db, g.ID.Hex(), &g); err != nil {
				t.Fatalf("\t%s\tShould be able to update the gallery : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to update the gallery.", tests.Success)

			if g.Headline != newHeadline {
				t.Fatalf("\t%s\tShould update the headline on the returned gallery : Expected \"%s\", got \"%s\"", tests.Failed, newHeadline, g.Headline)
			}
			t.Logf("\t%s\tShould update the headline on the returned gallery.", tests.Success)

			rg, err := gallery.Retrieve(tests.Context, db, g.ID.Hex())
			if err != nil {
				t.Fatalf("\t%s\tShould be able to retrieve the gallery : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to retrieve the gallery.", tests.Success)

			if rg.Headline != newHeadline {
				t.Fatalf("\t%s\tShould update the headline on the retrieved gallery : Expected \"%s\", got \"%s\"", tests.Failed, newHeadline, rg.Headline)
			}
			t.Logf("\t%s\tShould update the headline on the retrieved gallery.", tests.Success)
		}
	}
}
