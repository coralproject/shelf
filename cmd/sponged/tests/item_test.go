// Package tests implements users tests for the API layer.
package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ardanlabs/kit/cfg"
	"github.com/ardanlabs/kit/tests"
	"github.com/cayleygraph/cayley"
	"github.com/cayleygraph/cayley/quad"
	cayleyshelf "github.com/coralproject/shelf/internal/platform/db/cayley"
	"github.com/coralproject/shelf/internal/sponge/item"
	"github.com/coralproject/shelf/internal/sponge/item/itemfix"
)

// itemPrefix is the base name for items.
const itemPrefix = "ITEST_"

// setup initializes for each indivdual test.
func setup(t *testing.T) *cayley.Handle {
	tests.ResetLog()

	store, err := cayleyshelf.New(cfg.MustURL("MONGO_URI").String(), nil)
	if err != nil {
		t.Fatalf("\t%s\tShould be able to connect to the cayley graph : %s", tests.Failed, err)
	}

	return store
}

// teardown deinitializes for each indivdual test.
func teardown(t *testing.T, store *cayley.Handle) {
	store.Close()
	tests.DisplayLog()
}

// TestRetrieveItems tests the retrieval of items.
func TestRetrieveItems(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	t.Log("Given the need get a set of items by IDs.")
	{
		url := "/v1/item/ITEST_6eaaa19f-da7a-4095-bbe3-cee7a7631dd4,ITEST_d16790f8-13e9-4cb4-b9ef-d82835589660"
		r := httptest.NewRequest("GET", url, nil)
		w := httptest.NewRecorder()

		a.ServeHTTP(w, r)

		t.Logf("\tWhen calling url : %s", url)
		{
			t.Log("\tWhen we use version v1 of the item endpoint.")
			if w.Code != http.StatusOK {
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

// TestUpsertItem tests the insert and update of an item.
func TestUpsertItem(t *testing.T) {
	store := setup(t)
	defer teardown(t, store)

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

		url := "/v1/item"
		r := httptest.NewRequest("PUT", url, bytes.NewBuffer(itemStrData))
		w := httptest.NewRecorder()

		a.ServeHTTP(w, r)

		t.Logf("\tWhen calling url to insert : %s", url)
		{
			if w.Code != http.StatusOK {
				t.Fatalf("\t%s\tShould be able to insert the item : %v", tests.Failed, w.Code)
			}
			t.Logf("\t%s\tShould be able to insert the item.", tests.Success)
		}

		//----------------------------------------------------------------------
		// Check the inferred relationship.

		p := cayley.StartPath(store, quad.String("ITEST_80aa936a-f618-4234-a7be-df59a14cf8de")).Out(quad.String("authored"))
		it, _ := p.BuildIterator().Optimize()
		defer it.Close()
		for it.Next() {
			token := it.Result()
			value := store.NameOf(token)
			if quad.NativeOf(value) != "ITEST_d1dfa366-d2f7-4a4a-a64f-af89d4c97d82" {
				t.Fatalf("\t%s\tShould be able to get the inferred relationships from the graph", tests.Failed)
			}
		}
		if err := it.Err(); err != nil {
			t.Fatalf("\t%s\tShould be able to get the inferred relationships from the graph : %s", tests.Failed, err)
		}
		it.Close()
		t.Logf("\t%s\tShould be able to get the inferred relationships from the graph.", tests.Success)

		//----------------------------------------------------------------------
		// Retrieve the item.

		url = "/v1/item/" + items[0].ID
		r = httptest.NewRequest("GET", url, nil)
		w = httptest.NewRecorder()

		a.ServeHTTP(w, r)

		t.Logf("\tWhen calling url to get : %s", url)
		{
			if w.Code != http.StatusOK {
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

		url = "/v1/item"
		r = httptest.NewRequest("PUT", url, bytes.NewBuffer(itemStrData))
		w = httptest.NewRecorder()

		a.ServeHTTP(w, r)

		t.Logf("\tWhen calling url to update : %s", url)
		{
			if w.Code != http.StatusOK {
				t.Fatalf("\t%s\tShould be able to update the item : %v", tests.Failed, w.Code)
			}
			t.Logf("\t%s\tShould be able to update the item.", tests.Success)
		}

		//----------------------------------------------------------------------
		// Retrieve the Item.

		url = "/v1/item/" + items[0].ID
		r = httptest.NewRequest("GET", url, nil)
		w = httptest.NewRecorder()

		a.ServeHTTP(w, r)

		t.Logf("\tWhen calling url to get : %s", url)
		{
			if w.Code != http.StatusOK {
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
	store := setup(t)
	defer teardown(t, store)

	t.Log("Given the need to delete an item.")
	{
		//----------------------------------------------------------------------
		// Delete the Item.

		url := "/v1/item/ITEST_d1dfa366-d2f7-4a4a-a64f-af89d4c97d82"
		r := httptest.NewRequest("DELETE", url, nil)
		w := httptest.NewRecorder()

		a.ServeHTTP(w, r)

		t.Logf("\tWhen calling url to delete : %s", url)
		{
			if w.Code != http.StatusNoContent {
				t.Fatalf("\t%s\tShould be able to delete the item : %v", tests.Failed, w.Code)
			}
			t.Logf("\t%s\tShould be able to delete the item.", tests.Success)
		}

		//----------------------------------------------------------------------
		// Retrieve the Item.

		url = "/v1/view/ITEST_d1dfa366-d2f7-4a4a-a64f-af89d4c97d82"
		r = httptest.NewRequest("GET", url, nil)
		w = httptest.NewRecorder()

		a.ServeHTTP(w, r)

		t.Logf("\tWhen calling url to get : %s", url)
		{
			if w.Code != 404 {
				t.Fatalf("\t%s\tShould not be able to retrieve the item : %v", tests.Failed, w.Code)
			}
			t.Logf("\t%s\tShould not be able to retrieve the item.", tests.Success)
		}

		//----------------------------------------------------------------------
		// Check the inferred relationships.

		p := cayley.StartPath(store, quad.String("ITEST_80aa936a-f618-4234-a7be-df59a14cf8de")).Out(quad.String("authored"))
		it, _ := p.BuildIterator().Optimize()
		defer it.Close()

		var count int
		for it.Next() {
			count++
		}
		if err := it.Err(); err != nil {
			t.Fatalf("\t%s\tShould be able to confirm removed relationships : %s", tests.Failed, err)
		}

		if count > 0 {
			t.Fatalf("\t%s\tShould be able to confirm removed relationships.", tests.Failed)
		}
		t.Logf("\t%s\tShould be able to confirm removed relationships.", tests.Success)

	}
}
