package form

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/ardanlabs/kit/log"
	"github.com/coralproject/shelf/internal/ask/form/submission"
	"github.com/coralproject/shelf/internal/platform/db"
	"github.com/coralproject/shelf/internal/platform/db/mongo"
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
var (
	ErrInvalidID            = errors.New("ID is not in it's proper form")
	ErrUpdatingAggregations = errors.New("Could not aggregate Form statistics.")
)

//==============================================================================

// Collection is the mongo collection where Form documents are saved.
const Collection = "forms"

// Widget describes a specific question being asked by the Form which is
// contained within a Step.
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

// Step is a collection of Widget's.
type Step struct {
	ID      string   `json:"id" bson:"_id"`
	Name    string   `json:"name" bson:"name"`
	Widgets []Widget `json:"widgets" bson:"widgets"`
}

// MCAnswerAggregation holds the count for selections of a single multiple
// choice answer.
type MCAnswerAggregation struct {
	Title string `json:"answer" bson:"answer"`
	Count int    `json:"count" bson:"count"`
}

// MCAggregation holds a multiple choice question and a map aggregated counts for
// each answer. The Answers map is keyed off an md5 of the answer as not better keys exist
type MCAggregation struct {
	Question  string                         `json:"question" bson:"question"`
	MCAnswers map[string]MCAnswerAggregation `json:"answers" bson:"answers"`
}

// TextAggregation holds the aggregated text based answers for a single question.
type TextAggregation struct {
	Question string   `json:"question" bson:"question"`
	Answers  []string `json:"answers" bson:"ansewrs"`
}

// Stats describes the statistics being recorded by a specific Form.
type Stats struct {
	Responses        int                        `json:"responses" bson:"responses"`
	MCAggregations   map[string]MCAggregation   `json:"mc_aggregations" bson:"mc_aggregations"`
	TextAggregations map[string]TextAggregation `json:"text_aggregations" bson:"text_aggregations"`
}

// Group defines a key for a multiple choice question / answer combo.
type Group struct {
	QuestionID string `json:"question_id"`
	Question   string `json:"question"`
	Answer     string `json:"answer"`
}

//==============================================================================

// Form contains the conatical representation of a Form, containing all the
// Steps, and help text relating to completing the Form.
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

// Validate checks the Form value for consistency.
func (f *Form) Validate() error {
	if err := validate.Struct(f); err != nil {
		return err
	}

	return nil
}

//==============================================================================

// Upsert upserts the provided form into the MongoDB database collection.
func Upsert(context interface{}, db *db.DB, form *Form) error {
	log.Dev(context, "Upsert", "Started")

	var isNewForm bool

	// If there is no ID probided, we should set one as this is an Upsert
	// operation. It is also important to remember if this was a new form or not
	// because we need to update the stats if this wasn't a new form.
	if form.ID.Hex() == "" {
		form.ID = bson.NewObjectId()
		isNewForm = true
		form.DateCreated = time.Now()
	}

	form.DateUpdated = time.Now()

	if err := form.Validate(); err != nil {
		log.Error(context, "Upsert", err, "Completed")
		return err
	}

	f := func(c *mgo.Collection) error {
		log.Dev(context, "Upsert", "MGO : db.%s.upsert(%s, %s)", c.Name, mongo.Query(form.ID.Hex()), mongo.Query(form))
		_, err := c.UpsertId(form.ID, form)
		return err
	}

	if err := db.ExecuteMGO(context, Collection, f); err != nil {
		log.Error(context, "Upsert", err, "Completed")
		return err
	}

	// New forms don't have any stats so don't bother updating it.
	if !isNewForm {
		if _, err := UpdateStats(context, db, form.ID.Hex()); err != nil {
			log.Error(context, "Upsert", err, "Completed")
			return err
		}
	}

	log.Dev(context, "Upsert", "Completed")
	return nil
}

// UpdateStats updates the Stats on a given Form.
func UpdateStats(context interface{}, db *db.DB, id string) (*Stats, error) {
	log.Dev(context, "UpdateStats", "Started : Form[%s]", id)

	if !bson.IsObjectIdHex(id) {
		log.Error(context, "UpdateStats", ErrInvalidID, "Completed")
		return nil, ErrInvalidID
	}

	objectID := bson.ObjectIdHex(id)

	// Find the number of submissions on this form
	count, err := submission.Count(context, db, id)
	if err != nil {
		log.Error(context, "UpdateStats", ErrInvalidID, "Completed")
		return nil, err
	}

	stats := Stats{
		Responses: count,
	}

	f := func(c *mgo.Collection) error {
		u := bson.M{
			"$set": bson.M{
				"stats":        stats,
				"date_updated": time.Now(),
			},
		}
		log.Dev(context, "UpdateStats", "MGO : db.%s.update(%s, %s)", c.Name, mongo.Query(objectID), mongo.Query(u))
		return c.UpdateId(objectID, u)
	}

	if err := db.ExecuteMGO(context, Collection, f); err != nil {
		log.Error(context, "UpdateStats", err, "Completed")
		return nil, err
	}

	log.Dev(context, "UpdateStats", "Completed")
	return &stats, nil
}

// UpdateStatus updates the forms status and returns the updated form from
// the MongodB database collection.
func UpdateStatus(context interface{}, db *db.DB, id, status string) (*Form, error) {
	log.Dev(context, "UpdateStatus", "Started : Form[%s]", id)

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

	form, err := Retrieve(context, db, id)
	if err != nil {
		log.Error(context, "UpdateStatus", err, "Completed")
		return nil, err
	}

	log.Dev(context, "UpdateStatus", "Completed")
	return form, nil
}

// Retrieve retrieves the form from the MongodB database collection.
func Retrieve(context interface{}, db *db.DB, id string) (*Form, error) {
	log.Dev(context, "Retrieve", "Started : Form[%s]", id)

	if !bson.IsObjectIdHex(id) {
		log.Error(context, "Retrieve", ErrInvalidID, "Completed")
		return nil, ErrInvalidID
	}

	objectID := bson.ObjectIdHex(id)

	var form Form
	f := func(c *mgo.Collection) error {
		log.Dev(context, "Retrieve", "MGO : db.%s.find(%s)", c.Name, mongo.Query(objectID))
		return c.FindId(objectID).One(&form)
	}

	if err := db.ExecuteMGO(context, Collection, f); err != nil {
		log.Error(context, "Retrieve", err, "Completed")
		return nil, err
	}

	log.Dev(context, "Retrieve", "Completed")
	return &form, nil
}

// List retrieves a list of forms from the MongodB database collection.
func List(context interface{}, db *db.DB, limit, skip int) ([]Form, error) {
	log.Dev(context, "List", "Started")

	var forms = make([]Form, 0)
	f := func(c *mgo.Collection) error {
		log.Dev(context, "List", "MGO : db.%s.find().limit(%d).skip(%d)", c.Name, limit, skip)
		return c.Find(nil).Limit(limit).Skip(skip).All(&forms)
	}

	if err := db.ExecuteMGO(context, Collection, f); err != nil {
		log.Error(context, "List", err, "Completed")
		return nil, err
	}

	log.Dev(context, "List", "Completed")
	return forms, nil
}

// Delete removes the document matching the id provided from the MongoDB
// database collection.
func Delete(context interface{}, db *db.DB, id string) error {
	log.Dev(context, "Delete", "Started : Form[%s]", id)

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

	log.Dev(context, "Delete", "Completed")
	return nil
}

// MCAggregate calculates statistics on all multiple choice questions across submission
// on a form. It only considers qustions in the current form as currently there is no
// way to track how questions and answers change if the admin updates the form mid flight.
// Aggregate returns a map[_question_]map[_choice_]int_count datastructure.
func MCAggregate(context interface{}, db *db.DB, id string, subs []submission.Submission) (map[string]MCAggregation, error) {
	log.Dev(context, "Aggregate", "Started : Submission[%s]", id)

	if !bson.IsObjectIdHex(id) {
		log.Error(context, "Aggregate", ErrInvalidID, "Completed")
		return nil, ErrInvalidID
	}

	// Get the form in question.
	form, err := Retrieve(context, db, id)
	if err != nil {
		return nil, err
	}

	// Create a container for the aggregations: [question][option]count.
	aggs := make(map[string]MCAggregation)

	// Find the MultipleChoice widgets and add them to the aggs map
	for _, step := range form.Steps {
		for _, widget := range step.Widgets {
			if widget.Component == "MultipleChoice" {
				aggs[widget.ID] = MCAggregation{
					Question: widget.Title,
				}
			}
		}
	}

	// In this section we are looking through all submissions for answers to multiple choice
	// questions that are still active in the form and counting question/answer pairs.

	// Look at all submisisons.
	for _, sub := range subs {

		// Then at every anwer.
		for _, ans := range sub.Answers {

			// Skip nil answers.
			if ans.Answer == nil {
				continue
			}

			// The following section points to the need to refactor form submissions / introduce
			// stronger typing.

			// Unpack the answer object.
			a := ans.Answer.(bson.M)

			options := a["options"]

			// Options == nil points to a non MultipleChoice answer.
			if options == nil {
				continue
			}

			// This map of interfaces represent each checkbox the user clicked.
			opts := options.([]interface{})
			for _, opt := range opts {

				// Unpack the option.
				op := opt.(bson.M)

				// Use the title of the option as the map key.
				selection := op["title"].(string)

				// Hash the ansewr text for a unique key, as no actual key exists.
				hasher := md5.New()
				hasher.Write([]byte(op["title"].(string)))
				optKeyStr := hex.EncodeToString(hasher.Sum(nil))

				// If this question is not in the map then we can skip as it is not a current answer.
				if _, ok := aggs[ans.WidgetID]; !ok {
					continue
				}

				// If this is the first answer for this question, make a map for it.
				if aggs[ans.WidgetID].MCAnswers == nil {
					tmp := aggs[ans.WidgetID]
					tmp.MCAnswers = make(map[string]MCAnswerAggregation)
					aggs[ans.WidgetID] = tmp
				}

				// If this is the first time we've seen this answer, init the agg struct for it.
				if _, ok := aggs[ans.WidgetID].MCAnswers[optKeyStr]; !ok {
					aggs[ans.WidgetID].MCAnswers[optKeyStr] = MCAnswerAggregation{
						Title: selection,
						Count: 0,
					}
				}

				// Increment the counter for this question/answer pair.
				tmp := aggs[ans.WidgetID].MCAnswers[optKeyStr]
				tmp.Count++
				aggs[ans.WidgetID].MCAnswers[optKeyStr] = tmp

			}
		}
	}

	log.Dev(context, "Aggregate", "Completed : Submission[%s]", id)
	return aggs, nil
}

// TextAggregate returns all text answers flagged with includeInGroup.
func TextAggregate(context interface{}, subs []submission.Submission) (map[string]TextAggregation, error) {

	// Create a container for the aggregations: [question][option]count.
	aggs := make(map[string]TextAggregation)

	// Scan all the submissions and answers.
	for _, sub := range subs {
		for _, ans := range sub.Answers {

			// Skip nil answers.
			if ans.Answer == nil {
				continue
			}

			// Only include answers of questions flagged with "includeInGroups".
			props := ans.Props.(bson.M)
			include, ok := props["includeInGroups"]
			if include == nil || !ok {
				continue
			}

			// Unpack the answer and append it to the aggregation slice for that WidgetId.
			a := ans.Answer.(bson.M)
			tmp := aggs[ans.WidgetID]
			tmp.Question = ans.Question
			tmp.Answers = append(tmp.Answers, a["text"].(string))
			aggs[ans.WidgetID] = tmp

		}
	}

	log.Dev(context, "Text Aggregate", "Completed : Submission")
	return aggs, nil
}

// GroupSubmissions organizes submissions by Group. It looks for questions with the group by flag
// and creates Group structs.
func GroupSubmissions(context interface{}, db *db.DB, id string) (map[Group][]submission.Submission, error) {

	if !bson.IsObjectIdHex(id) {
		log.Error(context, "TextAggregate", ErrInvalidID, "Completed")
		return nil, ErrInvalidID
	}

	// Get the submissions for the form.Collection
	subs, err := submission.Search(context, db, id, 0, 0, submission.SearchOpts{})
	if err != nil {
		return nil, err
	}

	groups := make(map[Group][]submission.Submission)

	// Scan all the submissions and answers.
	for _, sub := range subs.Submissions {

		// Add all submissions to the [all,all] group
		group := Group{
			Question: "all",
			Answer:   "all",
		}
		tmp := groups[group]
		tmp = append(tmp, sub)
		groups[group] = tmp

		// Look for answers that may contain sub groups.
		for _, ans := range sub.Answers {

			// Skip nil answers.
			if ans.Answer == nil {
				continue
			}

			// Only include answers of questions flagged with "includeInGroups".
			props := ans.Props.(bson.M)
			include, ok := props["groupSubmissions"]
			if include == nil || !ok || include == false {
				continue
			}

			fmt.Printf("\n\n%#v", props["groupSubmissions"])

			// Unpack the answer object.
			a := ans.Answer.(bson.M)

			options := a["options"]

			// Options == nil points to a non MultipleChoice answer.
			if options == nil {
				continue
			}

			// This map of interfaces represent each checkbox the user clicked.
			opts := options.([]interface{})
			for _, opt := range opts {

				// Unpack the option.
				op := opt.(bson.M)

				// Use the title of the option as the map key.
				selection := op["title"].(string)

				// Add the submission to this subgroup
				group := Group{
					QuestionID: ans.WidgetID,
					Question:   ans.Question,
					Answer:     selection,
				}
				tmp := groups[group]
				tmp = append(tmp, sub)
				groups[group] = tmp

			}

		}
	}

	return groups, nil
}
