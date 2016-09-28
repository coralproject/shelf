// Package tests implements users tests for the API layer.
package tests

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ardanlabs/kit/tests"
	"github.com/coralproject/shelf/internal/sponge/item"
)

// TestItemsGET sample test for the GET call.
func TestItemsGET(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	t.Log("Given the need to test GET item call.")
	{
		url := "/v1/view/all/1/query/High_quality_commenters_in_Politics"
		r := tests.NewRequest("GET", url, nil)
		w := httptest.NewRecorder()

		a.ServeHTTP(w, r)

		t.Logf("\tWhen calling url : %s", url)
		{
			t.Log("\tWhen we user version v1 of the items endpoint.")
			if w.Code != http.StatusOK {
				t.Fatalf("\t%s\tShould be able to retrieve all the items for the query set High_quality_commenters_in_Politics: %v", tests.Failed, w.Code)
			}
			t.Logf("\t%s\tShould be able to retrieve the items for the query set.", tests.Success)

			var results map[string]interface{}

			err := json.NewDecoder(w.Body).Decode(&results)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to unmarshal the response : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to unmarshal the response.", tests.Success)

			items, ok := results["results"].([]interface{})[0].(map[string]interface{})["Docs"].([]interface{})
			if !ok {
				t.Errorf("\t%s\tShould have the correct type.", tests.Failed)
			}

			total := 4
			if len(items) != total {
				t.Log("GOT :", len(items))
				t.Log("WANT:", total)
				t.Errorf("\t%s\tShould have the correct amount of items.", tests.Failed)
			} else {
				t.Logf("\t%s\tShould have the correct amount of items.", tests.Success)
			}

			want := "56ce763b1fefce879aa0bb75"
			if items[0].(map[string]interface{})["_id"] != want {
				t.Log("GOT :", items[0].(map[string]interface{})["_id"])
				t.Log("WANT:", want)
				t.Errorf("\t%s\tShould have the correct id.", tests.Failed)
			} else {
				t.Logf("\t%s\tShould have the correct id.", tests.Success)
			}
		}
	}
}

// TestItemPUT sample test for the PUT call.
func TestItemPut(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	t.Log("Given the need to test PUT item call.")
	{
		url := "/v1/item"
		i := item.Item{Type: "comments", Version: 1, Data: map[string]interface{}{"content": "This is a new comment", "author_id": 1}}
		payload, err := json.Marshal(i)
		if err != nil {
			log.Fatal("\t%s\tShould be able to marshal the item: %v", tests.Failed, err)
		}
		r := tests.NewRequest("PUT", url, bytes.NewReader(payload))
		w := httptest.NewRecorder()

		a.ServeHTTP(w, r)

		t.Logf("\tWhen calling url : %s", url)
		{
			t.Log("\tWhen we user version v1 of the items endpoint.")
			if w.Code != http.StatusOK {
				t.Fatalf("\t%s\tShould be able to save the item: %v", tests.Failed, w.Code)
			}
			t.Logf("\t%s\tShould be able to save the item.", tests.Success)

			var result item.Item

			err := json.NewDecoder(w.Body).Decode(&result)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to unmarshal the response : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to unmarshal the response.", tests.Success)

			itemID := result.ID
			if itemID == "" {
				t.Errorf("\t%s\tShould have the correct type.", tests.Failed)
			}
		}
	}
}
