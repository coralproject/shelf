// Package tests implements users tests for the API layer.
package tests

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/ardanlabs/kit/cfg"
	"github.com/ardanlabs/kit/tests"
	"github.com/coralproject/shelf/cmd/sponged/routes"
	"github.com/coralproject/shelf/internal/platform/db"
	"github.com/coralproject/shelf/internal/sponge/item/itemfix"
	"github.com/coralproject/shelf/internal/wire/pattern/patternfix"
)

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

	if err := loadItems(tests.Context, db); err != nil {
		fmt.Println("Could not load items")
		return 1
	}
	defer itemfix.Remove(tests.Context, db, itemPrefix)

	if err := loadPatterns(tests.Context, db); err != nil {
		fmt.Println("Could not load patterns")
		return 1
	}
	defer patternfix.Remove(tests.Context, db, patternPrefix)

	return m.Run()
}

// loadItems adds items to run tests.
func loadItems(context interface{}, db *db.DB) error {
	items, err := itemfix.Get()
	if err != nil {
		return err
	}

	if err := itemfix.Add(context, db, items[1:]); err != nil {
		return err
	}

	return nil
}

// loadPatterns adds patterns to run tests.
func loadPatterns(context interface{}, db *db.DB) error {
	ps, _, err := patternfix.Get()
	if err != nil {
		return err
	}

	if err := patternfix.Add(context, db, ps[0:2]); err != nil {
		return err
	}

	return nil
}
