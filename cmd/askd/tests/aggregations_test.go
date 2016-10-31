// Package tests implements users tests for the API layer.
package tests

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/ardanlabs/kit/tests"
	"github.com/coralproject/shelf/cmd/askd/handlers"
	"github.com/coralproject/shelf/internal/ask/form"
)

// TestDigest tests the returning a form's digest
func TestDigest(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	t.Log("Given the need aggregate across a form's submissions.")
	{
		url := "/v1/form/580627b42600e2035218509f/digest"
		r := httptest.NewRequest("GET", url, nil)
		w := httptest.NewRecorder()

		a.ServeHTTP(w, r)

		t.Logf("\tWhen calling aggregate endpoint: %s", url)
		{
			t.Log("\tWhen we user version v1 of the aggregate endpoint.")
			if w.Code != 200 {
				t.Fatalf("\t%s\tShould be able to get the aggregation keys : %v", tests.Failed, w.Code)
			}
			t.Logf("\t%s\tShould be able to get the aggregation keys .", tests.Success)

			var fm handlers.FormDigest
			if err := json.Unmarshal(w.Body.Bytes(), &fm); err != nil {
				t.Fatalf("\t%s\tShould be able to unmarshal the results : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to unmarshal the results.", tests.Success)

			if len(fm.Questions) != 5 {
				t.Fatalf("\t%s\tShould be able to return 1 question instead of %v.", tests.Failed, len(fm.Questions))
			}
			t.Logf("\t%s\tShould be able to return 1 question.", tests.Success)
		}
	}
}

// TestAggregate tests the aggregation across a form's submissions
func TestAggregate(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	t.Log("Given the need aggregate across a form's submissions.")
	{
		url := "/v1/form/580627b42600e2035218509f/aggregate"
		r := httptest.NewRequest("GET", url, nil)
		w := httptest.NewRecorder()

		a.ServeHTTP(w, r)

		t.Logf("\tWhen calling aggregate endpoint: %s", url)
		{
			t.Log("\tWhen we user version v1 of the aggregate endpoint.")
			if w.Code != 200 {
				t.Fatalf("\t%s\tShould be able to get the aggregation keys : %v", tests.Failed, w.Code)
			}
			t.Logf("\t%s\tShould be able to get the aggregation keys .", tests.Success)

			var ak handlers.AggregationKeys
			if err := json.Unmarshal(w.Body.Bytes(), &ak); err != nil {
				t.Fatalf("\t%s\tShould be able to unmarshal the results : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to unmarshal the results.", tests.Success)

			if ak.Aggregations["all"].Count != 9 {
				t.Fatalf("\t%s\tShould be able to return 9 in Count.", tests.Failed)
			}
			t.Logf("\t%s\tShould be able to return 9 in Count.", tests.Success)
		}
	}
}

func TestAggregateGroup(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	t.Log("Given the need aggregate across a form's submissions.")
	{
		url := "/v1/form/580627b42600e2035218509f/aggregate/all"
		r := httptest.NewRequest("GET", url, nil)
		w := httptest.NewRecorder()

		a.ServeHTTP(w, r)

		t.Logf("\tWhen calling aggregate endpoint: %s", url)
		{
			t.Log("\tWhen we user version v1 of the aggregate endpoint.")
			if w.Code != 200 {
				t.Fatalf("\t%s\tShould be able to get the aggregation keys : %v", tests.Failed, w.Code)
			}
			t.Logf("\t%s\tShould be able to get the aggregation keys .", tests.Success)

			var ag form.Aggregation
			if err := json.Unmarshal(w.Body.Bytes(), &ag); err != nil {
				t.Fatalf("\t%s\tShould be able to unmarshal the results : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to unmarshal the results.", tests.Success)

			if ag.Count != 9 {
				t.Fatalf("\t%s\tShould have only one aggregation instead of %v.", tests.Failed, ag.Count)
			}
			t.Logf("\t%s\tShould have only one aggregation.", tests.Success)
		}
	}
}

func TestAggregateGroupSubmission(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	t.Log("Given the need aggregate across a form's submissions.")
	{
		url := "/v1/form/580627b42600e2035218509f/aggregate/d452b94d-e650-41c6-80af-c56091315c90/submission"
		r := httptest.NewRequest("GET", url, nil)
		w := httptest.NewRecorder()

		a.ServeHTTP(w, r)

		t.Logf("\tWhen calling aggregate endpoint: %s", url)
		{
			t.Log("\tWhen we user version v1 of the aggregate endpoint.")
			if w.Code != 200 {
				t.Fatalf("\t%s\tShould be able to get the aggregation keys : %v", tests.Failed, w.Code)
			}
			t.Logf("\t%s\tShould be able to get the aggregation keys .", tests.Success)

			var ta []form.TextAggregation
			if err := json.Unmarshal(w.Body.Bytes(), &ta); err != nil {
				t.Fatalf("\t%s\tShould be able to unmarshal the results : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to unmarshal the results.", tests.Success)
		}
	}
}
