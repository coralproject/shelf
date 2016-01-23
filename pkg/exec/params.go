package exec

import (
	"fmt"
	"strings"

	"github.com/coralproject/xenia/pkg/query"

	"github.com/ardanlabs/kit/log"
)

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
