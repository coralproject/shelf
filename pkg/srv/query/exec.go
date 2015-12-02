package query

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/coralproject/shelf/pkg/db"
	"github.com/coralproject/shelf/pkg/db/mongo"
	"github.com/coralproject/shelf/pkg/log"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// ExecuteQuerySet executes the giving query set, If found by its name and
// if the Set is a pipeline Type.
func ExecuteQuerySet(context interface{}, db *db.DB, name string) ([]Result, error) {
	log.Dev(context, "ExecuteQuerySet", "Started : Query Set[%s]", name)
	var res []Result

	set, err := GetSetByName(context, db, name)
	if err != nil {
		log.Error(context, "ExecuteQuerySet", err, "Completed : Query Set[%s]", name)
		return nil, err
	}

	// TODO: there is more to do here than just iterating the query list,
	// for now assume we are not needing a feedback of last result, but
	// just packing results from each query in the list.
	for _, q := range set.Queries {
		switch strings.ToLower(q.Type) {
		// TODO: do we need to standardize some type for these?
		case "pipeline":
			result, err := ExecuteQueryPipeline(context, db, q, set)
			if err != nil {
				// TODO: do we return the error or just pack in the result and set Result.Error
				// to true?
				// For now, lets return error.
				log.Error(context, "ExecuteQuerySet", err, "Completed : Query Set[%s]", name)
				return res, err
			}

			res = append(res, result)
		case "template":
			// No implementation yet
			continue
		}
	}

	log.Dev(context, "ExecuteQuerySet", "Completed : Query Set[%s]", name)
	return res, nil
}

// ExecuteQueryPipeline executes the giving PipelineType query from a supplied Set.
// Returns a Result type if successfull else returns an non-nil error.
func ExecuteQueryPipeline(context interface{}, db *db.DB, q Query, set *Set) (Result, error) {
	log.Dev(context, "ExecuteQueryPipeline", "Started : Query Set[%s] : Collect[%s]", set.Name, q.ScriptOptions.Collection)
	var res Result

	// TODO: Decide if we really need to check this?

	res.FeedName = set.Name
	res.QueryType = q.Type
	res.Collection = q.ScriptOptions.Collection

	// build parameter map for variable swapping.
	params := buildParamMap(set.Params)

	var renderedScripts []bson.M

	// render and transform into a map
	for _, script := range q.Scripts {
		rendered := renderScript(script, params)
		section := make(bson.M)

		//TODO: We need to sort out ScriptOptions.HasDate and ScriptOptions.HasObjectID

		if err := json.Unmarshal([]byte(rendered), &section); err != nil {
			res.Error = true
			log.Error(context, "ExecuteQueryPipeline", err, "Completed : Query Set[%s] : Collect[%s]", set.Name, q.ScriptOptions.Collection)
			return res, err
		}

		renderedScripts = append(renderedScripts, section)
	}

	var response []bson.M
	err := db.ExecuteMGO(context, q.ScriptOptions.Collection, func(c *mgo.Collection) error {
		return c.Pipe(&renderedScripts).All(&response)
	})

	if err != nil {
		res.Error = true
		log.Error(context, "ExecuteQueryPipeline", err, "Completed : Query Set[%s] : Collect[%s]", set.Name, q.ScriptOptions.Collection)
		return res, err
	}

	res.Results = renderBSONList(response, q.VarOptions)

	log.Dev(context, "ExecuteQueryPipeline", "Completed : Query Set[%s] : Collect[%s] : %s", set.Name, q.ScriptOptions.Collection, fmt.Sprintf("\n%s", mongo.Query(res)))
	return res, nil
}

//==============================================================================

// buildParamMap returns a map of parameter keys and default values from a list
// of Param.
func buildParamMap(params []SetParam) map[string]string {
	parameters := make(map[string]string)
	for _, param := range params {
		parameters[param.Name] = param.Default
	}
	return parameters
}

// variableMarkersRegExp represents the variable marker regular expression
// used in matching against a script src.
var variableMarkersRegExp = regexp.MustCompile(`#(.*?)#`)

// renderBSONList returns a list of maps with each index corresponding to the
// bson within the provided list.
func renderBSONList(data []bson.M, options *VarOption) []map[string]string {
	var mapped []map[string]string
	for _, block := range data {
		mapped = append(mapped, renderBSONMap(block, options))
	}

	return mapped
}

// renderBSONMap converts a bson.M response into a map of key values pairs.
func renderBSONMap(data bson.M, options *VarOption) map[string]string {
	rendered := make(map[string]string)

	for key, value := range data {
		switch v := value.(type) {
		case string, int, uint, float64, time.Time:
			rendered[key] = toString(v)
		case bson.ObjectId:
			rendered[key] = objectIDToString(v, (options != nil && !!options.ObjectID))
		case []string:
			var stringed []string
			for _, s := range v {
				stringed = append(stringed, "\""+s+"\"")
			}

			rendered[key] = fmt.Sprintf("[%s]", strings.Join(stringed, ","))
		case []int:
			var stringed []string
			for _, d := range v {
				stringed = append(stringed, fmt.Sprintf("%d", d))
			}

			rendered[key] = fmt.Sprintf("[%s]", strings.Join(stringed, ","))
		case []interface{}:
			var stringed []string
			for _, u := range v {
				stringed = append(stringed, toString(u))
			}

			rendered[key] = fmt.Sprintf("[%s]", strings.Join(stringed, ","))
		case interface{}:
			//since we can determine specifics, relay on toString() to handle
			// case, worst case will lead to using json.Marshal.
			rendered[key] = toString(v)
		}
	}

	return rendered
}

// toString returns a interface value as a string.
func toString(val interface{}) string {
	switch v := val.(type) {
	case string:
		return v
	case int:
		return strconv.FormatInt(int64(v), 10)
	case uint:
		return strconv.FormatUint(uint64(v), 10)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case time.Time:
		return dateToISO(v)
	}

	// attempt a json marshal
	json, err := json.Marshal(val)
	if err != nil {
		return ""
	}

	return string(json)
}

// objectIDToString returns a hex string version of a bson.ObjectId.
// wrap: determines if it returns a simple hex string or wraps the hex string in
// a "ObjectId(%s)" string.
func objectIDToString(id bson.ObjectId, wrap bool) string {
	if wrap {
		return "\"ObjectId('" + id.Hex() + "')\""
	}
	return "\"" + id.Hex() + "\""
}

// dateToISO returns a time.Time as a ISO string
func dateToISO(t time.Time) string {
	return t.Format("2006-01-02T15:04:05.000Z")
}

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
		// TODO: do we need to perform some transformation for dates and ObjectIDs?
		script = strings.Replace(script, match, val, 1)
	}

	return script
}
