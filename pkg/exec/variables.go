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
func ProcessVariables(context interface{}, commands map[string]interface{}, vars map[string]string, data map[string]interface{}) error {
	for key, value := range commands {

		// Test for the type of value we have.
		switch doc := value.(type) {

		// We have another document.
		case map[string]interface{}:
			if err := ProcessVariables(context, doc, vars, data); err != nil {
				return err
			}

		// We have a string value so check it.
		case string:
			if doc != "" && doc[0] == '#' {
				if err := varSub(context, key, doc, commands, vars, data); err != nil {
					return err
				}
			}

		// We have an array of values.
		case []interface{}:

			// Iterate over the array of values.
			for _, subDoc := range doc {

				// What type of subDoc is this array made of.
				switch arrDoc := subDoc.(type) {

				// We have another document.
				case map[string]interface{}:
					if err := ProcessVariables(context, arrDoc, vars, data); err != nil {
						return err
					}

				// We have a string value so check it.
				case string:
					if arrDoc != "" && arrDoc[0] == '#' {
						if err := varSub(context, key, arrDoc, commands, vars, data); err != nil {
							return err
						}
					}
				}
			}
		}
	}

	return nil
}

// varSub provides branching to execute the correct substitution.
func varSub(context interface{}, key, variable string, commands map[string]interface{}, vars map[string]string, data map[string]interface{}) error {

	// key:"field"  variable:"#cmd:variable_name"

	// Remove the # characters from the left.
	value := variable[1:]

	// Find the first instance of the separator.
	idx := strings.Index(value, ":")
	if idx == -1 {
		err := fmt.Errorf("Invalid variable format %q", variable)
		log.Error(context, "parseVar", err, "Parsing variable")
		return err
	}

	// Split the key and variable apart.
	cmd := value[0:idx]
	vari := value[idx+1:]

	switch key {
	case "$in":
		if len(cmd) == 6 && cmd[0:4] == "data" {
			v, err := fieldData(context, cmd[5:6], vari, data)
			if err != nil {
				return err
			}
			commands[key] = v
		}

	default:
		v, err := fieldVars(context, cmd, vari, vars, data)
		if err != nil {
			return err
		}
		commands[key] = v
	}

	return nil
}

// fieldVars focuses on replacing variables where the key is a field name.
func fieldVars(context interface{}, cmd, variable string, vars map[string]string, data map[string]interface{}) (interface{}, error) {

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

	// Let's perform the right action per command.
	switch cmd {
	case "number":
		i, err := strconv.Atoi(param)
		if err != nil {
			err = fmt.Errorf("Parameter %q is not a number", param)
			log.Error(context, "fieldData", err, "Index conversion")
			return nil, err
		}
		return i, nil

	case "string":
		return param, nil

	case "date":
		return isoDate(context, param)

	case "objid":
		return objID(context, param)

	default:
		if len(cmd) == 6 && cmd[0:4] == "data" {
			return fieldData(context, cmd[5:6], param, data)
		}

		if cmd == "data" {
			return nil, errors.New("Data command is missing the operator")
		}

		return nil, fmt.Errorf("Unknown command %q", cmd)
	}
}

// fieldData locates the data from the data map for the specified field.
func fieldData(context interface{}, cmdOp, variable string, data map[string]interface{}) (interface{}, error) {

	// We always want an array to be subsituted.					// We select the index and subtitue a single value.
	// Before : {"field" : {"$in": "#data.*:list.station_id"}}}		// Before : {"field" : "#data.0:list.station_id"}
	// After  : {"field" : {"$in": ["42021"]}}						// After  : {"field" : "42021"}
	//  	variable : "list.station_id"							//  	variable : "list.station_id"
	//  	data     : {"list": [{"station_id":"42021"}]}			//  	data     : {"list": [{"station_id":"42021"}, {"station_id":"23567"}]}

	// Find the data based on the variable and the field lookup.
	values, field, err := findData(context, variable, data)
	if err != nil {
		return variable, err
	}

	// How many values do we have.
	l := len(values)

	// If there are no values just use the literal value for the variable.
	if l == 0 {
		err := errors.New("No values returned")
		log.Error(context, "fieldData", err, "Checking length")
		return variable, err
	}

	// Do we need to return an array.
	if cmdOp == "*" {

		// We have more than one value so return an array of these values.
		var array []interface{}
		for _, doc := range values {

			// We have to find the value for the specified field.
			fldValue, exists := doc[field]
			if !exists {
				err := fmt.Errorf("Field %q not found", field)
				log.Error(context, "fieldData", err, "Map field lookup")
				return variable, err
			}

			// Append the value to the array.
			array = append(array, fldValue)
		}

		return array, nil
	}

	// Convert the index position to an int.
	index, err := strconv.Atoi(cmdOp)
	if err != nil {
		err = fmt.Errorf("Invalid operator command operator %q", cmdOp)
		log.Error(context, "fieldData", err, "Index conversion")
		return variable, err
	}

	// We can't ask for a position we don't have.
	if index > l-1 {
		err := fmt.Errorf("Index \"%d\" out of range, total \"%d\"", index, l)
		log.Error(context, "fieldData", err, "Index range check")
		return variable, err
	}

	fldValue, exists := values[index][field]
	if !exists {
		err := fmt.Errorf("Field %q not found at index \"%q\"", field, index)
		log.Error(context, "fieldData", err, "Map field lookup")
		return variable, err
	}

	return fldValue, nil
}

// findData process the variable and lookups up the data.
func findData(context interface{}, variable string, data map[string]interface{}) ([]bson.M, string, error) {

	// Before: "station.station_id"  After: {Field Data}, station_id

	// Split the variable into the data key and document field key.
	idx := strings.Index(variable, ".")
	if idx == -1 {
		err := fmt.Errorf("Invalid formated variable %q", variable)
		log.Error(context, "findData", err, "Parsing variable")
		return nil, "", err
	}

	// Extract the key and field.
	key := variable[0:idx]
	field := variable[idx+1:]

	// Find the results.
	results, exists := data[key]
	if !exists {
		err := fmt.Errorf("Key %q not found in saved results", key)
		log.Error(context, "findData", err, "Finding results")
		return nil, "", err
	}

	// Extract the concrete type from the interface.
	values, ok := results.([]bson.M)
	if !ok {
		err := errors.New("** FATAL : Expected the result to be an array of documents")
		log.Error(context, "findData", err, "Type assert results : %T", results)
		return nil, "", err
	}

	return values, field, nil
}

// isoDate is a helper function to convert the internal extension for dates
// into a BSON date. We convert the following string
// ISODate('2013-01-16T00:00:00.000Z') to a Go time value.
func isoDate(context interface{}, value string) (time.Time, error) {
	var parse string

	switch len(value) {
	case 10:
		parse = "2006-01-02"
	case 24:
		parse = "2006-01-02T15:04:05.999Z"
	case 23:
		parse = "2006-01-02T15:04:05.999"
	default:
		err := fmt.Errorf("Invalid date value %q", value)
		log.Error(context, "isoDate", err, "Selecting date parse string")
		return time.Time{}, err
	}

	dateTime, err := time.Parse(parse, value)
	if err != nil {
		log.Error(context, "isoDate", err, "Parsing date string")
		return time.Time{}, err
	}

	return dateTime, nil
}

// objID is a helper function to convert a string that represents a Mongo
// Object Id into a bson ObjectId type.
func objID(context interface{}, value string) (bson.ObjectId, error) {
	if !bson.IsObjectIdHex(value) {
		err := fmt.Errorf("Objectid %q is invalid", value)
		log.Error(context, "objID", err, "Checking obj validity")
		return bson.ObjectId(""), err
	}

	return bson.ObjectIdHex(value), nil
}
