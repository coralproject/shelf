package xenia_test

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/ardanlabs/kit/cfg"
	"github.com/ardanlabs/kit/tests"
	"github.com/coralproject/shelf/internal/platform/db"
	"github.com/coralproject/shelf/internal/sponge"
	"github.com/coralproject/shelf/internal/sponge/item/itemfix"
	"github.com/coralproject/shelf/internal/wire/pattern/patternfix"
	"github.com/coralproject/shelf/internal/wire/relationship/relationshipfix"
	"github.com/coralproject/shelf/internal/wire/view/viewfix"
	"github.com/coralproject/shelf/internal/xenia"
	"github.com/coralproject/shelf/internal/xenia/mask/mfix"
	"github.com/coralproject/shelf/internal/xenia/query"
	"github.com/coralproject/shelf/internal/xenia/regex/rfix"
	"github.com/coralproject/shelf/internal/xenia/script"
	"github.com/coralproject/shelf/internal/xenia/script/sfix"
	"github.com/coralproject/shelf/tstdata"
	"gopkg.in/mgo.v2/bson"
)

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

// setup initializes for each indivdual test.
func setup(t *testing.T) *db.DB {
	tests.ResetLog()

	db, err := db.NewMGO(tests.Context, tests.TestSession)
	if err != nil {
		t.Fatalf("%s\tShould be able to get a Mongo session : %v", tests.Failed, err)
	}

	if err := db.NewCayley(tests.Context, tests.TestSession); err != nil {
		t.Fatalf("%s\tShould be able to get a Cayley connection : %v", tests.Failed, err)
	}

	loadTestData(t, db)

	if err := loadRegex(db, "number.json"); err != nil {
		t.Fatalf("\t%s\tShould be able to load regex fixture : %v", tests.Failed, err)
	}
	if err := loadRegex(db, "email.json"); err != nil {
		t.Fatalf("\t%s\tShould be able to load regex fixture : %v", tests.Failed, err)
	}

	if err := loadRelationships("context", db); err != nil {
		t.Fatalf("\t%s\tShould be able to load relationship fixture : %v", tests.Failed, err)
	}

	if err := loadPatterns("context", db); err != nil {
		t.Fatalf("\t%s\tShould be able to load pattern fixture : %v", tests.Failed, err)
	}

	if err := loadViews("context", db); err != nil {
		t.Fatalf("\t%s\tShould be able to load view fixture : %v", tests.Failed, err)
	}

	if err := loadItems("context", db); err != nil {
		t.Fatalf("\t%s\tShould be able to load items : %v", tests.Failed, err)
	}

	return db
}

// loadItems adds items to run tests.
func loadItems(context interface{}, db *db.DB) error {

	store, err := db.GraphHandle(context)
	if err != nil {
		return err
	}

	items, err := itemfix.Get()
	if err != nil {
		return err
	}

	for _, itm := range items {
		if err := sponge.Import(context, db, store, &itm); err != nil {
			return err
		}
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

// unloadItems removes items from the items collection and the graph.
func unloadItems(context interface{}, db *db.DB) error {

	store, err := db.GraphHandle(context)
	if err != nil {
		return err
	}

	items, err := itemfix.Get()
	if err != nil {
		return err
	}

	for _, itm := range items {
		if err := sponge.Remove(context, db, store, itm.ID); err != nil {
			return err
		}
	}

	return nil
}

// loadRegex adds regex to run tests.
func loadRegex(db *db.DB, file string) error {
	rg, err := rfix.Get(file)
	if err != nil {
		return err
	}

	if err := rfix.Add(db, rg); err != nil {
		return err
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

// teardown deinitializes for each indivdual test.
func teardown(t *testing.T, db *db.DB) {
	relationshipfix.Remove(tests.Context, db, "RTEST_")
	viewfix.Remove(tests.Context, db, "VTEST_")
	rfix.Remove(db, "RTEST_")
	unloadItems(tests.Context, db)
	unloadTestData(t, db)
	db.CloseMGO(tests.Context)
	db.CloseCayley(tests.Context)
	tests.DisplayLog()
}

//==============================================================================

// TestExecuteSet tests the execution of different Sets that should succeed.
func TestExecuteSet(t *testing.T) {
	db := setup(t)
	defer teardown(t, db)

	// Build our table of the different test sets.
	execSets := []struct {
		typ string
		set []execSet
	}{
		{typ: "Positive", set: getPosExecSet()},
		{typ: "Negative", set: getNegExecSet()},
	}

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
