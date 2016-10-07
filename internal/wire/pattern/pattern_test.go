package pattern_test

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/ardanlabs/kit/cfg"
	"github.com/ardanlabs/kit/tests"
	"github.com/coralproject/shelf/internal/platform/db"
	"github.com/coralproject/shelf/internal/wire/pattern"
	"github.com/coralproject/shelf/internal/wire/pattern/patternfix"
)

// prefix is what we are looking to delete after the test.
const prefix = "PTEST_"

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
func setup(t *testing.T) ([]pattern.Pattern, *db.DB) {
	tests.ResetLog()

	patterns, _, err := patternfix.Get()
	if err != nil {
		t.Fatalf("%s\tShould load pattern records from the fixture file : %v", tests.Failed, err)
	}
	t.Logf("%s\tShould load pattern records from the fixture file.", tests.Success)

	db, err := db.NewMGO(tests.Context, tests.TestSession)
	if err != nil {
		t.Fatalf("%s\tShould be able to get a Mongo session : %v", tests.Failed, err)
	}

	return patterns, db
}

// teardown deinitializes for each indivdual test.
func teardown(t *testing.T, db *db.DB) {
	if err := patternfix.Remove(tests.Context, db, prefix); err != nil {
		t.Fatalf("%s\tShould be able to remove the pattern records : %v", tests.Failed, err)
	}
	t.Logf("%s\tShould be able to remove the pattern records.", tests.Success)

	db.CloseMGO(tests.Context)

	tests.DisplayLog()
}

//==============================================================================

// TestUpsertDelete tests if we can add/remove a pattern to/from the db.
func TestUpsertDelete(t *testing.T) {
	patterns, db := setup(t)
	defer teardown(t, db)

	t.Log("Given the need to upsert and delete patterns.")
	{
		t.Log("\tWhen starting from an empty patterns collection")
		{

			//----------------------------------------------------------------------
			// Upsert the pattern.

			if err := pattern.Upsert(tests.Context, db, &patterns[0]); err != nil {
				t.Fatalf("\t%s\tShould be able to upsert a pattern : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to upsert a pattern.", tests.Success)

			//----------------------------------------------------------------------
			// Get the pattern.

			pat, err := pattern.GetByType(tests.Context, db, patterns[0].Type)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to get the pattern by type : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to get the pattern by type.", tests.Success)

			//----------------------------------------------------------------------
			// Check that we got the relationship we expected.

			if !reflect.DeepEqual(patterns[0], *pat) {
				t.Logf("\t%+v", patterns[0])
				t.Logf("\t%+v", pat)
				t.Fatalf("\t%s\tShould be able to get back the same pattern.", tests.Failed)
			}
			t.Logf("\t%s\tShould be able to get back the same pattern.", tests.Success)

			//----------------------------------------------------------------------
			// Delete the pattern.

			if err := pattern.Delete(tests.Context, db, patterns[0].Type); err != nil {
				t.Fatalf("\t%s\tShould be able to delete the pattern : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to delete the pattern.", tests.Success)

			//----------------------------------------------------------------------
			// Get the pattern.

			_, err = pattern.GetByType(tests.Context, db, patterns[0].Type)
			if err == nil {
				t.Fatalf("\t%s\tShould generate an error when getting a pattern with the deleted type : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould generate an error when getting a pattern with the deleted type.", tests.Success)
		}
	}
}

// TestGetAll tests if we can get all patterns from the db.
func TestGetAll(t *testing.T) {
	patterns, db := setup(t)
	defer teardown(t, db)

	t.Log("Given the need to get all the patterns in the database.")
	{
		t.Log("\tWhen starting from an empty patterns collection")
		{

			for _, pat := range patterns {
				if err := pattern.Upsert(tests.Context, db, &pat); err != nil {
					t.Fatalf("\t%s\tShould be able to upsert patterns : %s", tests.Failed, err)
				}
			}
			t.Logf("\t%s\tShould be able to upsert patterns.", tests.Success)

			patternsBack, err := pattern.GetAll(tests.Context, db)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to get all patterns : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to get all patterns.", tests.Success)

			var filteredPats []pattern.Pattern
			for _, pat := range patternsBack {
				if pat.Type[0:len(prefix)] == prefix {
					filteredPats = append(filteredPats, pat)
				}
			}

			if !reflect.DeepEqual(patterns, filteredPats) {
				t.Logf("\t%+v", patterns)
				t.Logf("\t%+v", filteredPats)
				t.Fatalf("\t%s\tShould be able to get back the same patterns.", tests.Failed)
			}
			t.Logf("\t%s\tShould be able to get back the same patterns.", tests.Success)
		}
	}
}
