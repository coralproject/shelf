package xenia_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/ardanlabs/kit/cfg"
	"github.com/ardanlabs/kit/db"
	"github.com/ardanlabs/kit/db/mongo"
	"github.com/ardanlabs/kit/tests"
	"github.com/coralproject/shelf/internal/xenia"
	"github.com/coralproject/shelf/internal/xenia/mask/mfix"
	"github.com/coralproject/shelf/internal/xenia/query"
	"github.com/coralproject/shelf/internal/xenia/script"
	"github.com/coralproject/shelf/internal/xenia/script/sfix"
	"github.com/coralproject/shelf/tstdata"
	"gopkg.in/mgo.v2/bson"
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

//==============================================================================

// TestExecuteSet tests the execution of different Sets that should succeed.
func TestExecuteSet(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	// Build our table of the different test sets.
	execSets := []struct {
		typ string
		set []execSet
	}{
		{typ: "Positive", set: getPosExecSet()},
		{typ: "Negative", set: getNegExecSet()},
	}

	db, err := db.NewMGO(tests.Context, tests.TestSession)
	if err != nil {
		t.Fatalf("\t%s\tShould be able to get a Mongo session : %v", tests.Failed, err)
	}
	defer db.CloseMGO(tests.Context)

	t.Log("Given the need to load the test data.")
	{
		loadTestData(t, db)
	}

	defer func() {
		t.Log("Given the need to unload the test data.")
		{
			unloadTestData(t, db)
		}
	}()

	// Iterate over all the different test sets.
	for _, execSet := range execSets {

		t.Logf("Given the need to execute %s mongo tests.", execSet.typ)
		{
			for _, es := range execSet.set {

				// Setup a sub-test for each item.
				tf := func(t *testing.T) {
					t.Logf("\tWhen using Execute Set %s", es.set.Name)
					{
						result := xenia.Exec(tests.Context, db, es.set, es.vars)

						data, err := json.Marshal(result)
						if err != nil {
							t.Errorf("\t%s\tShould be able to marshal the result : %s", tests.Failed, err)
							return
						}
						t.Logf("\t%s\tShould be able to marshal the result.", tests.Success)

						var res query.Result
						if err := json.Unmarshal(data, &res); err != nil {
							t.Errorf("\t%s\tShould be able to unmarshal the result : %s", tests.Failed, err)
							return
						}
						t.Logf("\t%s\tShould be able to unmarshal the result.", tests.Success)

						// This support allowing the test to provide multiple documents
						// to check when data value order can be underterminstic.
						var found bool
						for _, rslt := range es.results {

							// We just need to find the string inside the result.
							if strings.HasPrefix(rslt, "#find:") {
								if strings.Contains(string(data), rslt[6:]) {
									found = true
									break
								}
								continue
							}

							// Compare the entire result.
							if string(data) == rslt {
								found = true
								break
							}
						}

						if !found {
							t.Log("Exp:", string(data))
							for _, rslt := range es.results {
								t.Log("Rsl:", rslt)
							}
							t.Errorf("\t%s\tShould have the correct result.", tests.Failed)
							return
						}
						t.Logf("\t%s\tShould have the correct result", tests.Success)
					}
				}

				t.Run(es.set.Name, tf)
			}
		}
	}
}

//==============================================================================

// loadTestData adds all the test data into the database.
func loadTestData(t *testing.T, db *db.DB) {
	t.Log("\tWhen loading data for the tests")
	{
		err := tstdata.Generate(db)
		if err != nil {
			t.Fatalf("\t%s\tShould be able to load system with test data : %v", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to load system with test data.", tests.Success)

		scripts := []string{
			"basic_script_pre.json",
			"basic_script_pst.json",
		}

		for _, file := range scripts {
			scr, err := sfix.Get(file)
			if err != nil {
				t.Fatalf("\t%s\tShould load script document from file : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould load script document from file.", tests.Success)

			// We need these scripts loaded under another name to allow tests
			// to run in parallel.
			scr.Name = strings.Replace(scr.Name, "STEST_O", "STEST_T", 1)

			if err := script.Upsert(tests.Context, db, scr); err != nil {
				t.Fatalf("\t%s\tShould be able to create a script : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to create a script.", tests.Success)
		}

		masks, err := mfix.Get("basic.json")
		if err != nil {
			t.Fatalf("\t%s\tShould load mask documents from file : %v", tests.Failed, err)
		}
		t.Logf("\t%s\tShould load mask documents from file.", tests.Success)

		for _, msk := range masks {
			if err := mfix.Add(db, msk); err != nil {
				t.Fatalf("\t%s\tShould be able to create a mask : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to create a mask.", tests.Success)
		}
	}
}

// unloadTestData removes all the test data from the database.
func unloadTestData(t *testing.T, db *db.DB) {
	t.Log("\tWhen unloading data for the tests")
	{
		tstdata.Drop(db)

		if err := sfix.Remove(db, "STEST_T"); err != nil {
			t.Fatalf("\t%s\tShould be able to remove the scripts : %v", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to remove the scripts.", tests.Success)

		if err := mfix.Remove(db, "test_xenia_data"); err != nil {
			t.Fatalf("\t%s\tShould be able to remove the masks : %v", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to remove the masks.", tests.Success)
	}
}

//==============================================================================

// execSet represents the table for the table test of execution tests.
type execSet struct {
	fail    bool
	set     *query.Set
	vars    map[string]string
	results []string
}

// docs represents what a user will receive after
// excuting a successful set.
type docs struct {
	Name string
	Docs []bson.M
}
