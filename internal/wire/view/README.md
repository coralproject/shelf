

# view
`import "github.com/coralproject/shelf/internal/wire/view"`

* [Overview](#pkg-overview)
* [Index](#pkg-index)
* [Subdirectories](#pkg-subdirectories)

## <a name="pkg-overview">Overview</a>



## <a name="pkg-index">Index</a>
* [Constants](#pkg-constants)
* [Variables](#pkg-variables)
* [func Delete(context interface{}, db *db.DB, name string) error](#Delete)
* [func GetAll(context interface{}, db *db.DB) ([]View, error)](#GetAll)
* [func Upsert(context interface{}, db *db.DB, view *View) error](#Upsert)
* [type Path](#Path)
  * [func (path *Path) Validate() error](#Path.Validate)
* [type PathSegment](#PathSegment)
  * [func (ps *PathSegment) Validate() error](#PathSegment.Validate)
* [type PathSegments](#PathSegments)
  * [func (slice PathSegments) Len() int](#PathSegments.Len)
  * [func (slice PathSegments) Less(i, j int) bool](#PathSegments.Less)
  * [func (slice PathSegments) Swap(i, j int)](#PathSegments.Swap)
* [type View](#View)
  * [func GetByName(context interface{}, db *db.DB, name string) (*View, error)](#GetByName)
  * [func (v *View) Validate() error](#View.Validate)


#### <a name="pkg-files">Package files</a>
[model.go](/src/github.com/coralproject/shelf/internal/wire/view/model.go) [view.go](/src/github.com/coralproject/shelf/internal/wire/view/view.go) 


## <a name="pkg-constants">Constants</a>
``` go
const Collection = "views"
```
Collection is the Mongo collection containing view metadata.


## <a name="pkg-variables">Variables</a>
``` go
var ErrNotFound = errors.New("View Not found")
```
ErrNotFound is an error variable thrown when no results are returned from a Mongo query.



## <a name="Delete">func</a> [Delete](/src/target/view.go?s=2505:2567#L81)
``` go
func Delete(context interface{}, db *db.DB, name string) error
```
Delete removes a view from from Mongo.



## <a name="GetAll">func</a> [GetAll](/src/target/view.go?s=1265:1324#L36)
``` go
func GetAll(context interface{}, db *db.DB) ([]View, error)
```
GetAll retrieves the current views from Mongo.



## <a name="Upsert">func</a> [Upsert](/src/target/view.go?s=534:595#L10)
``` go
func Upsert(context interface{}, db *db.DB, view *View) error
```
Upsert upserts a view to the collection of currently utilized views.




## <a name="Path">type</a> [Path](/src/target/model.go?s=1642:1822#L47)
``` go
type Path struct {
    StrictPath bool         `bson:"strict_path" json:"strict_path"`
    Segments   PathSegments `bson:"path_segments" json:"path_segments" validate:"required,min=1"`
}
```
Path includes information defining one or multiple graph paths,
along with a boolean choice for whether or not the path is a strict graph path.










### <a name="Path.Validate">func</a> (\*Path) [Validate](/src/target/model.go?s=1882:1916#L53)
``` go
func (path *Path) Validate() error
```
Validate checks the pathsegment value for consistency.




## <a name="PathSegment">type</a> [PathSegment](/src/target/model.go?s=517:838#L12)
``` go
type PathSegment struct {
    Level     int    `bson:"level" json:"level" validate:"required,min=1"`
    Direction string `bson:"direction" json:"direction" validate:"required,min=2"`
    Predicate string `bson:"predicate" json:"predicate" validate:"required,min=1"`
    Tag       string `bson:"tag,omitempty" json:"tag,omitempty"`
}
```
PathSegment contains metadata about a segment of a path,
which path partially defines a View.










### <a name="PathSegment.Validate">func</a> (\*PathSegment) [Validate](/src/target/model.go?s=898:937#L20)
``` go
func (ps *PathSegment) Validate() error
```
Validate checks the pathsegment value for consistency.




## <a name="PathSegments">type</a> [PathSegments](/src/target/model.go?s=1066:1097#L28)
``` go
type PathSegments []PathSegment
```
PathSegments is a slice of PathSegment values.










### <a name="PathSegments.Len">func</a> (PathSegments) [Len](/src/target/model.go?s=1150:1185#L31)
``` go
func (slice PathSegments) Len() int
```
Len is required to sort a slice of PathSegment.




### <a name="PathSegments.Less">func</a> (PathSegments) [Less](/src/target/model.go?s=1262:1307#L36)
``` go
func (slice PathSegments) Less(i, j int) bool
```
Less is required to sort a slice of PathSegment.




### <a name="PathSegments.Swap">func</a> (PathSegments) [Swap](/src/target/model.go?s=1405:1445#L41)
``` go
func (slice PathSegments) Swap(i, j int)
```
Swap is required to sort a slice of PathSegment.




## <a name="View">type</a> [View](/src/target/model.go?s=2037:2446#L61)
``` go
type View struct {
    Name       string `bson:"name" json:"name" validate:"required,min=3"`
    Collection string `bson:"collection" json:"collection" validate:"required,min=2"`
    StartType  string `bson:"start_type" json:"start_type" validate:"required,min=3"`
    ReturnRoot bool   `bson:"return_root,omitempty" json:"return_root,omitempty"`
    Paths      []Path `bson:"paths" json:"paths" validate:"required,min=1"`
}
```
View contains metadata about a view.







### <a name="GetByName">func</a> [GetByName](/src/target/view.go?s=1848:1922#L58)
``` go
func GetByName(context interface{}, db *db.DB, name string) (*View, error)
```
GetByName retrieves a view by name from Mongo.





### <a name="View.Validate">func</a> (\*View) [Validate](/src/target/model.go?s=2499:2530#L70)
``` go
func (v *View) Validate() error
```
Validate checks the View value for consistency.








- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)
