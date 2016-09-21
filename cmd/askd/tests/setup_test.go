// Package tests implements users tests for the API layer.
package tests

import (
	"fmt"
	"os"
	"testing"

	"github.com/ardanlabs/kit/cfg"
	"github.com/ardanlabs/kit/db"
	"github.com/ardanlabs/kit/tests"
	"github.com/ardanlabs/kit/web/app"
	"github.com/coralproject/shelf/cmd/askd/routes"
	"github.com/coralproject/shelf/internal/ask/form/submission/submissionfix"

	"github.com/coralproject/shelf/tstdata"
)

var a *app.App

func init() {
	// The call to API will force the init() function to initialize
	// cfg, log and mongodb.
	a = routes.API().(*app.App)
}

//==============================================================================

// TestMain helps to clean up the test data.
func TestMain(m *testing.M) {
	os.Exit(runTest(m))
}

// runTest initializes the environment for the tests and allows for
// the proper return code if the test fails or succeeds.
func runTest(m *testing.M) int {

	// In order to get a Mongo session we need the name of the database we
	// are using. The web framework middleware is using this by convention.
	dbName, err := cfg.String("MONGO_DB")
	if err != nil {
		fmt.Println("MongoDB is not configured")
		return 1
	}

	db, err := db.NewMGO("context", dbName)
	if err != nil {
		fmt.Println("Unable to get Mongo session")
		return 1
	}

	defer db.CloseMGO("context")

	tstdata.Generate(db)
	defer tstdata.Drop(db)

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
