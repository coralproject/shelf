
# submission
    import "github.com/coralproject/shelf/internal/ask/form/submission"




## Constants
``` go
const Collection = "form_submissions"
```
Collection is the mongo collection where Submission
documents are saved.


## Variables
``` go
var ErrInvalidID = errors.New("ID is not in it's proper form")
```
ErrInvalidID occurs when an ID is not in a valid form.


## func Count
``` go
func Count(context interface{}, db *db.DB, formID string) (int, error)
```
Count returns the count of current submissions for a given
form id in the Form Submissions MongoDB database collection.


## func Create
``` go
func Create(context interface{}, db *db.DB, formID string, submission *Submission) error
```
Create adds a new Submission based on a given Form into
the MongoDB database collection.


## func Delete
``` go
func Delete(context interface{}, db *db.DB, id string) error
```
Delete removes a given Form Submission from the MongoDB
database collection.


## func EnsureIndexes
``` go
func EnsureIndexes(context interface{}, db *db.DB) error
```
EnsureIndexes perform index create commands against Mongo for the indexes
needed for the ask package to run.


## func RetrieveMany
``` go
func RetrieveMany(context interface{}, db *db.DB, ids []string) ([]Submission, error)
```
RetrieveMany retrieves a list of Submission's from the MongoDB database collection.



## type Answer
``` go
type Answer struct {
    WidgetID     string      `json:"widget_id" bson:"widget_id" validate:"required,len=24"`
    Identity     bool        `json:"identity" bson:"identity"`
    Answer       interface{} `json:"answer" bson:"answer"`
    EditedAnswer interface{} `json:"edited" bson:"edited"`
    Question     interface{} `json:"question" bson:"question"`
    Props        interface{} `json:"props" bson:"props"`
}
```
Answer describes an answer submitted for a specific Form widget
with the specific question asked included as well.











## type AnswerInput
``` go
type AnswerInput struct {
    WidgetID string      `json:"widget_id" validate:"required"`
    Answer   interface{} `json:"answer" validate:"exists"`
}
```
AnswerInput describes the input accepted for a new submission
answer.











### func (\*AnswerInput) Validate
``` go
func (f *AnswerInput) Validate() error
```
Validate checks the AnswerInput value for consistency.



## type SearchOpts
``` go
type SearchOpts struct {
    DscOrder bool
    Query    string
    FilterBy string
}
```
SearchOpts is the options used to perform a search accross a
given forms submissions.











## type SearchResultCounts
``` go
type SearchResultCounts struct {
    SearchByFlag     map[string]int `json:"search_by_flag"`
    TotalSearch      int            `json:"total_search"`
    TotalSubmissions int            `json:"total_submissions"`
}
```
SearchResultCounts is a structured type containing the counts of results.











## type SearchResults
``` go
type SearchResults struct {
    Counts      SearchResultCounts `json:"counts"`
    Submissions []Submission
}
```
SearchResults is a structured type returning the results
expected from searching for submissions based on a form id.









### func Search
``` go
func Search(context interface{}, db *db.DB, formID string, limit, skip int, opts SearchOpts) (*SearchResults, error)
```
Search searches through form submissions for a given form
using the provided search options.




## type Submission
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









### func AddFlag
``` go
func AddFlag(context interface{}, db *db.DB, id, flag string) (*Submission, error)
```
AddFlag adds, and de-duplicates a flag to a given
Submission in the MongoDB database collection.


### func RemoveFlag
``` go
func RemoveFlag(context interface{}, db *db.DB, id, flag string) (*Submission, error)
```
RemoveFlag removes a flag from a given Submission in
the MongoDB database collection.


### func Retrieve
``` go
func Retrieve(context interface{}, db *db.DB, id string) (*Submission, error)
```
Retrieve retrieves a Submission from the MongoDB database
collection.


### func UpdateAnswer
``` go
func UpdateAnswer(context interface{}, db *db.DB, id string, answer AnswerInput) (*Submission, error)
```
UpdateAnswer updates the edited answer if it could find it
inside the MongoDB database collection atomically.


### func UpdateStatus
``` go
func UpdateStatus(context interface{}, db *db.DB, id, status string) (*Submission, error)
```
UpdateStatus updates a form submissions status inside the MongoDB database
collection.




### func (\*Submission) Validate
``` go
func (s *Submission) Validate() error
```
Validate checks the Submission value for consistency.









- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)