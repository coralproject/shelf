
# ask
    import "github.com/coralproject/shelf/internal/ask"




## Constants
``` go
const FormCollection = "forms"
```
FormCollection is the mongo collection where Form documents are saved.

``` go
const FormGalleryCollection = "form_galleries"
```
FormGalleryCollection is the mongo collection where FormGallery documents are
saved.

``` go
const FormSubmissionsCollection = "form_submissions"
```
FormSubmissionsCollection is the mongo collection where FormSubmission
documents are saved.


## Variables
``` go
var ErrInvalidID = errors.New("ID is not in it's proper form")
```
ErrInvalidID occurs when an ID is not in a valid form.


## func CountFormSubmissions
``` go
func CountFormSubmissions(context interface{}, db *db.DB, formID string) (int, error)
```
CountFormSubmissions returns the count of current submissions for a given
form id in the Form Submissions MongoDB database collection.


## func DeleteForm
``` go
func DeleteForm(context interface{}, db *db.DB, id string) error
```
DeleteForm removes the document matching the id provided from the MongoDB
database collection.


## func DeleteFormSubmission
``` go
func DeleteFormSubmission(context interface{}, db *db.DB, id, formID string) error
```
DeleteFormSubmission removes a given Form Submission from the MongoDB
database collection.


## func HydrateFormGalleries
``` go
func HydrateFormGalleries(context interface{}, db *db.DB, galleries []FormGallery) error
```
HydrateFormGalleries loads an array of form galleries with form submissions
from the MongoDB database collection.


## func HydrateFormGallery
``` go
func HydrateFormGallery(context interface{}, db *db.DB, gallery *FormGallery) error
```
HydrateFormGallery loads a FormGallery with form submissions from the MongoDB
database collection.


## func MergeSubmissionsIntoGalleryAnswers
``` go
func MergeSubmissionsIntoGalleryAnswers(gallery *FormGallery, submissions []FormSubmission)
```
MergeSubmissionsIntoGalleryAnswers associates the array of submissions onto
matching gallery answers.


## func RetrieveFormGalleriesForForm
``` go
func RetrieveFormGalleriesForForm(context interface{}, db *db.DB, formID string) ([]FormGallery, error)
```
RetrieveFormGalleriesForForm retrives the form galleries for a given form
from the MongoDB database collection.


## func RetrieveFormSubmissions
``` go
func RetrieveFormSubmissions(context interface{}, db *db.DB, ids []string) ([]FormSubmission, error)
```
RetrieveFormSubmissions retrieves a list of FormSubmission's from the MongoDB
database collection.


## func RetrieveManyForms
``` go
func RetrieveManyForms(context interface{}, db *db.DB, limit, skip int) ([]Form, error)
```
RetrieveManyForms retrieves a list of forms from the MongodB database
collection.


## func UpdateFormGallery
``` go
func UpdateFormGallery(context interface{}, db *db.DB, id string, gallery *FormGallery) error
```
UpdateFormGallery updates the form gallery in the MongoDB database
collection.


## func Upsert
``` go
func Upsert(context interface{}, db *db.DB, form *Form) error
```
Upsert upserts the provided form into the MongoDB database collection.


## func UpsertForm
``` go
func UpsertForm(context interface{}, db *db.DB, form *Form) error
```
UpsertForm upserts the provided form into the MongoDB database collection.



## type Form
``` go
type Form struct {
    ID             bson.ObjectId `json:"id" bson:"_id" validate:"required,len=24"`
    Status         string        `json:"status" bson:"status"`
    Theme          interface{}   `json:"theme" bson:"theme"`
    Settings       interface{}   `json:"settings" bson:"settings"`
    Header         interface{}   `json:"header" bson:"header"`
    Footer         interface{}   `json:"footer" bson:"footer"`
    FinishedScreen interface{}   `json:"finishedScreen" bson:"finishedScreen"`
    Steps          []FormStep    `json:"steps" bson:"steps"`
    Stats          FormStats     `json:"stats" bson:"stats"`
    CreatedBy      interface{}   `json:"created_by" bson:"created_by"`
    UpdatedBy      interface{}   `json:"updated_by" bson:"updated_by"`
    DeletedBy      interface{}   `json:"deleted_by" bson:"deleted_by"`
    DateCreated    time.Time     `json:"date_created,omitempty" bson:"date_created,omitempty"`
    DateUpdated    time.Time     `json:"date_updated,omitempty" bson:"date_updated,omitempty"`
    DateDeleted    time.Time     `json:"date_deleted,omitempty" bson:"date_deleted,omitempty"`
}
```
Form contains the conatical representation of a Form, containing all the
Steps, and help text relating to completing the Form.









### func RetrieveForm
``` go
func RetrieveForm(context interface{}, db *db.DB, id string) (*Form, error)
```
RetrieveForm retrieves the form from the MongodB database collection.


### func UpdateFormStatus
``` go
func UpdateFormStatus(context interface{}, db *db.DB, id, status string) (*Form, error)
```
UpdateFormStatus updates the forms status and returns the updated form from
the MongodB database collection.




### func (\*Form) Validate
``` go
func (f *Form) Validate() error
```
Validate checks the Form value for consistency.



## type FormGallery
``` go
type FormGallery struct {
    ID          bson.ObjectId          `json:"id" bson:"_id" validate:"required,len=24"`
    FormID      bson.ObjectId          `json:"form_id" bson:"form_id" validate:"required,len=24"`
    Headline    string                 `json:"headline" bson:"headline"`
    Description string                 `json:"description" bson:"description"`
    Config      map[string]interface{} `json:"config" bson:"config"`
    Answers     []FormGalleryAnswer    `json:"answers" bson:"answers"`
    DateCreated time.Time              `json:"date_created,omitempty" bson:"date_created,omitempty"`
    DateUpdated time.Time              `json:"date_updated,omitempty" bson:"date_updated,omitempty"`
}
```
FormGallery is a Form that has been moved to a shared space.









### func AddFormGalleryAnswer
``` go
func AddFormGalleryAnswer(context interface{}, db *db.DB, id, submissionID, answerID string) (*FormGallery, error)
```
AddFormGalleryAnswer adds an answer to a form gallery. Duplicated answers
are de-duplicated automatically and will not return an error.


### func CreateFormGallery
``` go
func CreateFormGallery(context interface{}, db *db.DB, formID string) (*FormGallery, error)
```
CreateFormGallery adds a form gallery based on the form id provided into the
MongoDB database collection.


### func RemoveFormGalleryAnswer
``` go
func RemoveFormGalleryAnswer(context interface{}, db *db.DB, id, submissionID, answerID string) (*FormGallery, error)
```
RemoveFormGalleryAnswer adds an answer to a form gallery. Duplicated answers
are de-duplicated automatically and will not return an error.


### func RetrieveFormGallery
``` go
func RetrieveFormGallery(context interface{}, db *db.DB, id string) (*FormGallery, error)
```
RetrieveFormGallery retrieves a form gallery from the MongoDB database
collection as well as hydrating the form gallery with form submissions.




### func (\*FormGallery) Validate
``` go
func (fg *FormGallery) Validate() error
```
Validate checks the FormGallery value for consistency.



## type FormGalleryAnswer
``` go
type FormGalleryAnswer struct {
    SubmissionID    bson.ObjectId          `json:"submission_id" bson:"submission_id" validate:"required,len=24"`
    AnswerID        string                 `json:"answer_id" bson:"answer_id" validate:"required,len=24"`
    Answer          FormSubmissionAnswer   `json:"answer,omitempty" bson:"-"`
    IdentityAnswers []FormSubmissionAnswer `json:"identity_answers,omitempty" bson:"-"`
}
```
FormGalleryAnswer describes an answer from a form which has been added to a
FormGallery.











## type FormStats
``` go
type FormStats struct {
    Responses int `json:"responses" bson:"responses"`
}
```
FormStats describes the statistics being recorded by a specific Form.









### func UpdateFormStats
``` go
func UpdateFormStats(context interface{}, db *db.DB, id string) (*FormStats, error)
```
UpdateFormStats updates the FormStats on a given Form.




## type FormStep
``` go
type FormStep struct {
    ID      string       `json:"id" bson:"_id"`
    Name    string       `json:"name" bson:"name"`
    Widgets []FormWidget `json:"widgets" bson:"widgets"`
}
```
FormStep is a collection of FormWidget's.











## type FormSubmission
``` go
type FormSubmission struct {
    ID             bson.ObjectId          `json:"id" bson:"_id"`
    FormID         bson.ObjectId          `json:"form_id" bson:"form_id"`
    Number         int                    `json:"number" bson:"number"`
    Status         string                 `json:"status" bson:"status"`
    Answers        []FormSubmissionAnswer `json:"replies" bson:"replies"`
    Flags          []string               `json:"flags" bson:"flags"` // simple, flexible string flagging
    Header         interface{}            `json:"header" bson:"header"`
    Footer         interface{}            `json:"footer" bson:"footer"`
    FinishedScreen interface{}            `json:"finishedScreen" bson:"finishedScreen"`
    CreatedBy      interface{}            `json:"created_by" bson:"created_by"` // Todo, decide how to represent ownership here
    UpdatedBy      interface{}            `json:"updated_by" bson:"updated_by"` // Todo, decide how to represent ownership here
    DateCreated    time.Time              `json:"date_created,omitempty" bson:"date_created,omitempty"`
    DateUpdated    time.Time              `json:"date_updated,omitempty" bson:"date_updated,omitempty"`
}
```
FormSubmission contains all the answers submitted for a specific Form as well
as any other details about the Form that were present at the time of the Form
submission.









### func AddFlagToFormSubmission
``` go
func AddFlagToFormSubmission(context interface{}, db *db.DB, id, flag string) (*FormSubmission, error)
```
AddFlagToFormSubmission adds, and de-duplicates a flag to a given
FormSubmission in the MongoDB database collection.


### func CreateFormSubmission
``` go
func CreateFormSubmission(context interface{}, db *db.DB, formID string, answers []FormSubmissionAnswerInput) (*FormSubmission, error)
```
CreateFormSubmission adds a new FormSubmission based on a given Form into
the MongoDB database collection.


### func RemoveFlagFromFormSubmission
``` go
func RemoveFlagFromFormSubmission(context interface{}, db *db.DB, id, flag string) (*FormSubmission, error)
```
RemoveFlagFromFormSubmission removes a flag from a given FormSubmission in
the MongoDB database collection.


### func RetrieveFormSubmission
``` go
func RetrieveFormSubmission(context interface{}, db *db.DB, id string) (*FormSubmission, error)
```
RetrieveFormSubmission retrieves a FormSubmission from the MongoDB database
collection.


### func UpdateFormSubmissionAnswer
``` go
func UpdateFormSubmissionAnswer(context interface{}, db *db.DB, id, answerID string, editedAnswer interface{}) (*FormSubmission, error)
```
UpdateFormSubmissionAnswer updates the edited answer if it could find it
inside the MongoDB database collection atomically.


### func UpdateFormSubmissionStatus
``` go
func UpdateFormSubmissionStatus(context interface{}, db *db.DB, id, status string) (*FormSubmission, error)
```
UpdateFormSubmissionStatus updates a form submissions status inside the
MongoDB database collection.




## type FormSubmissionAnswer
``` go
type FormSubmissionAnswer struct {
    WidgetID     string      `json:"widget_id" bson:"widget_id" validate:"required,len=24"`
    Identity     bool        `json:"identity" bson:"identity"`
    Answer       interface{} `json:"answer" bson:"answer"`
    EditedAnswer interface{} `json:"edited" bson:"edited"`
    Question     interface{} `json:"question" bson:"question"`
    Props        interface{} `json:"props" bson:"props"`
}
```
FormSubmissionAnswer describes an answer submitted for a specific Form widget
with the specific question asked included as well.











## type FormSubmissionAnswerInput
``` go
type FormSubmissionAnswerInput struct {
    WidgetID string      `json:"widget_id" validate:"required,len=24"`
    Answer   interface{} `json:"answer" validate:"exists"`
}
```
FormSubmissionAnswerInput describes the input accepted for a new submission
answer.











### func (\*FormSubmissionAnswerInput) Validate
``` go
func (f *FormSubmissionAnswerInput) Validate() error
```
Validate checks the FormSubmissionAnswerInput value for consistency.



## type FormSubmissionSearchOpts
``` go
type FormSubmissionSearchOpts struct {
    DscOrder bool
    Query    string
    FilterBy string
}
```
FormSubmissionSearchOpts is the options used to perform a search accross a
given forms submissions.











## type FormSubmissionSearchResults
``` go
type FormSubmissionSearchResults struct {
    Counts struct {
        SearchByFlag     map[string]int `json:"search_by_flag"`
        TotalSearch      int            `json:"total_search"`
        TotalSubmissions int            `json:"total_submissions"`
    } `json:"counts"`
    Submissions []FormSubmission
}
```
FormSubmissionSearchResults is a structured type returning the results
expected from searching for submissions based on a form id.









### func SearchFormSubmissions
``` go
func SearchFormSubmissions(context interface{}, db *db.DB, formID string, limit, skip int, opts FormSubmissionSearchOpts) (*FormSubmissionSearchResults, error)
```
SearchFormSubmissions searches through form submissions for a given form
using the provided search options.




## type FormWidget
``` go
type FormWidget struct {
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
FormWidget describes a specific question being asked by the Form which is
contained within a FormStep.

















- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)