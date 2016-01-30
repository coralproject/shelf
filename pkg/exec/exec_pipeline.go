// Package exec provides support for executing Sets and their different types
// of commands.
package exec

import (
	"errors"
	"fmt"

	"github.com/coralproject/xenia/pkg/query"

	"github.com/ardanlabs/kit/db"
	"github.com/ardanlabs/kit/db/mongo"
	"github.com/ardanlabs/kit/log"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// execPipeline executes the sepcified pipeline query.
func execPipeline(context interface{}, db *db.DB, q *query.Query, vars map[string]string, data map[string]interface{}) (docs, error) {

	// Validate we have scripts to run.
	if len(q.Commands) == 0 {
		return docs{}, errors.New("Invalid pipeline script")
	}

	// We need to check the last command for the extended $save command.
	commands := q.Commands
	l := len(q.Commands) - 1

	// If the last command is a $save, capture its value and remove
	// it from the pipeline.
	var save map[string]interface{}
	if v, exists := q.Commands[l]["$save"]; exists {
		if cmd, ok := v.(map[string]interface{}); ok {
			save = cmd
		}

		commands = q.Commands[0:l]
	}

	var agg string
	var pipeline []bson.M

	// Iterate over the commands and build the pipeline.
	for _, command := range commands {

		// Do we have variables to be substitued.
		if vars != nil {
			ProcessVariables(context, command, vars, data)
		}

		// Add the operation to the slice for the pipeline.
		pipeline = append(pipeline, command)

		// Build a logable version of this pipeline.
		agg += mongo.Query(command) + ",\n"
	}

	// Build the pipeline function for the execution.
	var results []bson.M
	f := func(c *mgo.Collection) error {
		log.Dev(context, "executePipeline", "MGO :\ndb.%s.aggregate([\n%s])", c.Name, agg)
		return c.Pipe(pipeline).All(&results)
	}

	// Execute the pipeline.
	if err := db.ExecuteMGO(context, q.Collection, f); err != nil {
		return docs{}, err
	}

	// If there were no results, return an empty array.
	if results == nil {
		results = []bson.M{}
	}

	// Do we need to save the result.
	if save != nil {
		if err := saveResult(context, save, results, data); err != nil {
			return docs{}, err
		}
	}

	return docs{q.Name, results}, nil
}

// saveResult processes the $save command for this result.
func saveResult(context interface{}, save map[string]interface{}, results []bson.M, data map[string]interface{}) error {

	// {"$map": "list"}

	// Capture the key and value and process the save.
	for cmd, value := range save {
		name, ok := value.(string)
		if !ok {
			return fmt.Errorf("Invalid map key provided : %v", value)
		}

		switch cmd {

		// Save the results into the map under the specified key.
		case "$map":
			log.Dev(context, "saveResult", "Saving result to map[%s]", name)
			data[name] = results
		}

		// Just process the first key because we should only have one.
		return nil
	}

	return nil
}
