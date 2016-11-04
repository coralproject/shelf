package form

import (
	"crypto/md5"
	"encoding/hex"

	"github.com/ardanlabs/kit/log"
	"github.com/coralproject/shelf/internal/ask/form/submission"
	"github.com/coralproject/shelf/internal/platform/db"
	"gopkg.in/mgo.v2/bson"
)

//==============================================================================

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

// TextAggregation holds the aggregated text based answers for a single question
// marked with the Incude in Aggregations tag, orderd by [question_id][answer].
type TextAggregation map[string]string

// Aggregation holds the various aggregations and stats collected.
type Aggregation struct {
	Group Group                    `json:"group" bson:"group"`
	Count int                      `json:"count" bson:"count"`
	MC    map[string]MCAggregation `json:"MultipleChoice" bson:"MultipleChoice"` // Capitalization matches widget type MultipleChoice
}

// Group defines a key for a multiple choice question / answer combo to be used
// to define slices of submissions to be aggregated.
type Group struct {
	ID       string `json:"group_id" bson:"group_id"`
	Question string `json:"question" bson:"question"`
	Answer   string `json:"answer" bson:"answer"`
}

//==============================================================================

// AggregateFormSubmissions retrieves the submissions for a form, groups them then
// runs aggregations and counts for each one.
func AggregateFormSubmissions(context interface{}, db *db.DB, id string) (map[string]Aggregation, error) {

	// Aggregations are hardcoded to only work with bookmarked submissions.
	opts := submission.SearchOpts{
		FilterBy: "bookmarked",
	}

	// Group the submissions.
	groupedSubmissions, err := GroupSubmissions(context, db, id, 0, 0, opts)
	if err != nil {
		return nil, err
	}

	// Create a container for the grouped aggregations.
	groupAggregations := make(map[string]Aggregation)

	// Loop through the groups of submissions.
	for group, submissions := range groupedSubmissions {

		// Perform the multiple choice aggregations.
		mcAggregations, err := mcAggregate(context, db, id, submissions)
		if err != nil {
			return nil, err
		}

		// Including the aggregation of the group that is being aggregated is redundant, remove.
		for key, agg := range mcAggregations {
			if agg.Question == group.Question {
				delete(mcAggregations, key)
			}
		}

		// Pack them in an aggregation along with the submission count and group.
		agg := Aggregation{
			Group: group,
			Count: len(submissions),
			MC:    mcAggregations,
		}

		// If this is the "all" group, set the key to all as it is not an answer.
		groupKey := ""
		if group.Question == "all" {
			groupKey = "all"
		}

		// Groups are ultimately based on a chosen answer. We do not have any keys for answers
		// so hash the answer text for a unique key.
		if group.Question != "all" {
			hasher := md5.New()
			hasher.Write([]byte(group.Answer))
			groupKey = hex.EncodeToString(hasher.Sum(nil))
		}

		// Add the aggregation to the map.
		groupAggregations[groupKey] = agg

	}

	return groupAggregations, nil
}

//==============================================================================

// SubmissionGroup is a transport that defines the transport structure for a submission group.
type SubmissionGroup struct {
	Submissions map[Group][]submission.Submission `json:"submissions" bson:"submissions"`
}

// GroupSubmissions organizes submissions by Group. It looks for questions with the group by flag
// and creates Group structs.
func GroupSubmissions(context interface{}, db *db.DB, formID string, limit int, skip int, opts submission.SearchOpts) (map[Group][]submission.Submission, error) {

	// Validate the formID.
	if !bson.IsObjectIdHex(formID) {
		log.Error(context, "TextAggregate", ErrInvalidID, "Completed")
		return nil, ErrInvalidID
	}

	// Load the form so that we can find which questions are flagged groupBy.
	form, err := Retrieve(context, db, formID)
	if err != nil {
		return nil, err
	}

	// Create a container for questions flagged as groupSubmissions.
	groupBys := make(map[string]bool)

	// Find the widgets to groupBy and add them to the map.
	for _, step := range form.Steps {
		for _, widget := range step.Widgets {

			// Type the props value.
			props, ok := widget.Props.(bson.M)
			if !ok {
				continue
			}

			// If groupSubmissions is set, add the ID to the map of questions to include.
			if props["groupSubmissions"] == true {
				groupBys[widget.ID] = true
			}
		}
	}

	// Get the submissions for the form.Collection.
	subs, err := submission.Search(context, db, formID, limit, skip, opts)
	if err != nil {
		return nil, err
	}

	groups := make(map[Group][]submission.Submission)

	// Scan all the submissions and answers.
	for _, sub := range subs.Submissions {

		// Add all submissions to the [all,all] group.
		group := Group{
			ID:       "all",
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
			_, ok := groupBys[ans.WidgetID]
			if !ok {
				continue
			}

			// Unpack the answer object.
			a, ok := ans.Answer.(bson.M)
			if !ok {
				continue
			}

			// Retrieve the options param from the answer.
			options, ok := a["options"]

			// Options == nil points to a non MultipleChoice answer.
			if !ok || options == nil {
				continue
			}

			// This map of interfaces represent each checkbox the user clicked.
			opts, ok := options.([]interface{})
			if !ok {
				continue
			}
			for _, opt := range opts {

				// Unpack the option.
				op := opt.(bson.M)

				// Use the title of the option as the map key.
				selection, ok := op["title"].(string)
				if !ok {
					continue
				}

				// Hash the answer text for a unique key, as no actual key exists.
				hasher := md5.New()
				hasher.Write([]byte(selection))
				optKeyStr := hex.EncodeToString(hasher.Sum(nil))

				// Add the submission to this subgroup.
				group := Group{
					ID:       optKeyStr,
					Question: ans.Question,
					Answer:   selection,
				}
				tmp := groups[group]
				tmp = append(tmp, sub)
				groups[group] = tmp
			}
		}
	}

	return groups, nil
}

//==============================================================================

// Aggregation functions take arrays of submissions and aggregate certain field
// types based on the parameters embedded in the form.

// TextAggregate returns all text answers flagged with includeInGroup.
func TextAggregate(context interface{}, db *db.DB, formID string, subs []submission.Submission) ([]TextAggregation, error) {

	// We load the form to look up wich fields to includeInGroups.

	// Ensure that the id passed is a valid bson IdHex.
	if !bson.IsObjectIdHex(formID) {
		log.Error(context, "Aggregate", ErrInvalidID, "Completed")
		return nil, ErrInvalidID
	}

	// Load the form so that we can find which questions are flagged to include
	// in the submission groups.
	form, err := Retrieve(context, db, formID)
	if err != nil {
		return nil, err
	}

	// Create a container for questions to include in the groups.
	includes := make(map[string]bool)

	// Find the widgets to includeInGroups and add them to the map.
	for _, step := range form.Steps {
		for _, widget := range step.Widgets {

			// Type the props value.
			props, ok := widget.Props.(bson.M)
			if !ok {
				continue
			}

			// If includeInGroups is set, add the ID to the map of questions to include.
			if props["includeInGroups"] == true {
				includes[widget.ID] = true
			}
		}
	}

	// Create a container for the aggregations: [question][option]count.
	textAggregations := []TextAggregation{}

	// Scan all the submissions and answers.
	for _, sub := range subs {

		// Create an aggregation and set the ID.
		textAggregation := TextAggregation{}
		textAggregation["id"] = sub.ID.Hex()

		for _, ans := range sub.Answers {

			// Skip nil answers.
			if ans.Answer == nil {
				continue
			}

			// Only include answers of questions flagged with "includeInGroups" that
			// we added to the includes map above.
			_, ok := includes[ans.WidgetID]
			if !ok {
				continue
			}

			// Options == nil points to a non MultipleChoice answer.
			var answer string

			a, ok := ans.Answer.(bson.M)
			if !ok {
				continue
			}

			options := a["options"]
			if options == nil {
				// Unpack the answer and add it to the map at the widgetID
				a, ok := ans.Answer.(bson.M)
				if !ok {
					continue
				}
				answer, ok = a["text"].(string)
				if !ok {
					continue
				}
			}

			// If we have multiple choice, use the first selection.
			if options != nil {
				opts, ok := options.([]interface{})
				if !ok {
					continue
				}

				// Unpack the option.
				op, ok := opts[0].(bson.M)
				if !ok {
					continue
				}

				// Use the title of the option as the map key.
				answer, ok = op["title"].(string)
				if !ok {
					continue
				}

			}

			textAggregation[ans.WidgetID] = answer

		}

		textAggregations = append(textAggregations, textAggregation)
	}

	log.Dev(context, "Text Aggregate", "Completed : Submission")
	return textAggregations, nil
}

// MCAggregate calculates statistics on all multiple choice questions.
func mcAggregate(context interface{}, db *db.DB, formID string, subs []submission.Submission) (map[string]MCAggregation, error) {
	log.Dev(context, "Aggregate", "Started : Submission[%s]", formID)

	// We load the form so that only the multiple choice questions currently in the form
	// will be included in the aggregation.

	// Ensure that the id passed is a valid bson IdHex.
	if !bson.IsObjectIdHex(formID) {
		log.Error(context, "Aggregate", ErrInvalidID, "Completed")
		return nil, ErrInvalidID
	}

	// Get the form in question.
	form, err := Retrieve(context, db, formID)
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

		// Then at every anwer, where an answer is to a question.
		for _, ans := range sub.Answers {

			// Skip nil answers.
			if ans.Answer == nil {
				continue
			}

			// The following section points to the need to refactor form submissions / introduce
			// stronger typing.

			// Unpack the answer object.
			a, ok := ans.Answer.(bson.M)
			if !ok {
				continue
			}

			// Retrieve the options param from the answer.
			options, ok := a["options"]

			// Options == nil points to a non MultipleChoice answer.
			if !ok || options == nil {
				continue
			}

			// This map of interfaces represent each checkbox the user clicked.
			opts, ok := options.([]interface{})
			if !ok {
				continue
			}
			for _, opt := range opts {

				// Unpack the option.
				op, ok := opt.(bson.M)
				if !ok {
					continue
				}

				// Use the title of the option as the map key.
				selection, ok := op["title"].(string)
				if !ok {
					continue
				}

				// Hash the ansewr text for a unique key, as no actual key exists.
				hasher := md5.New()
				title, ok := op["title"].(string)
				if !ok {
					continue
				}
				hasher.Write([]byte(title))
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

	log.Dev(context, "Aggregate", "Completed : Submission[%s]", formID)
	return aggs, nil
}
