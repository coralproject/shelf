package query

import (
	"encoding/json"
	"strings"
	"time"

	"gopkg.in/mgo.v2/bson"
)

// umarshalMongoScript converts a JSON Mongo commands into a BSON map.
func umarshalMongoScript(script string, so *ScriptOption) (bson.M, error) {
	query := []byte(script)

	var op bson.M
	if err := json.Unmarshal(query, &op); err != nil {
		return nil, err
	}

	// We have the HasDate and HasObjectID to prevent us from
	// trying to process these things when it is not necessary.
	if so != nil && (so.HasDate || so.HasObjectID) {
		op = mongoExtensions(op, so)
	}

	return op, nil
}

// mongoExtensions searches for our extensions that need to be converted
// from JSON into BSON, such as dates.
func mongoExtensions(op bson.M, so *ScriptOption) bson.M {
	for key, value := range op {
		// Recurse through the map if provided.
		if doc, ok := value.(map[string]interface{}); ok {
			mongoExtensions(doc, so)
		}

		// Is the value a string.
		if script, ok := value.(string); ok == true {
			if so.HasDate && strings.HasPrefix(script, "ISODate") {
				op[key] = isoDate(script)
			}

			if so.HasObjectID && strings.HasPrefix(script, "ObjectId") {
				op[key] = bson.ObjectIdHex(script[10:34])
			}
		}

		// Is the value an array.
		if array, ok := value.([]interface{}); ok {
			for _, item := range array {
				// Recurse through the map if provided.
				if doc, ok := item.(map[string]interface{}); ok {
					mongoExtensions(doc, so)
				}

				// Is the value a string.
				if script, ok := value.(string); ok == true {
					if so.HasDate && strings.HasPrefix(script, "ISODate") {
						op[key] = isoDate(script)
					}

					if so.HasObjectID && strings.HasPrefix(script, "ObjectId") {
						op[key] = bson.ObjectIdHex(script[10:34])
					}
				}
			}
		}
	}

	return op
}

// isoDate is a helper function to convert the internal extension for dates
// into a BSON date. We convert the following string
// ISODate('2013-01-16T00:00:00.000Z') to a Go time value.
func isoDate(script string) time.Time {
	dateTime, err := time.Parse("2006-01-02T15:04:05.999Z", script[9:len(script)-2])
	if err != nil {
		return time.Now().UTC()
	}

	return dateTime
}
