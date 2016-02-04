package exec_test

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/coralproject/xenia/pkg/exec"
	"github.com/coralproject/xenia/pkg/query"
	"github.com/coralproject/xenia/pkg/script"
	"github.com/coralproject/xenia/pkg/script/sfix"
	"github.com/coralproject/xenia/tstdata"

	"github.com/ardanlabs/kit/cfg"
	"github.com/ardanlabs/kit/db"
	"github.com/ardanlabs/kit/db/mongo"
	"github.com/ardanlabs/kit/tests"
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

// TestPreProcessing tests the ability to preprocess json documents.
func TestPreProcessing(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	time1, _ := time.Parse("2006-01-02T15:04:05.999Z", "2013-01-16T00:00:00.000Z")
	time2, _ := time.Parse("2006-01-02", "2013-01-16")

	commands := []struct {
		doc   map[string]interface{}
		vars  map[string]string
		after map[string]interface{}
	}{
		{
			map[string]interface{}{"field_name": "#string:name"},
			map[string]string{"name": "bill"},
			map[string]interface{}{"field_name": "bill"},
		},
		{
			map[string]interface{}{"field_name": "#number:value"},
			map[string]string{"value": "10"},
			map[string]interface{}{"field_name": 10},
		},
		{
			map[string]interface{}{"field_name": "#date:value"},
			map[string]string{"value": "2013-01-16T00:00:00.000Z"},
			map[string]interface{}{"field_name": time1},
		},
		{
			map[string]interface{}{"field_name": "#date:value"},
			map[string]string{"value": "2013-01-16"},
			map[string]interface{}{"field_name": time2},
		},
		{
			map[string]interface{}{"field_name": "#date:2013-01-16T00:00:00.000Z"},
			map[string]string{},
			map[string]interface{}{"field_name": time1},
		},
		{
			map[string]interface{}{"field_name": "#objid:value"},
			map[string]string{"value": "5660bc6e16908cae692e0593"},
			map[string]interface{}{"field_name": bson.ObjectIdHex("5660bc6e16908cae692e0593")},
		},
	}

	t.Logf("Given the need to preprocess commands.")
	{
		for _, cmd := range commands {
			t.Logf("\tWhen using %+v with %+v", cmd.doc, cmd.vars)
			{
				exec.ProcessVariables("", cmd.doc, cmd.vars, nil)

				if eq := compareBson(cmd.doc, cmd.after); !eq {
					t.Log(cmd.doc)
					t.Log(cmd.after)
					t.Errorf("\t%s\tShould get back the expected document.", tests.Failed)
					continue
				}
				t.Logf("\t%s\tShould get back the expected document.", tests.Success)
			}
		}
	}
}

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
				t.Logf("\tWhen using Execute Set %s", es.set.Name)
				{
					result := exec.Exec(tests.Context, db, es.set, es.vars)

					data, err := json.Marshal(result)
					if err != nil {
						t.Errorf("\t%s\tShould be able to marshal the result : %s", tests.Failed, err)
						continue
					}
					t.Logf("\t%s\tShould be able to marshal the result.", tests.Success)

					var res query.Result
					if err := json.Unmarshal(data, &res); err != nil {
						t.Errorf("\t%s\tShould be able to unmarshal the result : %s", tests.Failed, err)
						continue
					}
					t.Logf("\t%s\tShould be able to unmarshal the result.", tests.Success)

					// This support allowing the test to provide multiple documents
					// to check when data value order can be underterminstic.
					var found bool
					for _, rslt := range es.results {
						if string(data) == rslt {
							found = true
							break
						}
					}

					if !found {
						t.Log(string(data))
						for _, rslt := range es.results {
							t.Log(rslt)
						}
						t.Errorf("\t%s\tShould have the correct result.", tests.Failed)
						continue
					}
					t.Logf("\t%s\tShould have the correct result", tests.Success)
				}
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
				t.Fatalf("\t%s\tShould load script record from file : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould load script record from file.", tests.Success)

			// We need these scripts loaded under another name to allow tests
			// to run in parallel.
			scr.Name = strings.Replace(scr.Name, "STEST_O", "STEST_T", 1)

			if err := script.Upsert(tests.Context, db, scr); err != nil {
				t.Fatalf("\t%s\tShould be able to create a script : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to create a script.", tests.Success)
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

//==============================================================================

// compareBson compares two bson maps for equivalence.
func compareBson(m1 bson.M, m2 bson.M) bool {
	if len(m1) != len(m2) {
		return false
	}

	for k, v := range m1 {
		if m2[k] != v {
			return false
		}
	}

	for k, v := range m2 {
		if m1[k] != v {
			return false
		}
	}

	return true
}
