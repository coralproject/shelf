package query

// ScriptOption contains options for processing the scripts.
type ScriptOption struct {
	Collection  string `bson:"collection,omitempty" json:"collection,omitempty"`     // Name of the collection to use for processing the query.
	HasDate     bool   `bson:"has_date,omitempty" json:"has_date,omitempty"`         // Indicates there is a date to be pre-processed in the scripts.
	HasObjectID bool   `bson:"has_objectid,omitempty" json:"has_objectid,omitempty"` // Indicates there is an ObjectId to be pre-processed in the scripts.
}

// SaveOption contains options for saving results.
type SaveOption struct {
	SaveAs    string `bson:"save_as,omitempty" json:"save_as,omitempty"`     // Name of the memory variable to store the result into.
	Variables bool   `bson:"variables,omitempty" json:"variables,omitempty"` // Indicates if the result should be saved into the variables.
	ToJSON    bool   `bson:"to_json,omitempty" json:"to_json,omitempty"`     // Convert the string result to JSON. Template oriented.
}

// VarOption contains options for processing variables.
type VarOption struct {
	ObjectID bool `bson:"object_id,omitempty" json:"object_id,omitempty"` // Indicates to save ObjectId values with ObjectId tag.
}

// Query contains the configuration details for a query.
// Options use a pointer so they can be excluded when not in use.
type Query struct {
	Description   string        `bson:"desc,omitempty" json:"desc,omitempty"`                     // Description of this specific query.
	Type          string        `bson:"type" json:"type"`                                         // variable, inventory, pipeline, template
	Continue      bool          `bson:"continue,omitempty" json:"continue,omitempty"`             // Indicates that on failure to process the next query.
	ScriptOptions *ScriptOption `bson:"script_options,omitempty" json:"script_options,omitempty"` // Options associated with script processing.
	SaveOptions   *SaveOption   `bson:"save_options,omitempty" json:"save_options,omitempty"`     // Options associated with saving the result.
	VarOptions    *VarOption    `bson:"var_options,omitempty" json:"var_options,omitempty"`       // Options associated with variable processing.
	Scripts       []string      `bson:"scripts" json:"scripts"`                                   // Scripts to process for the query.
}

// SetParam contains meta-data about a required parameter for the query.
type SetParam struct {
	Name    string `bson:"name" json:"name"`       // Name of the parameter.
	Default string `bson:"default" json:"default"` // Default value for the parameter.
	Desc    string `bson:"desc" json:"desc"`       // Description about the parameter.
}

// Set contains the configuration details for a rule set.
type Set struct {
	Name        string     `bson:"name" json:"name"`       // Name of the query set.
	Description string     `bson:"desc" json:"desc"`       // Description of the query set.
	Enabled     bool       `bson:"enabled" json:"enabled"` // If the query set is enabled to run.
	Params      []SetParam `bson:"params" json:"params"`   // Collection of parameters.
	Queries     []Query    `bson:"queries" json:"queries"` // Collection of queries.
}

// Result contains the result of an query set execution.
type Result struct {
	FeedName   string      `json:"feed_name"`
	Collection string      `json:"collection"`
	QueryType  string      `json:"query_type"`
	Results    interface{} `json:"results"`
	Valid      bool        `json:"valid"`
	Error      bool        `json:"-"`
}
