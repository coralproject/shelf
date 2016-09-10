

# gallery
`import "github.com/coralproject/shelf/internal/ask/form/gallery"`

* [Overview](#pkg-overview)
* [Index](#pkg-index)
* [Subdirectories](#pkg-subdirectories)

## <a name="pkg-overview">Overview</a>



## <a name="pkg-index">Index</a>
* [Constants](#pkg-constants)
* [Variables](#pkg-variables)
* [func Create(context interface{}, db *db.DB, gallery *Gallery) error](#Create)
* [func Delete(context interface{}, db *db.DB, id string) error](#Delete)
* [func List(context interface{}, db *db.DB, formID string) ([]Gallery, error)](#List)
* [func Update(context interface{}, db *db.DB, id string, gallery *Gallery) error](#Update)
* [type Answer](#Answer)
  * [func (a *Answer) Validate() error](#Answer.Validate)
* [type Gallery](#Gallery)
  * [func AddAnswer(context interface{}, db *db.DB, id, submissionID, answerID string) (*Gallery, error)](#AddAnswer)
  * [func RemoveAnswer(context interface{}, db *db.DB, id, submissionID, answerID string) (*Gallery, error)](#RemoveAnswer)
  * [func Retrieve(context interface{}, db *db.DB, id string) (*Gallery, error)](#Retrieve)
  * [func (fg *Gallery) Validate() error](#Gallery.Validate)


#### <a name="pkg-files">Package files</a>
[gallery.go](/src/github.com/coralproject/shelf/internal/ask/form/gallery/gallery.go) 


## <a name="pkg-constants">Constants</a>
``` go
const Collection = "form_galleries"
```
Collection is the mongo collection where Gallery documents are
saved.


## <a name="pkg-variables">Variables</a>
``` go
var ErrInvalidID = errors.New("ID is not in it's proper form")
```
ErrInvalidID occurs when an ID is not in a valid form.



## <a name="Create">func</a> [Create](/src/target/gallery.go?s=2584:2651#L67)
``` go
func Create(context interface{}, db *db.DB, gallery *Gallery) error
```
Create adds a form gallery based on the form id provided into the
MongoDB database collection.



## <a name="Delete">func</a> [Delete](/src/target/gallery.go?s=12217:12277#L397)
``` go
func Delete(context interface{}, db *db.DB, id string) error
```
Delete removes the given Gallery with the ID provided.



## <a name="List">func</a> [List](/src/target/gallery.go?s=10432:10507#L330)
``` go
func List(context interface{}, db *db.DB, formID string) ([]Gallery, error)
```
List retrives the form galleries for a given form from the MongoDB database
collection.



## <a name="Update">func</a> [Update](/src/target/gallery.go?s=11348:11426#L365)
``` go
func Update(context interface{}, db *db.DB, id string, gallery *Gallery) error
```
Update updates the form gallery in the MongoDB database
collection.




## <a name="Answer">type</a> [Answer](/src/target/gallery.go?s=1037:1420#L28)
``` go
type Answer struct {
    SubmissionID    bson.ObjectId       `json:"submission_id" bson:"submission_id" validate:"required"`
    AnswerID        string              `json:"answer_id" bson:"answer_id" validate:"required"`
    Answer          submission.Answer   `json:"answer,omitempty" bson:"-" validate:"-"`
    IdentityAnswers []submission.Answer `json:"identity_answers,omitempty" bson:"-"`
}
```
Answer describes an answer from a form which has been added to a
Gallery.










### <a name="Answer.Validate">func</a> (\*Answer) [Validate](/src/target/gallery.go?s=1474:1507#L36)
``` go
func (a *Answer) Validate() error
```
Validate checks the Anser value for consistency.




## <a name="Gallery">type</a> [Gallery](/src/target/gallery.go?s=1646:2312#L45)
``` go
type Gallery struct {
    ID          bson.ObjectId          `json:"id" bson:"_id" validate:"required"`
    FormID      bson.ObjectId          `json:"form_id" bson:"form_id" validate:"required"`
    Headline    string                 `json:"headline" bson:"headline"`
    Description string                 `json:"description" bson:"description"`
    Config      map[string]interface{} `json:"config" bson:"config"`
    Answers     []Answer               `json:"answers" bson:"answers"`
    DateCreated time.Time              `json:"date_created,omitempty" bson:"date_created,omitempty"`
    DateUpdated time.Time              `json:"date_updated,omitempty" bson:"date_updated,omitempty"`
}
```
Gallery is a Form that has been moved to a shared space.







### <a name="AddAnswer">func</a> [AddAnswer](/src/target/gallery.go?s=7659:7758#L226)
``` go
func AddAnswer(context interface{}, db *db.DB, id, submissionID, answerID string) (*Gallery, error)
```
AddAnswer adds an answer to a form gallery. Duplicated answers
are de-duplicated automatically and will not return an error.


### <a name="RemoveAnswer">func</a> [RemoveAnswer](/src/target/gallery.go?s=9054:9156#L278)
``` go
func RemoveAnswer(context interface{}, db *db.DB, id, submissionID, answerID string) (*Gallery, error)
```
RemoveAnswer adds an answer to a form gallery. Duplicated answers
are de-duplicated automatically and will not return an error.


### <a name="Retrieve">func</a> [Retrieve](/src/target/gallery.go?s=3375:3449#L92)
``` go
func Retrieve(context interface{}, db *db.DB, id string) (*Gallery, error)
```
Retrieve retrieves a form gallery from the MongoDB database
collection as well as hydrating the form gallery with form submissions.





### <a name="Gallery.Validate">func</a> (\*Gallery) [Validate](/src/target/gallery.go?s=2368:2403#L57)
``` go
func (fg *Gallery) Validate() error
```
Validate checks the Gallery value for consistency.








- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)
