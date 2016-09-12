// Package tests implements users tests for the API layer.
package tests

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/ardanlabs/kit/tests"
	"github.com/coralproject/shelf/internal/xenia/query"
	"github.com/coralproject/shelf/internal/xenia/query/qfix"
)

// qPrefix is the base name for everything.
const qPrefix = "QTEST_O"

// TestQuerySets tests the retrieval of query sets.
func TestQuerySets(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	t.Log("Given the need get a set of query sets.")
	{
		url := "/1.0/query"
		r := tests.NewRequest("GET", url, nil)
		w := httptest.NewRecorder()

		a.ServeHTTP(w, r)

		t.Logf("\tWhen calling url : %s", url)
		{
			t.Log("\tWhen we user version 1.0 of the query endpoint.")
			if w.Code != 200 {
				t.Fatalf("\t%s\tShould be able to retrieve the query set list : %v", tests.Failed, w.Code)
			}
			t.Logf("\t%s\tShould be able to retrieve the query set list.", tests.Success)

			var sets []query.Set
			if err := json.Unmarshal(w.Body.Bytes(), &sets); err != nil {
				t.Fatalf("\t%s\tShould be able to unmarshal the results : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to unmarshal the results.", tests.Success)

			var count int
			for _, set := range sets {
				if len(set.Name) > len(qPrefix) && set.Name[0:len(qPrefix)] == qPrefix {
					count++
				}
			}

			if count != 2 {
				t.Fatalf("\t%s\tShould have two query sets : %d", tests.Failed, count)
			}
			t.Logf("\t%s\tShould have two query sets.", tests.Success)
		}
	}
}

// TestQueryByName tests the retrieval of a specific query.
func TestQueryByName(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	t.Log("Given the need to get a specific query.")
	{
		url := "/1.0/query/" + qPrefix + "_basic"
		r := tests.NewRequest("GET", url, nil)
		w := httptest.NewRecorder()

		a.ServeHTTP(w, r)

		t.Logf("\tWhen calling url : %s", url)
		{
			if w.Code != 200 {
				t.Fatalf("\t%s\tShould be able to retrieve the query : %v", tests.Failed, w.Code)
			}
			t.Logf("\t%s\tShould be able to retrieve the query.", tests.Success)

			var set query.Set
			if err := json.Unmarshal(w.Body.Bytes(), &set); err != nil {
				t.Fatalf("\t%s\tShould be able to unmarshal the results : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to unmarshal the results.", tests.Success)

			if set.Name != qPrefix+"_basic" {
				t.Fatalf("\t%s\tShould have the correct set : %s", tests.Failed, set.Name)
			}
			t.Logf("\t%s\tShould have the correct set.", tests.Success)
		}
	}
}

// TestQueryUpsert tests the insert and update of a set.
func TestQueryUpsert(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	t.Log("Given the need to insert and then update a set.")
	{
		//----------------------------------------------------------------------
		// Get the fixture.

		qs, err := qfix.Get("upsert.json")
		if err != nil {
			t.Fatalf("\t%s\tShould be able to retrieve the fixture : %v", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to retrieve the fixture.", tests.Success)

		qsStrData, err := json.Marshal(&qs)
		if err != nil {
			t.Fatalf("\t%s\tShould be able to marshal the fixture : %v", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to marshal the fixture.", tests.Success)

		//----------------------------------------------------------------------
		// Insert the Set.

		url := "/1.0/query"
		r := tests.NewRequest("PUT", url, bytes.NewBuffer(qsStrData))
		w := httptest.NewRecorder()

		a.ServeHTTP(w, r)

		t.Logf("\tWhen calling url to insert : %s", url)
		{
			if w.Code != 204 {
				t.Fatalf("\t%s\tShould be able to insert the set : %v", tests.Failed, w.Code)
			}
			t.Logf("\t%s\tShould be able to insert the set.", tests.Success)
		}

		//----------------------------------------------------------------------
		// Ensure the indexes.

		url = "/1.0/index/QTEST_O_upsert"
		r = tests.NewRequest("PUT", url, nil)
		w = httptest.NewRecorder()

		a.ServeHTTP(w, r)

		t.Logf("\tWhen calling url to ensure indexes : %s", url)
		{
			if w.Code != 204 {
				t.Fatalf("\t%s\tShould be able to ensure indexes for the set : %v", tests.Failed, w.Code)
			}
			t.Logf("\t%s\tShould be able to ensure indexes for the set.", tests.Success)
		}

		//----------------------------------------------------------------------
		// Retrieve the Set.

		url = "/1.0/query/" + qPrefix + "_upsert"
		r = tests.NewRequest("GET", url, nil)
		w = httptest.NewRecorder()

		a.ServeHTTP(w, r)

		t.Logf("\tWhen calling url to get : %s", url)
		{
			if w.Code != 200 {
				t.Fatalf("\t%s\tShould be able to retrieve the set : %v", tests.Failed, w.Code)
			}
			t.Logf("\t%s\tShould be able to retrieve the set.", tests.Success)

			recv := tests.IndentJSON(w.Body.String())
			resp := tests.IndentJSON(`{"name":"` + qPrefix + `_upsert","desc":"","pre_script":"","pst_script":"","params":[],"queries":[{"name":"Upsert","type":"pipeline","collection":"test_xenia_data","commands":[{"$match":{"station.d":"42021"}},{"$project":{"_id":0,"name":1}}],"indexes":[{"key":["station_id"],"unique":true}],"return":true}],"enabled":true,"explain":false}`)

			if resp != recv {
				t.Log(resp)
				t.Log(w.Body.String())
				t.Fatalf("\t%s\tShould get the expected result.", tests.Failed)
			}
			t.Logf("\t%s\tShould get the expected result.", tests.Success)
		}

		//----------------------------------------------------------------------
		// Update the Set.

		qs.Description = "C"

		qsStrData, err = json.Marshal(&qs)
		if err != nil {
			t.Fatalf("\t%s\tShould be able to marshal the changed fixture : %v", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to marshal the changed fixture.", tests.Success)

		url = "/1.0/query"
		r = tests.NewRequest("PUT", url, bytes.NewBuffer(qsStrData))
		w = httptest.NewRecorder()

		a.ServeHTTP(w, r)

		t.Logf("\tWhen calling url to update : %s", url)
		{
			if w.Code != 204 {
				t.Fatalf("\t%s\tShould be able to update the set : %v", tests.Failed, w.Code)
			}
			t.Logf("\t%s\tShould be able to update the set.", tests.Success)
		}

		//----------------------------------------------------------------------
		// Retrieve the Set.

		url = "/1.0/query/" + qPrefix + "_upsert"
		r = tests.NewRequest("GET", url, nil)
		w = httptest.NewRecorder()

		a.ServeHTTP(w, r)

		t.Logf("\tWhen calling url to get : %s", url)
		{
			if w.Code != 200 {
				t.Fatalf("\t%s\tShould be able to retrieve the set : %v", tests.Failed, w.Code)
			}
			t.Logf("\t%s\tShould be able to retrieve the set.", tests.Success)

			recv := tests.IndentJSON(w.Body.String())
			resp := tests.IndentJSON(`{"name":"` + qPrefix + `_upsert","desc":"C","pre_script":"","pst_script":"","params":[],"queries":[{"name":"Upsert","type":"pipeline","collection":"test_xenia_data","commands":[{"$match":{"station.d":"42021"}},{"$project":{"_id":0,"name":1}}],"indexes":[{"key":["station_id"],"unique":true}],"return":true}],"enabled":true,"explain":false}`)

			if resp != recv {
				t.Log(resp)
				t.Log(w.Body.String())
				t.Fatalf("\t%s\tShould get the expected result.", tests.Failed)
			}
			t.Logf("\t%s\tShould get the expected result.", tests.Success)
		}
	}
}

// TestQueryDelete tests the insert and deletion of a set.
func TestQueryDelete(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	t.Log("Given the need to insert and then delete a set.")
	{
		//----------------------------------------------------------------------
		// Get the fixture.

		qs, err := qfix.Get("upsert.json")
		if err != nil {
			t.Fatalf("\t%s\tShould be able to retrieve the fixture : %v", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to retrieve the fixture.", tests.Success)

		qsStrData, err := json.Marshal(&qs)
		if err != nil {
			t.Fatalf("\t%s\tShould be able to marshal the fixture : %v", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to marshal the fixture.", tests.Success)

		//----------------------------------------------------------------------
		// Insert the Set.

		url := "/1.0/query"
		r := tests.NewRequest("PUT", url, bytes.NewBuffer(qsStrData))
		w := httptest.NewRecorder()

		a.ServeHTTP(w, r)

		t.Logf("\tWhen calling url to insert : %s", url)
		{
			if w.Code != 204 {
				t.Fatalf("\t%s\tShould be able to insert the set : %v", tests.Failed, w.Code)
			}
			t.Logf("\t%s\tShould be able to insert the set.", tests.Success)
		}

		//----------------------------------------------------------------------
		// Retrieve the Set.

		url = "/1.0/query/" + qPrefix + "_upsert"
		r = tests.NewRequest("GET", url, nil)
		w = httptest.NewRecorder()

		a.ServeHTTP(w, r)

		t.Logf("\tWhen calling url to get : %s", url)
		{
			if w.Code != 200 {
				t.Fatalf("\t%s\tShould be able to retrieve the set : %v", tests.Failed, w.Code)
			}
			t.Logf("\t%s\tShould be able to retrieve the set.", tests.Success)

			recv := tests.IndentJSON(w.Body.String())
			resp := tests.IndentJSON(`{"name":"` + qPrefix + `_upsert","desc":"","pre_script":"","pst_script":"","params":[],"queries":[{"name":"Upsert","type":"pipeline","collection":"test_xenia_data","commands":[{"$match":{"station.d":"42021"}},{"$project":{"_id":0,"name":1}}],"indexes":[{"key":["station_id"],"unique":true}],"return":true}],"enabled":true,"explain":false}`)

			if resp != recv {
				t.Log(resp)
				t.Log(w.Body.String())
				t.Fatalf("\t%s\tShould get the expected result.", tests.Failed)
			}
			t.Logf("\t%s\tShould get the expected result.", tests.Success)
		}

		//----------------------------------------------------------------------
		// Delete the Set.

		url = "/1.0/query/" + qPrefix + "_upsert"
		r = tests.NewRequest("DELETE", url, nil)
		w = httptest.NewRecorder()

		a.ServeHTTP(w, r)

		t.Logf("\tWhen calling url to delete : %s", url)
		{
			if w.Code != 204 {
				t.Fatalf("\t%s\tShould be able to delete the set : %v", tests.Failed, w.Code)
			}
			t.Logf("\t%s\tShould be able to delete the set.", tests.Success)
		}

		//----------------------------------------------------------------------
		// Retrieve the Set.

		url = "/1.0/query/" + qPrefix + "_upsert"
		r = tests.NewRequest("GET", url, nil)
		w = httptest.NewRecorder()

		a.ServeHTTP(w, r)

		t.Logf("\tWhen calling url to get : %s", url)
		{
			if w.Code != 404 {
				t.Fatalf("\t%s\tShould Not be able to retrieve the set : %v", tests.Failed, w.Code)
			}
			t.Logf("\t%s\tShould Not be able to retrieve the set.", tests.Success)
		}
	}
}
