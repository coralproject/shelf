// Package tests implements users tests for the API layer.
package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ardanlabs/kit/tests"
)

// TestFormsGET sample test for the GET call.
func TestFormsGET(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	t.Log("Given the need to test GET form call.")
	{
		url := "/v1/form"
		r := tests.NewRequest("GET", url, nil)
		w := httptest.NewRecorder()

		a.ServeHTTP(w, r)

		t.Logf("\tWhen calling url : %s", url)
		{
			t.Log("\tWhen we user version 1.0 of the forms endpoint.")
			if w.Code != http.StatusOK {
				t.Fatalf("\t%s\tShould be able to retrieve the forms list : %v", tests.Failed, w.Code)
			}
			t.Logf("\t%s\tShould be able to retrieve the forms list.", tests.Success)

			var forms []struct {
				ID string `json:"id"`
			}

			err := json.NewDecoder(w.Body).Decode(&forms)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to unmarshal the response : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to unmarshal the response.", tests.Success)

			total := 141
			if len(forms) != total {
				t.Log("GOT :", len(forms))
				t.Log("WANT:", total)
				t.Errorf("\t%s\tShould have the correct amount of forms.", tests.Failed)
			} else {
				t.Logf("\t%s\tShould have the correct amount of forms.", tests.Success)
			}

			want := "5790f40c6413f60007228586"
			if forms[0].ID != want {
				t.Log("GOT :", forms[0].ID)
				t.Log("WANT:", want)
				t.Errorf("\t%s\tShould have the correct id.", tests.Failed)
			} else {
				t.Logf("\t%s\tShould have the correct id.", tests.Success)
			}
		}
	}
}
