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
// This function is exported because it is accessed by the tstdata package
// and tests.
func ProcessVariables(context interface{}, commands map[string]interface{}, vars map[string]string, results map[string]interface{}) error {

	// commands: Contains the mongodb pipeline with any extenstions.
	// vars    : Key/Value pairs passed into the set execution for variable substituion.
	// results : Any result from previous sets that have been saved.

	for key, value := range commands {

		// Test for the type of value we have.
		switch doc := value.(type) {

		// We have another document.
		case map[string]interface{}:
			if err := ProcessVariables(context, doc, vars, results); err != nil {
				return err
			}

		// We have a string value so check it.
		case string:
			if doc != "" && doc[0] == '#' {
				if err := varSub(context, key, doc, commands, vars, results); err != nil {
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
					if err := ProcessVariables(context, arrDoc, vars, results); err != nil {
						return err
					}

				// We have a string value so check it.
				case string:
					if arrDoc != "" && arrDoc[0] == '#' {
						if err := varSub(context, key, arrDoc, commands, vars, results); err != nil {
							return err
						}
					}
				}
			}
		}
	}

	return nil
}

// varSub replaces variables inside the command set with values.
func varSub(context interface{}, key, variable string, commands map[string]interface{}, vars map[string]string, results map[string]interface{}) error {

	// Before: {"field": "#number:variable_name"}  After: {"field": 1234}
	// 		key:"field"  variable:"#cmd:variable_name"

	// Remove the # characters from the left.
	value := variable[1:]

	// Find the first instance of the separator.
	idx := strings.Index(value, ":")
	if idx == -1 {
		err := fmt.Errorf("Invalid variable format %q, missing :", variable)
		log.Error(context, "varSub", err, "Parsing variable")
		return err
	}

	// Split the key and variable apart.
	cmd := value[0:idx]
	vari := value[idx+1:]

	switch key {
	case "$in":
		if len(cmd) != 6 || cmd[0:4] != "data" {
			err := fmt.Errorf("Invalid $in command %q, missing \"data\" keyword or malformed", cmd)
			log.Error(context, "varSub", err, "$in command processing")
			return err
		}

		v, err := fieldData(context, cmd[5:6], vari, results)
		if err != nil {
			return err
		}

		commands[key] = v
		return nil

	default:
		v, err := fieldVars(context, cmd, vari, vars, results)
		if err != nil {
			return err
		}

		commands[key] = v
		return nil
	}
}

// fieldVars looks up variables and returns their values as the specified type.
func fieldVars(context interface{}, cmd, variable string, vars map[string]string, results map[string]interface{}) (interface{}, error) {

	// {"field": "#cmd:variable"}
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
			log.Error(context, "fieldVars", err, "Index conversion")
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
			return fieldData(context, cmd[5:6], param, results)
		}

		if cmd == "data" {
			err := errors.New("Data command is missing the operator")
			log.Error(context, "fieldVars", err, "Checking cmd is data")
			return nil, err
		}

		err := fmt.Errorf("Unknown command %q", cmd)
		log.Error(context, "fieldVars", err, "Checking cmd is data")
		return nil, err
	}
}

// fieldData locates the data from the results based on the data operation
// and the lookup value.
func fieldData(context interface{}, dataOp, lookup string, results map[string]interface{}) (interface{}, error) {

	// We always want an array to be subsituted.					// We select the index and subtitue a single value.
	// Before: {"field" : {"$in": "#data.*:list.station_id"}}}		// Before: {"field" : "#data.0:list.station_id"}
	// After : {"field" : {"$in": ["42021"]}}						// After : {"field" : "42021"}
	//      dataOp : "*"                                         	//      dataOp : 0
	//  	lookup : "list.station_id"							    //  	lookup : "list.station_id"
	//  	results: {"list": [{"station_id":"42021"}]}				//  	results: {"list": [{"station_id":"42021"}, {"station_id":"23567"}]}

	// Find the result data based on the lookup and the field lookup.
	data, field, err := findResultData(context, lookup, results)
	if err != nil {
		return "", err
	}

	// How many documents do we have.
	l := len(data)

	// If there are no data just use the literal value for the lookup.
	if l == 0 {
		err := errors.New("The results contain no documents")
		log.Error(context, "fieldData", err, "Checking length")
		return "", err
	}

	// Do we need to return an array.
	if dataOp == "*" {

		// We need to create an array of the values.
		var array []interface{}
		for _, doc := range data {

			// We have to find the value for the specified field.
			fldValue, exists := doc[field]
			if !exists {
				err := fmt.Errorf("Field %q not found", field)
				log.Error(context, "fieldData", err, "Map field lookup")
				return "", err
			}

			// Append the value to the array.
			array = append(array, fldValue)
		}

		return array, nil
	}

	// Convert the index position to an int.
	index, err := strconv.Atoi(dataOp)
	if err != nil {
		err = fmt.Errorf("Invalid operator command operator %q", dataOp)
		log.Error(context, "fieldData", err, "Index conversion")
		return "", err
	}

	// We can't ask for a position we don't have.
	if index > l-1 {
		err := fmt.Errorf("Index \"%d\" out of range, total \"%d\"", index, l)
		log.Error(context, "fieldData", err, "Index range check")
		return "", err
	}

	// Extract the value for the specified index.
	fldValue, exists := data[index][field]
	if !exists {
		err := fmt.Errorf("Field %q not found at index \"%q\"", field, index)
		log.Error(context, "fieldData", err, "Map field lookup")
		return "", err
	}

	return fldValue, nil
}

// findResultData process the lookup against the results. Returns the result if
// found and the field name for location the field from the results later.
func findResultData(context interface{}, lookup string, results map[string]interface{}) ([]bson.M, string, error) {

	// lookup: "station.station_id"		lookup: "list.condition.wind_string"
	// 		key  :   station				key  : list
	//		field: station_id				field: condition.wind_string

	// Split the lookup into the data key and document field key.
	idx := strings.Index(lookup, ".")
	if idx == -1 {
		err := fmt.Errorf("Invalid formated lookup %q", lookup)
		log.Error(context, "findResultData", err, "Parsing lookup")
		return nil, "", err
	}

	// Extract the key and field.
	key := lookup[0:idx]
	field := lookup[idx+1:]

	// Find the result the user is looking for.
	data, exists := results[key]
	if !exists {
		err := fmt.Errorf("Key %q not found in saved results", key)
		log.Error(context, "findResultData", err, "Finding results")
		return nil, "", err
	}

	// Extract the concrete type from the interface.
	values, ok := data.([]bson.M)
	if !ok {
		err := errors.New("** FATAL : Expected the result to be an array of documents")
		log.Error(context, "findResultData", err, "Type assert results : %T", data)
		return nil, "", err
	}

	return values, field, nil
}

// isoDate is a helper function to convert the internal extension for dates
// into a BSON date. We convert the following string
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
