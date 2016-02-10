// Package tests implements users tests for the API layer.
package tests

import (
	"fmt"
	"testing"

	"github.com/coralproject/xenia/cmd/xeniad/routes"
	"github.com/coralproject/xenia/pkg/query/qfix"
	"github.com/coralproject/xenia/pkg/script/sfix"
	"github.com/coralproject/xenia/tstdata"

	"github.com/ardanlabs/kit/cfg"
	"github.com/ardanlabs/kit/db"
	"github.com/ardanlabs/kit/web/app"
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

	// In order to get a Mongo session we need the name of the database we
	// are using. The web framework middleware is using this by convention.
	dbName, err := cfg.String("MONGO_DB")
	if err != nil {
		fmt.Println("MongoDB is not configured")
		return
	}

	db, err := db.NewMGO("context", dbName)
	if err != nil {
		fmt.Println("Unable to get Mongo session")
		return
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

	m.Run()
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
