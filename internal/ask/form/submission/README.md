

# submission
`import "github.com/coralproject/shelf/internal/ask/form/submission"`

* [Overview](#pkg-overview)
* [Index](#pkg-index)
* [Subdirectories](#pkg-subdirectories)

## <a name="pkg-overview">Overview</a>



## <a name="pkg-index">Index</a>
* [Constants](#pkg-constants)
* [Variables](#pkg-variables)
* [func Count(context interface{}, db *db.DB, formID string) (int, error)](#Count)
* [func Create(context interface{}, db *db.DB, formID string, submission *Submission) error](#Create)
* [func Delete(context interface{}, db *db.DB, id string) error](#Delete)
* [func EnsureIndexes(context interface{}, db *db.DB) error](#EnsureIndexes)
* [func RetrieveMany(context interface{}, db *db.DB, ids []string) ([]Submission, error)](#RetrieveMany)
* [type Answer](#Answer)
* [type AnswerInput](#AnswerInput)
  * [func (f *AnswerInput) Validate() error](#AnswerInput.Validate)
* [type SearchOpts](#SearchOpts)
* [type SearchResultCounts](#SearchResultCounts)
* [type SearchResults](#SearchResults)
  * [func Search(context interface{}, db *db.DB, formID string, limit, skip int, opts SearchOpts) (*SearchResults, error)](#Search)
* [type Submission](#Submission)
  * [func AddFlag(context interface{}, db *db.DB, id, flag string) (*Submission, error)](#AddFlag)
  * [func RemoveFlag(context interface{}, db *db.DB, id, flag string) (*Submission, error)](#RemoveFlag)
  * [func Retrieve(context interface{}, db *db.DB, id string) (*Submission, error)](#Retrieve)
  * [func UpdateAnswer(context interface{}, db *db.DB, id string, answer AnswerInput) (*Submission, error)](#UpdateAnswer)
  * [func UpdateStatus(context interface{}, db *db.DB, id, status string) (*Submission, error)](#UpdateStatus)
  * [func (s *Submission) Validate() error](#Submission.Validate)


#### <a name="pkg-files">Package files</a>
[submission.go](/src/github.com/coralproject/shelf/internal/ask/form/submission/submission.go) 


## <a name="pkg-constants">Constants</a>
``` go
const Collection = "form_submissions"
```
Collection is the mongo collection where Submission
documents are saved.


## <a name="pkg-variables">Variables</a>
``` go
var ErrInvalidID = errors.New("ID is not in it's proper form")
```
ErrInvalidID occurs when an ID is not in a valid form.



## <a name="Count">func</a> [Count](/src/target/submission.go?s=10182:10252#L321)
``` go
func Count(context interface{}, db *db.DB, formID string) (int, error)
```
Count returns the count of current submissions for a given
form id in the Form Submissions MongoDB database collection.



## <a name="Create">func</a> [Create](/src/target/submission.go?s=4873:4961#L134)
``` go
func Create(context interface{}, db *db.DB, formID string, submission *Submission) error
```
Create adds a new Submission based on a given Form into
the MongoDB database collection.



## <a name="Delete">func</a> [Delete](/src/target/submission.go?s=16282:16342#L545)
``` go
func Delete(context interface{}, db *db.DB, id string) error
```
Delete removes a given Form Submission from the MongoDB
database collection.



## <a name="EnsureIndexes">func</a> [EnsureIndexes](/src/target/submission.go?s=1076:1132#L29)
``` go
func EnsureIndexes(context interface{}, db *db.DB) error
```
EnsureIndexes perform index create commands against Mongo for the indexes
needed for the ask package to run.



## <a name="RetrieveMany">func</a> [RetrieveMany](/src/target/submission.go?s=6716:6801#L197)
``` go
func RetrieveMany(context interface{}, db *db.DB, ids []string) ([]Submission, error)
```
RetrieveMany retrieves a list of Submission's from the MongoDB database collection.




## <a name="Answer">type</a> [Answer](/src/target/submission.go?s=3002:3399#L95)
``` go
type Answer struct {
    WidgetID     string      `json:"widget_id" bson:"widget_id" validate:"required,len=24"`
    Identity     bool        `json:"identity" bson:"identity"`
    Answer       interface{} `json:"answer" bson:"answer"`
    EditedAnswer interface{} `json:"edited" bson:"edited"`
    Question     string      `json:"question" bson:"question"`
    Props        interface{} `json:"props" bson:"props"`
}
```
Answer describes an answer submitted for a specific Form widget
with the specific question asked included as well.










## <a name="AnswerInput">type</a> [AnswerInput](/src/target/submission.go?s=2560:2704#L79)
``` go
type AnswerInput struct {
    WidgetID string      `json:"widget_id" validate:"required"`
    Answer   interface{} `json:"answer" validate:"exists"`
}
```
AnswerInput describes the input accepted for a new submission
answer.










### <a name="AnswerInput.Validate">func</a> (\*AnswerInput) [Validate](/src/target/submission.go?s=2764:2802#L85)
``` go
func (f *AnswerInput) Validate() error
```
Validate checks the AnswerInput value for consistency.




## <a name="SearchOpts">type</a> [SearchOpts](/src/target/submission.go?s=2407:2482#L71)
``` go
type SearchOpts struct {
    DscOrder bool
    Query    string
    FilterBy string
}
```
SearchOpts is the options used to perform a search accross a
given forms submissions.










## <a name="SearchResultCounts">type</a> [SearchResultCounts](/src/target/submission.go?s=1803:2009#L55)
``` go
type SearchResultCounts struct {
    SearchByFlag     map[string]int `json:"search_by_flag"`
    TotalSearch      int            `json:"total_search"`
    TotalSubmissions int            `json:"total_submissions"`
}
```
SearchResultCounts is a structured type containing the counts of results.










## <a name="SearchResults">type</a> [SearchResults](/src/target/submission.go?s=2134:2313#L63)
``` go
type SearchResults struct {
    Counts      SearchResultCounts `json:"counts"`
    Submissions []Submission       `json:"submissions"`
    CSVURL      string             `json:"csv_url"`
}
```
SearchResults is a structured type returning the results
expected from searching for submissions based on a form id.







### <a name="Search">func</a> [Search](/src/target/submission.go?s=11021:11137#L354)
``` go
func Search(context interface{}, db *db.DB, formID string, limit, skip int, opts SearchOpts) (*SearchResults, error)
```
Search searches through form submissions for a given form
using the provided search options.





## <a name="Submission">type</a> [Submission](/src/target/submission.go?s=3574:4603#L107)
``` go
type Submission struct {
    ID             bson.ObjectId `json:"id" bson:"_id"`
    FormID         bson.ObjectId `json:"form_id" bson:"form_id"`
    Number         int           `json:"number" bson:"number"`
    Status         string        `json:"status" bson:"status"`
    Answers        []Answer      `json:"replies" bson:"replies"`
    Flags          []string      `json:"flags" bson:"flags"` // simple, flexible string flagging
    Header         interface{}   `json:"header" bson:"header"`
    Footer         interface{}   `json:"footer" bson:"footer"`
    FinishedScreen interface{}   `json:"finishedScreen" bson:"finishedScreen"`
    CreatedBy      interface{}   `json:"created_by" bson:"created_by"` // Todo, decide how to represent ownership here
    UpdatedBy      interface{}   `json:"updated_by" bson:"updated_by"` // Todo, decide how to represent ownership here
    DateCreated    time.Time     `json:"date_created,omitempty" bson:"date_created,omitempty"`
    DateUpdated    time.Time     `json:"date_updated,omitempty" bson:"date_updated,omitempty"`
}
```
Submission contains all the answers submitted for a specific Form as well
as any other details about the Form that were present at the time of the Form
submission.







### <a name="AddFlag">func</a> [AddFlag](/src/target/submission.go?s=14308:14390#L471)
``` go
func AddFlag(context interface{}, db *db.DB, id, flag string) (*Submission, error)
```
AddFlag adds, and de-duplicates a flag to a given
Submission in the MongoDB database collection.


### <a name="RemoveFlag">func</a> [RemoveFlag](/src/target/submission.go?s=15291:15376#L508)
``` go
func RemoveFlag(context interface{}, db *db.DB, id, flag string) (*Submission, error)
```
RemoveFlag removes a flag from a given Submission in
the MongoDB database collection.


### <a name="Retrieve">func</a> [Retrieve](/src/target/submission.go?s=5925:6002#L171)
``` go
func Retrieve(context interface{}, db *db.DB, id string) (*Submission, error)
```
Retrieve retrieves a Submission from the MongoDB database
collection.


### <a name="UpdateAnswer">func</a> [UpdateAnswer](/src/target/submission.go?s=8740:8841#L271)
``` go
func UpdateAnswer(context interface{}, db *db.DB, id string, answer AnswerInput) (*Submission, error)
```
UpdateAnswer updates the edited answer if it could find it
inside the MongoDB database collection atomically.


### <a name="UpdateStatus">func</a> [UpdateStatus](/src/target/submission.go?s=7660:7749#L233)
``` go
func UpdateStatus(context interface{}, db *db.DB, id, status string) (*Submission, error)
```
UpdateStatus updates a form submissions status inside the MongoDB database
collection.





### <a name="Submission.Validate">func</a> (\*Submission) [Validate](/src/target/submission.go?s=4662:4699#L124)
``` go
func (s *Submission) Validate() error
```
Validate checks the Submission value for consistency.








- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)
