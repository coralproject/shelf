// Package tests implements users tests for the API layer.
package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/ardanlabs/kit/tests"

	"github.com/coralproject/shelf/internal/sponge/item"
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

		//----------------------------------------------------------------------
		// Retrieve the item.

		var items []item.Item
		url = "/1.0/item/" + fmt.Sprintf("%s_%v", typ, dat["id"])
		r = tests.NewRequest("GET", url, nil)
		w = httptest.NewRecorder()

		a.ServeHTTP(w, r)

		t.Logf("\tWhen calling url to get : %s", url)
		{
			if w.Code != 200 {
				t.Fatalf("\t%s\tShould be able to retrieve the item : %v", tests.Failed, w.Code)
			}
			t.Logf("\t%s\tShould be able to retrieve the item.", tests.Success)

			var itemsBack []item.Item
			if err := json.Unmarshal(w.Body.Bytes(), &itemsBack); err != nil {
				t.Fatalf("\t%s\tShould be able to unmarshal the results : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to unmarshal the results.", tests.Success)

			if itemsBack[0].ID != items[0].ID || itemsBack[0].Type != items[0].Type {
				t.Logf("\t%+v", items[0])
				t.Logf("\t%+v", itemsBack[0])
				t.Fatalf("\t%s\tShould be able to get back the same item.", tests.Failed)
			}
			t.Logf("\t%s\tShould be able to get back the same item.", tests.Success)
		}
	}
}
