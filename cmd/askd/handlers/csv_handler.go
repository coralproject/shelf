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

// build the header for the CSV.
// it is returning a slice to get the columns in order when building the CSV
func buildHeader(questions map[string]string) []map[string]string {

	header := []map[string]string{
		{
			"title":    "ID",
			"widgetID": "",
		},
		{
			"title":    "FormID",
			"widgetID": "",
		},
		{
			"title":    "Number",
			"widgetID": "",
		},
		{
			"title":    "Status",
			"widgetID": "",
		},
		{
			"title":    "Flags",
			"widgetID": "",
		},
		{
			"title":    "CreatedBy",
			"widgetID": "",
		},
		{
			"title":    "UpdatedBy",
			"widgetID": "",
		},
		{
			"title":    "DateCreated",
			"widgetID": "",
		},
		{
			"title":    "DateUpdated",
			"widgetID": "",
		},
	}

	for k, v := range questions {
		header = append(header, map[string]string{"title": k, "widgetID": v})
	}

	return header
}

// gets the data associated to the header from the submission and build the row
func buildRow(header []map[string]string, submission submission.Submission) ([]string, error) {

	var row []string
	var err error

	for i := 0; i < len(header); i++ {

		widgetID := header[i]["widgetID"]
		title := header[i]["title"]

		if widgetID != "" { // if the column is one of the questions.

			row = append(row, findAnswerToQuestion(submission, widgetID))

		} else { // if the column is other information (not a question) about the submission.

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
	}

	return row, nil
}

// convert a bson.M into string to display in the CSV.
// This is quite complicated code to be able to deal
// and convert any type of data that comes up in the fields. For example, multiple option/multiple choice.
func convertToString(m bson.M) string {
	var s string
	for _, val := range m {
		switch t := val.(type) {
		case string:
			if s == "" {
				s = t
			} else {
				s = fmt.Sprintf("%s, %s ", s, val)
			}
		case bson.M:
			s = fmt.Sprintf("%s, %s ", s, convertToString(t))
		case []interface{}: //map[  options: [map[index:2 title:Clarinet]] ]
			for _, option := range val.([]interface{}) {
				switch o := option.(type) {
				case bson.M:
					if _, ok := o["title"]; !ok {
						s = fmt.Sprintf("%s, %v", s, o)
					} else {
						if s == "" {
							s = o["title"].(string)
						} else {
							s = fmt.Sprintf("%s, %s", s, o["title"])
						}
					}
				default:
					s = fmt.Sprintf("%v", option)
				}

			}
		default:
			s = fmt.Sprintf("%v", val)
		}
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

// get the titles of the header
func getTitles(header []map[string]string) []string {

	var titles []string

	for v := range header {
		titles = append(titles, header[v]["title"])
	}

	return titles
}
