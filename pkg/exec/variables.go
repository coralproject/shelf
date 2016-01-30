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

	// Remove the # characters from the left.
	value := variable[1:]

	// Find the first instance of the separator.
	idx := strings.Index(value, ":")
	if idx == -1 {
		err := fmt.Errorf("Invalid format : %s", variable)
		log.Error(context, "parseVar", err, "Parsing variable")
		return
	}

	// Split the key and variable apart.
	typ := value[0:idx]
	vari := value[idx+1:]

	switch key {
	case "$in":
		if typ == "data" {
			commands[key] = fieldData(context, vari, data)
		}

	default:
		commands[key] = fieldVars(context, typ, vari, vars, data)
	}
}

// fieldVars focuses on replacing variables where the key is a field name.
func fieldVars(context interface{}, typ, variable string, vars map[string]string, data map[string]interface{}) interface{} {

	// Before: {"field": "#number:variable_name"}  After: {"field": 1234}
	// Before: {"field": "#string:variable_name"}  After: {"field": "value"}
	// Before: {"field": "#date:variable_name"}    After: {"field": time.Time}
	// Before: {"field": "#objid:variable_name"}   After: {"field": mgo.ObjectId}
	// Before: {"field": "#data:doc.station_id"}   After: {"field": "23453"}

	// If the variable does not exist, use the variable straight up.
	param, exists := vars[variable]
	if !exists {
		param = variable
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

	case "data":
		return fieldData(context, param, data)
	}

	return variable
}

// fieldData locates the data from the data map for the specified field.
func fieldData(context interface{}, variable string, data map[string]interface{}) interface{} {

	// Before : {"field" : {"$in": "#data:list.station_id"}}}
	// After  : {"field" : {"$in": ["12345", 45678"]}}
	//  	variable : "list.station_id"
	//  	data     : {"list": [{"station_id":"42021"}, {"station_id":"45098"]}]}

	// Before : {"field" : "#data:station.station_id"}
	// After  : {"field" : "12345"}
	//  	variable : "station.station_id"
	//  	data     : {"station": [{"station_id":"42021"}]}

	// Find the data based on the variable and the field lookup.
	values, field, err := findData(context, variable, data)
	if err != nil {
		return variable
	}

	// How many values do we have.
	l := len(values)

	// If there are no values just use the literal value
	// for the variable.
	if l == 0 {
		return variable
	}

	// If there is onlu one value, use it.
	if l == 1 {
		fldValue, exists := values[0][field]
		if !exists {
			log.Error(context, "dataLookup", fmt.Errorf("Field not found : %s", field), "Map field lookup")
			return variable
		}

		return fldValue
	}

	// We have more than one value so return an array of these values.
	var array []interface{}
	for _, doc := range values {

		// We have to find the value for the specified field.
		fldValue, exists := doc[field]
		if !exists {
			log.Error(context, "dataLookup", fmt.Errorf("Field not found : %s", field), "Map field lookup")
			return variable
		}

		// Append the value to the array.
		array = append(array, fldValue)
	}

	return array
}

// findData process the variable and lookups up the data.
func findData(context interface{}, variable string, data map[string]interface{}) ([]bson.M, string, error) {

	// Before: "station.station_id"  After: Data, station_id

	// Split the variable into the data key and document field key.
	idx := strings.Index(variable, ".")
	if idx == -1 {
		err := fmt.Errorf("Invalid format : %s", variable)
		log.Error(context, "findData", err, "Parsing variable")
		return nil, "", err
	}

	// Extract the key and field.
	key := variable[0:idx]
	field := variable[idx+1:]

	// Find the results.
	results, exists := data[key]
	if !exists {
		err := fmt.Errorf("Key not found : %s", key)
		log.Error(context, "findData", err, "Finding results")
		return nil, "", err
	}

	// Extract the concrete type from the interface.
	values, ok := results.([]bson.M)
	if !ok {
		err := errors.New("Expected an array to exist")
		log.Error(context, "findData", err, "Type assert results : %T", results)
		return nil, "", err
	}

	return values, field, nil
}

// isoDate is a helper function to convert the internal extension for dates
// into a BSON date. We convert the following string
// ISODate('2013-01-16T00:00:00.000Z') to a Go time value.
func isoDate(value string) time.Time {
	var parse string

	switch len(value) {
	case 10:
		parse = "2006-01-02"
	case 24:
		parse = "2006-01-02T15:04:05.999Z"
	case 23:
		parse = "2006-01-02T15:04:05.999"
	default:
		return time.Now().UTC()
	}

	dateTime, err := time.Parse(parse, value)
	if err != nil {
		return time.Now().UTC()
	}

	return dateTime
}

// objID is a helper function to convert a string that represents a Mongo
// Object Id into a bson ObjectId type.
func objID(value string) bson.ObjectId {
	if len(value) > 24 {
		return bson.ObjectId("")
	}

	return bson.ObjectIdHex(value)
}
