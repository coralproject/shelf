package view_test

import (
	"reflect"
	"testing"

	"github.com/ardanlabs/kit/cfg"
	"github.com/ardanlabs/kit/db"
	"github.com/ardanlabs/kit/db/mongo"
	"github.com/ardanlabs/kit/tests"
	"github.com/coralproject/xenia/internal/shelf/view"
	"github.com/coralproject/xenia/internal/shelf/view/viewfix"
)

func init() {
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
}

// prefix is what we are looking to delete after the test.
const prefix = "VTEST_"

// TestUpsertDelete tests if we can add/remove a view to/from the db.
func TestUpsertDelete(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	db, err := db.NewMGO(tests.Context, tests.TestSession)
	if err != nil {
		t.Fatalf("\t%s\tShould be able to get a Mongo session : %v", tests.Failed, err)
	}
	defer db.CloseMGO(tests.Context)

	defer func() {
		if err := viewfix.Remove(tests.Context, db, prefix); err != nil {
			t.Fatalf("\t%s\tShould be able to remove the views : %v", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to remove the views.", tests.Success)
	}()

	t.Log("Given the need to upsert and delete views.")
	{
		t.Log("\tWhen starting from an empty views collection")
		{

			//----------------------------------------------------------------------
			// Get the fixture.

			views, err := viewfix.Get()
			if err != nil {
				t.Fatalf("\t%s\tShould be able retrieve view fixture : %s", tests.Failed, err)
			}

			//----------------------------------------------------------------------
			// Upsert the view.

			if err := view.Upsert(tests.Context, db, &views[0]); err != nil {
				t.Fatalf("\t%s\tShould be able to upsert a view : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to upsert a view.", tests.Success)

			//----------------------------------------------------------------------
			// Get the view.

			v, err := view.GetByName(tests.Context, db, views[0].Name)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to get the view by name : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to get the view by name.", tests.Success)

			//----------------------------------------------------------------------
			// Check that we got the expected view.

			if !reflect.DeepEqual(views[0], *v) {
				t.Logf("\t%+v", views[0])
				t.Logf("\t%+v", v)
				t.Fatalf("\t%s\tShould be able to get back the same view.", tests.Failed)
			}
			t.Logf("\t%s\tShould be able to get back the same view.", tests.Success)

			//----------------------------------------------------------------------
			// Delete the view.

			if err := view.Delete(tests.Context, db, views[0].Name); err != nil {
				t.Fatalf("\t%s\tShould be able to delete the view : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to delete the view.", tests.Success)

			//----------------------------------------------------------------------
			// Get the view.

			v, err = view.GetByName(tests.Context, db, views[0].Name)
			if err == nil {
				t.Fatalf("\t%s\tShould generate an error when getting a view with the deleted name : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould generate an error when getting a view with the deleted name.", tests.Success)
		}
	}
}

// TestGetAll tests if we can get all views from the db.
func TestGetAll(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	db, err := db.NewMGO(tests.Context, tests.TestSession)
	if err != nil {
		t.Fatalf("\t%s\tShould be able to get a Mongo session : %v", tests.Failed, err)
	}
	defer db.CloseMGO(tests.Context)

	defer func() {
		if err := viewfix.Remove(tests.Context, db, prefix); err != nil {
			t.Fatalf("\t%s\tShould be able to remove the views : %v", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to remove the views.", tests.Success)
	}()

	t.Log("Given the need to get all the views in the database.")
	{
		t.Log("\tWhen starting from an empty views collection")
		{
			views1, err := viewfix.Get()
			if err != nil {
				t.Fatalf("\t%s\tShould be able retrieve view fixture : %s", tests.Failed, err)
			}

			for _, v := range views1 {
				if err := view.Upsert(tests.Context, db, &v); err != nil {
					t.Fatalf("\t%s\tShould be able to upsert views : %s", tests.Failed, err)
				}
			}
			t.Logf("\t%s\tShould be able to upsert views.", tests.Success)

			views2, err := view.GetAll(tests.Context, db)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to get all views : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to get all views.", tests.Success)

			if !reflect.DeepEqual(views1, views2) {
				t.Logf("\t%+v", views1)
				t.Logf("\t%+v", views2)
				t.Fatalf("\t%s\tShould be able to get back the same views.", tests.Failed)
			}
			t.Logf("\t%s\tShould be able to get back the same views.", tests.Success)
		}
	}
}
