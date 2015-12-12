// Package endpoint implements users tests for the API layer.
package tests

import (
	"net/http/httptest"
	"testing"

	"github.com/coralproject/shelf/app/xenia/routes"
	"github.com/coralproject/shelf/pkg/query"

	"github.com/ardanlabs/kit/db"
	"github.com/ardanlabs/kit/tests"
	"github.com/ardanlabs/kit/web/app"
)

var a *app.App

func init() {
	tests.Init("SHELF")
	tests.InitMongo()

	a = routes.API().(*app.App)
}

//==============================================================================

// TestQueryNames tests the retrieval of query names.
func TestQueryNames(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	loadQuery()
	defer removeQuery()

	r := tests.NewRequest("GET", "/1.0/query", nil)
	w := httptest.NewRecorder()

	a.ServeHTTP(w, r)

	t.Log("Given the need get a list of query names.")
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

	loadQuery()
	defer removeQuery()

	r := tests.NewRequest("GET", "/1.0/query/QTEST_basic", nil)
	w := httptest.NewRecorder()

	a.ServeHTTP(w, r)

	t.Log("Given the need get a specific query.")
	{
		t.Log("\tWhen we user version 1.0 of the query/basic endpoint.")
		if w.Code != 200 {
			t.Fatalf("\t%s\tShould be able to retrieve the query : %v", tests.Failed, w.Code)
		}
		t.Logf("\t%s\tShould be able to retrieve the query.", tests.Success)
	}
}

//==============================================================================

// loadQuery adds queries to run tests.
func loadQuery() error {
	const fixture = "basic.json"
	qs1, err := query.GetFixture(fixture)
	if err != nil {
		return err
	}

	db := db.NewMGO()
	defer db.CloseMGO()

	if err := query.AddTestSet(db, qs1); err != nil {
		return err
	}

	return nil
}

// removeQuery removes the queries for the tests.
func removeQuery() error {
	db := db.NewMGO()
	defer db.CloseMGO()

	return query.RemoveTestSets(db)
}
