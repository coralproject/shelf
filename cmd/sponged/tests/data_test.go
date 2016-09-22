// Package tests implements users tests for the API layer.
package tests

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/ardanlabs/kit/tests"
	"github.com/coralproject/shelf/internal/sponge/item/itemfix"
)

// TestUpsertItem tests the insert and update of an item.
func DataUpsertItem(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	t.Log("Given the need to insert items received as data.")
	{
		//----------------------------------------------------------------------
		// Get the fixture.

		dat, err := itemfix.GetData("data.json")
		if err != nil {
			t.Fatalf("\t%s\tShould be able to retrieve the fixture : %v", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to retrieve the fixture.", tests.Success)

		itemStrData, err := json.Marshal(dat)
		if err != nil {
			t.Fatalf("\t%s\tShould be able to marshal the fixture : %v", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to marshal the fixture.", tests.Success)

		//----------------------------------------------------------------------
		// Insert the Item.

		typ := "test_type"
		url := "/1.0/data/" + typ
		r := tests.NewRequest("POST", url, bytes.NewBuffer(itemStrData))
		w := httptest.NewRecorder()

		a.ServeHTTP(w, r)

		t.Logf("\tWhen calling url to insert : %s", url)
		{
			if w.Code != 204 {
				t.Fatalf("\t%s\tShould be able to insert the data : %v", tests.Failed, w.Code)
			}
			t.Logf("\t%s\tShould be able to insert the data.", tests.Success)
		}
	}
}
