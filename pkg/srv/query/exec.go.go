package query

import (
	"regexp"
	"strings"

	"gopkg.in/mgo.v2/bson"
)

// variableMarkersRegExp represents the variable marker regular expression
// used in matching against a script src.
var variableMarkersRegExp = regexp.MustCompile(`#(.*?)#`)

// renderScript evalues and replaces all variable markers (#variable#) in the
// source script string with the approrpiate value from the map if found.
func renderScript(script string, variables map[string]string) string {
	// Fetches all markers matching a (#.*?#) regexp.
	matches := variableMarkersRegExp.FindAllString(script, -1)

	if len(matches) == 0 {
		return script
	}

	for _, match := range matches {
		varName := strings.Trim(match, "#")
		val, ok := variables[varName]
		if !ok {
			continue
		}
		script = strings.Replace(script, match, val, 1)
	}

	return script
}

// RenderResult converts a bson.M response into a map of key values pairs.
func RenderResult(data bson.M, options VarOption) (map[string]string, error) {
	return nil, nil
}
