package query

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

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

// ExecuteSet executes the specified query set by name.
func ExecuteSet(context interface{}, db *db.DB, set *Set, vars map[string]string) *Result {
	log.Dev(context, "ExecuteSet", "Started : Name[%s]", set.Name)

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
	r := Result{
		Results: results,
	}

	log.Dev(context, "ExecuteSet", "Completed : \n%s\n", mongo.Query(results))
	return &r
}

// errResult creates a result value with the error.
func errResult(context interface{}, err error) *Result {
	r := Result{
		Results: bson.M{"error": err.Error()},
		Error:   true,
	}

	log.Error(context, "ExecuteSet", err, "Completed")
	return &r
}

// validateParameters validates the variables against the query string
// of parameters. Plus it loads default values.
func validateParameters(context interface{}, set *Set, vars map[string]string) error {

	// Do we not have parameters.
	if len(set.Params) == 0 {
		return nil
	}

	// Do we not have variables, load the default values.
	if len(vars) == 0 {
		for _, p := range set.Params {
			if p.Default != "" {
				log.Dev(context, "validateParameters", "Adding : Name[%s] Default[%s]", p.Name, p.Default)
				vars[p.Name] = p.Default
			}
		}
	}

	var missing []string

	// Validate each know parameter is represented in the variable list.
	for _, p := range set.Params {
		if _, ok := vars[p.Name]; !ok {

			// The variable was not provided but we have a
			// default value for this so use it.
			if p.Default != "" {
				log.Dev(context, "validateParameters", "Adding : Name[%s] Default[%s]", p.Name, p.Default)
				vars[p.Name] = p.Default
				continue
			}

			// We are missing the parameter.
			missing = append(missing, p.Name)
		}
	}

	// Were there missing parameters.
	if missing == nil {
		return nil
	}

	return fmt.Errorf("Variables [%s] were not included with the call", strings.Join(missing, ","))
}

// loadPrePostScripts updates each query script slice with pre/post commands.
func loadPrePostScripts(context interface{}, db *db.DB, set *Set) error {
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
	for _, q := range set.Queries {
		if set.PreScript != "" {
			scripts[0].Commands = append(scripts[0].Commands, q.Scripts...)
			q.Scripts = scripts[0].Commands
		}

		if set.PstScript != "" {
			q.Scripts = append(q.Scripts, scripts[1].Commands...)
		}
	}

	return nil
}

// executePipeline executes the sepcified pipeline query.
func executePipeline(context interface{}, db *db.DB, q *Query, vars map[string]string) (docs, error) {

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
		op, err := UmarshalMongoScript(script, q)
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

//==============================================================================
// MongoDB specific functions.

// UmarshalMongoScript converts a JSON Mongo commands into a BSON map.
func UmarshalMongoScript(script string, q *Query) (bson.M, error) {
	query := []byte(script)

	var op bson.M
	if err := json.Unmarshal(query, &op); err != nil {
		return nil, err
	}

	// We have the HasDate and HasObjectID to prevent us from
	// trying to process these things when it is not necessary.
	if q != nil && (q.HasDate || q.HasObjectID) {
		op = mongoExtensions(op, q)
	}

	return op, nil
}

// mongoExtensions searches for our extensions that need to be converted
// from JSON into BSON, such as dates.
func mongoExtensions(op bson.M, q *Query) bson.M {
	for key, value := range op {

		// Recurse through the map if provided.
		if doc, ok := value.(map[string]interface{}); ok {
			mongoExtensions(doc, q)
		}

		// Is the value a string.
		if script, ok := value.(string); ok == true {
			if q.HasDate && strings.HasPrefix(script, "ISODate") {
				op[key] = isoDate(script)
			}

			if q.HasObjectID && strings.HasPrefix(script, "ObjectId") {
				op[key] = bson.ObjectIdHex(script[10:34])
			}
		}

		// Is the value an array.
		if array, ok := value.([]interface{}); ok {
			for _, item := range array {

				// Recurse through the map if provided.
				if doc, ok := item.(map[string]interface{}); ok {
					mongoExtensions(doc, q)
				}

				// Is the value a string.
				if script, ok := value.(string); ok == true {
					if q.HasDate && strings.HasPrefix(script, "ISODate") {
						op[key] = isoDate(script)
					}

					if q.HasObjectID && strings.HasPrefix(script, "ObjectId") {
						op[key] = objID(script)
					}
				}
			}
		}
	}

	return op
}

// objID is a helper function to convert a string that represents a Mongo
// Object Id into a bson ObjectId type.
func objID(script string) bson.ObjectId {
	if len(script) > 34 {
		return bson.ObjectId("")
	}

	return bson.ObjectIdHex(script[10:34])
}

// isoDate is a helper function to convert the internal extension for dates
// into a BSON date. We convert the following string
// ISODate('2013-01-16T00:00:00.000Z') to a Go time value.
func isoDate(script string) time.Time {
	var parse string
	if len(script) == 21 {
		parse = "2006-01-02"
	} else {
		parse = "2006-01-02T15:04:05.999Z"
	}

	dateTime, err := time.Parse(parse, script[9:len(script)-2])
	if err != nil {
		return time.Now().UTC()
	}

	return dateTime
}
