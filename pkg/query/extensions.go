package query

import (
	"strings"
	"time"

	"gopkg.in/mgo.v2/bson"
)

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
