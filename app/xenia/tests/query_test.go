// Package endpoint implements users tests for the API layer.
package tests

import (
	"net/http/httptest"
	"testing"

	"github.com/coralproject/xenia/app/xenia/routes"
	"github.com/coralproject/xenia/pkg/query"

	"github.com/ardanlabs/kit/db"
	"github.com/ardanlabs/kit/tests"
	"github.com/ardanlabs/kit/web/app"
)

var a *app.App

func init() {
	tests.Init("XENIA")
	tests.InitMongo()

	a = routes.API().(*app.App)
}

// TestMain helps to clean up the test data.
func TestMain(m *testing.M) {
	db := db.NewMGO()
	defer db.CloseMGO()

	query.GenerateTestData(db)
	defer query.DropTestData()

	loadQuery(db, "basic.json")
	loadQuery(db, "basic_var.json")
	defer query.RemoveTestSets(db)

	m.Run()
}

//==============================================================================

// TestQueryNames tests the retrieval of query names.
func TestQueryNames(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	url := "/1.0/query"
	r := tests.NewRequest("GET", "/1.0/query", nil)
	w := httptest.NewRecorder()

	a.ServeHTTP(w, r)

	t.Logf("\tWhen calling url : %s", url)
	{
		t.Log("\tWhen we user version 1.0 of the query endpoint.")
		if w.Code != 200 {
			t.Fatalf("\t%s\tShould be able to retrieve the query list : %v", tests.Failed, w.Code)
		}
		t.Logf("\t%s\tShould be able to retrieve the query list.", tests.Success)
	}
}

// TestQueryByName tests the retrieval of a specific query.
func TestQueryByName(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	url := "/1.0/query/QTEST_basic"
	r := tests.NewRequest("GET", "/1.0/query/QTEST_basic", nil)
	w := httptest.NewRecorder()

	a.ServeHTTP(w, r)

	t.Log("Given the need get a specific query.")
	{
		t.Logf("\tWhen calling url : %s", url)
		if w.Code != 200 {
			t.Fatalf("\t%s\tShould be able to retrieve the query : %v", tests.Failed, w.Code)
		}
		t.Logf("\t%s\tShould be able to retrieve the query.", tests.Success)

		resp := `{"name":"QTEST_basic","desc":"","enabled":true,"params":[],"queries":[{"name":"Basic","type":"pipeline","collection":"test_query","save":true,"scripts":["{\"$match\": {\"station_id\" : \"42021\"}}","{\"$project\": {\"_id\": 0, \"name\": 1}}"]}]}`
		if resp[0:245] != w.Body.String()[0:245] {
			t.Log(resp)
			t.Log(w.Body.String())
			t.Fatalf("\t%s\tShould get the expected result.", tests.Failed)
		}
		t.Logf("\t%s\tShould get the expected result.", tests.Success)
	}
}

// TestQueryExec tests the execution of a specific query.
func TestQueryExec(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	url := "/1.0/query/QTEST_basic/exec?station_id=42021"
	r := tests.NewRequest("GET", url, nil)
	w := httptest.NewRecorder()

	a.ServeHTTP(w, r)

	t.Log("Given the need get a specific query.")
	{
		t.Logf("\tWhen calling url : %s", url)
		if w.Code != 200 {
			t.Fatalf("\t%s\tShould be able to retrieve the query : %v", tests.Failed, w.Code)
		}
		t.Logf("\t%s\tShould be able to retrieve the query.", tests.Success)

		resp := `{"results":[{"Name":"Basic","Docs":[{"name":"C14 - Pasco County Buoy, FL"}]}],"error":false}`
		if resp[0:92] != w.Body.String()[0:92] {
			t.Log(resp)
			t.Log(w.Body.String())
			t.Fatalf("\t%s\tShould get the expected result.", tests.Failed)
		}
		t.Logf("\t%s\tShould get the expected result.", tests.Success)
	}
}

// TestQueryExecJSONP tests the execution of a specific query using JSONP.
func TestQueryExecJSONP(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	url := "/1.0/query/QTEST_basic/exec?station_id=42021&callback=handle_data"
	r := tests.NewRequest("GET", url, nil)
	w := httptest.NewRecorder()

	a.ServeHTTP(w, r)

	t.Log("Given the need get a specific query.")
	{
		t.Logf("\tWhen calling url : %s", url)
		if w.Code != 200 {
			t.Fatalf("\t%s\tShould be able to retrieve the query : %v", tests.Failed, w.Code)
		}
		t.Logf("\t%s\tShould be able to retrieve the query.", tests.Success)

		resp := `handle_data({"results":[{"Name":"Basic","Docs":[{"name":"C14 - Pasco County Buoy, FL"}]}],"error":false})`
		if resp[0:92] != w.Body.String()[0:92] {
			t.Log(resp)
			t.Log(w.Body.String())
			t.Fatalf("\t%s\tShould get the expected result.", tests.Failed)
		}
		t.Logf("\t%s\tShould get the expected result.", tests.Success)
	}
}

//==============================================================================

// loadQuery adds queries to run tests.
func loadQuery(db *db.DB, fixture string) error {
	qs1, err := query.GetFixture(fixture)
	if err != nil {
		return err
	}

	if err := query.AddTestSet(db, qs1); err != nil {
		return err
	}

	return nil
}
