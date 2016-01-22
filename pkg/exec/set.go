package exec

import (
	"errors"
	"strings"

	"github.com/coralproject/xenia/pkg/query"

	"github.com/ardanlabs/kit/db"
	"github.com/ardanlabs/kit/db/mongo"
	"github.com/ardanlabs/kit/log"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Set executes the specified query set by name.
func Set(context interface{}, db *db.DB, set *query.Set, vars map[string]string) *query.Result {
	log.Dev(context, "Set", "Started : Name[%s]", set.Name)

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

	log.Dev(context, "Set", "Completed : \n%s\n", mongo.Query(results))
	return &r
}

//==============================================================================

// executePipeline executes the sepcified pipeline query.
func executePipeline(context interface{}, db *db.DB, q *query.Query, vars map[string]string) (docs, error) {

	// Validate we have scripts to run.
	if len(q.Scripts) == 0 {
		return docs{}, errors.New("Invalid pipeline script")
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
		op, err := q.UmarshalMongoScript(script)
		if err != nil {
			return docs{}, err
		}

		// Add the operation to the slice for the pipeline.
		pipeline = append(pipeline, op)
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
