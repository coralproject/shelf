package xenia_test

import (
	"github.com/coralproject/shelf/internal/xenia/query"
	"github.com/coralproject/shelf/tstdata"
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
		basicParamDefault(),
		basicVarRegex(),
		basicSaveIn(),
		basicSaveInObjectID(),
		basicSaveVar(),
		multiFieldLookup(),
		mongoRegex(),
		masking(),
		withAdjTime(),
		fieldReplace(),
		explain(),
		basicView(),
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
			`{"results":[{"Name":"NoResults","Docs":[]}]}`,
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
			`{"results":[{"Name":"Basic","Docs":[{"name":"C14 - Pasco County Buoy, FL"}]}]}`,
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
			`{"results":[{"Name":"Basic Array","Docs":[{"name":"C14 - Pasco County Buoy, FL"},{"name":"NANTUCKET 54NM Southeast of Nantucket"}]}]}`,
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
			PreScript: "STEST_T_basic_script_pre",
			PstScript: "STEST_T_basic_script_pst",
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
			`{"results":[{"Name":"Basic PrePost","Docs":[{"name":"C14 - Pasco County Buoy, FL"}]}]}`,
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
			`{"results":[{"Name":"Time","Docs":[{"name":"C14 - Pasco County Buoy, FL"},{"name":"GULF OF MAINE 78 NM EAST OF PORTSMOUTH,NH"}]}]}`,
			`{"results":[{"Name":"Time","Docs":[{"name":"GULF OF MAINE 78 NM EAST OF PORTSMOUTH,NH"},{"name":"NANTUCKET 54NM Southeast of Nantucket"}]}]}`,
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
			`{"results":[{"Name":"Short Time","Docs":[{"name":"C14 - Pasco County Buoy, FL"},{"name":"GULF OF MAINE 78 NM EAST OF PORTSMOUTH,NH"}]}]}`,
			`{"results":[{"Name":"Short Time","Docs":[{"name":"GULF OF MAINE 78 NM EAST OF PORTSMOUTH,NH"},{"name":"NANTUCKET 54NM Southeast of Nantucket"}]}]}`,
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
			`{"results":[{"Name":"Basic","Docs":[{"name":"C14 - Pasco County Buoy, FL"}]},{"Name":"Time","Docs":[{"name":"C14 - Pasco County Buoy, FL"},{"name":"GULF OF MAINE 78 NM EAST OF PORTSMOUTH,NH"}]}]}`,
			`{"results":[{"Name":"Basic","Docs":[{"name":"C14 - Pasco County Buoy, FL"}]},{"Name":"Time","Docs":[{"name":"GULF OF MAINE 78 NM EAST OF PORTSMOUTH,NH"},{"name":"NANTUCKET 54NM Southeast of Nantucket"}]}]}`,
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
			`{"results":[{"Name":"Vars","Docs":[{"name":"C14 - Pasco County Buoy, FL"}]}]}`,
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
			`{"results":[{"Name":"Vars","Docs":[{"name":"C14 - Pasco County Buoy, FL"}]}]}`,
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
				{Name: "station_id", RegexName: "RTEST_number"},
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
			`{"results":[{"Name":"Basic Var Regex","Docs":[{"name":"C14 - Pasco County Buoy, FL"}]}]}`,
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
			`{"results":[{"Name":"Get Documents","Docs":[{"name":"C14 - Pasco County Buoy, FL"},{"name":"GULF OF MAINE 78 NM EAST OF PORTSMOUTH,NH"},{"name":"NANTUCKET 54NM Southeast of Nantucket"},{"name":"GEORGES BANK 170 NM East of Hyannis, MA"},{"name":"SE Cape Cod 30NM East of Nantucket, MA"}]}]}`,
		},
	}
}

// basicSaveInObjectID performs a simple query where the result of the first query
// is used in an $In statement with ObjectId's.
func basicSaveInObjectID() execSet {
	return execSet{
		fail: false,
		set: &query.Set{
			Name:    "Basic Save In ObjectID",
			Enabled: true,
			Queries: []query.Query{
				{
					Name:       "Get Object Ids",
					Type:       "pipeline",
					Collection: tstdata.CollectionExecTest,
					Return:     false,
					Commands: []map[string]interface{}{
						{"$project": map[string]interface{}{"_id": 1, "station_id": 1}},
						{"$limit": 5},
						{"$save": map[string]interface{}{"$map": "list"}},
					},
				},
				{
					Name:       "Get Documents ObjectId",
					Type:       "pipeline",
					Collection: tstdata.CollectionExecTest,
					Return:     true,
					Commands: []map[string]interface{}{
						{"$match": map[string]interface{}{"_id": map[string]interface{}{"$in": "#data.*:list._id"}}},
						{"$project": map[string]interface{}{"_id": 0, "name": 1}},
					},
				},
			},
		},
		results: []string{
			`{"results":[{"Name":"Get Documents ObjectId","Docs":[{"name":"C14 - Pasco County Buoy, FL"},{"name":"GULF OF MAINE 78 NM EAST OF PORTSMOUTH,NH"},{"name":"NANTUCKET 54NM Southeast of Nantucket"},{"name":"GEORGES BANK 170 NM East of Hyannis, MA"},{"name":"SE Cape Cod 30NM East of Nantucket, MA"}]}]}`,
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
			`{"results":[{"Name":"Get Documents","Docs":[{"name":"C14 - Pasco County Buoy, FL"}]}]}`,
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
			`{"results":[{"Name":"Get Documents","Docs":[{"name":"C14 - Pasco County Buoy, FL"}]}]}`,
		},
	}
}

// mongoRegex performs a Mongo regex inside the pipeline.
func mongoRegex() execSet {
	return execSet{
		fail: false,
		set: &query.Set{
			Name:    "Mongo Regex",
			Enabled: true,
			Queries: []query.Query{
				{
					Name:       "Mongo Regex 1",
					Type:       "pipeline",
					Collection: tstdata.CollectionExecTest,
					Return:     true,
					Commands: []map[string]interface{}{
						{"$match": map[string]interface{}{"name": "#regex:/east/i"}},
						{"$group": map[string]interface{}{"_id": "station_id", "count": map[string]interface{}{"$sum": 1}}},
					},
				},
				{
					Name:       "Mongo Regex 2",
					Type:       "pipeline",
					Collection: tstdata.CollectionExecTest,
					Return:     true,
					Commands: []map[string]interface{}{
						{"$match": map[string]interface{}{"name": "#regex:/East/"}},
						{"$group": map[string]interface{}{"_id": "station_id", "count": map[string]interface{}{"$sum": 1}}},
					},
				},
			},
		},
		results: []string{
			`{"results":[{"Name":"Mongo Regex 1","Docs":[{"_id":"station_id","count":5}]},{"Name":"Mongo Regex 2","Docs":[{"_id":"station_id","count":3}]}]}`,
		},
	}
}

// masking projects fields that are configured to be masked.
func masking() execSet {
	return execSet{
		fail: false,
		set: &query.Set{
			Name:    "Masking",
			Enabled: true,
			Queries: []query.Query{
				{
					Name:       "Masking",
					Type:       "pipeline",
					Collection: tstdata.CollectionExecTest,
					Return:     true,
					Commands: []map[string]interface{}{
						{"$match": map[string]interface{}{"station_id": "42021"}},
						{"$project": map[string]interface{}{"_id": 0, "name": 1, "condition.observation_time": 1, "condition.pressure_string": 1}},
					},
				},
			},
		},
		results: []string{
			`{"results":[{"Name":"Masking","Docs":[{"condition":{"pressure_string":"******"},"name":"C14 - Pasco County Buoy, FL"}]}]}`,
			`{"results":[{"Name":"Masking","Docs":[{"condition":{"observation_time":"Last Updated on Oct 30 2012, 11:00 am CDT","pressure_string":"1014.0 mb"},"name":"C14 - Pasco County Buoy, FL"}]}]}`,
			`{"results":[{"Name":"Masking","Docs":[{"condition":{"pressure_string":"1014.0 mb"},"name":"C14 - Pasco County Buoy, FL"}]}]}`,
		},

		// NOT SURE WHAT TO DO. When tests are run in parallel the masks may be
		// gone. I can't fudge this because it is tied to the collection we
		// are running the query again. So I have both results for now and I will
		// check for both. One will be right :)
	}
}

// withAdjTime creates a simple query using the time command.
func withAdjTime() execSet {
	return execSet{
		fail: false,
		set: &query.Set{
			Name:    "Since",
			Enabled: true,
			Queries: []query.Query{
				{
					Name:       "Since",
					Type:       "pipeline",
					Collection: tstdata.CollectionExecTest,
					Return:     true,
					Commands: []map[string]interface{}{
						{"$match": map[string]interface{}{"condition.date": map[string]interface{}{"$gt": "#time:-87600h"}}},
						{"$project": map[string]interface{}{"_id": 0, "name": 1}},
						{"$limit": 2},
					},
				},
			},
		},
		results: []string{
			`{"results":[{"Name":"Since","Docs":[{"name":"C14 - Pasco County Buoy, FL"},{"name":"GULF OF MAINE 78 NM EAST OF PORTSMOUTH,NH"}]}]}`,
			`{"results":[{"Name":"Since","Docs":[{"name":"GULF OF MAINE 78 NM EAST OF PORTSMOUTH,NH"},{"name":"NANTUCKET 54NM Southeast of Nantucket"}]}]}`,
		},
	}
}

// fieldReplace tests the replacement of fields.
func fieldReplace() execSet {
	return execSet{
		fail: false,
		vars: map[string]string{"cond": "condition", "dt": "date"},
		set: &query.Set{
			Name:    "Find Replace",
			Enabled: true,
			Queries: []query.Query{
				{
					Name:       "Find Replace",
					Type:       "pipeline",
					Collection: tstdata.CollectionExecTest,
					Return:     true,
					Commands: []map[string]interface{}{
						{"$match": map[string]interface{}{"{cond}.{dt}": map[string]interface{}{"$gt": "#date:2013-01-01T00:00:00.000Z"}}},
						{"$project": map[string]interface{}{"_id": 0, "name": 1}},
						{"$limit": 2},
					},
				},
			},
		},
		results: []string{
			`{"results":[{"Name":"Find Replace","Docs":[{"name":"C14 - Pasco County Buoy, FL"},{"name":"GULF OF MAINE 78 NM EAST OF PORTSMOUTH,NH"}]}]}`,
			`{"results":[{"Name":"Find Replace","Docs":[{"name":"GULF OF MAINE 78 NM EAST OF PORTSMOUTH,NH"},{"name":"NANTUCKET 54NM Southeast of Nantucket"}]}]}`,
		},
	}
}

// explain tests the use of the explain output.
func explain() execSet {
	return execSet{
		fail: false,
		set: &query.Set{
			Name:    "Explain",
			Enabled: true,
			Explain: true,
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
			`#find:queryPlanner`,
		},
	}
}

// basicView performs simple query on a view.
func basicView() execSet {
	return execSet{
		fail: false,
		vars: map[string]string{
			"view":             "VTEST_thread",
			"item":             "ITEST_c1b2bbfe-af9f-4903-8777-bd47c4d5b20a",
			"item_of_interest": "ITEST_d1dfa366-d2f7-4a4a-a64f-af89d4c97d82",
		},
		set: &query.Set{
			Name:    "Basic View",
			Enabled: true,
			Params: []query.Param{
				{Name: "item_of_interest"},
			},
			Queries: []query.Query{
				{
					Name:       "ViewVars",
					Type:       "pipeline",
					Collection: "view",
					Return:     true,
					Commands: []map[string]interface{}{
						{"$match": map[string]interface{}{"item_id": "#string:item_of_interest"}},
						{"$project": map[string]interface{}{"_id": 0, "item_id": 1}},
					},
				},
			},
		},
		results: []string{
			`{"results":[{"Name":"ViewVars","Docs":[{"item_id":"ITEST_d1dfa366-d2f7-4a4a-a64f-af89d4c97d82"}]}]}`,
		},
	}
}
