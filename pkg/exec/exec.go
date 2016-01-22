// Package exec provides support for executing Sets and their different types
// of commands.
package exec

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/coralproject/xenia/pkg/query"
	"github.com/coralproject/xenia/pkg/script"

	"github.com/ardanlabs/kit/db"
	"github.com/ardanlabs/kit/log"
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

// validateParameters validates the variables against the query string
// of parameters. Plus it loads default values.
func validateParameters(context interface{}, set *query.Set, vars map[string]string) error {

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
			scripts[0].Commands = append(scripts[0].Commands, set.Queries[i].Scripts...)
			set.Queries[i].Scripts = scripts[0].Commands
		}

		if set.PstScript != "" {
			set.Queries[i].Scripts = append(set.Queries[i].Scripts, scripts[1].Commands...)
		}
	}

	return nil
}
