// Package tests implements users tests for the API layer.
package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/ardanlabs/kit/tests"
	"github.com/coralproject/shelf/internal/wire/pattern"
	"github.com/coralproject/shelf/internal/wire/pattern/patternfix"
)

// patternPrefix is the base name for everything.
const patternPrefix = "PTEST_"

// TestListPatterns tests the retrieval of patterns.
func TestListPatterns(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	t.Log("Given the need get a set of patterns.")
	{
		url := "/1.0/pattern"
		r := tests.NewRequest("GET", url, nil)
		w := httptest.NewRecorder()

		a.ServeHTTP(w, r)

		t.Logf("\tWhen calling url : %s", url)
		{
			t.Log("\tWhen we use version 1.0 of the pattern endpoint.")
			if w.Code != 200 {
				t.Fatalf("\t%s\tShould be able to retrieve the list of patterns : %v", tests.Failed, w.Code)
			}
			t.Logf("\t%s\tShould be able to retrieve the list of patterns.", tests.Success)

			var ps []pattern.Pattern
			if err := json.Unmarshal(w.Body.Bytes(), &ps); err != nil {
				t.Fatalf("\t%s\tShould be able to unmarshal the results : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to unmarshal the results.", tests.Success)

			var count int
			for _, p := range ps {
				if len(p.Type) > len(patternPrefix) && p.Type[0:len(patternPrefix)] == patternPrefix {
					count++
				}
			}

			if count != 2 {
				fmt.Println(ps)
				t.Fatalf("\t%s\tShould have two patterns : %d", tests.Failed, count)
			}
			t.Logf("\t%s\tShould have two patterns.", tests.Success)
		}
	}
}

// TestRetrievePattern tests the retrieval of a specific pattern.
func TestRetrievePattern(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	t.Log("Given the need to get a specific pattern.")
	{
		url := "/1.0/pattern/" + patternPrefix + "comment"
		r := tests.NewRequest("GET", url, nil)
		w := httptest.NewRecorder()

		a.ServeHTTP(w, r)

		t.Logf("\tWhen calling url : %s", url)
		{
			if w.Code != 200 {
				t.Fatalf("\t%s\tShould be able to retrieve the pattern : %v", tests.Failed, w.Code)
			}
			t.Logf("\t%s\tShould be able to retrieve the pattern.", tests.Success)

			var p pattern.Pattern
			if err := json.Unmarshal(w.Body.Bytes(), &p); err != nil {
				t.Fatalf("\t%s\tShould be able to unmarshal the results : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to unmarshal the results.", tests.Success)

			if p.Type != patternPrefix+"comment" {
				t.Fatalf("\t%s\tShould have the correct pattern : %s", tests.Failed, p.Type)
			}
			t.Logf("\t%s\tShould have the correct pattern.", tests.Success)
		}
	}
}

// TestUpsertPattern tests the insert and update of a pattern.
func TestUpsertPattern(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	t.Log("Given the need to insert and then update a pattern.")
	{
		//----------------------------------------------------------------------
		// Get the fixture.

		ps, _, err := patternfix.Get()
		if err != nil {
			t.Fatalf("\t%s\tShould be able to retrieve the fixture : %v", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to retrieve the fixture.", tests.Success)

		pStrData, err := json.Marshal(&ps[2])
		if err != nil {
			t.Fatalf("\t%s\tShould be able to marshal the fixture : %v", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to marshal the fixture.", tests.Success)

		//----------------------------------------------------------------------
		// Insert the Pattern.

		url := "/1.0/pattern"
		r := tests.NewRequest("PUT", url, bytes.NewBuffer(pStrData))
		w := httptest.NewRecorder()

		a.ServeHTTP(w, r)

		t.Logf("\tWhen calling url to insert : %s", url)
		{
			if w.Code != 204 {
				t.Fatalf("\t%s\tShould be able to insert the pattern : %v", tests.Failed, w.Code)
			}
			t.Logf("\t%s\tShould be able to insert the pattern.", tests.Success)
		}

		//----------------------------------------------------------------------
		// Retrieve the Pattern.

		url = "/1.0/pattern/" + patternPrefix + "asset"
		r = tests.NewRequest("GET", url, nil)
		w = httptest.NewRecorder()

		a.ServeHTTP(w, r)

		t.Logf("\tWhen calling url to get : %s", url)
		{
			if w.Code != 200 {
				t.Fatalf("\t%s\tShould be able to retrieve the pattern : %v", tests.Failed, w.Code)
			}
			t.Logf("\t%s\tShould be able to retrieve the pattern.", tests.Success)

			var p pattern.Pattern
			if err := json.Unmarshal(w.Body.Bytes(), &p); err != nil {
				t.Fatalf("\t%s\tShould be able to unmarshal the results : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to unmarshal the results.", tests.Success)

			if !reflect.DeepEqual(ps[2], p) {
				t.Logf("\t%+v", ps[2])
				t.Logf("\t%+v", p)
				t.Fatalf("\t%s\tShould be able to get back the same pattern.", tests.Failed)
			}
			t.Logf("\t%s\tShould be able to get back the same relationship.", tests.Success)
		}

		//----------------------------------------------------------------------
		// Update the Pattern.

		ps[2].Type = patternPrefix + "article"

		pStrData, err = json.Marshal(ps[2])
		if err != nil {
			t.Fatalf("\t%s\tShould be able to marshal the changed fixture : %v", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to marshal the changed fixture.", tests.Success)

		url = "/1.0/pattern"
		r = tests.NewRequest("PUT", url, bytes.NewBuffer(pStrData))
		w = httptest.NewRecorder()

		a.ServeHTTP(w, r)

		t.Logf("\tWhen calling url to update : %s", url)
		{
			if w.Code != 204 {
				t.Fatalf("\t%s\tShould be able to update the pattern : %v", tests.Failed, w.Code)
			}
			t.Logf("\t%s\tShould be able to update the pattern.", tests.Success)
		}

		//----------------------------------------------------------------------
		// Retrieve the Pattern.

		url = "/1.0/pattern/" + patternPrefix + "article"
		r = tests.NewRequest("GET", url, nil)
		w = httptest.NewRecorder()

		a.ServeHTTP(w, r)

		t.Logf("\tWhen calling url to get : %s", url)
		{
			if w.Code != 200 {
				t.Fatalf("\t%s\tShould be able to retrieve the pattern : %v", tests.Failed, w.Code)
			}
			t.Logf("\t%s\tShould be able to retrieve the pattern.", tests.Success)

			var pUpdated pattern.Pattern
			if err := json.Unmarshal(w.Body.Bytes(), &pUpdated); err != nil {
				t.Fatalf("\t%s\tShould be able to unmarshal the results : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to unmarshal the results.", tests.Success)

			if pUpdated.Type != patternPrefix+"article" {
				t.Log(w.Body.String())
				t.Fatalf("\t%s\tShould get the expected result.", tests.Failed)
			}
			t.Logf("\t%s\tShould get the expected result.", tests.Success)
		}
	}
}

// TestDeletePattern tests the insert and deletion of a pattern.
func TestDeletePattern(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	t.Log("Given the need to delete a pattern.")
	{
		//----------------------------------------------------------------------
		// Delete the Pattern.

		url := "/1.0/pattern/" + patternPrefix + "comment"
		r := tests.NewRequest("DELETE", url, nil)
		w := httptest.NewRecorder()

		a.ServeHTTP(w, r)

		t.Logf("\tWhen calling url to delete : %s", url)
		{
			if w.Code != 204 {
				t.Fatalf("\t%s\tShould be able to delete the pattern : %v", tests.Failed, w.Code)
			}
			t.Logf("\t%s\tShould be able to delete the pattern.", tests.Success)
		}

		//----------------------------------------------------------------------
		// Retrieve the Pattern.

		url = "/1.0/pattern/" + patternPrefix + "comment"
		r = tests.NewRequest("GET", url, nil)
		w = httptest.NewRecorder()

		a.ServeHTTP(w, r)

		t.Logf("\tWhen calling url to get : %s", url)
		{
			if w.Code != 404 {
				t.Fatalf("\t%s\tShould not be able to retrieve the pattern : %v", tests.Failed, w.Code)
			}
			t.Logf("\t%s\tShould not be able to retrieve the pattern.", tests.Success)
		}
	}
}
