// Package tests implements users tests for the API layer.
package tests

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/ardanlabs/kit/tests"
	"github.com/coralproject/shelf/internal/xenia/mask"
	"github.com/coralproject/shelf/internal/xenia/mask/mfix"
)

// mCollection is the collection to use for everything.
const mCollection = "test_xenia_data"

// TestMasks tests the retrieval of masks.
func TestMasks(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	t.Log("Given the need get a set of masks.")
	{
		url := "/1.0/mask"
		r := tests.NewRequest("GET", url, nil)
		w := httptest.NewRecorder()

		a.ServeHTTP(w, r)

		t.Logf("\tWhen calling url : %s", url)
		{
			t.Log("\tWhen we user version 1.0 of the script endpoint.")
			if w.Code != 200 {
				t.Fatalf("\t%s\tShould be able to retrieve the mask list : %v", tests.Failed, w.Code)
			}
			t.Logf("\t%s\tShould be able to retrieve the mask list.", tests.Success)

			var masks map[string]mask.Mask
			if err := json.Unmarshal(w.Body.Bytes(), &masks); err != nil {
				t.Fatalf("\t%s\tShould be able to unmarshal the results : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to unmarshal the results.", tests.Success)

			var count int
			for _, scr := range masks {
				if scr.Collection == mCollection {
					count++
				}
			}

			if count != 2 {
				t.Fatalf("\t%s\tShould have two masks : %d", tests.Failed, count)
			}
			t.Logf("\t%s\tShould have two masks.", tests.Success)
		}
	}
}

// TestMaskByName tests the retrieval of a specific mask.
func TestMaskByName(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	t.Log("Given the need to get a specific mask.")
	{
		url := "/1.0/mask/" + mCollection + "/observation_time"
		r := tests.NewRequest("GET", url, nil)
		w := httptest.NewRecorder()

		a.ServeHTTP(w, r)

		t.Logf("\tWhen calling url : %s", url)
		{
			if w.Code != 200 {
				t.Fatalf("\t%s\tShould be able to retrieve the mask : %v", tests.Failed, w.Code)
			}
			t.Logf("\t%s\tShould be able to retrieve the mask.", tests.Success)

			var msk mask.Mask
			if err := json.Unmarshal(w.Body.Bytes(), &msk); err != nil {
				t.Fatalf("\t%s\tShould be able to unmarshal the results : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to unmarshal the results.", tests.Success)

			if msk.Collection != mCollection && msk.Field != "observation_time" {
				t.Fatalf("\t%s\tShould have the correct mask : %s - %s", tests.Failed, msk.Collection, msk.Field)
			}
			t.Logf("\t%s\tShould have the correct mask.", tests.Success)
		}
	}
}

// TestMaskUpsert tests the insert and update of a mask.
func TestMaskUpsert(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	t.Log("Given the need to insert and then update a mask.")
	{
		//----------------------------------------------------------------------
		// Get the fixture.

		masks, err := mfix.Get("basic.json")
		if err != nil {
			t.Fatalf("\t%s\tShould be able to retrieve the fixture : %v", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to retrieve the fixture.", tests.Success)

		masks[0].Field = "test_insert"
		masks[0].Type = "left"

		mskStrData, err := json.Marshal(masks[0])
		if err != nil {
			t.Fatalf("\t%s\tShould be able to marshal the fixture : %v", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to marshal the fixture.", tests.Success)

		//----------------------------------------------------------------------
		// Insert the Mask.

		url := "/1.0/mask"
		r := tests.NewRequest("PUT", url, bytes.NewBuffer(mskStrData))
		w := httptest.NewRecorder()

		a.ServeHTTP(w, r)

		t.Logf("\tWhen calling url to insert : %s", url)
		{
			if w.Code != 204 {
				t.Fatalf("\t%s\tShould be able to insert the mask : %v", tests.Failed, w.Code)
			}
			t.Logf("\t%s\tShould be able to insert the mask.", tests.Success)
		}

		//----------------------------------------------------------------------
		// Retrieve the Mask.

		url = "/1.0/mask/" + mCollection + "/test_insert"
		r = tests.NewRequest("GET", url, nil)
		w = httptest.NewRecorder()

		a.ServeHTTP(w, r)

		t.Logf("\tWhen calling url to get : %s", url)
		{
			if w.Code != 200 {
				t.Fatalf("\t%s\tShould be able to retrieve the mask : %v", tests.Failed, w.Code)
			}
			t.Logf("\t%s\tShould be able to retrieve the mask.", tests.Success)

			recv := tests.IndentJSON(w.Body.String())
			resp := tests.IndentJSON(`{"collection":"` + mCollection + `","field":"test_insert","type":"left"}`)

			if resp != recv {
				t.Log(resp)
				t.Log(w.Body.String())
				t.Fatalf("\t%s\tShould get the expected result.", tests.Failed)
			}
			t.Logf("\t%s\tShould get the expected result.", tests.Success)
		}

		//----------------------------------------------------------------------
		// Update the Mask.

		masks[0].Type = "right"

		mskStrData, err = json.Marshal(masks[0])
		if err != nil {
			t.Fatalf("\t%s\tShould be able to marshal the changed fixture : %v", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to marshal the changed fixture.", tests.Success)

		url = "/1.0/mask"
		r = tests.NewRequest("PUT", url, bytes.NewBuffer(mskStrData))
		w = httptest.NewRecorder()

		a.ServeHTTP(w, r)

		t.Logf("\tWhen calling url to update : %s", url)
		{
			if w.Code != 204 {
				t.Fatalf("\t%s\tShould be able to update the mask : %v", tests.Failed, w.Code)
			}
			t.Logf("\t%s\tShould be able to update the mask.", tests.Success)
		}

		//----------------------------------------------------------------------
		// Retrieve the Mask.

		url = "/1.0/mask/" + mCollection + "/test_insert"
		r = tests.NewRequest("GET", url, nil)
		w = httptest.NewRecorder()

		a.ServeHTTP(w, r)

		t.Logf("\tWhen calling url to get : %s", url)
		{
			if w.Code != 200 {
				t.Fatalf("\t%s\tShould be able to retrieve the mask : %v", tests.Failed, w.Code)
			}
			t.Logf("\t%s\tShould be able to retrieve the mask.", tests.Success)

			recv := tests.IndentJSON(w.Body.String())
			resp := tests.IndentJSON(`{"collection":"` + mCollection + `","field":"test_insert","type":"right"}`)

			if resp != recv {
				t.Log(resp)
				t.Log(w.Body.String())
				t.Fatalf("\t%s\tShould get the expected result.", tests.Failed)
			}
			t.Logf("\t%s\tShould get the expected result.", tests.Success)
		}
	}
}

// TestMaskDelete tests the insert and deletion of a mask.
func TestMaskDelete(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	t.Log("Given the need to insert and then delete a mask.")
	{
		//----------------------------------------------------------------------
		// Get the fixture.

		masks, err := mfix.Get("basic.json")
		if err != nil {
			t.Fatalf("\t%s\tShould be able to retrieve the fixture : %v", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to retrieve the fixture.", tests.Success)

		masks[0].Field = "test_delete"
		masks[0].Type = "left"

		mskStrData, err := json.Marshal(masks[0])
		if err != nil {
			t.Fatalf("\t%s\tShould be able to marshal the fixture : %v", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to marshal the fixture.", tests.Success)

		//----------------------------------------------------------------------
		// Insert the Mask.

		url := "/1.0/mask"
		r := tests.NewRequest("PUT", url, bytes.NewBuffer(mskStrData))
		w := httptest.NewRecorder()

		a.ServeHTTP(w, r)

		t.Logf("\tWhen calling url to insert : %s", url)
		{
			if w.Code != 204 {
				t.Fatalf("\t%s\tShould be able to insert the mask : %v", tests.Failed, w.Code)
			}
			t.Logf("\t%s\tShould be able to insert the mask.", tests.Success)
		}

		//----------------------------------------------------------------------
		// Retrieve the Mask.

		url = "/1.0/mask/" + mCollection + "/test_delete"
		r = tests.NewRequest("GET", url, nil)
		w = httptest.NewRecorder()

		a.ServeHTTP(w, r)

		t.Logf("\tWhen calling url to get : %s", url)
		{
			if w.Code != 200 {
				t.Fatalf("\t%s\tShould be able to retrieve the mask : %v", tests.Failed, w.Code)
			}
			t.Logf("\t%s\tShould be able to retrieve the mask.", tests.Success)

			recv := tests.IndentJSON(w.Body.String())
			resp := tests.IndentJSON(`{"collection":"` + mCollection + `","field":"test_delete","type":"left"}`)

			if resp != recv {
				t.Log(resp)
				t.Log(w.Body.String())
				t.Fatalf("\t%s\tShould get the expected result.", tests.Failed)
			}
			t.Logf("\t%s\tShould get the expected result.", tests.Success)
		}

		//----------------------------------------------------------------------
		// Delete the Mask.

		url = "/1.0/mask/" + mCollection + "/test_delete"
		r = tests.NewRequest("DELETE", url, nil)
		w = httptest.NewRecorder()

		a.ServeHTTP(w, r)

		t.Logf("\tWhen calling url to delete : %s", url)
		{
			if w.Code != 204 {
				t.Fatalf("\t%s\tShould be able to delete the mask : %v", tests.Failed, w.Code)
			}
			t.Logf("\t%s\tShould be able to delete the mask.", tests.Success)
		}

		//----------------------------------------------------------------------
		// Retrieve the Mask.

		url = "/1.0/mask/" + mCollection + "/test_delete"
		r = tests.NewRequest("GET", url, nil)
		w = httptest.NewRecorder()

		a.ServeHTTP(w, r)

		t.Logf("\tWhen calling url to get : %s", url)
		{
			if w.Code != 404 {
				t.Fatalf("\t%s\tShould Not be able to retrieve the mask : %v", tests.Failed, w.Code)
			}
			t.Logf("\t%s\tShould Not be able to retrieve the mask.", tests.Success)
		}
	}
}
