package exec

import (
	"strconv"
	"strings"
	"time"

	"gopkg.in/mgo.v2/bson"
)

// PreProcess walks the document preprocessing it for use.
func PreProcess(commands map[string]interface{}, vars map[string]string) {
	for key, value := range commands {

		// Test for the type of value we have.
		switch doc := value.(type) {

		// We have another document.
		case map[string]interface{}:
			PreProcess(doc, vars)

		// We have a string value.
		case string:
			if doc != "" && doc[0] == '#' {
				commands[key] = parse(doc, vars)
			}

		// We have an array of values.
		case []interface{}:

			// Iterate over the array of values.
			for _, subDoc := range doc {

				// What type of subDoc is this array made of.
				switch arrDoc := subDoc.(type) {

				// We have another document.
				case map[string]interface{}:
					PreProcess(arrDoc, vars)

				// We have a string subDoc.
				case string:
					if arrDoc != "" && arrDoc[0] == '#' {
						commands[key] = parse(arrDoc, vars)
					}
				}
			}
		}
	}
}

// renderCommand replaces variables inside of a query command.
func parse(varsub string, vars map[string]string) interface{} {

	// We now have the following expressions.
	// {"field": "#number:variable_name"}
	// {"field": "#string:variable_name"}
	// {"field": "#date:variable_name"}
	// {"field": "#objid:variable_name"}

	// Remove the # characters from the left.
	data := varsub[1:]

	// Find the first instance of the separator.
	idx := strings.Index(data, ":")
	if idx == -1 {
		return varsub
	}

	// Split the two parts.
	typ := data[0:idx]
	variable := data[idx+1:]

	// If the variable does not exist, use the value straight up.
	value, exists := vars[variable]
	if !exists {
		value = variable
	}

	// Let's perform the right action per type.
	switch typ {
	case "number":
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}

	case "string":
		return value

	case "date":
		return isoDate(value)

	case "objid":
		return objID(value)
	}

	return varsub
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
