

# form
`import "github.com/coralproject/shelf/internal/ask/form"`

* [Overview](#pkg-overview)
* [Index](#pkg-index)
* [Subdirectories](#pkg-subdirectories)

## <a name="pkg-overview">Overview</a>



## <a name="pkg-index">Index</a>
* [Constants](#pkg-constants)
* [Variables](#pkg-variables)
* [func Delete(context interface{}, db *db.DB, id string) error](#Delete)
* [func List(context interface{}, db *db.DB, limit, skip int) ([]Form, error)](#List)
* [func Upsert(context interface{}, db *db.DB, form *Form) error](#Upsert)
* [type Form](#Form)
  * [func Retrieve(context interface{}, db *db.DB, id string) (*Form, error)](#Retrieve)
  * [func UpdateStatus(context interface{}, db *db.DB, id, status string) (*Form, error)](#UpdateStatus)
  * [func (f *Form) Validate() error](#Form.Validate)
* [type Stats](#Stats)
  * [func UpdateStats(context interface{}, db *db.DB, id string) (*Stats, error)](#UpdateStats)
* [type Step](#Step)
* [type Widget](#Widget)


#### <a name="pkg-files">Package files</a>
[form.go](/src/github.com/coralproject/shelf/internal/ask/form/form.go) 


## <a name="pkg-constants">Constants</a>
``` go
const Collection = "forms"
```
Collection is the mongo collection where Form documents are saved.


## <a name="pkg-variables">Variables</a>
``` go
var ErrInvalidID = errors.New("ID is not in it's proper form")
```
ErrInvalidID occurs when an ID is not in a valid form.



## <a name="Delete">func</a> [Delete](/src/target/form.go?s=8245:8305#L250)
``` go
func Delete(context interface{}, db *db.DB, id string) error
```
Delete removes the document matching the id provided from the MongoDB
database collection.



## <a name="List">func</a> [List](/src/target/form.go?s=7623:7697#L230)
``` go
func List(context interface{}, db *db.DB, limit, skip int) ([]Form, error)
```
List retrieves a list of forms from the MongodB database collection.



## <a name="Upsert">func</a> [Upsert](/src/target/form.go?s=3584:3645#L84)
``` go
func Upsert(context interface{}, db *db.DB, form *Form) error
```
Upsert upserts the provided form into the MongoDB database collection.




## <a name="Form">type</a> [Form](/src/target/form.go?s=2063:3265#L54)
``` go
type Form struct {
    ID             bson.ObjectId          `json:"id" bson:"_id" validate:"required"`
    Status         string                 `json:"status" bson:"status"`
    Theme          interface{}            `json:"theme" bson:"theme"`
    Settings       map[string]interface{} `json:"settings" bson:"settings"`
    Header         interface{}            `json:"header" bson:"header"`
    Footer         interface{}            `json:"footer" bson:"footer"`
    FinishedScreen interface{}            `json:"finishedScreen" bson:"finishedScreen"`
    Steps          []Step                 `json:"steps" bson:"steps"`
    Stats          Stats                  `json:"stats" bson:"stats"`
    CreatedBy      interface{}            `json:"created_by" bson:"created_by"`
    UpdatedBy      interface{}            `json:"updated_by" bson:"updated_by"`
    DeletedBy      interface{}            `json:"deleted_by" bson:"deleted_by"`
    DateCreated    time.Time              `json:"date_created,omitempty" bson:"date_created,omitempty"`
    DateUpdated    time.Time              `json:"date_updated,omitempty" bson:"date_updated,omitempty"`
    DateDeleted    time.Time              `json:"date_deleted,omitempty" bson:"date_deleted,omitempty"`
}
```
Form contains the conatical representation of a Form, containing all the
Steps, and help text relating to completing the Form.







### <a name="Retrieve">func</a> [Retrieve](/src/target/form.go?s=6881:6952#L204)
``` go
func Retrieve(context interface{}, db *db.DB, id string) (*Form, error)
```
Retrieve retrieves the form from the MongodB database collection.


### <a name="UpdateStatus">func</a> [UpdateStatus](/src/target/form.go?s=5872:5955#L167)
``` go
func UpdateStatus(context interface{}, db *db.DB, id, status string) (*Form, error)
```
UpdateStatus updates the forms status and returns the updated form from
the MongodB database collection.





### <a name="Form.Validate">func</a> (\*Form) [Validate](/src/target/form.go?s=3318:3349#L73)
``` go
func (f *Form) Validate() error
```
Validate checks the Form value for consistency.




## <a name="Stats">type</a> [Stats](/src/target/form.go?s=1774:1846#L46)
``` go
type Stats struct {
    Responses int `json:"responses" bson:"responses"`
}
```
Stats describes the statistics being recorded by a specific Form.







### <a name="UpdateStats">func</a> [UpdateStats](/src/target/form.go?s=4776:4851#L126)
``` go
func UpdateStats(context interface{}, db *db.DB, id string) (*Stats, error)
```
UpdateStats updates the Stats on a given Form.





## <a name="Step">type</a> [Step](/src/target/form.go?s=1548:1703#L39)
``` go
type Step struct {
    ID      string   `json:"id" bson:"_id"`
    Name    string   `json:"name" bson:"name"`
    Widgets []Widget `json:"widgets" bson:"widgets"`
}
```
Step is a collection of Widget's.










## <a name="Widget">type</a> [Widget](/src/target/form.go?s=1040:1509#L27)
``` go
type Widget struct {
    ID          string      `json:"id" bson:"_id"`
    Type        string      `json:"type" bson:"type"`
    Identity    bool        `json:"identity" bson:"identity"`
    Component   string      `json:"component" bson:"component"`
    Title       string      `json:"title" bson:"title"`
    Description string      `json:"description" bson:"description"`
    Wrapper     interface{} `json:"wrapper" bson:"wrapper"`
    Props       interface{} `json:"props" bson:"props"`
}
```
Widget describes a specific question being asked by the Form which is
contained within a Step.














- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)
