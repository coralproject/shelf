package comment

type Taxonomy struct {
	Name       	string 			`json:"name" bson:"name"`
	Value    	string 			`json:"value" bson:"value"`
}
type TaxonomyArray []Taxonomy

type Asset struct {
	Id       	string 			`json:"id" bson:"_id"`
	Url    		string 			`json:"url" bson:"url"`
	Taxonomy 	TaxonomyArray 	`json:"taxonomy" bson:"taxonomy"`
}

type AssetArray []Asset
