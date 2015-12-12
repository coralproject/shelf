
# query
    import "github.com/coralproject/shelf/pkg/query"

Package query provides API's for managing querysets which will be used in
executing different aggregation tests against their respective data collection.

QuerySet
In query, records are required to follow specific formatting and are at this
point, only allowed to be in a json serializable format which meet the query.Set
structure.

The query set execution supports the following types:

- Pipeline


	  Pipeline query set types take advantage of MongoDB's aggregation API
	(the currently supported data backend), which allows insightful use of its
	internal query language, in providing context against data sets within the database.

QuerySet Sample:

```json
{


	"name":"spending_advice",
	"description":"tests against user spending rate and provides adequate advice on saving more",
	"enabled": true,
	"params":[
	  {
	    "name":"user_id",
	    "default":"396bc782-6ac6-4183-a671-6e75ca5989a5",
	    "desc":"provides the user_id to check against the collection"
	  }
	],
	"rules": [
	{
	  "desc":"match spending rate over 20 dollars",
	  "type":"pipeline",
	  "continue": true,
	  "script_options": {
	    "collection":"demo_user_transactions",
	    "has_date":false,
	    "has_objectid": false
	  },
	  "save_options": {
	    "save_as":"high_dollar_users",
	    "variables": true,
	    "to_json": true
	  },
	  "var_options":{},
	  "scripts":[
	    "{ \"$match\" : { \"user_id\" : \"#userId#\", \"category\" : \"gas\" }}",
	    "{ \"$group\" : { \"_id\" : { \"category\" : \"$category\" }, \"amount\" : { \"$sum\" : \"$amount\" }}}",
	    "{ \"$match\" : { \"amount\" : { \"$gt\" : 20.00} }}"
	  ]
	 }]

}
```




## Constants
``` go
const (
    TypePipeline = "pipeline"
    TypeTemplate = "template"
)
```
Set of query types we expect to receive

``` go
const (
    Collection        = "query_sets"
    CollectionHistory = "query_sets_history"
)
```
Contains the name of Mongo collections.



## func DeleteSet
``` go
func DeleteSet(context interface{}, db *db.DB, name string) error
```
DeleteSet is used to remove an existing Set document.


## func GetSetNames
``` go
func GetSetNames(context interface{}, db *db.DB) ([]string, error)
```
GetSetNames retrieves a list of rule names.


## func UmarshalMongoScript
``` go
func UmarshalMongoScript(script string, q *Query) (bson.M, error)
```
UmarshalMongoScript converts a JSON Mongo commands into a BSON map.


## func UpsertSet
``` go
func UpsertSet(context interface{}, db *db.DB, qs *Set) error
```
UpsertSet is used to create or update an existing Set document.



## type Query
``` go
type Query struct {
    // Unique name per Set where results are stored.
    Name string `bson:"name" json:"name"`

    // Description of this specific query.
    Description string `bson:"desc,omitempty" json:"desc,omitempty"`

    // TypePipeline, TypeTemplate
    Type string `bson:"type" json:"type"`

    // Name of the collection to use for processing the query.
    Collection string `bson:"collection,omitempty" json:"collection,omitempty"`

    // Indicates that on failure to process the next query.
    Continue bool `bson:"continue,omitempty" json:"continue,omitempty"`

    // Return the results back to the user with Name as the key.
    Return bool `bson:"save" json:"save"`

    // Indicates there is a date to be pre-processed in the scripts.
    HasDate bool `bson:"has_date,omitempty" json:"has_date,omitempty"`

    // Indicates there is an ObjectId to be pre-processed in the scripts.
    HasObjectID bool `bson:"has_objectid,omitempty" json:"has_objectid,omitempty"`

    // Scripts to process for the query.
    Scripts []string `bson:"scripts" json:"scripts"`
}
```
Query contains the configuration details for a query.











## type Result
``` go
type Result struct {
    Results interface{} `json:"results"`
    Error   bool        `json:"error"`
}
```
Result contains the result of an query set execution.









### func ExecuteSet
``` go
func ExecuteSet(context interface{}, db *db.DB, set *Set, vars map[string]string) *Result
```
ExecuteSet executes the specified query set by name.




## type Set
``` go
type Set struct {
    Name        string     `bson:"name" json:"name"`       // Name of the query set.
    Description string     `bson:"desc" json:"desc"`       // Description of the query set.
    Enabled     bool       `bson:"enabled" json:"enabled"` // If the query set is enabled to run.
    Params      []SetParam `bson:"params" json:"params"`   // Collection of parameters.
    Queries     []Query    `bson:"queries" json:"queries"` // Collection of queries.
}
```
Set contains the configuration details for a rule set.









### func GetLastSetHistoryByName
``` go
func GetLastSetHistoryByName(context interface{}, db *db.DB, name string) (*Set, error)
```
GetLastSetHistoryByName gets the last written Set within the query_history
collection and returns the last one else returns a non-nil error if it fails.


### func GetSetByName
``` go
func GetSetByName(context interface{}, db *db.DB, name string) (*Set, error)
```
GetSetByName retrieves the configuration for the specified Set.




## type SetParam
``` go
type SetParam struct {
    Name    string `bson:"name" json:"name"`       // Name of the parameter.
    Default string `bson:"default" json:"default"` // Default value for the parameter.
    Desc    string `bson:"desc" json:"desc"`       // Description about the parameter.
}
```
SetParam contains meta-data about a required parameter for the query.

















- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)