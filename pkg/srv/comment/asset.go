package comment

// Taxonomy holds all name-value pairs.
type Taxonomy struct {
	Name  string `json:"name" bson:"name"`
	Value string `json:"value" bson:"value"`
}

// Asset denotes an asset in the system e.g. an article or a blog etc.
type Asset struct {
	ID         string     `json:"id" bson:"_id"`
	SourceID   string     `json:"src_id" bson:"src_id"`
	URL        string     `json:"url" bson:"url"`
	Taxonomies []Taxonomy `json:"taxonomies" bson:"taxonomies"`
}
