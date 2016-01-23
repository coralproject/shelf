// Package exec provides support for executing Sets and their different types
// of commands.
package exec

import (
	"errors"
	"strings"

	"github.com/coralproject/xenia/pkg/query"
	"github.com/coralproject/xenia/pkg/script"

	"github.com/ardanlabs/kit/db"
	"github.com/ardanlabs/kit/db/mongo"
	"github.com/ardanlabs/kit/log"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// docs represents what a user will receive after
// excuting a successful set.
type docs struct {
	Name string
	Docs []bson.M
}

// emptyResult is for returning empty runs.
var emptyResult []docs

//==============================================================================

// errResult creates a result value with the error.
func errResult(context interface{}, err error) *query.Result {
	r := query.Result{
		Results: bson.M{"error": err.Error()},
		Error:   true,
	}

	log.Error(context, "errResult", err, "Completed")
	return &r
}

//==============================================================================

// loadPrePostScripts updates each query script slice with pre/post commands.
func loadPrePostScripts(context interface{}, db *db.DB, set *query.Set) error {
	if set.PreScript == "" && set.PstScript == "" {
		return nil
	}

	// Load the set of scripts we need to fetch.
	fetchScripts := make([]string, 2)

	if set.PreScript != "" {
		fetchScripts[0] = set.PreScript
	}

	if set.PstScript != "" {
		fetchScripts[1] = set.PstScript
	}

	// Pull all the script documents we need.
	scripts, err := script.GetByNames(context, db, fetchScripts)
	if err != nil {
		return err
	}

	// Add the commands to the query scripts. Since order of the
	// pre/post scripts is maintained, this is simplified.
	for i := range set.Queries {
		if set.PreScript != "" {
			scripts[0].Commands = append(scripts[0].Commands, set.Queries[i].Commands...)
			set.Queries[i].Commands = scripts[0].Commands
		}

		if set.PstScript != "" {
			set.Queries[i].Commands = append(set.Queries[i].Commands, scripts[1].Commands...)
		}
	}

	return nil
}

//==============================================================================

// Exec executes the specified query set by name.
func Exec(context interface{}, db *db.DB, set *query.Set, vars map[string]string) *query.Result {
	log.Dev(context, "Exec", "Started : Name[%s]", set.Name)

	// Validate the set that is provided.
	if err := set.Validate(); err != nil {
		return errResult(context, err)
	}

	// Is the rule enabled.
	if !set.Enabled {
		return errResult(context, errors.New("Set disabled"))
	}

	// If we have been provided a nil map, make one.
	if vars == nil {
		vars = make(map[string]string)
	}

	// Did we get everything we need. Also load defaults.
	if err := validateParameters(context, set, vars); err != nil {
		return errResult(context, err)
	}

	// Load the pre/post scripts.
	if err := loadPrePostScripts(context, db, set); err != nil {
		return errResult(context, err)
	}

	// Final results of running the set of queries.
	var results []docs

	// Iterate of the set of queries.
	for _, q := range set.Queries {
		var result docs
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

				// Go execute the next query starting over.
				continue
			}

			// We need to return an error result.
			return errResult(context, err)
		}

		// Append these results to the final set.
		if q.Return {
			results = append(results, result)
		}
	}

	// Setup the result we will return.
	r := query.Result{
		Results: results,
	}

	log.Dev(context, "Exec", "Completed : \n%s\n", mongo.Query(results))
	return &r
}

//==============================================================================

// executePipeline executes the sepcified pipeline query.
func executePipeline(context interface{}, db *db.DB, q *query.Query, vars map[string]string) (docs, error) {

	// Validate we have scripts to run.
	if len(q.Commands) == 0 {
		return docs{}, errors.New("Invalid pipeline script")
	}

	var pipeline []bson.M

	// Iterate over the scripts building the pipeline.
	for _, command := range q.Commands {

		// Do we have variables to be substitued.
		if vars != nil {
			command = PreProcess(command, vars)
		}

		// Add the operation to the slice for the pipeline.
		pipeline = append(pipeline, command)
	}

	collName := q.Collection

	// Build the pipeline function for the execution.
	var results []bson.M
	f := func(c *mgo.Collection) error {
		var ops string
		for _, op := range pipeline {
			ops += mongo.Query(op) + ",\n"
		}

		log.Dev(context, "executePipeline", "MGO :\ndb.%s.aggregate([\n%s])", c.Name, ops)
		return c.Pipe(pipeline).All(&results)
	}

	// Execute the pipeline.
	if err := db.ExecuteMGO(context, collName, f); err != nil {
		return docs{}, err
	}

	// If there were no results, return an empty array.
	if results == nil {
		results = []bson.M{}
	}

	return docs{q.Name, results}, nil
}
