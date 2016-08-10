package xenia

import (
	"errors"
	"fmt"
	"strings"

	"github.com/ardanlabs/kit/db"
	"github.com/ardanlabs/kit/log"
	"github.com/coralproject/xenia/internal/xenia/query"
	"github.com/coralproject/xenia/internal/xenia/regex"
)

// processParams validates the variables against the query string of parameters.
// It also loads default values and processes parameter regexes.
func processParams(context interface{}, db *db.DB, set *query.Set, vars map[string]string) error {

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

	var errs []string

	// Validate each known parameter is represented in the variable list.
	for _, p := range set.Params {
		if _, exists := vars[p.Name]; !exists {

			// The variable was not provided but we have a
			// default value for this so use it.
			if p.Default != "" {
				log.Dev(context, "validateParameters", "Adding : Name[%s] Default[%s]", p.Name, p.Default)
				vars[p.Name] = p.Default
			} else {

				// We are missing the parameter.
				errs = append(errs, "Missing["+p.Name+"]")
			}
		}

		// Is there a regex to validate against?
		if p.RegexName != "" {
			value := vars[p.Name]
			if err := validateRegex(context, db, value, p.RegexName); err != nil {
				errs = append(errs, "Invalid["+value+":"+p.RegexName+":"+err.Error()+"]")
			}
		}
	}

	// Were there any errors.
	if errs != nil {
		return errors.New(strings.Join(errs, ","))
	}

	return nil
}

// validateRegex compares the value to the configured regex.
func validateRegex(context interface{}, db *db.DB, value string, name string) error {
	rgx, err := regex.GetByName(context, db, name)
	if err != nil {
		return err
	}

	if rgx.Compile == nil {
		err := errors.New("FATAL ERROR: Regex is not pre-compiled")
		log.Error(context, "validateRegex", err, "Check compiled")
		return err
	}

	if !rgx.Compile.MatchString(value) {
		err := fmt.Errorf("Value %q does not match %q expression", value, rgx.Name)
		log.Error(context, "validateRegex", err, "Preform match")
		return err
	}

	return nil
}
