// Package tests implements users tests for the API layer.
package tests

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/ardanlabs/kit/cfg"
	"github.com/ardanlabs/kit/tests"
	"github.com/coralproject/shelf/cmd/askd/routes"
	"github.com/coralproject/shelf/internal/ask/form/submission/submissionfix"
	"github.com/coralproject/shelf/internal/platform/db"

	"github.com/coralproject/shelf/tstdata"
)

// a is the global App reference which actually is a http.Handler.
var a http.Handler

//==============================================================================

// TestMain helps to clean up the test data.
func TestMain(m *testing.M) {
	os.Exit(runTest(m))
}

// runTest initializes the environment for the tests and allows for
// the proper return code if the test fails or succeeds.
func runTest(m *testing.M) int {

	// Initialize MongoDB using the `tests.TestSession` as the name of the
	// master session.
	if err := db.RegMasterSession(tests.Context, tests.TestSession, cfg.MustURL("MONGO_URI").String(), 0); err != nil {
		fmt.Println("Can't register master session: " + err.Error())
		return 1
	}

	// Setup the app for performing tests.
	a = routes.API()

	// Snatch the mongo session so we can create some test data.
	db, err := db.NewMGO(tests.Context, tests.TestSession)
	if err != nil {
		fmt.Println("Unable to get Mongo session")
		return 1
	}
	defer db.CloseMGO(tests.Context)

	// Generate the test data.
	tstdata.Generate(db)
	defer tstdata.Drop(db)

	// Load in the submissions from the fixture.
	if err = loadSubmissions(db, "submission.json"); err != nil {
		fmt.Println("Unable to load submissions: ", err)
	}
	defer submissionfix.Remove(tests.Context, db, subPrefix)

	return m.Run()
}

// loadSubmissions adds submissions to run tests.
func loadSubmissions(db *db.DB, file string) error {
	submissions, err := submissionfix.GetMany(file)
	if err != nil {
		return err
	}

	if err := submissionfix.Add(tests.Context, db, submissions); err != nil {
		return err
	}

	return nil
}
