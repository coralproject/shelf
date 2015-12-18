package query_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/coralproject/xenia/pkg/query"

	"github.com/ardanlabs/kit/db"
	"github.com/ardanlabs/kit/tests"

	"gopkg.in/mgo.v2/bson"
)

func init() {
	tests.Init("XENIA")
	tests.InitMongo()
}

//==============================================================================

// TestUmarshalMongoScript tests the ability to convert string based Mongo
// commands into a bson map for processing.
func TestUmarshalMongoScript(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	scripts := []struct {
		text string
		qry  *query.Query
		cmp  bson.M
	}{
		{
			`{"name":"bill"}`,
			nil,
			bson.M{"name": "bill"},
		},
		{
			`{"date":"ISODate('2013-01-16T00:00:00.000Z')"}`,
			&query.Query{HasDate: true},
			bson.M{"date": time.Date(2013, time.January, 16, 0, 0, 0, 0, time.UTC)},
		},
		{
			`{"_id":"ObjectId(\"5660bc6e16908cae692e0593\")"}`,
			&query.Query{HasObjectID: true},
			bson.M{"_id": bson.ObjectIdHex("5660bc6e16908cae692e0593")},
		},
	}

	t.Logf("Given the need to convert mongo commands.")
	{
		for _, script := range scripts {
			t.Logf("\tWhen using %s with %+v", script.text, script.qry)
			{
				b, err := query.UmarshalMongoScript(script.text, script.qry)
				if err != nil {
					t.Errorf("\t%s\tShould be able to convert without an error : %v", tests.Failed, err)
					continue
				}
				t.Logf("\t%s\tShould be able to convert without an error.", tests.Success)

				if eq := compareBson(b, script.cmp); !eq {
					t.Log(b)
					t.Log(script.cmp)
					t.Errorf("\t%s\tShould get back the expected bson document.", tests.Failed)
					continue
				}
				t.Logf("\t%s\tShould get back the expected bson document.", tests.Success)
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
		err := query.GenerateTestData(db)
		if err != nil {
			t.Fatalf("\t%s\tShould be able to load system with test data : %v", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to load system with test data.", tests.Success)

		defer query.DropTestData()

		for _, es := range execSet {
			t.Logf("\tWhen using Execute Set %s", es.set.Name)
			{
				result := query.ExecuteSet(tests.Context, db, es.set, es.vars)
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

				if string(data) != es.result {
					t.Log(string(data))
					t.Log(es.result)
					t.Errorf("\t%s\tShould have the correct result.", tests.Failed)
					continue
				}
				t.Logf("\t%s\tShould have the correct result", tests.Success)
			}
		}
	}
}

//==============================================================================

// execSet represents the table for the table test of execution tests.
type execSet struct {
	fail   bool
	set    *query.Set
	vars   map[string]string
	result string
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
		querySetWithTime(),
		querySetWithMultiResults(),
		querySetNoResults(),
		querySetMalformed(),
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
					Collection: query.CollectionExecTest,
					Return:     true,
					Scripts: []string{
						`{"$match": {"station_id" : "42021"}}`,
						`{"$project": {"_id": 0, "name": 1}}`,
					},
				},
			},
		},
		result: `{"results":[{"Name":"Basic","Docs":[{"name":"C14 - Pasco County Buoy, FL"}]}],"error":false}`,
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
					Collection: query.CollectionExecTest,
					Return:     true,
					HasDate:    true,
					Scripts: []string{
						`{"$match": {"condition.date" : {"$gt": "ISODate(\"2013-01-01T00:00:00.000Z\")"}}}`,
						`{"$project": {"_id": 0, "name": 1}}`,
						`{"$limit": 2}`,
					},
				},
			},
		},
		result: `{"results":[{"Name":"Time","Docs":[{"name":"C14 - Pasco County Buoy, FL"},{"name":"GULF OF MAINE 78 NM EAST OF PORTSMOUTH,NH"}]}],"error":false}`,
	}
}

// querySetWithMultiResults creates a simple query set using time.
func querySetWithMultiResults() execSet {
	return execSet{
		fail: false,
		set: &query.Set{
			Name:    "MultiResults",
			Enabled: true,
			Queries: []query.Query{
				{
					Name:       "Basic",
					Type:       "pipeline",
					Collection: query.CollectionExecTest,
					Return:     true,
					Scripts: []string{
						`{"$match": {"station_id" : "42021"}}`,
						`{"$project": {"_id": 0, "name": 1}}`,
					},
				},
				{
					Name:       "Time",
					Type:       "pipeline",
					Collection: query.CollectionExecTest,
					Return:     true,
					HasDate:    true,
					Scripts: []string{
						`{"$match": {"condition.date" : {"$gt": "ISODate(\"2013-01-01T00:00:00.000Z\")"}}}`,
						`{"$project": {"_id": 0, "name": 1}}`,
						`{"$limit": 2}`,
					},
				},
			},
		},
		result: `{"results":[{"Name":"Basic","Docs":[{"name":"C14 - Pasco County Buoy, FL"}]},{"Name":"Time","Docs":[{"name":"C14 - Pasco County Buoy, FL"},{"name":"GULF OF MAINE 78 NM EAST OF PORTSMOUTH,NH"}]}],"error":false}`,
	}
}

// querySetNoResults starts with a simple query set with no results.
func querySetNoResults() execSet {
	return execSet{
		fail: false,
		set: &query.Set{
			Name:    "NoResults",
			Enabled: true,
			Queries: []query.Query{
				{
					Name:       "NoResults",
					Type:       "pipeline",
					Collection: query.CollectionExecTest,
					Return:     true,
					Scripts: []string{
						`{"$match": {"station_id" : "XXXXXX"}}`,
						`{"$project": {"_id": 0, "name": 1}}`,
					},
				},
			},
		},
		result: `{"results":[{"Name":"NoResults","Docs":[]}],"error":false}`,
	}
}

// querySetMalformed creates a query set with a malformed query.
func querySetMalformed() execSet {
	return execSet{
		fail: true,
		set: &query.Set{
			Name:    "Malformed",
			Enabled: true,
			Queries: []query.Query{
				{
					Name:       "Malformed",
					Type:       "pipeline",
					Collection: query.CollectionExecTest,
					Return:     true,
					Scripts: []string{
						`{"$match": {"station_id" : "XXXXXX"`,
						`{"$project": {"_id": 0, "name": 1}}`,
					},
				},
			},
		},
		result: `{"results":{"error":"unexpected end of JSON input"},"error":true}`,
	}
}

// querySetBasicVars performs simple query with variables.
func querySetBasicVars() execSet {
	return execSet{
		fail: false,
		vars: map[string]string{"station_id": "42021"},
		set: &query.Set{
			Name:    "BasicVars",
			Enabled: true,
			Params: []query.Param{
				{Name: "station_id"},
			},
			Queries: []query.Query{
				{
					Name:       "Vars",
					Type:       "pipeline",
					Collection: query.CollectionExecTest,
					Return:     true,
					Scripts: []string{
						`{"$match": {"station_id" : "#station_id#"}}`,
						`{"$project": {"_id": 0, "name": 1}}`,
					},
				},
			},
		},
		result: `{"results":[{"Name":"Vars","Docs":[{"name":"C14 - Pasco County Buoy, FL"}]}],"error":false}`,
	}
}

// querySetBasicMissingVars performs simple query with missing parameters.
func querySetBasicMissingVars() execSet {
	return execSet{
		fail: true,
		set: &query.Set{
			Name:    "MissingVars",
			Enabled: true,
			Params: []query.Param{
				{Name: "station_id"},
			},
			Queries: []query.Query{
				{
					Name:       "Vars",
					Type:       "pipeline",
					Collection: query.CollectionExecTest,
					Return:     true,
					Scripts: []string{
						`{"$match": {"station_id" : "#station_id#"}}`,
						`{"$project": {"_id": 0, "name": 1}}`,
					},
				},
			},
		},
		result: `{"results":{"error":"Variables [station_id] were not included with the call"},"error":true}`,
	}
}

// querySetBasicParamDefault performs simple query with a default parameters.
func querySetBasicParamDefault() execSet {
	return execSet{
		fail: false,
		set: &query.Set{
			Name:    "ParamDefault",
			Enabled: true,
			Params: []query.Param{
				{Name: "station_id", Default: "42021"},
			},
			Queries: []query.Query{
				{
					Name:       "Vars",
					Type:       "pipeline",
					Collection: query.CollectionExecTest,
					Return:     true,
					Scripts: []string{
						`{"$match": {"station_id" : "#station_id#"}}`,
						`{"$project": {"_id": 0, "name": 1}}`,
					},
				},
			},
		},
		result: `{"results":[{"Name":"Vars","Docs":[{"name":"C14 - Pasco County Buoy, FL"}]}],"error":false}`,
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
