// Package endpoint implements users tests for the API layer.
package tests

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/coralproject/xenia/app/xenia/routes"
	"github.com/coralproject/xenia/pkg/query"
	"github.com/coralproject/xenia/pkg/query/qfix"
	"github.com/coralproject/xenia/tstdata"

	"github.com/ardanlabs/kit/cfg"
	"github.com/ardanlabs/kit/db"
	"github.com/ardanlabs/kit/db/mongo"
	"github.com/ardanlabs/kit/tests"
	"github.com/ardanlabs/kit/web/app"
)

var a *app.App

func init() {
	tests.Init("XENIA")

	cfg := mongo.Config{
		Host:     cfg.MustString("MONGO_HOST"),
		AuthDB:   cfg.MustString("MONGO_AUTHDB"),
		DB:       cfg.MustString("MONGO_DB"),
		User:     cfg.MustString("MONGO_USER"),
		Password: cfg.MustString("MONGO_PASS"),
	}
	tests.InitMongo(cfg)

	a = routes.API().(*app.App)
}

//==============================================================================

// TestMain helps to clean up the test data.
func TestMain(m *testing.M) {
	db := db.NewMGO()
	defer db.CloseMGO()

	tstdata.Generate(db)
	defer tstdata.Drop()

	loadQuery(db, "basic.json")
	loadQuery(db, "basic_var.json")
	defer qfix.Remove(db)

	m.Run()
}

// loadQuery adds queries to run tests.
func loadQuery(db *db.DB, file string) error {
	qs1, err := qfix.Get(file)
	if err != nil {
		return err
	}

	if err := qfix.Add(db, qs1); err != nil {
		return err
	}

	return nil
}

//==============================================================================

// TestQuerySets tests the retrieval of query sets.
func TestQuerySets(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	t.Log("Given the need get a set of query sets.")
	{
		url := "/1.0/query"
		r := tests.NewRequest("GET", url, nil)
		w := httptest.NewRecorder()

		a.ServeHTTP(w, r)

		t.Logf("\tWhen calling url : %s", url)
		{
			t.Log("\tWhen we user version 1.0 of the query endpoint.")
			if w.Code != 200 {
				t.Fatalf("\t%s\tShould be able to retrieve the query set list : %v", tests.Failed, w.Code)
			}
			t.Logf("\t%s\tShould be able to retrieve the query set list.", tests.Success)

			var sets []query.Set
			if err := json.Unmarshal(w.Body.Bytes(), &sets); err != nil {
				t.Fatalf("\t%s\tShould be able to unmarshal the results : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to unmarshal the results.", tests.Success)

			var count int
			for _, set := range sets {
				if set.Name[0:5] == "QTEST" {
					count++
				}
			}

			if count != 2 {
				t.Fatalf("\t%s\tShould have two query sets : %d", tests.Failed, count)
			}
			t.Logf("\t%s\tShould have two query sets.", tests.Success)
		}
	}
}

// TestQueryByName tests the retrieval of a specific query.
func TestQueryByName(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	t.Log("Given the need to get a specific query.")
	{
		url := "/1.0/query/QTEST_basic"
		r := tests.NewRequest("GET", url, nil)
		w := httptest.NewRecorder()

		a.ServeHTTP(w, r)

		t.Logf("\tWhen calling url : %s", url)
		{
			if w.Code != 200 {
				t.Fatalf("\t%s\tShould be able to retrieve the query : %v", tests.Failed, w.Code)
			}
			t.Logf("\t%s\tShould be able to retrieve the query.", tests.Success)

			var set query.Set
			if err := json.Unmarshal(w.Body.Bytes(), &set); err != nil {
				t.Fatalf("\t%s\tShould be able to unmarshal the results : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to unmarshal the results.", tests.Success)

			if set.Name != "QTEST_basic" {
				t.Fatalf("\t%s\tShould have the correct set : %s", tests.Failed, set.Name)
			}
			t.Logf("\t%s\tShould have the correct set.", tests.Success)
		}
	}
}

// TestQueryExec tests the execution of a specific query.
func TestQueryExec(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	t.Log("Given the need to execute a specific query.")
	{
		url := "/1.0/exec/QTEST_basic?station_id=42021"
		r := tests.NewRequest("GET", url, nil)
		w := httptest.NewRecorder()

		a.ServeHTTP(w, r)

		t.Logf("\tWhen calling url : %s", url)
		{
			if w.Code != 200 {
				t.Fatalf("\t%s\tShould be able to retrieve the query : %v", tests.Failed, w.Code)
			}
			t.Logf("\t%s\tShould be able to retrieve the query.", tests.Success)

			recv := w.Body.String()
			resp := `{"results":[{"Name":"Basic","Docs":[{"name":"C14 - Pasco County Buoy, FL"}]}],"error":false}`

			if resp != recv {
				t.Log(resp)
				t.Log(recv)
				t.Fatalf("\t%s\tShould get the expected result.", tests.Failed)
			}
			t.Logf("\t%s\tShould get the expected result.", tests.Success)
		}
	}
}

// TestQueryExecCustom tests the execution of a custom query.
func TestQueryExecCustom(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	t.Log("Given the need to execute a custom query.")
	{
		qs, err := qfix.Get("basic.json")
		if err != nil {
			t.Fatalf("\t%s\tShould be able to retrieve the fixture : %v", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to retrieve the fixture.", tests.Success)

		qsStrData, err := json.Marshal(&qs)
		if err != nil {
			t.Fatalf("\t%s\tShould be able to marshal the fixture : %v", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to marshal the fixture.", tests.Success)

		url := "/1.0/exec"
		r := tests.NewRequest("POST", url, bytes.NewBuffer(qsStrData))
		w := httptest.NewRecorder()

		a.ServeHTTP(w, r)

		t.Logf("\tWhen calling url : %s", url)
		{
			if w.Code != 200 {
				t.Fatalf("\t%s\tShould be able to retrieve the query : %v", tests.Failed, w.Code)
			}
			t.Logf("\t%s\tShould be able to retrieve the query.", tests.Success)

			recv := w.Body.String()
			resp := `{"results":[{"Name":"Basic","Docs":[{"name":"C14 - Pasco County Buoy, FL"}]}],"error":false}`

			if resp != recv {
				t.Log(resp)
				t.Log(recv)
				t.Fatalf("\t%s\tShould get the expected result.", tests.Failed)
			}
			t.Logf("\t%s\tShould get the expected result.", tests.Success)
		}
	}
}

// TestQueryExecJSONP tests the execution of a specific query using JSONP.
func TestQueryExecJSONP(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	t.Log("Given the need to execute a specific query with JSONP output.")
	{
		url := "/1.0/exec/QTEST_basic?station_id=42021&callback=handle_data"
		r := tests.NewRequest("GET", url, nil)
		w := httptest.NewRecorder()

		a.ServeHTTP(w, r)

		t.Logf("\tWhen calling url : %s", url)
		{
			if w.Code != 200 {
				t.Fatalf("\t%s\tShould be able to retrieve the query : %v", tests.Failed, w.Code)
			}
			t.Logf("\t%s\tShould be able to retrieve the query.", tests.Success)

			recv := w.Body.String()
			resp := `handle_data({"results":[{"Name":"Basic","Docs":[{"name":"C14 - Pasco County Buoy, FL"}]}],"error":false})`

			if resp != recv {
				t.Log(resp)
				t.Log(recv)
				t.Fatalf("\t%s\tShould get the expected result.", tests.Failed)
			}
			t.Logf("\t%s\tShould get the expected result.", tests.Success)
		}
	}
}

// TestQueryUpsert tests the insert and update of a set.
func TestQueryUpsert(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	t.Log("Given the need to insert and then update a set.")
	{
		//----------------------------------------------------------------------
		// Get the fixture.

		qs, err := qfix.Get("upsert.json")
		if err != nil {
			t.Fatalf("\t%s\tShould be able to retrieve the fixture : %v", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to retrieve the fixture.", tests.Success)

		qsStrData, err := json.Marshal(&qs)
		if err != nil {
			t.Fatalf("\t%s\tShould be able to marshal the fixture : %v", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to marshal the fixture.", tests.Success)

		//----------------------------------------------------------------------
		// Insert the Set.

		url := "/1.0/query"
		r := tests.NewRequest("PUT", url, bytes.NewBuffer(qsStrData))
		w := httptest.NewRecorder()

		a.ServeHTTP(w, r)

		t.Logf("\tWhen calling url to insert : %s", url)
		{
			if w.Code != 204 {
				t.Fatalf("\t%s\tShould be able to insert the set : %v", tests.Failed, w.Code)
			}
			t.Logf("\t%s\tShould be able to insert the set.", tests.Success)
		}

		//----------------------------------------------------------------------
		// Retrieve the Set.

		url = "/1.0/query/QTEST_upsert"
		r = tests.NewRequest("GET", url, nil)
		w = httptest.NewRecorder()

		a.ServeHTTP(w, r)

		t.Logf("\tWhen calling url to get : %s", url)
		{
			if w.Code != 200 {
				t.Fatalf("\t%s\tShould be able to retrieve the set : %v", tests.Failed, w.Code)
			}
			t.Logf("\t%s\tShould be able to retrieve the set.", tests.Success)

			recv := w.Body.String()
			resp := `{"name":"QTEST_upsert","desc":"","pre_script":"","pst_script":"","params":[],"queries":[{"name":"Upsert","type":"pipeline","collection":"test_xenia_data","scripts":["{\"$match\": {\"station_id\" : \"42021\"}}","{\"$project\": {\"_id\": 0, \"name\": 1}}"],"return":true}],"enabled":true}`

			if resp != recv {
				t.Log(resp)
				t.Log(w.Body.String())
				t.Fatalf("\t%s\tShould get the expected result.", tests.Failed)
			}
			t.Logf("\t%s\tShould get the expected result.", tests.Success)
		}

		//----------------------------------------------------------------------
		// Update the Set.

		qs.Description = "C"

		qsStrData, err = json.Marshal(&qs)
		if err != nil {
			t.Fatalf("\t%s\tShould be able to marshal the changed fixture : %v", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to marshal the changed fixture.", tests.Success)

		url = "/1.0/query"
		r = tests.NewRequest("PUT", url, bytes.NewBuffer(qsStrData))
		w = httptest.NewRecorder()

		a.ServeHTTP(w, r)

		t.Logf("\tWhen calling url to update : %s", url)
		{
			if w.Code != 204 {
				t.Fatalf("\t%s\tShould be able to update the set : %v", tests.Failed, w.Code)
			}
			t.Logf("\t%s\tShould be able to update the set.", tests.Success)
		}

		//----------------------------------------------------------------------
		// Retrieve the Set.

		url = "/1.0/query/QTEST_upsert"
		r = tests.NewRequest("GET", url, nil)
		w = httptest.NewRecorder()

		a.ServeHTTP(w, r)

		t.Logf("\tWhen calling url to get : %s", url)
		{
			if w.Code != 200 {
				t.Fatalf("\t%s\tShould be able to retrieve the set : %v", tests.Failed, w.Code)
			}
			t.Logf("\t%s\tShould be able to retrieve the set.", tests.Success)

			recv := w.Body.String()
			resp := `{"name":"QTEST_upsert","desc":"C","pre_script":"","pst_script":"","params":[],"queries":[{"name":"Upsert","type":"pipeline","collection":"test_xenia_data","scripts":["{\"$match\": {\"station_id\" : \"42021\"}}","{\"$project\": {\"_id\": 0, \"name\": 1}}"],"return":true}],"enabled":true}`

			if resp != recv {
				t.Log(resp)
				t.Log(w.Body.String())
				t.Fatalf("\t%s\tShould get the expected result.", tests.Failed)
			}
			t.Logf("\t%s\tShould get the expected result.", tests.Success)
		}
	}
}

// TestQueryDelete tests the insert and deletion of a set.
func TestQueryDelete(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	t.Log("Given the need to insert and then delete a set.")
	{
		//----------------------------------------------------------------------
		// Get the fixture.

		qs, err := qfix.Get("upsert.json")
		if err != nil {
			t.Fatalf("\t%s\tShould be able to retrieve the fixture : %v", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to retrieve the fixture.", tests.Success)

		qsStrData, err := json.Marshal(&qs)
		if err != nil {
			t.Fatalf("\t%s\tShould be able to marshal the fixture : %v", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to marshal the fixture.", tests.Success)

		//----------------------------------------------------------------------
		// Insert the Set.

		url := "/1.0/query"
		r := tests.NewRequest("PUT", url, bytes.NewBuffer(qsStrData))
		w := httptest.NewRecorder()

		a.ServeHTTP(w, r)

		t.Logf("\tWhen calling url to insert : %s", url)
		{
			if w.Code != 204 {
				t.Fatalf("\t%s\tShould be able to insert the set : %v", tests.Failed, w.Code)
			}
			t.Logf("\t%s\tShould be able to insert the set.", tests.Success)
		}

		//----------------------------------------------------------------------
		// Retrieve the Set.

		url = "/1.0/query/QTEST_upsert"
		r = tests.NewRequest("GET", url, nil)
		w = httptest.NewRecorder()

		a.ServeHTTP(w, r)

		t.Logf("\tWhen calling url to get : %s", url)
		{
			if w.Code != 200 {
				t.Fatalf("\t%s\tShould be able to retrieve the set : %v", tests.Failed, w.Code)
			}
			t.Logf("\t%s\tShould be able to retrieve the set.", tests.Success)

			recv := w.Body.String()
			resp := `{"name":"QTEST_upsert","desc":"","pre_script":"","pst_script":"","params":[],"queries":[{"name":"Upsert","type":"pipeline","collection":"test_xenia_data","scripts":["{\"$match\": {\"station_id\" : \"42021\"}}","{\"$project\": {\"_id\": 0, \"name\": 1}}"],"return":true}],"enabled":true}`

			if resp != recv {
				t.Log(resp)
				t.Log(w.Body.String())
				t.Fatalf("\t%s\tShould get the expected result.", tests.Failed)
			}
			t.Logf("\t%s\tShould get the expected result.", tests.Success)
		}

		//----------------------------------------------------------------------
		// Delete the Set.

		url = "/1.0/query/QTEST_upsert"
		r = tests.NewRequest("DELETE", url, nil)
		w = httptest.NewRecorder()

		a.ServeHTTP(w, r)

		t.Logf("\tWhen calling url to delete : %s", url)
		{
			if w.Code != 204 {
				t.Fatalf("\t%s\tShould be able to delete the set : %v", tests.Failed, w.Code)
			}
			t.Logf("\t%s\tShould be able to delete the set.", tests.Success)
		}

		//----------------------------------------------------------------------
		// Retrieve the Set.

		url = "/1.0/query/QTEST_upsert"
		r = tests.NewRequest("GET", url, nil)
		w = httptest.NewRecorder()

		a.ServeHTTP(w, r)

		t.Logf("\tWhen calling url to get : %s", url)
		{
			if w.Code != 404 {
				t.Fatalf("\t%s\tShould Not be able to retrieve the set : %v", tests.Failed, w.Code)
			}
			t.Logf("\t%s\tShould Not be able to retrieve the set.", tests.Success)
		}
	}
}
