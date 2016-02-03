package exec_test

import (
	"github.com/coralproject/xenia/pkg/query"
	"github.com/coralproject/xenia/tstdata"
)

// getPosExecSet returns the table for the testing.
func getPosExecSet() []execSet {
	return []execSet{
		noResults(),
		basic(),
		basicArray(),
		basicPrePost(),
		withTime(),
		withShortTime(),
		withMultiResults(),
		basicVars(),
		basicMissingVars(),
		basicParamDefault(),
		basicVarRegex(),
		basicVarRegexFail(),
		basicVarRegexMissing(),
		basicSaveIn(),
		basicSaveVar(),
		multiFieldLookup(),
	}
}

// noResults starts with a simple query set with no results.
func noResults() execSet {
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

// basic starts with a simple query set.
func basic() execSet {
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

// basicArray validates hard coded arrays work.
func basicArray() execSet {
	return execSet{
		fail: false,
		set: &query.Set{
			Name:    "Basic Array",
			Enabled: true,
			Queries: []query.Query{
				{
					Name:       "Basic Array",
					Type:       "pipeline",
					Collection: tstdata.CollectionExecTest,
					Return:     true,
					Commands: []map[string]interface{}{
						{"$match": map[string]interface{}{"station_id": map[string]interface{}{"$in": []string{"42021", "44008"}}}},
						{"$project": map[string]interface{}{"_id": 0, "name": 1}},
					},
				},
			},
		},
		results: []string{
			`{"results":[{"Name":"Basic Array","Docs":[{"name":"C14 - Pasco County Buoy, FL"},{"name":"NANTUCKET 54NM Southeast of Nantucket"}]}],"error":false}`,
		},
	}
}

// basicPrePost executes a simple query with pre/post commands.
func basicPrePost() execSet {
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

// withTime creates a simple query set using time.
func withTime() execSet {
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

// withShortTime creates a simple query set using short time.
func withShortTime() execSet {
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

// withMultiResults creates a simple query set using time.
func withMultiResults() execSet {
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

// basicVars performs simple query with variables.
func basicVars() execSet {
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

// basicMissingVars performs simple query with missing parameters.
func basicMissingVars() execSet {
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
			`{"results":{"error":"Missing[station_id]"},"error":true}`,
		},
	}
}

// basicParamDefault performs simple query with a default parameters.
func basicParamDefault() execSet {
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

// basicVarRegex performs simple query with variables and regex validation.
func basicVarRegex() execSet {
	return execSet{
		fail: false,
		vars: map[string]string{"station_id": "42021"},
		set: &query.Set{
			Name:    "Basic Var Regex",
			Enabled: true,
			Params: []query.Param{
				{Name: "station_id", RegexName: "number"},
			},
			Queries: []query.Query{
				{
					Name:       "Basic Var Regex",
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
			`{"results":[{"Name":"Basic Var Regex","Docs":[{"name":"C14 - Pasco County Buoy, FL"}]}],"error":false}`,
		},
	}
}

// basicVarRegexFail performs simple query with variables and an
// invalid regex validation.
func basicVarRegexFail() execSet {
	return execSet{
		fail: true,
		vars: map[string]string{"station_id": "42021"},
		set: &query.Set{
			Name:    "Basic Var Regex Fail",
			Enabled: true,
			Params: []query.Param{
				{Name: "station_id", RegexName: "email"},
			},
			Queries: []query.Query{
				{
					Name:       "Basic Var Regex Fail",
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
			`{"results":{"error":"Invalid[42021:email:Value \"42021\" does not match \"email\" expression]"},"error":true}`,
		},
	}
}

// basicVarRegexMissing performs simple query with variables and a
// missing regex validation.
func basicVarRegexMissing() execSet {
	return execSet{
		fail: true,
		vars: map[string]string{"station_id": "42021"},
		set: &query.Set{
			Name:    "Basic Var Regex Missing",
			Enabled: true,
			Params: []query.Param{
				{Name: "station_id", RegexName: "numbers"},
			},
			Queries: []query.Query{
				{
					Name:       "Basic Var Regex Missing",
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
			`{"results":{"error":"Invalid[42021:numbers:Regex Not found]"},"error":true}`,
		},
	}
}

// basicSaveIn performs a simple query where the result of the first query
// is used in an $In statement.
func basicSaveIn() execSet {
	return execSet{
		fail: false,
		set: &query.Set{
			Name:    "Basic Save In",
			Enabled: true,
			Queries: []query.Query{
				{
					Name:       "Get Ids",
					Type:       "pipeline",
					Collection: tstdata.CollectionExecTest,
					Return:     false,
					Commands: []map[string]interface{}{
						{"$project": map[string]interface{}{"_id": 0, "station_id": 1}},
						{"$limit": 5},
						{"$save": map[string]interface{}{"$map": "list"}},
					},
				},
				{
					Name:       "Get Documents",
					Type:       "pipeline",
					Collection: tstdata.CollectionExecTest,
					Return:     true,
					Commands: []map[string]interface{}{
						{"$match": map[string]interface{}{"station_id": map[string]interface{}{"$in": "#data.*:list.station_id"}}},
						{"$project": map[string]interface{}{"_id": 0, "name": 1}},
					},
				},
			},
		},
		results: []string{
			`{"results":[{"Name":"Get Documents","Docs":[{"name":"C14 - Pasco County Buoy, FL"},{"name":"GULF OF MAINE 78 NM EAST OF PORTSMOUTH,NH"},{"name":"NANTUCKET 54NM Southeast of Nantucket"},{"name":"GEORGES BANK 170 NM East of Hyannis, MA"},{"name":"SE Cape Cod 30NM East of Nantucket, MA"}]}],"error":false}`,
		},
	}
}

// basicSaveVar performs a simple query where the result of the first query
// is used in a variable replacement.
func basicSaveVar() execSet {
	return execSet{
		fail: false,
		set: &query.Set{
			Name:    "Basic Save Var",
			Enabled: true,
			Queries: []query.Query{
				{
					Name:       "Get Ids",
					Type:       "pipeline",
					Collection: tstdata.CollectionExecTest,
					Return:     false,
					Commands: []map[string]interface{}{
						{"$project": map[string]interface{}{"_id": 0, "station_id": 1}},
						{"$limit": 1},
						{"$save": map[string]interface{}{"$map": "station"}},
					},
				},
				{
					Name:       "Get Documents",
					Type:       "pipeline",
					Collection: tstdata.CollectionExecTest,
					Return:     true,
					Commands: []map[string]interface{}{
						{"$match": map[string]interface{}{"station_id": "#data.0:station.station_id"}},
						{"$project": map[string]interface{}{"_id": 0, "name": 1}},
					},
				},
			},
		},
		results: []string{
			`{"results":[{"Name":"Get Documents","Docs":[{"name":"C14 - Pasco County Buoy, FL"}]}],"error":false}`,
		},
	}
}

func multiFieldLookup() execSet {
	return execSet{
		fail: false,
		set: &query.Set{
			Name:    "Multi Field Lookup",
			Enabled: true,
			Queries: []query.Query{
				{
					Name:       "Get Document",
					Type:       "pipeline",
					Collection: tstdata.CollectionExecTest,
					Return:     false,
					Commands: []map[string]interface{}{
						{"$limit": 1},
						{"$save": map[string]interface{}{"$map": "station"}},
					},
				},
				{
					Name:       "Get Documents",
					Type:       "pipeline",
					Collection: tstdata.CollectionExecTest,
					Return:     true,
					Commands: []map[string]interface{}{
						{"$match": map[string]interface{}{"location.type": "#data.0:station.location.type"}},
						{"$match": map[string]interface{}{"station_id": "#data.0:station.station_id"}},
						{"$project": map[string]interface{}{"_id": 0, "name": 1}},
					},
				},
			},
		},
		results: []string{
			`{"results":[{"Name":"Get Documents","Docs":[{"name":"C14 - Pasco County Buoy, FL"}]}],"error":false}`,
		},
	}
}
