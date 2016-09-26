package handlers

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/coralproject/shelf/internal/ask/form/submission"
	"gopkg.in/mgo.v2/bson"
)

// encodeSubmissionsToCSV retrieves the questions and answers on a specific submission.
// we are returning a [][]string to use the resolut in returing the data in CSV.
func encodeSubmissionsToCSV(submissions []submission.Submission) ([]byte, error) {

	// Get all the questions for the submissions to the form. It returns a map that contains the widget_id and question.
	questions := make(map[string]string, 0)

	for _, s := range submissions {
		for _, r := range s.Answers {
			questions[r.Question] = r.WidgetID
		}
	}

	// Build the header with the columns that we need.
	header := buildHeader(questions)

	var rows [][]string

	// Add the titles of the header to the CSV.
	rows = append(rows, getTitles(header))

	// Build the rows.
	for _, s := range submissions {

		row, err := buildRow(header, s)
		if err != nil {
			return nil, err
		}

		// add the row to the CSV data
		rows = append(rows, row)
	}

	// Marshal the data into a CSV string.
	var buf bytes.Buffer
	w := csv.NewWriter(&buf)

	w.WriteAll(rows)

	if err := w.Error(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// buildHeader builds the header for the CSV.
// It is returning a slice to get the columns in order when building the CSV.
func buildHeader(questions map[string]string) []map[string]string {
	header := []map[string]string{
		{"title": "FormID"},
		{"title": "ID"},
		{"title": "Status"},
		{"title": "Flags"},
		{"title": "CreatedBy"},
		{"title": "UpdatedBy"},
		{"title": "DateCreated"},
		{"title": "DateUpdated"},
	}

	// Adds the questions to the map header.
	for k, v := range questions {
		header = append(header, map[string]string{"title": k, "widgetID": v})
	}

	return header
}

// buildRow gets the data associated to the header from the submission and build the row of the CSV.
// It returns the row to add tot he CSV.
func buildRow(header []map[string]string, submission submission.Submission) ([]string, error) {
	var row []string
	var err error

	for i := 0; i < len(header); i++ {
		widgetID := header[i]["widgetID"]
		title := header[i]["title"]

		if widgetID != "" { // If the column is one of the questions.
			row = append(row, findAnswerToQuestion(submission, widgetID))
			continue
		}

		// If the column is other information (not a question) about the submission.
		var value string
		v := reflect.ValueOf(submission)
		switch t := reflect.Indirect(v).FieldByName(title).Interface().(type) {
		case string:
			value = t
		case int:
			value = strconv.Itoa(t)
		case []string:
			value = strings.Join(t, ", ")
		case time.Time:
			value = t.String()
		case bson.ObjectId:
			value = t.Hex()
		case nil:
			value = ""
		default:
			err = fmt.Errorf("Type not found for field %v. Value: %v", title, t)
			return nil, err
		}
		row = append(row, value)
	}

	return row, nil
}

// convertToString convert a bson.M into string to display in the CSV.
// This is quite complicated code to be able to deal
// and convert any type of data that comes up in the fields. For example,
// multiple option/multiple choice.
func convertToString(doc bson.M) string {

	// Construct an array of strings to collect our bson properties
	var strs []string

	// Loop over the document.
	for _, val := range doc {

		switch docValue := val.(type) {

		// If the value is a string, just append it.
		case string:
			strs = append(strs, docValue)

		// If the value is another bson.M document, then recurse.
		case bson.M:
			strs = append(strs, convertToString(docValue))

		// If the value is an array of documents, then range of that one.
		case []interface{}:
			// Loop over the
			for _, subDoc := range docValue {

				switch subDocValue := subDoc.(type) {

				// If this is another doc, then walk into it.
				case bson.M:

					// If the doc has a property called title, we need to select that one
					// instead of another field.
					if title, ok := subDocValue["title"]; ok {
						if titleString, ok := title.(string); ok {
							strs = append(strs, titleString)
							continue
						}
					}

					// The inner object is a bson.M object, we should keep extracting it.
					strs = append(strs, convertToString(subDocValue))

				default:
					strs = append(strs, fmt.Sprintf("%v", subDocValue))
				}
			}

		default:
			strs = append(strs, fmt.Sprintf("%v", val))
		}
	}

	return strings.Join(strs, ", ")
}

// findAnswerToQuestion finds the answer for the specific question in the submission.
// It returns an empty string if it does not find it.
func findAnswerToQuestion(s submission.Submission, widgetID string) string {
	for _, r := range s.Answers {
		if r.WidgetID == widgetID {
			switch t := r.Answer.(type) {
			case bson.M:
				return convertToString(t)
			case string:
				return t
			default:
				return fmt.Sprintf("%v", t)
			}
		}
	}

	return ""
}

// getTitles get the titles of the header
func getTitles(header []map[string]string) []string {
	var titles []string
	for v := range header {
		titles = append(titles, header[v]["title"])
	}

	return titles
}
