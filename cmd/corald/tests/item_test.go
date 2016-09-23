// Package tests implements users tests for the API layer.
package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ardanlabs/kit/tests"
)

// TestItemsGET sample test for the GET call.
func TestItemsGET(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	t.Log("Given the need to test GET form call.")
	{
		url := "/v1/item/all/1/High_quality_commenters_in_Politics"
		r := tests.NewRequest("GET", url, nil)
		w := httptest.NewRecorder()

		a.ServeHTTP(w, r)

		t.Logf("\tWhen calling url : %s", url)
		{
			t.Log("\tWhen we user version 1.0 of the items endpoint.")
			if w.Code != http.StatusOK {
				t.Fatalf("\t%s\tShould be able to retrieve all the items for the query set High_quality_commenters_in_Politics.", tests.Failed)
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

// // TestItemPOST sample test for the POST call.
// func TestItemPOST(t *testing.T) {
// 	tests.ResetLog()
// 	defer tests.DisplayLog()
//
// 	t.Log("Given the need to test POST form call.")
// 	{
// 		url := "/v1/item"
// 		r := tests.NewRequest("POST", url, nil)
// 		w := httptest.NewRecorder()
//
// 		a.ServeHTTP(w, r)
//
// 		t.Logf("\tWhen calling url : %s", url)
// 		{
// 			t.Log("\tWhen we user version 1.0 of the forms endpoint.")
// 			if w.Code != http.StatusCreated {
// 				t.Fatalf("\t%s\tShould be able to retrieve the forms list : %v", tests.Failed, w.Code)
// 			}
// 			t.Logf("\t%s\tShould be able to retrieve the forms list.", tests.Success)
//
// 			var form struct {
// 				ID string `json:"id"`
// 			}
//
// 			err := json.NewDecoder(w.Body).Decode(&form)
// 			if err != nil {
// 				t.Fatalf("\t%s\tShould be able to unmarshal the response : %v", tests.Failed, err)
// 			}
// 			t.Logf("\t%s\tShould be able to unmarshal the response.", tests.Success)
//
// 			want := "57daa6680cefe53d4b2adce0"
// 			if form.ID != want {
// 				t.Log("GOT :", form.ID)
// 				t.Log("WANT:", want)
// 				t.Errorf("\t%s\tShould have the correct id.", tests.Failed)
// 			} else {
// 				t.Logf("\t%s\tShould have the correct id.", tests.Success)
// 			}
// 		}
// 	}
// }

// TestItemsPUT sample test for the PUT call.
func TestItemsPUT(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	t.Log("Given the need to test POST form call.")
	{
		url := "/v1/item/1"
		r := tests.NewRequest("PUT", url, nil)
		w := httptest.NewRecorder()

		a.ServeHTTP(w, r)

		t.Logf("\tWhen calling url : %s", url)
		{
			t.Log("\tWhen we user version 1.0 of the forms endpoint.")
			if w.Code != http.StatusNoContent {
				t.Fatalf("\t%s\tShould be able to retrieve the forms list : %v", tests.Failed, w.Code)
			}
			t.Logf("\t%s\tShould be able to retrieve the forms list.", tests.Success)
		}
	}
}
