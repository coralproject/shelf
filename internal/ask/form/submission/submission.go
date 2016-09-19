package submission

import (
	"errors"
	"strings"
	"time"

	"github.com/ardanlabs/kit/db"
	"github.com/ardanlabs/kit/db/mongo"
	"github.com/ardanlabs/kit/log"

	validator "gopkg.in/bluesuncorp/validator.v8"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//==============================================================================

// validate is used to perform model field validation.
var validate *validator.Validate

func init() {
	validate = validator.New(&validator.Config{TagName: "validate"})
}

//==============================================================================

// ErrInvalidID occurs when an ID is not in a valid form.
var ErrInvalidID = errors.New("ID is not in it's proper form")

//==============================================================================

// Collection is the mongo collection where Submission
// documents are saved.
const Collection = "form_submissions"

// EnsureIndexes perform index create commands against Mongo for the indexes
// needed for the ask package to run.
func EnsureIndexes(context interface{}, db *db.DB) error {
	log.Dev(context, "EnsureIndexes", "Started")

	f := func(c *mgo.Collection) error {
		index := mgo.Index{
			Key:        []string{"$text:$**"},
			Unique:     false,
			DropDups:   false,
			Background: false,
			Sparse:     true,
			Name:       "$**_text",
		}
		log.Dev(context, "EnsureIndexes", "MGO : db.%s.ensureIndex(%s)", c.Name, mongo.Query(index))
		return c.EnsureIndex(index)
	}

	if err := db.ExecuteMGO(context, Collection, f); err != nil {
		log.Error(context, "EnsureIndexes", err, "Completed")
		return err
	}

	log.Dev(context, "EnsureIndexes", "Completed")
	return nil
}

// SearchResultCounts is a structured type containing the counts of results.
type SearchResultCounts struct {
	SearchByFlag     map[string]int `json:"search_by_flag"`
	TotalSearch      int            `json:"total_search"`
	TotalSubmissions int            `json:"total_submissions"`
}

// SearchResults is a structured type returning the results
// expected from searching for submissions based on a form id.
type SearchResults struct {
	Counts      SearchResultCounts `json:"counts"`
	Submissions []Submission       `json:"submissions"`
	CSVURL      string             `json:"csv_url"`
}

// SearchOpts is the options used to perform a search accross a
// given forms submissions.
type SearchOpts struct {
	DscOrder bool
	Query    string
	FilterBy string
}

// AnswerInput describes the input accepted for a new submission
// answer.
type AnswerInput struct {
	WidgetID string      `json:"widget_id" validate:"required"`
	Answer   interface{} `json:"answer" validate:"exists"`
}

// Validate checks the AnswerInput value for consistency.
func (f *AnswerInput) Validate() error {
	if err := validate.Struct(f); err != nil {
		return err
	}

	return nil
}

// Answer describes an answer submitted for a specific Form widget
// with the specific question asked included as well.
type Answer struct {
	WidgetID     string      `json:"widget_id" bson:"widget_id" validate:"required,len=24"`
	Identity     bool        `json:"identity" bson:"identity"`
	Answer       interface{} `json:"answer" bson:"answer"`
	EditedAnswer interface{} `json:"edited" bson:"edited"`
	Question     string      `json:"question" bson:"question"`
	Props        interface{} `json:"props" bson:"props"`
}

// Submission contains all the answers submitted for a specific Form as well
// as any other details about the Form that were present at the time of the Form
// submission.
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

// Validate checks the Submission value for consistency.
func (s *Submission) Validate() error {
	if err := validate.Struct(s); err != nil {
		return err
	}

	return nil
}

// Create adds a new Submission based on a given Form into
// the MongoDB database collection.
func Create(context interface{}, db *db.DB, formID string, submission *Submission) error {
	log.Dev(context, "Create", "Started : Form[%s]", formID)

	if !bson.IsObjectIdHex(formID) {
		log.Error(context, "Create", ErrInvalidID, "Completed")
		return ErrInvalidID
	}

	if err := submission.Validate(); err != nil {
		return err
	}

	// FIXME: handle Number field maybe with https://docs.mongodb.com/v3.0/tutorial/create-an-auto-incrementing-field/ to resolve race condition
	count, err := Count(context, db, formID)
	if err != nil {
		log.Error(context, "Create", err, "Completed")
		return err
	}

	submission.Number = count + 1

	f := func(c *mgo.Collection) error {
		log.Dev(context, "Create", "MGO : db.%s.insert(%s)", c.Name, mongo.Query(submission))
		return c.Insert(submission)
	}

	if err := db.ExecuteMGO(context, Collection, f); err != nil {
		log.Error(context, "Create", err, "Completed")
		return err
	}

	log.Dev(context, "Create", "Completed")
	return nil
}

// Retrieve retrieves a Submission from the MongoDB database
// collection.
func Retrieve(context interface{}, db *db.DB, id string) (*Submission, error) {
	log.Dev(context, "Retrieve", "Started : Submission[%s]", id)

	if !bson.IsObjectIdHex(id) {
		log.Error(context, "Retrieve", ErrInvalidID, "Completed")
		return nil, ErrInvalidID
	}

	objectID := bson.ObjectIdHex(id)

	var submission Submission
	f := func(c *mgo.Collection) error {
		log.Dev(context, "Retrieve", "MGO : db.%s.find(%s)", c.Name, mongo.Query(objectID))
		return c.FindId(objectID).One(&submission)
	}

	if err := db.ExecuteMGO(context, Collection, f); err != nil {
		log.Error(context, "Retrieve", err, "Completed")
		return nil, err
	}

	log.Dev(context, "Retrieve", "Started")
	return &submission, nil
}

// RetrieveMany retrieves a list of Submission's from the MongoDB database collection.
func RetrieveMany(context interface{}, db *db.DB, ids []string) ([]Submission, error) {
	log.Dev(context, "RetrieveMany", "Started")

	var objectIDs = make([]bson.ObjectId, len(ids))

	for i, id := range ids {
		if !bson.IsObjectIdHex(id) {
			log.Error(context, "RetrieveMany", ErrInvalidID, "Completed")
			return nil, ErrInvalidID
		}

		objectIDs[i] = bson.ObjectIdHex(id)
	}

	var submissions []Submission
	f := func(c *mgo.Collection) error {
		q := bson.M{
			"_id": bson.M{
				"$in": objectIDs,
			},
		}
		log.Dev(context, "RetrieveMany", "MGO : db.%s.find(%s)", c.Name, mongo.Query(q))
		return c.Find(q).All(&submissions)
	}

	if err := db.ExecuteMGO(context, Collection, f); err != nil {
		log.Error(context, "RetrieveMany", err, "Completed")
		return nil, err
	}

	log.Dev(context, "RetrieveMany", "Started")
	return submissions, nil
}

// UpdateStatus updates a form submissions status inside the MongoDB database
// collection.
func UpdateStatus(context interface{}, db *db.DB, id, status string) (*Submission, error) {
	log.Dev(context, "UpdateStatus", "Started : Submission[%s]", id)

	if !bson.IsObjectIdHex(id) {
		log.Error(context, "UpdateStatus", ErrInvalidID, "Completed")
		return nil, ErrInvalidID
	}

	objectID := bson.ObjectIdHex(id)

	f := func(c *mgo.Collection) error {
		u := bson.M{
			"$set": bson.M{
				"status":       status,
				"date_updated": time.Now(),
			},
		}
		log.Dev(context, "UpdateStatus", "MGO : db.%s.update(%s, %s)", c.Name, mongo.Query(objectID), mongo.Query(u))
		return c.UpdateId(objectID, u)
	}

	if err := db.ExecuteMGO(context, Collection, f); err != nil {
		log.Error(context, "UpdateStatus", err, "Completed")
		return nil, err
	}

	submission, err := Retrieve(context, db, id)
	if err != nil {
		log.Error(context, "UpdateStatus", err, "Completed")
		return nil, err
	}

	log.Dev(context, "UpdateStatus", "Completed")
	return submission, nil
}

// UpdateAnswer updates the edited answer if it could find it
// inside the MongoDB database collection atomically.
func UpdateAnswer(context interface{}, db *db.DB, id string, answer AnswerInput) (*Submission, error) {
	log.Dev(context, "UpdateAnswer", "Started : Submission[%s]", id)

	if !bson.IsObjectIdHex(id) {
		log.Error(context, "UpdateAnswer", ErrInvalidID, "Completed")
		return nil, ErrInvalidID
	}

	if err := answer.Validate(); err != nil {
		log.Error(context, "UpdateAnswer", ErrInvalidID, "Completed")
		return nil, ErrInvalidID
	}

	objectID := bson.ObjectIdHex(id)

	f := func(c *mgo.Collection) error {
		q := bson.M{
			"_id":               objectID,
			"replies.widget_id": answer.WidgetID,
		}

		// Update the nested subdocument using the $ projection operator:
		// https://docs.mongodb.com/manual/reference/operator/update/positional/
		u := bson.M{
			"$set": bson.M{
				"replies.$.edited": answer.Answer,
			},
		}

		log.Dev(context, "UpdateAnswer", "MGO : db.%s.update(%s, %s)", c.Name, mongo.Query(q), mongo.Query(u))
		return c.Update(q, u)
	}

	if err := db.ExecuteMGO(context, Collection, f); err != nil {
		log.Error(context, "UpdateAnswer", err, "Completed")
		return nil, err
	}

	submission, err := Retrieve(context, db, id)
	if err != nil {
		log.Error(context, "UpdateAnswer", err, "Completed")
		return nil, err
	}

	log.Dev(context, "UpdateAnswer", "Completed")
	return submission, nil
}

// Count returns the count of current submissions for a given
// form id in the Form Submissions MongoDB database collection.
func Count(context interface{}, db *db.DB, formID string) (int, error) {
	log.Dev(context, "Count", "Completed : Form[%s]", formID)

	if !bson.IsObjectIdHex(formID) {
		log.Error(context, "Count", ErrInvalidID, "Completed")
		return 0, ErrInvalidID
	}

	formObjectID := bson.ObjectIdHex(formID)

	var count int
	f := func(c *mgo.Collection) error {
		var err error

		q := bson.M{
			"form_id": formObjectID,
		}
		log.Dev(context, "Count", "MGO : db.%s.find(%s).count()", c.Name, mongo.Query(q))
		count, err = c.Find(q).Count()
		return err
	}

	if err := db.ExecuteMGO(context, Collection, f); err != nil {
		log.Error(context, "Count", err, "Completed")
		return 0, err
	}

	log.Dev(context, "Count", "Completed")
	return count, nil
}

// Search searches through form submissions for a given form
// using the provided search options.
func Search(context interface{}, db *db.DB, formID string, limit, skip int, opts SearchOpts) (*SearchResults, error) {
	log.Dev(context, "Search", "Completed")

	if !bson.IsObjectIdHex(formID) {
		log.Error(context, "Search", ErrInvalidID, "Completed")
		return nil, ErrInvalidID
	}

	formObjectID := bson.ObjectIdHex(formID)

	var sort string
	if opts.DscOrder {
		sort = "-date_created"
	} else {
		sort = "date_created"
	}

	var results = SearchResults{
		Submissions: make([]Submission, 0),
		Counts: SearchResultCounts{
			SearchByFlag: make(map[string]int),
		},
	}
	f := func(c *mgo.Collection) error {
		var err error

		q := bson.M{
			"form_id": formObjectID,
		}

		log.Dev(context, "Search", "MGO : db.%s.find(%s).count()", c.Name, mongo.Query(q))
		results.Counts.TotalSubmissions, err = c.Find(q).Count()
		if err != nil {
			return err
		}

		// If the query or the filter is specificed, we do need to mutate the query
		// to include these terms. If that's the case, we also need to perform a
		// count based on the new search terms.
		if opts.Query != "" || opts.FilterBy != "" {
			if opts.Query != "" {
				// Search query includes the optional text query.
				q["$text"] = bson.M{
					"$search": opts.Query,
				}
			}

			if opts.FilterBy != "" {
				// This must be a tag based filter, so determine if the flag is a
				// negation or not and add the proper filter.
				if strings.HasPrefix(opts.FilterBy, "-") {
					notflag := strings.TrimLeft(opts.FilterBy, "-")
					q["flags"] = bson.M{"$nin": []string{notflag}}
				} else {
					q["flags"] = bson.M{"$in": []string{opts.FilterBy}}
				}
			}

			log.Dev(context, "Search", "MGO : db.%s.find(%s).count()", c.Name, mongo.Query(q))
			results.Counts.TotalSearch, err = c.Find(q).Count()
			if err != nil {
				return err
			}
		} else {
			// As there's no extra filtering criterion, we don't need to re-count the
			// total results as a result of the filtering because there wasn't any!
			results.Counts.TotalSearch = results.Counts.TotalSubmissions
		}

		log.Dev(context, "Search", "MGO : db.%s.find(%s).skip(%d).limit(%d).sort(%s)", c.Name, mongo.Query(q), skip, limit, sort)
		err = c.Find(q).Skip(skip).Limit(limit).Sort(sort).All(&results.Submissions)
		if err != nil {
			return err
		}

		// Instead of pulling all the form submissions ourself, we can just use
		// the MongoDB Query Aggregation.
		pipeline := []bson.M{
			{"$match": q},
			{"$unwind": "$flags"},
			{"$group": bson.M{
				"_id": "$flags",
				"count": bson.M{
					"$sum": 1,
				},
			}},
		}

		var flagBuckets []struct {
			Name  string `bson:"_id"`
			Count int    `bson:"count"`
		}
		log.Dev(context, "Search", "MGO : db.%s.aggregate(%s)", c.Name, mongo.Query(pipeline))
		err = c.Pipe(pipeline).All(&flagBuckets)
		if err != nil {
			return err
		}

		// Load the buckets into the flag aggregation.
		for _, bucket := range flagBuckets {
			results.Counts.SearchByFlag[bucket.Name] = bucket.Count
		}

		return nil
	}

	if err := db.ExecuteMGO(context, Collection, f); err != nil {
		log.Error(context, "Search", err, "Completed")
		return nil, err
	}

	log.Dev(context, "Search", "Completed")
	return &results, nil
}

// AddFlag adds, and de-duplicates a flag to a given
// Submission in the MongoDB database collection.
func AddFlag(context interface{}, db *db.DB, id, flag string) (*Submission, error) {
	log.Dev(context, "AddFlag", "Started : Submission[%s]", id)

	if !bson.IsObjectIdHex(id) {
		log.Error(context, "AddFlag", ErrInvalidID, "Completed")
		return nil, ErrInvalidID
	}

	objectID := bson.ObjectIdHex(id)

	f := func(c *mgo.Collection) error {
		u := bson.M{
			"$addToSet": bson.M{
				"flags": flag,
			},
		}
		log.Dev(context, "AddFlag", "MGO : db.%s.update(%s, %s)", c.Name, mongo.Query(objectID), mongo.Query(u))
		return c.UpdateId(objectID, u)
	}

	if err := db.ExecuteMGO(context, Collection, f); err != nil {
		log.Error(context, "AddFlag", err, "Completed")
		return nil, err
	}

	submission, err := Retrieve(context, db, id)
	if err != nil {
		log.Error(context, "AddFlag", err, "Completed")
		return nil, err
	}

	log.Dev(context, "AddFlag", "Completed")
	return submission, nil
}

// RemoveFlag removes a flag from a given Submission in
// the MongoDB database collection.
func RemoveFlag(context interface{}, db *db.DB, id, flag string) (*Submission, error) {
	log.Dev(context, "RemoveFlag", "Started : Submission[%s]", id)

	if !bson.IsObjectIdHex(id) {
		log.Error(context, "RemoveFlag", ErrInvalidID, "Completed")
		return nil, ErrInvalidID
	}

	objectID := bson.ObjectIdHex(id)

	f := func(c *mgo.Collection) error {
		u := bson.M{
			"$pull": bson.M{
				"flags": flag,
			},
		}
		log.Dev(context, "RemoveFlag", "MGO : db.%s.update(%s, %s)", c.Name, mongo.Query(objectID), mongo.Query(u))
		return c.UpdateId(objectID, u)
	}

	if err := db.ExecuteMGO(context, Collection, f); err != nil {
		log.Error(context, "RemoveFlag", err, "Completed")
		return nil, err
	}

	submission, err := Retrieve(context, db, id)
	if err != nil {
		log.Error(context, "RemoveFlag", err, "Completed")
		return nil, err
	}

	log.Dev(context, "RemoveFlag", "Completed")
	return submission, nil
}

// Delete removes a given Form Submission from the MongoDB
// database collection.
func Delete(context interface{}, db *db.DB, id string) error {
	log.Dev(context, "Delete", "Started : Submission[%s]", id)

	if !bson.IsObjectIdHex(id) {
		log.Error(context, "Delete", ErrInvalidID, "Completed")
		return ErrInvalidID
	}

	objectID := bson.ObjectIdHex(id)

	f := func(c *mgo.Collection) error {
		log.Dev(context, "Delete", "MGO : db.%s.remove(%s)", c.Name, mongo.Query(objectID))
		return c.RemoveId(objectID)
	}

	if err := db.ExecuteMGO(context, Collection, f); err != nil {
		log.Error(context, "Delete", err, "Completed")
		return err
	}

	log.Dev(context, "Delete", "Started")
	return nil
}
