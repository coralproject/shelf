package exec

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ardanlabs/kit/log"

	"gopkg.in/mgo.v2/bson"
)

// ProcessVariables walks the document performing variable substitutions.
//
// In some cases we are just replacing the variable name for what is in the map.
// If we have dates and objectid, that requires more work to convert to proper types.
// This function is also accessed by the tstdata package and exec_test.go
func ProcessVariables(context interface{}, commands map[string]interface{}, vars map[string]string, data map[string]interface{}) {
	for key, value := range commands {

		// Test for the type of value we have.
		switch doc := value.(type) {

		// We have another document.
		case map[string]interface{}:
			ProcessVariables(context, doc, vars, data)

		// We have a string value so check it.
		case string:
			if doc != "" && doc[0] == '#' {
				varSub(context, key, doc, commands, vars, data)
			}

		// We have an array of values.
		case []interface{}:

			// Iterate over the array of values.
			for _, subDoc := range doc {

				// What type of subDoc is this array made of.
				switch arrDoc := subDoc.(type) {

				// We have another document.
				case map[string]interface{}:
					ProcessVariables(context, arrDoc, vars, data)

				// We have a string value so check it.
				case string:
					if arrDoc != "" && arrDoc[0] == '#' {
						varSub(context, key, arrDoc, commands, vars, data)
					}
				}
			}
		}
	}
}

// varSub provides branching to execute the correct substitution.
func varSub(context interface{}, key string, variable string, commands map[string]interface{}, vars map[string]string, data map[string]interface{}) {
	switch key {
	case "$in":
		commands[key] = inSub(context, variable, data)

	default:
		commands[key] = fieldSub(context, variable, vars)
	}
}

// parseVar splits a variable in half and returns its parts.
func parseVar(context interface{}, variable string) (key string, value string, err error) {
	// Remove the # characters from the left.
	data := variable[1:]

	// Find the first instance of the separator.
	idx := strings.Index(data, ":")
	if idx == -1 {
		err := fmt.Errorf("Invalid format : %s", variable)
		log.Error(context, "parseVar", err, "Parsing variable")
		return "", "", err
	}

	// Split the key and value apart.
	return data[0:idx], data[idx+1:], nil
}

// inSub focuses on producing an array of values from the data.
func inSub(context interface{}, variable string, data map[string]interface{}) interface{} {

	// Before : {"$in": "#in:list.station_id"}
	//
	// The data map contains a key to a document of the results that were saved.
	// {"key": [{"station_id":"42021"}, {"station_id":"45098"]}]}
	//
	// After  : { $in: []string{"42021", "45098"} }

	// Split the variable into its parts.
	_, value, err := parseVar(context, variable)
	if err != nil {
		return variable
	}

	// Split the value into the data key and document field key.
	idx := strings.Index(value, ".")
	if idx == -1 {
		log.Error(context, "inSub", fmt.Errorf("Invalid format : %s", value), "Parsing value")
		return variable
	}

	// Extract the key and field.
	key := value[0:idx]
	field := value[idx+1:]

	// Find the results.
	results, exists := data[key]
	if !exists {
		log.Error(context, "inSub", fmt.Errorf("Key not found : %s", key), "Finding results")
		return variable
	}

	// Extract the concrete type from the interface.
	values, ok := results.([]bson.M)
	if !ok {
		log.Error(context, "inSub", errors.New("Expected an array to exist"), "Type assert results")
		return variable
	}

	// Iterate over the interface values which represent a document
	// and find the specified field in each document.
	var array []interface{}
	for _, doc := range values {

		// We have to find the value for the specified field.
		fldValue, exists := doc[field]
		if !exists {
			log.Error(context, "inSub", fmt.Errorf("Field not found : %s", field), "Map field lookup")
			return variable
		}

		// Append the value to the array.
		array = append(array, fldValue)
	}

	return array
}

// fieldSub focuses on replacing variables where the key is a field name.
func fieldSub(context interface{}, variable string, vars map[string]string) interface{} {

	// Before : {"field": "#number:variable_name"}  After : {"field": 1234}
	// Before : {"field": "#string:variable_name"}  After : {"field": "value"}
	// Before : {"field": "#date:variable_name"}    After : {"field": time.Time}
	// Before : {"field": "#objid:variable_name"}   After : {"field": mgo.ObjectId}

	// Split the variable into its parts.
	typ, value, err := parseVar(context, variable)
	if err != nil {
		return variable
	}

	// If the variable does not exist, use the value straight up.
	param, exists := vars[value]
	if !exists {
		param = value
	}

	// Let's perform the right action per type.
	switch typ {
	case "number":
		if i, err := strconv.Atoi(param); err == nil {
			return i
		}

	case "string":
		return param

	case "date":
		return isoDate(param)

	case "objid":
		return objID(param)
	}

	return variable
}

// isoDate is a helper function to convert the internal extension for dates
// into a BSON date. We convert the following string
// ISODate('2013-01-16T00:00:00.000Z') to a Go time value.
func isoDate(script string) time.Time {
	var parse string

	switch len(script) {
	case 10:
		parse = "2006-01-02"
	case 24:
		parse = "2006-01-02T15:04:05.999Z"
	case 23:
		parse = "2006-01-02T15:04:05.999"
	default:
		return time.Now().UTC()
	}

	dateTime, err := time.Parse(parse, script)
	if err != nil {
		return time.Now().UTC()
	}

	return dateTime
}

// objID is a helper function to convert a string that represents a Mongo
// Object Id into a bson ObjectId type.
func objID(script string) bson.ObjectId {
	if len(script) > 24 {
		return bson.ObjectId("")
	}

	return bson.ObjectIdHex(script)
}
