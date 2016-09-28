package xenia_test

import (
	"github.com/coralproject/shelf/internal/xenia/query"
	"github.com/coralproject/shelf/tstdata"
)

// getNegExecSet returns the table for the testing.
func getNegExecSet() []execSet {
	return []execSet{
		badTime(),
		badObjid(),
		dataMissingOperator(),
		dataMissingInvldOperator(),
		basicMissingVars(),
		dataMissingResults(),
		basicVarRegexFail(),
		basicVarRegexMissing(),
		dataInvldIndex(),
		dataInMalformed(),
		mongoRegexMalformed1(),
		mongoRegexMalformed2(),
	}
}

// withBadTime creates a simple query set using a bad time.
func badTime() execSet {
	return execSet{
		fail: true,
		set: &query.Set{
			Name:    "Bad Time",
			Enabled: true,
			Queries: []query.Query{
				{
					Name:       "Time",
					Type:       "pipeline",
					Collection: tstdata.CollectionExecTest,
					Return:     true,
					Commands: []map[string]interface{}{
						{"$match": map[string]interface{}{"condition.date": map[string]interface{}{"$gt": "#date:2000-1-1"}}},
						{"$project": map[string]interface{}{"_id": 0, "name": 1}},
						{"$limit": 1},
					},
				},
			},
		},
		results: []string{
			`{"results":{"commands":[{"$match":{"condition.date":{"$gt":"#date:2000-1-1"}}},{"$project":{"_id":0,"name":1}},{"$limit":1}],"error":"Invalid date value \"2000-1-1\""}}`,
		},
	}
}

// badObjid creates a simple query set using a bad object id.
func badObjid() execSet {
	return execSet{
		fail: true,
		set: &query.Set{
			Name:    "Bat Objectid",
			Enabled: true,
			Queries: []query.Query{
				{
					Name:       "Objectid",
					Type:       "pipeline",
					Collection: tstdata.CollectionExecTest,
					Return:     true,
					Commands: []map[string]interface{}{
						{"$match": map[string]interface{}{"station_id": "#objid:5660bc6e16908cae"}},
						{"$project": map[string]interface{}{"_id": 0, "name": 1}},
					},
				},
			},
		},
		results: []string{
			`{"results":{"commands":[{"$match":{"station_id":"#objid:5660bc6e16908cae"}},{"$project":{"_id":0,"name":1}}],"error":"Objectid \"5660bc6e16908cae\" is invalid"}}`,
		},
	}
}

// dataMissingOperator performs a test for when the data command is used but
// missing an operator like .* or .N.
func dataMissingOperator() execSet {
	return execSet{
		fail: true,
		set: &query.Set{
			Name:    "Data Missing Operator",
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
						{"$match": map[string]interface{}{"station_id": "#data:station.station_id"}},
						{"$project": map[string]interface{}{"_id": 0, "name": 1}},
					},
				},
			},
		},
		results: []string{
			`{"results":{"commands":[{"$match":{"station_id":"#data:station.station_id"}},{"$project":{"_id":0,"name":1}}],"error":"Data command is missing the operator"}}`,
		},
	}
}

// dataMissingInvldOperator performs a test for when the data command is used but
// the operator is unknown.
func dataMissingInvldOperator() execSet {
	return execSet{
		fail: true,
		set: &query.Set{
			Name:    "Data Missing Invalid Operator",
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
						{"$match": map[string]interface{}{"station_id": "#data.?:station.station_id"}},
						{"$project": map[string]interface{}{"_id": 0, "name": 1}},
					},
				},
			},
		},
		results: []string{
			`{"results":{"commands":[{"$match":{"station_id":"#data.?:station.station_id"}},{"$project":{"_id":0,"name":1}}],"error":"Invalid operator command operator \"?\""}}`,
		},
	}
}

// dataMissingResults performs a test for when the data command is used but
// the results are not found.
func dataMissingResults() execSet {
	return execSet{
		fail: true,
		set: &query.Set{
			Name:    "Data Missing Results",
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
						{"$match": map[string]interface{}{"station_id": "#data.0:list.station_id"}},
						{"$project": map[string]interface{}{"_id": 0, "name": 1}},
					},
				},
			},
		},
		results: []string{
			`{"results":{"commands":[{"$match":{"station_id":"#data.0:list.station_id"}},{"$project":{"_id":0,"name":1}}],"error":"Key \"list\" not found in saved results"}}`,
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
				{Name: "station_id", RegexName: "RTEST_email"},
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
			`{"results":{"error":"Invalid[42021:RTEST_email:Value \"42021\" does not match \"RTEST_email\" expression]"}}`,
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
			`{"results":{"error":"Invalid[42021:numbers:Regex Not found]"}}`,
		},
	}
}

// dataInvldIndex performs a test for when the data command is used but an
// invalid index is selected.
func dataInvldIndex() execSet {
	return execSet{
		fail: true,
		set: &query.Set{
			Name:    "Data Invalid Index",
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
						{"$match": map[string]interface{}{"station_id": "#data.8:station.station_id"}},
						{"$project": map[string]interface{}{"_id": 0, "name": 1}},
					},
				},
			},
		},
		results: []string{
			`{"results":{"commands":[{"$match":{"station_id":"#data.8:station.station_id"}},{"$project":{"_id":0,"name":1}}],"error":"Index \"8\" out of range, total \"1\""}}`,
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
			`{"results":{"error":"Missing[station_id]"}}`,
		},
	}
}

// dataInMalformed performs a test for when the $in command is malformed.
func dataInMalformed() execSet {
	return execSet{
		fail: true,
		set: &query.Set{
			Name:    "Data In Malformed",
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
						{"$match": map[string]interface{}{"station_id": map[string]interface{}{"$in": "#tada:list.station_id"}}},
						{"$project": map[string]interface{}{"_id": 0, "name": 1}},
					},
				},
			},
		},
		results: []string{
			`{"results":{"commands":[{"$match":{"station_id":{"$in":"#tada:list.station_id"}}},{"$project":{"_id":0,"name":1}}],"error":"Invalid $in command \"tada\", missing \"data\" keyword or malformed"}}`,
		},
	}
}

// mongoRegexMalformed1 performs a Mongo malformed regex inside the pipeline.
func mongoRegexMalformed1() execSet {
	return execSet{
		fail: true,
		set: &query.Set{
			Name:    "Mongo Regex Malformed 1",
			Enabled: true,
			Queries: []query.Query{
				{
					Name:       "Mongo Regex",
					Type:       "pipeline",
					Collection: tstdata.CollectionExecTest,
					Return:     true,
					Commands: []map[string]interface{}{
						{"$match": map[string]interface{}{"name": "#regex:east"}},
						{"$group": map[string]interface{}{"_id": "station_id", "count": map[string]interface{}{"$sum": 1}}},
					},
				},
			},
		},
		results: []string{
			`{"results":{"commands":[{"$match":{"name":"#regex:east"}},{"$group":{"_id":"station_id","count":{"$sum":1}}}],"error":"Parameter \"east\" is not a regular expression"}}`,
		},
	}
}

// mongoRegexMalformed2 performs a Mongo malformed regex inside the pipeline.
func mongoRegexMalformed2() execSet {
	return execSet{
		fail: true,
		set: &query.Set{
			Name:    "Mongo Regex Malformed 2",
			Enabled: true,
			Queries: []query.Query{
				{
					Name:       "Mongo Regex",
					Type:       "pipeline",
					Collection: tstdata.CollectionExecTest,
					Return:     true,
					Commands: []map[string]interface{}{
						{"$match": map[string]interface{}{"name": "#regex:/east"}},
						{"$group": map[string]interface{}{"_id": "station_id", "count": map[string]interface{}{"$sum": 1}}},
					},
				},
			},
		},
		results: []string{
			`{"results":{"commands":[{"$match":{"name":"#regex:/east"}},{"$group":{"_id":"station_id","count":{"$sum":1}}}],"error":"Parameter \"/east\" is not a regular expression"}}`,
		},
	}
}
