

# form
`import "github.com/coralproject/shelf/internal/ask/form"`

* [Overview](#pkg-overview)
* [Index](#pkg-index)
* [Subdirectories](#pkg-subdirectories)

## <a name="pkg-overview">Overview</a>



## <a name="pkg-index">Index</a>
* [Constants](#pkg-constants)
* [Variables](#pkg-variables)
* [func AggregateFormSubmissions(context interface{}, db *db.DB, id string) (map[string]Aggregation, error)](#AggregateFormSubmissions)
* [func Delete(context interface{}, db *db.DB, id string) error](#Delete)
* [func GroupSubmissions(context interface{}, db *db.DB, formID string, limit int, skip int, opts submission.SearchOpts) (map[Group][]submission.Submission, error)](#GroupSubmissions)
* [func List(context interface{}, db *db.DB, limit, skip int) ([]Form, error)](#List)
* [func TextAggregate(context interface{}, db *db.DB, formID string, subs []submission.Submission) ([]TextAggregation, error)](#TextAggregate)
* [func Upsert(context interface{}, db *db.DB, form *Form) error](#Upsert)
* [type Aggregation](#Aggregation)
* [type Form](#Form)
  * [func Retrieve(context interface{}, db *db.DB, id string) (*Form, error)](#Retrieve)
  * [func UpdateStatus(context interface{}, db *db.DB, id, status string) (*Form, error)](#UpdateStatus)
  * [func (f *Form) Validate() error](#Form.Validate)
* [type Group](#Group)
* [type MCAggregation](#MCAggregation)
* [type MCAnswerAggregation](#MCAnswerAggregation)
* [type Stats](#Stats)
  * [func UpdateStats(context interface{}, db *db.DB, id string) (*Stats, error)](#UpdateStats)
* [type Step](#Step)
* [type SubmissionGroup](#SubmissionGroup)
* [type TextAggregation](#TextAggregation)
* [type Widget](#Widget)


#### <a name="pkg-files">Package files</a>
[aggregation.go](/src/github.com/coralproject/shelf/internal/ask/form/aggregation.go) [form.go](/src/github.com/coralproject/shelf/internal/ask/form/form.go) 


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



## <a name="AggregateFormSubmissions">func</a> [AggregateFormSubmissions](/src/target/aggregation.go?s=1945:2049#L42)
``` go
func AggregateFormSubmissions(context interface{}, db *db.DB, id string) (map[string]Aggregation, error)
```
AggregateFormSubmissions retrieves the submissions for a form, groups them then
runs aggregations and counts for each one.



## <a name="Delete">func</a> [Delete](/src/target/form.go?s=8404:8464#L255)
``` go
func Delete(context interface{}, db *db.DB, id string) error
```
Delete removes the document matching the id provided from the MongoDB
database collection.



## <a name="GroupSubmissions">func</a> [GroupSubmissions](/src/target/aggregation.go?s=3928:4088#L107)
``` go
func GroupSubmissions(context interface{}, db *db.DB, formID string, limit int, skip int, opts submission.SearchOpts) (map[Group][]submission.Submission, error)
```
GroupSubmissions organizes submissions by Group. It looks for questions with the group by flag
and creates Group structs.



## <a name="List">func</a> [List](/src/target/form.go?s=7782:7856#L235)
``` go
func List(context interface{}, db *db.DB, limit, skip int) ([]Form, error)
```
List retrieves a list of forms from the MongodB database collection.



## <a name="TextAggregate">func</a> [TextAggregate](/src/target/aggregation.go?s=7067:7189#L233)
``` go
func TextAggregate(context interface{}, db *db.DB, formID string, subs []submission.Submission) ([]TextAggregation, error)
```
TextAggregate returns all text answers flagged with includeInGroup.



## <a name="Upsert">func</a> [Upsert](/src/target/form.go?s=3630:3691#L84)
``` go
func Upsert(context interface{}, db *db.DB, form *Form) error
```
Upsert upserts the provided form into the MongoDB database collection.




## <a name="Aggregation">type</a> [Aggregation](/src/target/aggregation.go?s=1149:1427#L24)
``` go
type Aggregation struct {
    Group Group                    `json:"group" bson:"group"`
    Count int                      `json:"count" bson:"count"`
    MC    map[string]MCAggregation `json:"MultipleChoice" bson:"MultipleChoice"` // Capitalization matches widget type MultipleChoice
}
```
Aggregation holds the various aggregations and stats collected.










## <a name="Form">type</a> [Form](/src/target/form.go?s=2109:3311#L54)
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







### <a name="Retrieve">func</a> [Retrieve](/src/target/form.go?s=7040:7111#L209)
``` go
func Retrieve(context interface{}, db *db.DB, id string) (*Form, error)
```
Retrieve retrieves the form from the MongodB database collection.


### <a name="UpdateStatus">func</a> [UpdateStatus](/src/target/form.go?s=6031:6114#L172)
``` go
func UpdateStatus(context interface{}, db *db.DB, id, status string) (*Form, error)
```
UpdateStatus updates the forms status and returns the updated form from
the MongodB database collection.





### <a name="Form.Validate">func</a> (\*Form) [Validate](/src/target/form.go?s=3364:3395#L73)
``` go
func (f *Form) Validate() error
```
Validate checks the Form value for consistency.




## <a name="Group">type</a> [Group](/src/target/aggregation.go?s=1562:1732#L32)
``` go
type Group struct {
    ID       string `json:"group_id" bson:"group_id"`
    Question string `json:"question" bson:"question"`
    Answer   string `json:"answer" bson:"answer"`
}
```
Group defines a key for a multiple choice question / answer combo to be used
to define slices of submissions to be aggregated.










## <a name="MCAggregation">type</a> [MCAggregation](/src/target/aggregation.go?s=700:879#L14)
``` go
type MCAggregation struct {
    Question  string                         `json:"question" bson:"question"`
    MCAnswers map[string]MCAnswerAggregation `json:"answers" bson:"answers"`
}
```
MCAggregation holds a multiple choice question and a map aggregated counts for
each answer. The Answers map is keyed off an md5 of the answer as not better keys exist










## <a name="MCAnswerAggregation">type</a> [MCAnswerAggregation](/src/target/aggregation.go?s=404:525#L7)
``` go
type MCAnswerAggregation struct {
    Title string `json:"answer" bson:"answer"`
    Count int    `json:"count" bson:"count"`
}
```
MCAnswerAggregation holds the count for selections of a single multiple
choice answer.










## <a name="Stats">type</a> [Stats](/src/target/form.go?s=1820:1892#L46)
``` go
type Stats struct {
    Responses int `json:"responses" bson:"responses"`
}
```
Stats describes the statistics being recorded by a specific Form.







### <a name="UpdateStats">func</a> [UpdateStats](/src/target/form.go?s=4886:4961#L129)
``` go
func UpdateStats(context interface{}, db *db.DB, id string) (*Stats, error)
```
UpdateStats updates the Stats on a given Form.





## <a name="Step">type</a> [Step](/src/target/form.go?s=1594:1749#L39)
``` go
type Step struct {
    ID      string   `json:"id" bson:"_id"`
    Name    string   `json:"name" bson:"name"`
    Widgets []Widget `json:"widgets" bson:"widgets"`
}
```
Step is a collection of Widget's.










## <a name="SubmissionGroup">type</a> [SubmissionGroup](/src/target/aggregation.go?s=3680:3798#L101)
``` go
type SubmissionGroup struct {
    Submissions map[Group][]submission.Submission `json:"submissions" bson:"submissions"`
}
```
SubmissionGroup is a transport that defines the transport structure for a submission group.










## <a name="TextAggregation">type</a> [TextAggregation](/src/target/aggregation.go?s=1042:1080#L21)
``` go
type TextAggregation map[string]string
```
TextAggregation holds the aggregated text based answers for a single question
marked with the Incude in Aggregations tag, orderd by [question_id][answer].










## <a name="Widget">type</a> [Widget](/src/target/form.go?s=1086:1555#L27)
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
