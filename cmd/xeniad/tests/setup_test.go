// Package tests implements users tests for the API layer.
package tests

import (
	"fmt"
	"os"
	"testing"

	"github.com/ardanlabs/kit/cfg"
	"github.com/ardanlabs/kit/db"
	"github.com/ardanlabs/kit/web/app"
	"github.com/coralproject/shelf/cmd/xeniad/routes"
	"github.com/coralproject/shelf/internal/wire/pattern/patternfix"
	"github.com/coralproject/shelf/internal/wire/relationship/relationshipfix"
	"github.com/coralproject/shelf/internal/wire/view/viewfix"
	"github.com/coralproject/shelf/internal/xenia/mask/mfix"
	"github.com/coralproject/shelf/internal/xenia/query/qfix"
	"github.com/coralproject/shelf/internal/xenia/script/sfix"
	"github.com/coralproject/shelf/tstdata"
)

var a *app.App

func init() {
	// The call to API will force the init() function to initialize
	// cfg, log and mongodb.
	a = routes.API(true).(*app.App)
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

	loadQuery(db, "basic.json")
	loadQuery(db, "basic_var.json")
	defer qfix.Remove(db, "QTEST_O")

	loadScript(db, "basic_script_pre.json")
	loadScript(db, "basic_script_pst.json")
	defer sfix.Remove(db, "STEST_O")

	loadMasks(db, "basic.json")
	defer mfix.Remove(db, "test_xenia_data")

	loadRelationships("context", db)
	defer relationshipfix.Remove("context", db, relPrefix)

	loadViews("context", db)
	defer viewfix.Remove("context", db, viewPrefix)

	loadPatterns("context", db)
	defer patternfix.Remove("context", db, patternPrefix)

	return m.Run()
}

// loadQuery adds queries to run tests.
func loadQuery(db *db.DB, file string) error {
	set, err := qfix.Get(file)
	if err != nil {
		return err
	}

	if err := qfix.Add(db, set); err != nil {
		return err
	}

	return nil
}

// loadScript adds scripts to run tests.
func loadScript(db *db.DB, file string) error {
	scr, err := sfix.Get(file)
	if err != nil {
		return err
	}

	if err := sfix.Add(db, scr); err != nil {
		return err
	}

	return nil
}

// loadMasks adds masks to run tests.
func loadMasks(db *db.DB, file string) error {
	masks, err := mfix.Get(file)
	if err != nil {
		return err
	}

	for _, msk := range masks {
		if err := mfix.Add(db, msk); err != nil {
			return err
		}
	}

	return nil
}

// loadRelationships adds relationships to run tests.
func loadRelationships(context interface{}, db *db.DB) error {
	rels, err := relationshipfix.Get()
	if err != nil {
		return err
	}

	if err := relationshipfix.Add(context, db, rels[0:2]); err != nil {
		return err
	}

	return nil
}

// loadViews adds views to run tests.
func loadViews(context interface{}, db *db.DB) error {
	views, err := viewfix.Get()
	if err != nil {
		return err
	}

	if err := viewfix.Add(context, db, views[0:2]); err != nil {
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
