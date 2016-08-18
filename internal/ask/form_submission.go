package ask

import (
	"strings"
	"time"

	"github.com/ardanlabs/kit/db"
	"github.com/ardanlabs/kit/db/mongo"
	"github.com/ardanlabs/kit/log"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// FormSubmissionsCollection is the mongo collection where FormSubmission
// documents are saved.
const FormSubmissionsCollection = "form_submissions"

// FormSubmissionSearchResults is a structured type returning the results
// expected from searching for submissions based on a form id.
type FormSubmissionSearchResults struct {
	Counts struct {
		SearchByFlag     map[string]int `json:"search_by_flag"`
		TotalSearch      int            `json:"total_search"`
		TotalSubmissions int            `json:"total_submissions"`
	} `json:"counts"`
	Submissions []FormSubmission
}

// FormSubmissionSearchOpts is the options used to perform a search accross a
// given forms submissions.
type FormSubmissionSearchOpts struct {
	DscOrder bool
	Query    string
	FilterBy string
}

// FormSubmissionAnswerInput describes the input accepted for a new submission
// answer.
type FormSubmissionAnswerInput struct {
	WidgetID string      `json:"widget_id"`
	Answer   interface{} `json:"answer"`
}

// FormSubmissionAnswer describes an answer submitted for a specific Form widget
// with the specific question asked included as well.
type FormSubmissionAnswer struct {
	WidgetID     string      `json:"widget_id" bson:"widget_id"`
	Identity     bool        `json:"identity" bson:"identity"`
	Answer       interface{} `json:"answer" bson:"answer"`
	EditedAnswer interface{} `json:"edited" bson:"edited"`
	Question     interface{} `json:"question" bson:"question"`
	Props        interface{} `json:"props" bson:"props"`
}

// FormSubmission contains all the answers submitted for a specific Form as well
// as any other details about the Form that were present at the time of the Form
// submission.
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

// CreateFormSubmission adds a new FormSubmission based on a given Form into
// the MongoDB database collection.
func CreateFormSubmission(context interface{}, db *db.DB, formID string, answers []FormSubmissionAnswerInput) (*FormSubmission, error) {
	log.Dev(context, "CreateFormSubmission", "Started")

	if !bson.IsObjectIdHex(formID) {
		log.Error(context, "CreateFormSubmission", ErrInvalidID, "Completed")
		return nil, ErrInvalidID
	}

	form, err := RetrieveForm(context, db, formID)
	if err != nil {
		log.Error(context, "CreateFormSubmission", err, "Completed")
		return nil, err
	}

	// create the new form submission
	fs := FormSubmission{
		ID:          bson.NewObjectId(),
		FormID:      bson.ObjectIdHex(formID),
		Header:      form.Header,
		Footer:      form.Footer,
		Answers:     make([]FormSubmissionAnswer, 0),
		DateCreated: time.Now(),
		DateUpdated: time.Now(),
	}

	// FIXME: handle Number field maybe with https://docs.mongodb.com/v3.0/tutorial/create-an-auto-incrementing-field/

	// for each answer
	for _, answer := range answers {
		var found bool

		// we must check each step of the form
		for _, step := range form.Steps {
			// and each widget
			for _, widget := range step.Widgets {
				// to see if we can find the matching widget for this answer
				if answer.WidgetID == widget.ID {
					// and push that answer into the form submission
					fs.Answers = append(fs.Answers, FormSubmissionAnswer{
						WidgetID: widget.ID,
						Answer:   answer,
						Identity: widget.Identity,
						Question: widget.Title,
						Props:    widget.Props,
					})

					// mark the answer as found
					found = true

					break
				}
			}

			// so that if the answer was already found...
			if found {
				// we can break out of this step loop
				break
			}
		}

	}

	f := func(c *mgo.Collection) error {
		log.Dev(context, "CreateFormSubmission", "MGO : db.%s.insert(%s)", c.Name, mongo.Query(fs))
		return c.Insert(fs)
	}

	if err := db.ExecuteMGO(context, FormSubmissionsCollection, f); err != nil {
		log.Error(context, "CreateFormSubmission", err, "Completed")
		return nil, err
	}

	if _, err := UpdateFormStats(context, db, formID); err != nil {
		log.Error(context, "CreateFormSubmission", err, "Completed")
		return nil, err
	}

	log.Dev(context, "CreateFormSubmission", "Completed")
	return &fs, nil
}

// RetrieveFormSubmission retrieves a FormSubmission from the MongoDB database
// collection.
func RetrieveFormSubmission(context interface{}, db *db.DB, id string) (*FormSubmission, error) {
	log.Dev(context, "RetrieveFormSubmission", "Started")

	if !bson.IsObjectIdHex(id) {
		log.Error(context, "RetrieveFormSubmission", ErrInvalidID, "Completed")
		return nil, ErrInvalidID
	}

	objectID := bson.ObjectIdHex(id)

	var submission FormSubmission
	f := func(c *mgo.Collection) error {
		log.Dev(context, "RetrieveFormSubmission", "MGO : db.%s.find(%s)", c.Name, mongo.Query(objectID.Hex()))
		return c.FindId(objectID).One(&submission)
	}

	if err := db.ExecuteMGO(context, FormSubmissionsCollection, f); err != nil {
		log.Error(context, "RetrieveFormSubmission", err, "Completed")
		return nil, err
	}

	log.Dev(context, "RetrieveFormSubmission", "Started")
	return &submission, nil
}

// UpdateFormSubmissionStatus updates a form submissions status inside the
// MongoDB database collection.
func UpdateFormSubmissionStatus(context interface{}, db *db.DB, id, status string) (*FormSubmission, error) {
	log.Dev(context, "UpdateFormSubmissionStatus", "Started")

	if !bson.IsObjectIdHex(id) {
		log.Error(context, "UpdateFormSubmissionStatus", ErrInvalidID, "Completed")
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
		log.Dev(context, "UpdateFormSubmissionStatus", "MGO : db.%s.update(%s, %s)", c.Name, mongo.Query(objectID.Hex()), mongo.Query(u))
		return c.UpdateId(objectID, u)
	}

	if err := db.ExecuteMGO(context, FormSubmissionsCollection, f); err != nil {
		log.Error(context, "UpdateFormSubmissionStatus", err, "Completed")
		return nil, err
	}

	submission, err := RetrieveFormSubmission(context, db, id)
	if err != nil {
		log.Error(context, "UpdateFormSubmissionStatus", err, "Completed")
		return nil, err
	}

	log.Dev(context, "UpdateFormSubmissionStatus", "Completed")
	return submission, nil
}

// UpdateFormSubmissionAnswer updates the edited answer if it could find it
// inside the MongoDB database collection atomically.
func UpdateFormSubmissionAnswer(context interface{}, db *db.DB, id, answerID string, editedAnswer interface{}) (*FormSubmission, error) {
	log.Dev(context, "UpdateFormSubmissionAnswer", "Started")

	if !bson.IsObjectIdHex(id) {
		log.Error(context, "UpdateFormSubmissionAnswer", ErrInvalidID, "Completed")
		return nil, ErrInvalidID
	}

	if !bson.IsObjectIdHex(answerID) {
		log.Error(context, "UpdateFormSubmissionAnswer", ErrInvalidID, "Completed")
		return nil, ErrInvalidID
	}

	objectID := bson.ObjectIdHex(id)

	f := func(c *mgo.Collection) error {
		q := bson.M{
			"_id":               objectID,
			"replies.widget_id": answerID,
		}

		// update the nested subdocument using the $ projection operator:
		// https://docs.mongodb.com/manual/reference/operator/update/positional/
		u := bson.M{
			"$set": bson.M{
				"replies.$.edited": editedAnswer,
			},
		}

		log.Dev(context, "UpdateFormSubmissionAnswer", "MGO : db.%s.update(%s, %s)", c.Name, mongo.Query(q), mongo.Query(u))
		return c.Update(q, u)
	}

	if err := db.ExecuteMGO(context, FormSubmissionsCollection, f); err != nil {
		log.Error(context, "UpdateFormSubmissionAnswer", err, "Completed")
		return nil, err
	}

	submission, err := RetrieveFormSubmission(context, db, id)
	if err != nil {
		log.Error(context, "UpdateFormSubmissionAnswer", err, "Completed")
		return nil, err
	}

	log.Dev(context, "UpdateFormSubmissionAnswer", "Completed")
	return submission, nil
}

// SearchFormSubmissions searches through form submissions for a given form
// using the provided search options.
func SearchFormSubmissions(context interface{}, db *db.DB, formID string, limit, skip int, opts FormSubmissionSearchOpts) (*FormSubmissionSearchResults, error) {
	log.Dev(context, "SearchFormSubmissions", "Completed")

	if !bson.IsObjectIdHex(formID) {
		log.Error(context, "SearchFormSubmissions", ErrInvalidID, "Completed")
		return nil, ErrInvalidID
	}

	formObjectID := bson.ObjectIdHex(formID)

	var sort string
	if opts.DscOrder {
		sort = "-date_created"
	} else {
		sort = "date_created"
	}

	var results = FormSubmissionSearchResults{
		Submissions: make([]FormSubmission, 0),
	}
	f := func(c *mgo.Collection) error {
		var err error

		q := bson.M{
			"form_id": formObjectID,
		}

		log.Dev(context, "SearchFormSubmissions", "MGO : db.%s.find(%s).count()", c.Name, mongo.Query(q))
		results.Counts.TotalSubmissions, err = c.Find(q).Count()
		if err != nil {
			return err
		}

		if opts.Query != "" || opts.FilterBy != "" {
			if opts.Query != "" {
				// search query including the optional text query
				q["$text"] = bson.M{
					"$search": opts.Query,
				}
			} else {
				// a flag based filter
				if strings.HasPrefix(opts.FilterBy, "-") {
					notflag := strings.TrimLeft(opts.FilterBy, "-")
					q["flags"] = bson.M{"$nin": []string{notflag}}
				} else {
					q["flags"] = bson.M{"$in": []string{opts.FilterBy}}
				}
			}

			log.Dev(context, "SearchFormSubmissions", "MGO : db.%s.find(%s).count()", c.Name, mongo.Query(q))
			results.Counts.TotalSearch, err = c.Find(q).Count()
			if err != nil {
				return err
			}
		} else {
			// as there's no extra filtering criterion, we don't need to re-count the
			// total results as a result of the filtering because there wasn't any!
			results.Counts.TotalSearch = results.Counts.TotalSubmissions
		}

		log.Dev(context, "SearchFormSubmissions", "MGO : db.%s.find(%s).skip(%d).limit(%d).sort(%s)", c.Name, mongo.Query(q), skip, limit, sort)
		err = c.Find(q).Skip(skip).Limit(limit).Sort(sort).All(&results.Submissions)
		if err != nil {
			return err
		}

		// instead of pulling all the form submissions ourself, we're
		pipeline := []bson.M{
			bson.M{"$match": q},
			bson.M{"$unwind": "$flags"},
			bson.M{"$group": bson.M{
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
		log.Dev(context, "SearchFormSubmissions", "MGO : db.%s.aggregate(%s)", c.Name, mongo.Query(pipeline))
		err = c.Pipe(pipeline).All(&flagBuckets)
		if err != nil {
			return err
		}

		// load the buckets into the flag aggregation
		for _, bucket := range flagBuckets {
			results.Counts.SearchByFlag[bucket.Name] = bucket.Count
		}

		return nil
	}

	if err := db.ExecuteMGO(context, FormSubmissionsCollection, f); err != nil {
		log.Error(context, "SearchFormSubmissions", err, "Completed")
		return nil, err
	}

	log.Dev(context, "SearchFormSubmissions", "Completed")
	return &results, nil
}

// AddFlagToFormSubmission adds, and de-duplicates a flag to a given
// FormSubmission in the MongoDB database collection.
func AddFlagToFormSubmission(context interface{}, db *db.DB, id, flag string) (*FormSubmission, error) {
	// TODO: implement
	return nil, nil
}

// RemoveFlagFromFormSubmission removes a flag from a given FormSubmission in
// the MongoDB database collection.
func RemoveFlagFromFormSubmission(context interface{}, db *db.DB, id, flag string) (*FormSubmission, error) {
	// TODO: implement
	return nil, nil
}

// DeleteFormSubmission removes a given Form Submission from the MongoDB
// database collection.
func DeleteFormSubmission(context interface{}, db *db.DB, id, formID string) error {
	log.Dev(context, "DeleteFormSubmission", "Started")

	if !bson.IsObjectIdHex(id) {
		log.Error(context, "DeleteFormSubmission", ErrInvalidID, "Completed")
		return ErrInvalidID
	}

	// FIXME: uncomment once old API has been deprecated.
	// if !bson.IsObjectIdHex(formID) {
	// 	log.Error(context, "DeleteFormSubmission", ErrInvalidID, "Completed")
	// 	return ErrInvalidID
	// }

	objectID := bson.ObjectIdHex(id)

	f := func(c *mgo.Collection) error {
		log.Dev(context, "DeleteFormSubmission", "MGO : db.%s.remove(%s)", c.Name, mongo.Query(objectID.Hex()))
		return c.RemoveId(objectID)
	}

	if err := db.ExecuteMGO(context, FormSubmissionsCollection, f); err != nil {
		log.Error(context, "DeleteFormSubmission", err, "Completed")
		return err
	}

	// FIXME: remove once old API has been deprecated.
	if formID != "" {
		if _, err := UpdateFormStats(context, db, formID); err != nil {
			log.Error(context, "DeleteFormSubmission", err, "Completed")
			return err
		}
	}

	log.Dev(context, "DeleteFormSubmission", "Started")
	return nil
}
