package exec_test

import (
	"github.com/coralproject/xenia/pkg/query"
	"github.com/coralproject/xenia/tstdata"
)

// getNegExecSet returns the table for the testing.
func getNegExecSet() []execSet {
	return []execSet{
		badTime(),
		badObjid(),
		dataMissingOperator(),
		dataMissingInvldOperator(),
		dataMissingResults(),
		dataInvldIndex(),
		dataInMalformed(),
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
			`{"results":{"commands":[{"$match":{"condition.date":{"$gt":"#date:2000-1-1"}}},{"$project":{"_id":0,"name":1}},{"$limit":1}],"error":"Invalid date value \"2000-1-1\""},"error":true}`,
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
			`{"results":{"commands":[{"$match":{"station_id":"#objid:5660bc6e16908cae"}},{"$project":{"_id":0,"name":1}}],"error":"Objectid \"5660bc6e16908cae\" is invalid"},"error":true}`,
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
			`{"results":{"commands":[{"$match":{"station_id":"#data:station.station_id"}},{"$project":{"_id":0,"name":1}}],"error":"Data command is missing the operator"},"error":true}`,
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
			`{"results":{"commands":[{"$match":{"station_id":"#data.?:station.station_id"}},{"$project":{"_id":0,"name":1}}],"error":"Invalid operator command operator \"?\""},"error":true}`,
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
			`{"results":{"commands":[{"$match":{"station_id":"#data.0:list.station_id"}},{"$project":{"_id":0,"name":1}}],"error":"Key \"list\" not found in saved results"},"error":true}`,
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
			`{"results":{"commands":[{"$match":{"station_id":"#data.8:station.station_id"}},{"$project":{"_id":0,"name":1}}],"error":"Index \"8\" out of range, total \"1\""},"error":true}`,
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
			`{"results":{"commands":[{"$match":{"station_id":{"$in":"#tada:list.station_id"}}},{"$project":{"_id":0,"name":1}}],"error":"Invalid $in command \"tada\", missing \"data\" keyword or malformed"},"error":true}`,
		},
	}
}
