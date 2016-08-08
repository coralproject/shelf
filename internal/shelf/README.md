

# shelf
`import "github.com/coralproject/xenia/internal/shelf"`

* [Overview](#pkg-overview)
* [Index](#pkg-index)
* [Subdirectories](#pkg-subdirectories)

## <a name="pkg-overview">Overview</a>



## <a name="pkg-index">Index</a>
* [Constants](#pkg-constants)
* [Variables](#pkg-variables)
* [func AddRelationship(context interface{}, db *db.DB, rel Relationship) (string, error)](#AddRelationship)
* [func AddView(context interface{}, db *db.DB, view View) (string, error)](#AddView)
* [func ClearRelManager(context interface{}, db *db.DB) error](#ClearRelManager)
* [func NewRelManager(context interface{}, db *db.DB, rm RelManager) error](#NewRelManager)
* [func RemoveRelationship(context interface{}, db *db.DB, relID string) error](#RemoveRelationship)
* [func RemoveView(context interface{}, db *db.DB, viewID string) error](#RemoveView)
* [func UpdateRelationship(context interface{}, db *db.DB, rel Relationship) error](#UpdateRelationship)
* [func UpdateView(context interface{}, db *db.DB, view View) error](#UpdateView)
* [type PathSegment](#PathSegment)
  * [func (ps PathSegment) Validate() error](#PathSegment.Validate)
* [type RelManager](#RelManager)
  * [func GetRelManager(context interface{}, db *db.DB) (RelManager, error)](#GetRelManager)
  * [func (rm RelManager) Validate() error](#RelManager.Validate)
* [type Relationship](#Relationship)
  * [func (r Relationship) Validate() error](#Relationship.Validate)
* [type View](#View)
  * [func (v View) Validate() error](#View.Validate)


#### <a name="pkg-files">Package files</a>
[manager.go](/src/target/manager.go) [model.go](/src/target/model.go) [relationships.go](/src/target/relationships.go) [utils.go](/src/target/utils.go) [views.go](/src/target/views.go) 


## <a name="pkg-constants">Constants</a>
``` go
const Collection = "relationship_manager"
```
Collection is the MongoDB collection housing metadata about relationships and views.


## <a name="pkg-variables">Variables</a>
``` go
var (
    ErrNotFound = errors.New("Set Not found")
)
```
Set of error variables.



## <a name="AddRelationship">func</a> [AddRelationship](/src/target/relationships.go?s=273:359#L6)
``` go
func AddRelationship(context interface{}, db *db.DB, rel Relationship) (string, error)
```
AddRelationship adds a relationship to the relationship manager.



## <a name="AddView">func</a> [AddView](/src/target/views.go?s=257:328#L6)
``` go
func AddView(context interface{}, db *db.DB, view View) (string, error)
```
AddView adds a view to the relationship manager.



## <a name="ClearRelManager">func</a> [ClearRelManager](/src/target/manager.go?s=1605:1663#L45)
``` go
func ClearRelManager(context interface{}, db *db.DB) error
```
ClearRelManager clears a current relationship manager from Mongo.



## <a name="NewRelManager">func</a> [NewRelManager](/src/target/manager.go?s=486:557#L12)
``` go
func NewRelManager(context interface{}, db *db.DB, rm RelManager) error
```
NewRelManager creates a new relationship manager, either with defaults
or based on a provided JSON config.



## <a name="RemoveRelationship">func</a> [RemoveRelationship](/src/target/relationships.go?s=1542:1617#L41)
``` go
func RemoveRelationship(context interface{}, db *db.DB, relID string) error
```
RemoveRelationship removes a relationship from the relationship manager.



## <a name="RemoveView">func</a> [RemoveView](/src/target/views.go?s=1857:1925#L53)
``` go
func RemoveView(context interface{}, db *db.DB, viewID string) error
```
RemoveView removes a view from the relationship manager.



## <a name="UpdateRelationship">func</a> [UpdateRelationship](/src/target/relationships.go?s=2821:2900#L79)
``` go
func UpdateRelationship(context interface{}, db *db.DB, rel Relationship) error
```
UpdateRelationship updates a relationship in the relationship manager.



## <a name="UpdateView">func</a> [UpdateView](/src/target/views.go?s=2434:2498#L71)
``` go
func UpdateView(context interface{}, db *db.DB, view View) error
```
UpdateView updates a view in the relationship manager.




## <a name="PathSegment">type</a> [PathSegment](/src/target/model.go?s=2725:3078#L81)
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










### <a name="PathSegment.Validate">func</a> (PathSegment) [Validate](/src/target/model.go?s=3138:3176#L89)
``` go
func (ps PathSegment) Validate() error
```
Validate checks the PathSegment value for consistency.




## <a name="RelManager">type</a> [RelManager](/src/target/model.go?s=532:794#L12)
``` go
type RelManager struct {
    ID            int            `bson:"id" json:"id"`
    Relationships []Relationship `bson:"relationships" json:"relationships" validate:"required,min=1"`
    Views         []View         `bson:"views" json:"views" validate:"required,min=1"`
}
```
RelManager contains metadata about what relationships and views are currenlty
being utilized in the system.







### <a name="GetRelManager">func</a> [GetRelManager](/src/target/manager.go?s=2143:2213#L61)
``` go
func GetRelManager(context interface{}, db *db.DB) (RelManager, error)
```
GetRelManager retrieves the current relationship manager from Mongo.





### <a name="RelManager.Validate">func</a> (RelManager) [Validate](/src/target/model.go?s=853:890#L19)
``` go
func (rm RelManager) Validate() error
```
Validate checks the RelManager value for consistency.




## <a name="Relationship">type</a> [Relationship](/src/target/model.go?s=1919:2447#L61)
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










### <a name="Relationship.Validate">func</a> (Relationship) [Validate](/src/target/model.go?s=2508:2546#L71)
``` go
func (r Relationship) Validate() error
```
Validate checks the Relationship value for consistency.




## <a name="View">type</a> [View](/src/target/model.go?s=3295:3631#L97)
``` go
type View struct {
    ID        string        `bson:"id" json:"id" validate:"required,min=1"`
    Name      string        `bson:"name" json:"name" validate:"required,min=3"`
    StartType string        `bson:"start_type" json:"start_type" validate:"required,min=3"`
    Path      []PathSegment `bson:"path" json:"path" validate:"required,min=1"`
}
```
View contains metadata about a view.










### <a name="View.Validate">func</a> (View) [Validate](/src/target/model.go?s=3684:3714#L105)
``` go
func (v View) Validate() error
```
Validate checks the View value for consistency.








- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)
