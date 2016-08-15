// Package tests implements users tests for the API layer.
package tests

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/ardanlabs/kit/tests"
	"github.com/coralproject/xenia/internal/wire/view"
	"github.com/coralproject/xenia/internal/wire/view/viewfix"
)

// viewPrefix is the base name for everything.
const viewPrefix = "VTEST_"

// TestListViews tests the retrieval of views.
func TestListViews(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	t.Log("Given the need get a set of views.")
	{
		url := "/1.0/view"
		r := tests.NewRequest("GET", url, nil)
		w := httptest.NewRecorder()

		a.ServeHTTP(w, r)

		t.Logf("\tWhen calling url : %s", url)
		{
			t.Log("\tWhen we use version 1.0 of the view endpoint.")
			if w.Code != 200 {
				t.Fatalf("\t%s\tShould be able to retrieve the list of views : %v", tests.Failed, w.Code)
			}
			t.Logf("\t%s\tShould be able to retrieve the list of views.", tests.Success)

			var views []view.View
			if err := json.Unmarshal(w.Body.Bytes(), &views); err != nil {
				t.Fatalf("\t%s\tShould be able to unmarshal the results : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to unmarshal the results.", tests.Success)

			var count int
			for _, v := range views {
				if len(v.Name) > len(viewPrefix) && v.Name[0:len(viewPrefix)] == viewPrefix {
					count++
				}
			}

			if count != 2 {
				t.Fatalf("\t%s\tShould have two views : %d", tests.Failed, count)
			}
			t.Logf("\t%s\tShould have two views.", tests.Success)
		}
	}
}

// TestRetrieveView tests the retrieval of a specific view.
func TestRetrieveView(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	t.Log("Given the need to get a specific view.")
	{
		url := "/1.0/view/" + viewPrefix + "thread"
		r := tests.NewRequest("GET", url, nil)
		w := httptest.NewRecorder()

		a.ServeHTTP(w, r)

		t.Logf("\tWhen calling url : %s", url)
		{
			if w.Code != 200 {
				t.Fatalf("\t%s\tShould be able to retrieve the view : %v", tests.Failed, w.Code)
			}
			t.Logf("\t%s\tShould be able to retrieve the view.", tests.Success)

			var v view.View
			if err := json.Unmarshal(w.Body.Bytes(), &v); err != nil {
				t.Fatalf("\t%s\tShould be able to unmarshal the results : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to unmarshal the results.", tests.Success)

			if v.Name != viewPrefix+"thread" {
				t.Fatalf("\t%s\tShould have the correct view : %s", tests.Failed, v.Name)
			}
			t.Logf("\t%s\tShould have the correct view.", tests.Success)
		}
	}
}

// TestUpsertView tests the insert and update of a view.
func TestUpsertView(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	t.Log("Given the need to insert and then update a view.")
	{
		//----------------------------------------------------------------------
		// Get the fixture.

		views, err := viewfix.Get()
		if err != nil {
			t.Fatalf("\t%s\tShould be able to retrieve the fixture : %v", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to retrieve the fixture.", tests.Success)

		viewStrData, err := json.Marshal(&views[2])
		if err != nil {
			t.Fatalf("\t%s\tShould be able to marshal the fixture : %v", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to marshal the fixture.", tests.Success)

		//----------------------------------------------------------------------
		// Insert the View.

		url := "/1.0/view"
		r := tests.NewRequest("PUT", url, bytes.NewBuffer(viewStrData))
		w := httptest.NewRecorder()

		a.ServeHTTP(w, r)

		t.Logf("\tWhen calling url to insert : %s", url)
		{
			if w.Code != 204 {
				t.Fatalf("\t%s\tShould be able to insert the view : %v", tests.Failed, w.Code)
			}
			t.Logf("\t%s\tShould be able to insert the view.", tests.Success)
		}

		//----------------------------------------------------------------------
		// Retrieve the View.

		url = "/1.0/view/" + viewPrefix + "comments%20from%20authors%20flagged%20by%20a%20user"
		r = tests.NewRequest("GET", url, nil)
		w = httptest.NewRecorder()

		a.ServeHTTP(w, r)

		t.Logf("\tWhen calling url to get : %s", url)
		{
			if w.Code != 200 {
				t.Fatalf("\t%s\tShould be able to retrieve the view : %v", tests.Failed, w.Code)
			}
			t.Logf("\t%s\tShould be able to retrieve the view.", tests.Success)

			var v view.View
			if err := json.Unmarshal(w.Body.Bytes(), &v); err != nil {
				t.Fatalf("\t%s\tShould be able to unmarshal the results : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to unmarshal the results.", tests.Success)

			if !reflect.DeepEqual(views[2], v) {
				t.Logf("\t%+v", views[2])
				t.Logf("\t%+v", v)
				t.Fatalf("\t%s\tShould be able to get back the same view.", tests.Failed)
			}
			t.Logf("\t%s\tShould be able to get back the same view.", tests.Success)
		}

		//----------------------------------------------------------------------
		// Update the View.

		views[2].Name = viewPrefix + "better_name"

		viewStrData, err = json.Marshal(views[2])
		if err != nil {
			t.Fatalf("\t%s\tShould be able to marshal the changed fixture : %v", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to marshal the changed fixture.", tests.Success)

		url = "/1.0/view"
		r = tests.NewRequest("PUT", url, bytes.NewBuffer(viewStrData))
		w = httptest.NewRecorder()

		a.ServeHTTP(w, r)

		t.Logf("\tWhen calling url to update : %s", url)
		{
			if w.Code != 204 {
				t.Fatalf("\t%s\tShould be able to update the view : %v", tests.Failed, w.Code)
			}
			t.Logf("\t%s\tShould be able to update the view.", tests.Success)
		}

		//----------------------------------------------------------------------
		// Retrieve the View.

		url = "/1.0/view/" + viewPrefix + "better_name"
		r = tests.NewRequest("GET", url, nil)
		w = httptest.NewRecorder()

		a.ServeHTTP(w, r)

		t.Logf("\tWhen calling url to get : %s", url)
		{
			if w.Code != 200 {
				t.Fatalf("\t%s\tShould be able to retrieve the view : %v", tests.Failed, w.Code)
			}
			t.Logf("\t%s\tShould be able to retrieve the view.", tests.Success)

			var vUpdated view.View
			if err := json.Unmarshal(w.Body.Bytes(), &vUpdated); err != nil {
				t.Fatalf("\t%s\tShould be able to unmarshal the results : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to unmarshal the results.", tests.Success)

			if vUpdated.Name != "VTEST_better_name" {
				t.Log(w.Body.String())
				t.Fatalf("\t%s\tShould get the expected result.", tests.Failed)
			}
			t.Logf("\t%s\tShould get the expected result.", tests.Success)
		}
	}
}

// TestDeleteView tests the insert and deletion of a view.
func TestDeleteView(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	t.Log("Given the need to delete a view.")
	{
		//----------------------------------------------------------------------
		// Delete the View.

		url := "/1.0/view/" + viewPrefix + "thread"
		r := tests.NewRequest("DELETE", url, nil)
		w := httptest.NewRecorder()

		a.ServeHTTP(w, r)

		t.Logf("\tWhen calling url to delete : %s", url)
		{
			if w.Code != 204 {
				t.Fatalf("\t%s\tShould be able to delete the view : %v", tests.Failed, w.Code)
			}
			t.Logf("\t%s\tShould be able to delete the view.", tests.Success)
		}

		//----------------------------------------------------------------------
		// Retrieve the View.

		url = "/1.0/view/" + viewPrefix + "thread"
		r = tests.NewRequest("GET", url, nil)
		w = httptest.NewRecorder()

		a.ServeHTTP(w, r)

		t.Logf("\tWhen calling url to get : %s", url)
		{
			if w.Code != 404 {
				t.Fatalf("\t%s\tShould not be able to retrieve the view : %v", tests.Failed, w.Code)
			}
			t.Logf("\t%s\tShould not be able to retrieve the view.", tests.Success)
		}
	}
}
