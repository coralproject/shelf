

# shelf
`import "github.com/coralproject/xenia/internal/shelf/"`

* [Overview](#pkg-overview)
* [Index](#pkg-index)
* [Subdirectories](#pkg-subdirectories)

## <a name="pkg-overview">Overview</a>



## <a name="pkg-index">Index</a>
* [Constants](#pkg-constants)
* [Variables](#pkg-variables)
* [func AddRelationship(context interface{}, db *db.DB, rel Relationship) (string, error)](#AddRelationship)
* [func AddView(context interface{}, db *db.DB, view View) (string, error)](#AddView)
* [func ClearRelsAndViews(context interface{}, db *db.DB) error](#ClearRelsAndViews)
* [func NewRelsAndViews(context interface{}, db *db.DB, rm RelsAndViews) error](#NewRelsAndViews)
* [func RemoveRelationship(context interface{}, db *db.DB, relID string) error](#RemoveRelationship)
* [func RemoveView(context interface{}, db *db.DB, viewID string) error](#RemoveView)
* [func UpdateRelationship(context interface{}, db *db.DB, rel Relationship) error](#UpdateRelationship)
* [func UpdateView(context interface{}, db *db.DB, view View) error](#UpdateView)
* [type PathSegment](#PathSegment)
  * [func (ps PathSegment) Validate() error](#PathSegment.Validate)
* [type Relationship](#Relationship)
  * [func (r Relationship) Validate() error](#Relationship.Validate)
* [type RelsAndViews](#RelsAndViews)
  * [func GetRelsAndViews(context interface{}, db *db.DB) (RelsAndViews, error)](#GetRelsAndViews)
  * [func (rm RelsAndViews) Validate() error](#RelsAndViews.Validate)
* [type View](#View)
  * [func (v View) Validate() error](#View.Validate)


#### <a name="pkg-files">Package files</a>
[model.go](/src/target/model.go) [relationships.go](/src/target/relationships.go) [shelf.go](/src/target/shelf.go) [utils.go](/src/target/utils.go) [views.go](/src/target/views.go) 


## <a name="pkg-constants">Constants</a>
``` go
const (
    // ViewCollection is the Mongo collection containing view metadata.
    ViewCollection = "views"
    // RelCollection is the Mongo collection containing relationship metadata.
    RelCollection = "relationships"
)
```

## <a name="pkg-variables">Variables</a>
``` go
var (
    ErrNotFound = errors.New("Set Not found")
)
```
Set of error variables.



## <a name="AddRelationship">func</a> [AddRelationship](/src/target/relationships.go?s=226:312#L4)
``` go
func AddRelationship(context interface{}, db *db.DB, rel Relationship) (string, error)
```
AddRelationship adds a relationship to the currently utilized relationships.



## <a name="AddView">func</a> [AddView](/src/target/views.go?s=191:262#L4)
``` go
func AddView(context interface{}, db *db.DB, view View) (string, error)
```
AddView adds a view to the current views.



## <a name="ClearRelsAndViews">func</a> [ClearRelsAndViews](/src/target/shelf.go?s=1550:1610#L55)
``` go
func ClearRelsAndViews(context interface{}, db *db.DB) error
```
ClearRelsAndViews clears a current relationships and views from Mongo.



## <a name="NewRelsAndViews">func</a> [NewRelsAndViews](/src/target/shelf.go?s=496:571#L15)
``` go
func NewRelsAndViews(context interface{}, db *db.DB, rm RelsAndViews) error
```
NewRelsAndViews creates new relationships and views, based on input JSON.



## <a name="RemoveRelationship">func</a> [RemoveRelationship](/src/target/relationships.go?s=1552:1627#L50)
``` go
func RemoveRelationship(context interface{}, db *db.DB, relID string) error
```
RemoveRelationship removes a relationship from the current relationships.



## <a name="RemoveView">func</a> [RemoveView](/src/target/views.go?s=1854:1922#L62)
``` go
func RemoveView(context interface{}, db *db.DB, viewID string) error
```
RemoveView removes a view from the current views.



## <a name="UpdateRelationship">func</a> [UpdateRelationship](/src/target/relationships.go?s=2661:2740#L88)
``` go
func UpdateRelationship(context interface{}, db *db.DB, rel Relationship) error
```
UpdateRelationship updates a relationship in the current relationships.



## <a name="UpdateView">func</a> [UpdateView](/src/target/views.go?s=2320:2384#L80)
``` go
func UpdateView(context interface{}, db *db.DB, view View) error
```
UpdateView updates a view in the current views.




## <a name="PathSegment">type</a> [PathSegment](/src/target/model.go?s=2679:3032#L80)
``` go
type PathSegment struct {
    Level          int    `bson:"level" json:"level" validate:"required,min=1"`
    Direction      string `bson:"direction" json:"direction" validate:"required,min=2"`
    RelationshipID string `bson:"relationship_id" json:"relationship_id" validate:"required,min=1"`
    Tag            string `bson:"tag,omitempty" json:"tag,omitempty"`
}
```
PathSegment contains metadata about a segment of a path,
which path partially defines a View.










### <a name="PathSegment.Validate">func</a> (PathSegment) [Validate](/src/target/model.go?s=3092:3130#L88)
``` go
func (ps PathSegment) Validate() error
```
Validate checks the PathSegment value for consistency.




## <a name="Relationship">type</a> [Relationship](/src/target/model.go?s=1873:2401#L60)
``` go
type Relationship struct {
    ID           string   `bson:"id" json:"id" validate:"required,min=1"`
    SubjectTypes []string `bson:"subject_types" json:"subject_types" validate:"required,min=1"`
    Predicate    string   `bson:"predicate" json:"predicate" validate:"required,min=2"`
    ObjectTypes  []string `bson:"object_types" json:"object_types" validate:"required,min=1"`
    InString     string   `bson:"in_string,omitempty" json:"in_string,omitempty"`
    OutString    string   `bson:"out_string,omitempty" json:"out_string,omitempty"`
}
```
Relationship contains metadata about a relationship.
Note, predicate should be unique.










### <a name="Relationship.Validate">func</a> (Relationship) [Validate](/src/target/model.go?s=2462:2500#L70)
``` go
func (r Relationship) Validate() error
```
Validate checks the Relationship value for consistency.




## <a name="RelsAndViews">type</a> [RelsAndViews](/src/target/model.go?s=534:746#L12)
``` go
type RelsAndViews struct {
    Relationships []Relationship `bson:"relationships" json:"relationships" validate:"required,min=1"`
    Views         []View         `bson:"views" json:"views" validate:"required,min=1"`
}
```
RelsAndViews contains metadata about what relationships and views are currently
being utilized in the system.







### <a name="GetRelsAndViews">func</a> [GetRelsAndViews](/src/target/shelf.go?s=2220:2294#L79)
``` go
func GetRelsAndViews(context interface{}, db *db.DB) (RelsAndViews, error)
```
GetRelsAndViews retrieves the current relationships and views from Mongo.





### <a name="RelsAndViews.Validate">func</a> (RelsAndViews) [Validate](/src/target/model.go?s=807:846#L18)
``` go
func (rm RelsAndViews) Validate() error
```
Validate checks the RelsAndViews value for consistency.




## <a name="View">type</a> [View](/src/target/model.go?s=3249:3585#L96)
``` go
type View struct {
    ID        string        `bson:"id" json:"id" validate:"required,min=1"`
    Name      string        `bson:"name" json:"name" validate:"required,min=3"`
    StartType string        `bson:"start_type" json:"start_type" validate:"required,min=3"`
    Path      []PathSegment `bson:"path" json:"path" validate:"required,min=1"`
}
```
View contains metadata about a view.










### <a name="View.Validate">func</a> (View) [Validate](/src/target/model.go?s=3638:3668#L104)
``` go
func (v View) Validate() error
```
Validate checks the View value for consistency.








- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)
