package view_test

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/ardanlabs/kit/cfg"
	"github.com/ardanlabs/kit/tests"
	"github.com/coralproject/shelf/internal/platform/db"
	"github.com/coralproject/shelf/internal/wire/view"
	"github.com/coralproject/shelf/internal/wire/view/viewfix"
)

// prefix is what we are looking to delete after the test.
const prefix = "VTEST_"

func TestMain(m *testing.M) {
	os.Exit(runTest(m))
}

// runTest initializes the environment for the tests and allows for
// the proper return code if the test fails or succeeds.
func runTest(m *testing.M) int {

	// Initialize the configuration and logging systems. Plus anything
	// else the web app layer needs.
	tests.Init("XENIA")

	// Initialize MongoDB using the `tests.TestSession` as the name of the
	// master session.
	if err := db.RegMasterSession(tests.Context, tests.TestSession, cfg.MustURL("MONGO_URI").String(), 0); err != nil {
		fmt.Println("Can't register master session: " + err.Error())
		return 1
	}

	return m.Run()
}

//==============================================================================

// setup initializes for each indivdual test.
func setup(t *testing.T) ([]view.View, *db.DB) {
	tests.ResetLog()

	views, err := viewfix.Get()
	if err != nil {
		t.Fatalf("%s\tShould load view records from file : %v", tests.Failed, err)
	}
	t.Logf("%s\tShould load view records from file.", tests.Success)

	db, err := db.NewMGO(tests.Context, tests.TestSession)
	if err != nil {
		t.Fatalf("%s\tShould be able to get a Mongo session : %v", tests.Failed, err)
	}

	return views, db
}

// teardown deinitializes for each indivdual test.
func teardown(t *testing.T, db *db.DB) {
	if err := viewfix.Remove(tests.Context, db, prefix); err != nil {
		t.Fatalf("%s\tShould be able to remove the view records : %v", tests.Failed, err)
	}
	t.Logf("%s\tShould be able to remove the view records.", tests.Success)

	db.CloseMGO(tests.Context)

	tests.DisplayLog()
}

//==============================================================================

// TestUpsertDelete tests if we can add/remove a view to/from the db.
func TestUpsertDelete(t *testing.T) {
	views, db := setup(t)
	defer teardown(t, db)

	t.Log("Given the need to upsert and delete views.")
	{
		t.Log("\tWhen starting from an empty views collection")
		{

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

			_, err = view.GetByName(tests.Context, db, views[0].Name)
			if err == nil {
				t.Fatalf("\t%s\tShould generate an error when getting a view with the deleted name : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould generate an error when getting a view with the deleted name.", tests.Success)
		}
	}
}

// TestGetAll tests if we can get all views from the db.
func TestGetAll(t *testing.T) {
	views1, db := setup(t)
	defer teardown(t, db)

	t.Log("Given the need to get all the views in the database.")
	{
		t.Log("\tWhen starting from an empty views collection")
		{

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

			var filteredViews []view.View
			for _, vw := range views2 {
				if vw.Name[0:len(prefix)] == prefix {
					filteredViews = append(filteredViews, vw)
				}
			}

			if !reflect.DeepEqual(views1, filteredViews) {
				t.Logf("\t%+v", views1)
				t.Logf("\t%+v", filteredViews)
				t.Fatalf("\t%s\tShould be able to get back the same views.", tests.Failed)
			}
			t.Logf("\t%s\tShould be able to get back the same views.", tests.Success)
		}
	}
}
