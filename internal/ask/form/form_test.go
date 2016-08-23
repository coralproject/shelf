package form_test

import (
	"log"
	"os"
	"testing"

	"github.com/ardanlabs/kit/cfg"
	"github.com/ardanlabs/kit/db"
	"github.com/ardanlabs/kit/db/mongo"
	"github.com/ardanlabs/kit/tests"
	"github.com/coralproject/shelf/internal/ask/form"
	"github.com/coralproject/shelf/internal/ask/form/formfix"
)

var dbSession *db.DB

// prefix is what we are looking to delete after the test.
const prefix = "FTEST_"

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

	db, err := db.NewMGO(tests.Context, tests.TestSession)
	if err != nil {
		log.Fatalf("\t%s\tShould be able to get a Mongo session : %v", tests.Failed, err)
	}
	defer db.CloseMGO(tests.Context)

	// set the package database handle
	dbSession = db

	os.Exit(m.Run())
}

func Test_UpsertDelete(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	defer func() {
		if err := formfix.Remove(tests.Context, dbSession, prefix); err != nil {
			t.Fatalf("\t%s\tShould be able to remove the forms : %v", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to remove the forms.", tests.Success)
	}()

	t.Log("Given the need to upsert and delete forms.")
	{
		t.Log("\tWhen starting from an empty forms collection")
		{
			//----------------------------------------------------------------------
			// Get the fixture.

			fms, err := formfix.Get()
			if err != nil {
				t.Fatalf("\t%s\tShould be able retrieve form fixture : %s", tests.Failed, err)
			}

			//----------------------------------------------------------------------
			// Upsert the form.

			if err := form.Upsert(tests.Context, dbSession, &fms[0]); err != nil {
				t.Fatalf("\t%s\tShould be able to upsert a form : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to upsert a form.", tests.Success)

			//----------------------------------------------------------------------
			// Get the form.

			fm, err := form.Retrieve(tests.Context, dbSession, fms[0].ID.Hex())
			if err != nil {
				t.Fatalf("\t%s\tShould be able to get the form by id : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to get the form by id.", tests.Success)

			//----------------------------------------------------------------------
			// Check that we got the form we expected.

			if fms[0].ID.Hex() != fm.ID.Hex() {
				t.Fatalf("\t%s\tShould be able to get back the same form.", tests.Failed)
			}
			t.Logf("\t%s\tShould be able to get back the same form.", tests.Success)

			//----------------------------------------------------------------------
			// Delete the form.

			if err := form.Delete(tests.Context, dbSession, fms[0].ID.Hex()); err != nil {
				t.Fatalf("\t%s\tShould be able to delete the form : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to delete the form.", tests.Success)

			//----------------------------------------------------------------------
			// Get the form.

			_, err = form.Retrieve(tests.Context, dbSession, fms[0].ID.Hex())
			if err == nil {
				t.Fatalf("\t%s\tShould generate an error when getting a form with the deleted id : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould generate an error when getting a form with the deleted id.", tests.Success)
		}
	}
}
