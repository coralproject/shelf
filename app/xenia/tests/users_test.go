// Package endpoint implements users tests for the API layer.
// USE THIS AS A MODEL FOR NOW.
package endpoint

import (
	"net/http/httptest"
	"testing"

	"github.com/coralproject/shelf/app/xenia/app"
	"github.com/coralproject/shelf/app/xenia/routes"

	"github.com/coralproject/shelf/pkg/db"
	"github.com/coralproject/shelf/pkg/tests"
)

func init() {
	tests.Init()
}

// TestQueryNames tests the retrieval of query names.
func TestQueryNames(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	c := &app.Context{
		DB:        db.NewMGO(),
		SessionID: "TESTING",
	}
	defer c.DB.CloseMGO()

	a := routes.API().(*app.App)

	r := tests.NewRequest("GET", "/1.0/query/names", nil)
	w := httptest.NewRecorder()

	a.ServeHTTP(w, r)

	t.Log("Given the need get a list of query names.")
	{
		t.Log("\tWhen we user version 1.0 of the query/names endpoint.")
		if w.Code == 404 {
			t.Fatalf("\t%s\tShould be able to retrieve the query list : %v", tests.Failed, w.Code)
		}
		t.Logf("\t%s\tShould be able to retrieve the query list.", tests.Success)
	}
}
