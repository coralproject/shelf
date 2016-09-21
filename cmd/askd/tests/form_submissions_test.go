// Package tests implements users tests for the API layer.
package tests

import (
	"encoding/csv"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"gopkg.in/mgo.v2/bson"

	"github.com/ardanlabs/kit/tests"
)

// subPrefix is the base name for everything.
const subPrefix = "FSTEST"

// TestExport tests the retrieval of a URL for a CSV file to download with the submissions.
func TextExport(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	t.Log("Given the need download submissions in a CSV format.")
	{
		url := "/v1/form/57be0437e65ada0851000002/submission/export"
		r := tests.NewRequest("GET", url, nil)
		w := httptest.NewRecorder()

		a.ServeHTTP(w, r)

		t.Logf("\tWhen calling url : %s", url)
		{
			t.Log("\tWhen we user version 1.0 of the export endpoint.")
			if w.Code != http.StatusOK {
				t.Fatalf("\t%s\tShould be able to get the URL of the file to download : %v", tests.Failed, w.Code)
			}
			t.Logf("\t%s\tShould be able to get the URL of the file to download.", tests.Success)

			var result bson.M
			if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
				t.Fatalf("\t%s\tShould be able to unmarshal the results : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to unmarshal the results.", tests.Success)

			expectedURL := "http:///v1/form/57be0437e65ada0851000002/submission/export?download=true"
			csvURL := result["csv_url"]

			if csvURL != expectedURL {
				t.Fatalf("\t%s\tShould have a different URL to download CSV : %s", tests.Failed, csvURL)
			}
			t.Logf("\t%s\tShould have a different URL to download CSV.", tests.Success)
		}
	}
}

// TestDownloadCSV tests the retrieval of a CSV file.
func TestDownloadCSV(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	t.Log("Given the need download submissions in a CSV format.")
	{
		url := "/v1/form/57be0437e65ada0851000002/submission/export?download=true"
		r := tests.NewRequest("GET", url, nil)
		w := httptest.NewRecorder()

		a.ServeHTTP(w, r)

		t.Logf("\tWhen calling CSV URL : %s", url)
		{
			t.Log("\tWhen we user version 1.0 of the export endpoint.")
			if w.Code != 200 {
				t.Fatalf("\t%s\tShould be able to get the URL of the file to download : %v", tests.Failed, w.Code)
			}
			t.Logf("\t%s\tShould be able to get the URL of the file to download.", tests.Success)

			isCSV := false
			for _, h := range w.HeaderMap["Content-Type"] {
				if h == "text/csv" {
					isCSV = true
					break
				}
			}
			if !isCSV {
				t.Fatalf("\t%s\tShould be able to get a CSV content-type file.", tests.Failed)
			}
			t.Logf("\t%s\tShould be able to get a CSV content-type file.", tests.Success)

			r := csv.NewReader(strings.NewReader(string(w.Body.Bytes())))
			records, err := r.ReadAll()
			if err != nil {
				t.Fatalf("\t%s\tShould be able to unmarshal the results : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to unmarshal the results.", tests.Success)

			expectedCount := 3
			if len(records) != expectedCount {
				t.Fatalf("\t%s\tShould have exactly %d rows but it has %d.", tests.Failed, expectedCount, len(records))
			}
			t.Logf("\t%s\tShould have exactly %d rows.", tests.Success, expectedCount)

		}
	}
}
