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

	// Add the values of the header to the CSV.
	rows = append(rows, getValues(header))

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

// build the header for the CSV.
func buildHeader(questions map[string]string) map[string]string {

	header := map[string]string{
		"ID":          "",
		"FormID":      "",
		"Number":      "",
		"Status":      "",
		"Flags":       "",
		"CreatedBy":   "",
		"UpdatedBy":   "",
		"DateCreated": "",
		"DateUpdated": "",
	}
	for k, v := range questions {
		header[k] = v
	}

	return header
}

// gets the data associated to the header from the submission and build the row
func buildRow(header map[string]string, submission submission.Submission) ([]string, error) {

	var row []string
	var err error

	for head, widgetID := range header {

		if widgetID != "" { // if the column is one of the questions.

			row = append(row, findAnswerToQuestion(submission, widgetID))

		} else { // if the column is other information (not a question) about the submission.

			var value string
			v := reflect.ValueOf(submission)
			switch t := reflect.Indirect(v).FieldByName(head).Interface().(type) {
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
				err = fmt.Errorf("Type not found for field %v. Value: %v", head, t)
				return nil, err
			}

			row = append(row, value)
		}
	}

	return row, nil
}

// convert a bson.M into string to display in the CSV.
func convertToString(m bson.M) string {

	var s string
	for _, val := range m {
		s = s + fmt.Sprintf("%v ", val)
	}

	return s
}

// find the Answer for the specific question in the submission. It returns an empty string if it does not find it.
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

// get values of a map.
func getValues(m map[string]string) []string {

	var values []string

	for v := range m {
		values = append(values, v)
	}

	return values
}
