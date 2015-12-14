package query

// Set of query types we expect to receive
const (
	TypePipeline = "pipeline"
	TypeTemplate = "template"
)

// Query contains the configuration details for a query.
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
	Return bool `bson:"return" json:"return"`

	// Indicates there is a date to be pre-processed in the scripts.
	HasDate bool `bson:"has_date,omitempty" json:"has_date,omitempty"`

	// Indicates there is an ObjectId to be pre-processed in the scripts.
	HasObjectID bool `bson:"has_objectid,omitempty" json:"has_objectid,omitempty"`

	// Scripts to process for the query.
	Scripts []string `bson:"scripts" json:"scripts"`
}

// Param contains meta-data about a required parameter for the query.
type Param struct {
	Name    string `bson:"name" json:"name"`       // Name of the parameter.
	Default string `bson:"default" json:"default"` // Default value for the parameter.
	Desc    string `bson:"desc" json:"desc"`       // Description about the parameter.
}

// Set contains the configuration details for a rule set.
type Set struct {
	Name        string  `bson:"name" json:"name"`       // Name of the query set.
	Description string  `bson:"desc" json:"desc"`       // Description of the query set.
	Enabled     bool    `bson:"enabled" json:"enabled"` // If the query set is enabled to run.
	Params      []Param `bson:"params" json:"params"`   // Collection of parameters.
	Queries     []Query `bson:"queries" json:"queries"` // Collection of queries.
}

// Result contains the result of an query set execution.
type Result struct {
	Results interface{} `json:"results"`
	Error   bool        `json:"error"`
}
