package query

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/coralproject/shelf/pkg/db"
	"github.com/coralproject/shelf/pkg/db/mongo"
	"github.com/coralproject/shelf/pkg/log"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// emptyResult is for returning empty runs.
var emptyResult = []bson.M{}

// ExecuteSet executes the specified query set by name.
func ExecuteSet(context interface{}, db *db.DB, set *Set, vars map[string]string) *Result {
	log.Dev(context, "ExecuteSet", "Started : Name[%s]", set.Name)

	// Setup the result we will return.
	r := Result{
		Results: emptyResult,
	}

	// Validate the variables against the meta-data.
	if len(set.Params) > 0 {
		if vars == nil {
			err := errors.New("Invalid Variables List.")
			r.Error = true
			r.Results = []bson.M{bson.M{"error": err.Error()}}
			log.Error(context, "ExecuteSet", err, "Completed")
			return &r
		}

		// Validate each know parameter is represented in the variable list.
		for _, p := range set.Params {
			if _, ok := vars[p.Name]; !ok {
				err := fmt.Errorf("Variable %s not included with the call.", p.Name)
				r.Error = true
				r.Results = []bson.M{bson.M{"error": err.Error()}}
				log.Error(context, "ExecuteSet", err, "Completed")
				return &r
			}
		}
	}

	// Final results of running the set of queries.
	var results []bson.M

	// Iterate of the set of queries.
	for _, q := range set.Queries {
		var result []bson.M
		var err error

		// We only have pipeline right now.
		switch strings.ToLower(q.Type) {
		case "pipeline":
			result, err = executePipeline(context, db, &q, vars)
		}

		// Was there an error processing the query.
		if err != nil {
			// Were we told to continue to the next one.
			if q.Continue {
				// Reset any existing result, it is invalid.
				r.Results = emptyResult
				r.Error = false

				// Go execute the next query starting over.
				continue
			}

			// We need to return an error result.
			r.Error = true
			r.Results = []bson.M{bson.M{"error": err.Error()}}
			log.Error(context, "ExecuteSet", err, "Completed")
			return &r
		}

		// Append these results to the final set.
		results = append(results, result...)
	}

	// Save the final results to be returned.
	r.Results = results

	log.Dev(context, "executePipeline", "Completed : \n%s\n", mongo.Query(results))
	return &r
}

// executePipeline executes the sepcified pipeline query.
func executePipeline(context interface{}, db *db.DB, q *Query, vars map[string]string) ([]bson.M, error) {
	// Validate we have scripts to run.
	if len(q.Scripts) == 0 {
		return nil, errors.New("Invalid pipeline script")
	}

	var pipeline []bson.M

	// Iterate over the scripts building the pipeline.
	for _, script := range q.Scripts {
		// This marker means to skip over this script.
		if strings.HasPrefix(script, "-") {
			continue
		}

		// Do we have variables to be substitued.
		if vars != nil {
			script = renderScript(script, vars)
		}

		// Unmarshal the script into a bson.M for use.
		op, err := umarshalMongoScript(script, q.ScriptOptions)
		if err != nil {
			return nil, err
		}

		// Add the operation to the slice for the pipeline.
		pipeline = append(pipeline, op)
	}

	collName := q.ScriptOptions.Collection

	// Build the pipeline function for the execution.
	var results []bson.M
	f := func(collection *mgo.Collection) error {
		var ops string
		for _, op := range pipeline {
			ops += mongo.Query(op) + ",\n"
		}

		log.Dev(context, "executePipeline", "MGO : db.%s.aggregate([\n%s])", collName, ops)
		return collection.Pipe(pipeline).All(&results)
	}

	// Execute the pipeline.
	if err := db.ExecuteMGO(context, collName, f); err != nil {
		return nil, err
	}

	// If there were not results, treat it as an error.
	if len(results) == 0 {
		return nil, errors.New("No result")
	}

	return results, nil
}

//==============================================================================

// reVarSub represents a regular expression for processing variables.
var reVarSub = regexp.MustCompile(`#(.*?)#`)

// renderScript replaces variables inside of a query script.
func renderScript(script string, vars map[string]string) string {
	matches := reVarSub.FindAllString(script, -1)
	if matches == nil {
		return script
	}

	for _, match := range matches {
		varName := strings.Trim(match, "#")
		if v, exists := vars[varName]; exists {
			script = strings.Replace(script, match, v, 1)
		}
	}

	return script
}
