// Package tests implements users tests for the API layer.
package tests

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/ardanlabs/kit/tests"
	"github.com/coralproject/xenia/internal/shelf/relationship"
	"github.com/coralproject/xenia/internal/shelf/relationship/relationshipfix"
)

// relPrefix is the base name for everything.
const relPrefix = "RTEST_"

// TestListRelationships tests the retrieval of relationships.
func TestListRelationships(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	t.Log("Given the need get a set of relationships.")
	{
		url := "/1.0/relationship"
		r := tests.NewRequest("GET", url, nil)
		w := httptest.NewRecorder()

		a.ServeHTTP(w, r)

		t.Logf("\tWhen calling url : %s", url)
		{
			t.Log("\tWhen we use version 1.0 of the relationship endpoint.")
			if w.Code != 200 {
				t.Fatalf("\t%s\tShould be able to retrieve the list of relationships : %v", tests.Failed, w.Code)
			}
			t.Logf("\t%s\tShould be able to retrieve the list of relationships.", tests.Success)

			var rels []relationship.Relationship
			if err := json.Unmarshal(w.Body.Bytes(), &rels); err != nil {
				t.Fatalf("\t%s\tShould be able to unmarshal the results : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to unmarshal the results.", tests.Success)

			var count int
			for _, rel := range rels {
				if len(rel.Predicate) > len(relPrefix) && rel.Predicate[0:len(relPrefix)] == relPrefix {
					count++
				}
			}

			if count != 2 {
				t.Fatalf("\t%s\tShould have two relationships : %d", tests.Failed, count)
			}
			t.Logf("\t%s\tShould have two relationships.", tests.Success)
		}
	}
}

// TestRetrieveRelationship tests the retrieval of a specific relationship.
func TestRetrieveRelationship(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	t.Log("Given the need to get a specific relationship.")
	{
		url := "/1.0/relationship/" + relPrefix + "authored"
		r := tests.NewRequest("GET", url, nil)
		w := httptest.NewRecorder()

		a.ServeHTTP(w, r)

		t.Logf("\tWhen calling url : %s", url)
		{
			if w.Code != 200 {
				t.Fatalf("\t%s\tShould be able to retrieve the relationship : %v", tests.Failed, w.Code)
			}
			t.Logf("\t%s\tShould be able to retrieve the relationship.", tests.Success)

			var rel relationship.Relationship
			if err := json.Unmarshal(w.Body.Bytes(), &rel); err != nil {
				t.Fatalf("\t%s\tShould be able to unmarshal the results : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to unmarshal the results.", tests.Success)

			if rel.Predicate != relPrefix+"authored" {
				t.Fatalf("\t%s\tShould have the correct relationship : %s", tests.Failed, rel.Predicate)
			}
			t.Logf("\t%s\tShould have the correct relationship.", tests.Success)
		}
	}
}

// TestUpsertRelationship tests the insert and update of a relationship.
func TestUpsertRelationship(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	t.Log("Given the need to insert and then update a relationship.")
	{
		//----------------------------------------------------------------------
		// Get the fixture.

		rels, err := relationshipfix.Get()
		if err != nil {
			t.Fatalf("\t%s\tShould be able to retrieve the fixture : %v", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to retrieve the fixture.", tests.Success)

		relStrData, err := json.Marshal(&rels[2])
		if err != nil {
			t.Fatalf("\t%s\tShould be able to marshal the fixture : %v", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to marshal the fixture.", tests.Success)

		//----------------------------------------------------------------------
		// Insert the Relationship.

		url := "/1.0/relationship"
		r := tests.NewRequest("PUT", url, bytes.NewBuffer(relStrData))
		w := httptest.NewRecorder()

		a.ServeHTTP(w, r)

		t.Logf("\tWhen calling url to insert : %s", url)
		{
			if w.Code != 204 {
				t.Fatalf("\t%s\tShould be able to insert the relationship : %v", tests.Failed, w.Code)
			}
			t.Logf("\t%s\tShould be able to insert the relationship.", tests.Success)
		}

		//----------------------------------------------------------------------
		// Retrieve the Relationship.

		url = "/1.0/relationship/" + relPrefix + "flagged"
		r = tests.NewRequest("GET", url, nil)
		w = httptest.NewRecorder()

		a.ServeHTTP(w, r)

		t.Logf("\tWhen calling url to get : %s", url)
		{
			if w.Code != 200 {
				t.Fatalf("\t%s\tShould be able to retrieve the relationship : %v", tests.Failed, w.Code)
			}
			t.Logf("\t%s\tShould be able to retrieve the relationship.", tests.Success)

			var rel relationship.Relationship
			if err := json.Unmarshal(w.Body.Bytes(), &rel); err != nil {
				t.Fatalf("\t%s\tShould be able to unmarshal the results : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to unmarshal the results.", tests.Success)

			if !reflect.DeepEqual(rels[2], rel) {
				t.Logf("\t%+v", rels[2])
				t.Logf("\t%+v", rel)
				t.Fatalf("\t%s\tShould be able to get back the same relationship.", tests.Failed)
			}
			t.Logf("\t%s\tShould be able to get back the same relationship.", tests.Success)
		}

		//----------------------------------------------------------------------
		// Update the Relationship.

		rels[2].InString = "blocked"

		relStrData, err = json.Marshal(rels[2])
		if err != nil {
			t.Fatalf("\t%s\tShould be able to marshal the changed fixture : %v", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to marshal the changed fixture.", tests.Success)

		url = "/1.0/relationship"
		r = tests.NewRequest("PUT", url, bytes.NewBuffer(relStrData))
		w = httptest.NewRecorder()

		a.ServeHTTP(w, r)

		t.Logf("\tWhen calling url to update : %s", url)
		{
			if w.Code != 204 {
				t.Fatalf("\t%s\tShould be able to update the relationship : %v", tests.Failed, w.Code)
			}
			t.Logf("\t%s\tShould be able to update the relationship.", tests.Success)
		}

		//----------------------------------------------------------------------
		// Retrieve the Relationship.

		url = "/1.0/relationship/" + relPrefix + "flagged"
		r = tests.NewRequest("GET", url, nil)
		w = httptest.NewRecorder()

		a.ServeHTTP(w, r)

		t.Logf("\tWhen calling url to get : %s", url)
		{
			if w.Code != 200 {
				t.Fatalf("\t%s\tShould be able to retrieve the relationship : %v", tests.Failed, w.Code)
			}
			t.Logf("\t%s\tShould be able to retrieve the relationship.", tests.Success)

			var relUpdated relationship.Relationship
			if err := json.Unmarshal(w.Body.Bytes(), &relUpdated); err != nil {
				t.Fatalf("\t%s\tShould be able to unmarshal the results : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to unmarshal the results.", tests.Success)

			if relUpdated.InString != "blocked" {
				t.Log(w.Body.String())
				t.Fatalf("\t%s\tShould get the expected result.", tests.Failed)
			}
			t.Logf("\t%s\tShould get the expected result.", tests.Success)
		}
	}
}

// TestDeleteRelationship tests the insert and deletion of a relationship.
func TestDeleteRelationship(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	t.Log("Given the need to delete a relationship.")
	{
		//----------------------------------------------------------------------
		// Delete the Relationship.

		url := "/1.0/relationship/" + relPrefix + "on"
		r := tests.NewRequest("DELETE", url, nil)
		w := httptest.NewRecorder()

		a.ServeHTTP(w, r)

		t.Logf("\tWhen calling url to delete : %s", url)
		{
			if w.Code != 204 {
				t.Fatalf("\t%s\tShould be able to delete the relationship : %v", tests.Failed, w.Code)
			}
			t.Logf("\t%s\tShould be able to delete the relationship.", tests.Success)
		}

		//----------------------------------------------------------------------
		// Retrieve the Relationship.

		url = "/1.0/relationship/" + relPrefix + "on"
		r = tests.NewRequest("GET", url, nil)
		w = httptest.NewRecorder()

		a.ServeHTTP(w, r)

		t.Logf("\tWhen calling url to get : %s", url)
		{
			if w.Code != 404 {
				t.Fatalf("\t%s\tShould not be able to retrieve the relationship : %v", tests.Failed, w.Code)
			}
			t.Logf("\t%s\tShould not be able to retrieve the relationship.", tests.Success)
		}
	}
}
