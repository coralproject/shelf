// Package tests implements users tests for the API layer.
package tests

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/ardanlabs/kit/tests"
	"github.com/coralproject/shelf/internal/sponge/item"
	"github.com/coralproject/shelf/internal/sponge/item/itemfix"
)

// itemPrefix is the base name for everything.
const itemPrefix = "ITEST_"

// TestRetrieveItems tests the retrieval of items.
func TestRetrieveItems(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	t.Log("Given the need get a set of items by IDs.")
	{
		url := "/1.0/item/ITEST_6eaaa19f-da7a-4095-bbe3-cee7a7631dd4,ITEST_d16790f8-13e9-4cb4-b9ef-d82835589660"
		r := tests.NewRequest("GET", url, nil)
		w := httptest.NewRecorder()

		a.ServeHTTP(w, r)

		t.Logf("\tWhen calling url : %s", url)
		{
			t.Log("\tWhen we use version 1.0 of the item endpoint.")
			if w.Code != 200 {
				t.Fatalf("\t%s\tShould be able to retrieve the items : %v", tests.Failed, w.Code)
			}
			t.Logf("\t%s\tShould be able to retrieve the items.", tests.Success)

			var items []item.Item
			if err := json.Unmarshal(w.Body.Bytes(), &items); err != nil {
				t.Fatalf("\t%s\tShould be able to unmarshal the results : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to unmarshal the results.", tests.Success)

			var count int
			for _, it := range items {
				if len(it.ID) > len(itemPrefix) && it.ID[0:len(itemPrefix)] == itemPrefix {
					count++
				}
			}

			if count != 2 {
				t.Fatalf("\t%s\tShould have retrieved two items : %d", tests.Failed, count)
			}
			t.Logf("\t%s\tShould have retrieved two items.", tests.Success)
		}
	}
}

// TestUpsertView tests the insert and update of an item.
func TestUpsertView(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	t.Log("Given the need to insert and then update an item.")
	{
		//----------------------------------------------------------------------
		// Get the fixture.

		items, err := itemfix.Get()
		if err != nil {
			t.Fatalf("\t%s\tShould be able to retrieve the fixture : %v", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to retrieve the fixture.", tests.Success)

		itemStrData, err := json.Marshal(&items[0])
		if err != nil {
			t.Fatalf("\t%s\tShould be able to marshal the fixture : %v", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to marshal the fixture.", tests.Success)

		//----------------------------------------------------------------------
		// Insert the Item.

		url := "/1.0/item"
		r := tests.NewRequest("PUT", url, bytes.NewBuffer(itemStrData))
		w := httptest.NewRecorder()

		a.ServeHTTP(w, r)

		t.Logf("\tWhen calling url to insert : %s", url)
		{
			if w.Code != 204 {
				t.Fatalf("\t%s\tShould be able to insert the item : %v", tests.Failed, w.Code)
			}
			t.Logf("\t%s\tShould be able to insert the item.", tests.Success)
		}

		//----------------------------------------------------------------------
		// Retrieve the item.

		url = "/1.0/item/" + items[0].ID
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

		//----------------------------------------------------------------------
		// Update the Item.

		items[0].Version = 2

		itemStrData, err = json.Marshal(items[0])
		if err != nil {
			t.Fatalf("\t%s\tShould be able to marshal the changed fixture : %v", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to marshal the changed fixture.", tests.Success)

		url = "/1.0/item"
		r = tests.NewRequest("PUT", url, bytes.NewBuffer(itemStrData))
		w = httptest.NewRecorder()

		a.ServeHTTP(w, r)

		t.Logf("\tWhen calling url to update : %s", url)
		{
			if w.Code != 204 {
				t.Fatalf("\t%s\tShould be able to update the item : %v", tests.Failed, w.Code)
			}
			t.Logf("\t%s\tShould be able to update the item.", tests.Success)
		}

		//----------------------------------------------------------------------
		// Retrieve the Item.

		url = "/1.0/item/" + items[0].ID
		r = tests.NewRequest("GET", url, nil)
		w = httptest.NewRecorder()

		a.ServeHTTP(w, r)

		t.Logf("\tWhen calling url to get : %s", url)
		{
			if w.Code != 200 {
				t.Fatalf("\t%s\tShould be able to retrieve the item : %v", tests.Failed, w.Code)
			}
			t.Logf("\t%s\tShould be able to retrieve the item.", tests.Success)

			var itUpdated []item.Item
			if err := json.Unmarshal(w.Body.Bytes(), &itUpdated); err != nil {
				t.Fatalf("\t%s\tShould be able to unmarshal the results : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to unmarshal the results.", tests.Success)

			if itUpdated[0].Version != 2 {
				t.Log(w.Body.String())
				t.Fatalf("\t%s\tShould get the expected result.", tests.Failed)
			}
			t.Logf("\t%s\tShould get the expected result.", tests.Success)
		}
	}
}

// TestDeleteItem tests the insert and deletion of a item.
func TestDeleteItem(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	t.Log("Given the need to delete an item.")
	{
		//----------------------------------------------------------------------
		// Delete the View.

		url := "/1.0/item/ITEST_6eaaa19f-da7a-4095-bbe3-cee7a7631dd4"
		r := tests.NewRequest("DELETE", url, nil)
		w := httptest.NewRecorder()

		a.ServeHTTP(w, r)

		t.Logf("\tWhen calling url to delete : %s", url)
		{
			if w.Code != 204 {
				t.Fatalf("\t%s\tShould be able to delete the item : %v", tests.Failed, w.Code)
			}
			t.Logf("\t%s\tShould be able to delete the item.", tests.Success)
		}

		//----------------------------------------------------------------------
		// Retrieve the Item.

		url = "/1.0/view/ITEST_6eaaa19f-da7a-4095-bbe3-cee7a7631dd4"
		r = tests.NewRequest("GET", url, nil)
		w = httptest.NewRecorder()

		a.ServeHTTP(w, r)

		t.Logf("\tWhen calling url to get : %s", url)
		{
			if w.Code != 404 {
				t.Fatalf("\t%s\tShould not be able to retrieve the item : %v", tests.Failed, w.Code)
			}
			t.Logf("\t%s\tShould not be able to retrieve the item.", tests.Success)
		}
	}
}
