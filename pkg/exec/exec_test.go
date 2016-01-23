package exec_test

import (
	"encoding/json"
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
	tests.Init("XENIA")

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
				exec.PreProcess(cmd.doc, cmd.vars)

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

	execSet := getExecSet()

	db := db.NewMGO()
	defer db.CloseMGO()

	t.Logf("Given the need to execute mongo commands.")
	{
		loadTestData(t, db)
		defer unloadTestData(t, db)

		for _, es := range execSet {
			t.Logf("\tWhen using Execute Set %s", es.set.Name)
			{
				result := exec.Exec(tests.Context, db, es.set, es.vars)
				if !es.fail {
					if result.Error {
						t.Errorf("\t%s\tShould be able to execute the query set : %+v", tests.Failed, result.Results)
						continue
					}
					t.Logf("\t%s\tShould be able to execute the query set.", tests.Success)
				} else {
					if !result.Error {
						t.Errorf("\t%s\tShould Not be able to execute the query set : %+v", tests.Failed, result.Results)
						continue
					}
					t.Logf("\t%s\tShould Not be able to execute the query set.", tests.Success)
				}

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
		tstdata.Drop()

		if err := sfix.Remove(db); err != nil {
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

// getExecSet returns the table for the testing.
func getExecSet() []execSet {
	return []execSet{
		querySetBasic(),
		querySetBasicPrePost(),
		querySetWithTime(),
		querySetWithShortTime(),
		querySetWithMultiResults(),
		querySetNoResults(),
		querySetBasicVars(),
		querySetBasicMissingVars(),
		querySetBasicParamDefault(),
	}
}

// querySetBasic starts with a simple query set.
func querySetBasic() execSet {
	return execSet{
		fail: false,
		set: &query.Set{
			Name:    "Basic",
			Enabled: true,
			Queries: []query.Query{
				{
					Name:       "Basic",
					Type:       "pipeline",
					Collection: tstdata.CollectionExecTest,
					Return:     true,
					Commands: []map[string]interface{}{
						{"$match": map[string]interface{}{"station_id": "42021"}},
						{"$project": map[string]interface{}{"_id": 0, "name": 1}},
					},
				},
			},
		},
		results: []string{
			`{"results":[{"Name":"Basic","Docs":[{"name":"C14 - Pasco County Buoy, FL"}]}],"error":false}`,
		},
	}
}

// querySetBasicPrePost executes a simple query with pre/post commands.
func querySetBasicPrePost() execSet {
	return execSet{
		fail: false,
		set: &query.Set{
			Name:      "Basic PrePost",
			Enabled:   true,
			PreScript: "STEST_basic_script_pre",
			PstScript: "STEST_basic_script_pst",
			Queries: []query.Query{
				{
					Name:       "Basic PrePost",
					Type:       "pipeline",
					Collection: tstdata.CollectionExecTest,
					Return:     true,
					Commands: []map[string]interface{}{
						{"$project": map[string]interface{}{"_id": 0, "name": 1}},
					},
				},
			},
		},
		results: []string{
			`{"results":[{"Name":"Basic PrePost","Docs":[{"name":"C14 - Pasco County Buoy, FL"}]}],"error":false}`,
		},
	}
}

// querySetWithTime creates a simple query set using time.
func querySetWithTime() execSet {
	return execSet{
		fail: false,
		set: &query.Set{
			Name:    "Time",
			Enabled: true,
			Queries: []query.Query{
				{
					Name:       "Time",
					Type:       "pipeline",
					Collection: tstdata.CollectionExecTest,
					Return:     true,
					Commands: []map[string]interface{}{
						{"$match": map[string]interface{}{"condition.date": map[string]interface{}{"$gt": "#date:2013-01-01T00:00:00.000Z"}}},
						{"$project": map[string]interface{}{"_id": 0, "name": 1}},
						{"$limit": 2},
					},
				},
			},
		},
		results: []string{
			`{"results":[{"Name":"Time","Docs":[{"name":"C14 - Pasco County Buoy, FL"},{"name":"GULF OF MAINE 78 NM EAST OF PORTSMOUTH,NH"}]}],"error":false}`,
			`{"results":[{"Name":"Time","Docs":[{"name":"GULF OF MAINE 78 NM EAST OF PORTSMOUTH,NH"},{"name":"NANTUCKET 54NM Southeast of Nantucket"}]}],"error":false}`,
		},
	}
}

// querySetWithShortTime creates a simple query set using short time.
func querySetWithShortTime() execSet {
	return execSet{
		fail: false,
		set: &query.Set{
			Name:    "Short Time",
			Enabled: true,
			Queries: []query.Query{
				{
					Name:       "Short Time",
					Type:       "pipeline",
					Collection: tstdata.CollectionExecTest,
					Return:     true,
					Commands: []map[string]interface{}{
						{"$match": map[string]interface{}{"condition.date": map[string]interface{}{"$gt": "#date:2013-01-01"}}},
						{"$project": map[string]interface{}{"_id": 0, "name": 1}},
						{"$limit": 2},
					},
				},
			},
		},
		results: []string{
			`{"results":[{"Name":"Short Time","Docs":[{"name":"C14 - Pasco County Buoy, FL"},{"name":"GULF OF MAINE 78 NM EAST OF PORTSMOUTH,NH"}]}],"error":false}`,
			`{"results":[{"Name":"Short Time","Docs":[{"name":"GULF OF MAINE 78 NM EAST OF PORTSMOUTH,NH"},{"name":"NANTUCKET 54NM Southeast of Nantucket"}]}],"error":false}`,
		},
	}
}

// querySetWithMultiResults creates a simple query set using time.
func querySetWithMultiResults() execSet {
	return execSet{
		fail: false,
		set: &query.Set{
			Name:    "Multi Results",
			Enabled: true,
			Queries: []query.Query{
				{
					Name:       "Basic",
					Type:       "pipeline",
					Collection: tstdata.CollectionExecTest,
					Return:     true,
					Commands: []map[string]interface{}{
						{"$match": map[string]interface{}{"station_id": "42021"}},
						{"$project": map[string]interface{}{"_id": 0, "name": 1}},
					},
				},
				{
					Name:       "Time",
					Type:       "pipeline",
					Collection: tstdata.CollectionExecTest,
					Return:     true,
					Commands: []map[string]interface{}{
						{"$match": map[string]interface{}{"condition.date": map[string]interface{}{"$gt": "#date:2013-01-01T00:00:00.000Z"}}},
						{"$project": map[string]interface{}{"_id": 0, "name": 1}},
						{"$limit": 2},
					},
				},
			},
		},
		results: []string{
			`{"results":[{"Name":"Basic","Docs":[{"name":"C14 - Pasco County Buoy, FL"}]},{"Name":"Time","Docs":[{"name":"C14 - Pasco County Buoy, FL"},{"name":"GULF OF MAINE 78 NM EAST OF PORTSMOUTH,NH"}]}],"error":false}`,
			`{"results":[{"Name":"Basic","Docs":[{"name":"C14 - Pasco County Buoy, FL"}]},{"Name":"Time","Docs":[{"name":"GULF OF MAINE 78 NM EAST OF PORTSMOUTH,NH"},{"name":"NANTUCKET 54NM Southeast of Nantucket"}]}],"error":false}`,
		},
	}
}

// querySetNoResults starts with a simple query set with no results.
func querySetNoResults() execSet {
	return execSet{
		fail: false,
		set: &query.Set{
			Name:    "No Results",
			Enabled: true,
			Queries: []query.Query{
				{
					Name:       "NoResults",
					Type:       "pipeline",
					Collection: tstdata.CollectionExecTest,
					Return:     true,
					Commands: []map[string]interface{}{
						{"$match": map[string]interface{}{"station_id": "XXXXXX"}},
						{"$project": map[string]interface{}{"_id": 0, "name": 1}},
					},
				},
			},
		},
		results: []string{
			`{"results":[{"Name":"NoResults","Docs":[]}],"error":false}`,
		},
	}
}

// querySetBasicVars performs simple query with variables.
func querySetBasicVars() execSet {
	return execSet{
		fail: false,
		vars: map[string]string{"station_id": "42021"},
		set: &query.Set{
			Name:    "Basic Vars",
			Enabled: true,
			Params: []query.Param{
				{Name: "station_id"},
			},
			Queries: []query.Query{
				{
					Name:       "Vars",
					Type:       "pipeline",
					Collection: tstdata.CollectionExecTest,
					Return:     true,
					Commands: []map[string]interface{}{
						{"$match": map[string]interface{}{"station_id": "#string:station_id"}},
						{"$project": map[string]interface{}{"_id": 0, "name": 1}},
					},
				},
			},
		},
		results: []string{
			`{"results":[{"Name":"Vars","Docs":[{"name":"C14 - Pasco County Buoy, FL"}]}],"error":false}`,
		},
	}
}

// querySetBasicMissingVars performs simple query with missing parameters.
func querySetBasicMissingVars() execSet {
	return execSet{
		fail: true,
		set: &query.Set{
			Name:    "Missing Vars",
			Enabled: true,
			Params: []query.Param{
				{Name: "station_id"},
			},
			Queries: []query.Query{
				{
					Name:       "Vars",
					Type:       "pipeline",
					Collection: tstdata.CollectionExecTest,
					Return:     true,
					Commands: []map[string]interface{}{
						{"$match": map[string]interface{}{"station_id": "#string:station_id"}},
						{"$project": map[string]interface{}{"_id": 0, "name": 1}},
					},
				},
			},
		},
		results: []string{
			`{"results":{"error":"Variables [station_id] were not included with the call"},"error":true}`,
		},
	}
}

// querySetBasicParamDefault performs simple query with a default parameters.
func querySetBasicParamDefault() execSet {
	return execSet{
		fail: false,
		set: &query.Set{
			Name:    "Param Default",
			Enabled: true,
			Params: []query.Param{
				{Name: "station_id", Default: "42021"},
			},
			Queries: []query.Query{
				{
					Name:       "Vars",
					Type:       "pipeline",
					Collection: tstdata.CollectionExecTest,
					Return:     true,
					Commands: []map[string]interface{}{
						{"$match": map[string]interface{}{"station_id": "#string:station_id"}},
						{"$project": map[string]interface{}{"_id": 0, "name": 1}},
					},
				},
			},
		},
		results: []string{
			`{"results":[{"Name":"Vars","Docs":[{"name":"C14 - Pasco County Buoy, FL"}]}],"error":false}`,
		},
	}
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
