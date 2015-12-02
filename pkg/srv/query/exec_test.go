package query_test

import (
	"testing"

	"github.com/coralproject/shelf/pkg/db"
	"github.com/coralproject/shelf/pkg/srv/query"
	"github.com/coralproject/shelf/pkg/tests"
)

// TestQueryExecution validates the process of executing a giving query record
// in the db.
func TestQueryExecution(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	const name = "QTEST_spending_advice"
	const fixture = "./fixtures/spending_advice.json"
	qs1, err := getFixture(fixture)
	if err != nil {
		t.Fatalf("\t%s\tShould load query record from file : %v", tests.Failed, err)
	}
	t.Logf("\t%s\tShould load query record from file.", tests.Success)

	db := db.NewMGO()
	defer db.CloseMGO()

	defer func() {
		if err := removeSets(db); err != nil {
			t.Fatalf("\t%s\tShould be able to remove the query set : %v", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to remove the query set.", tests.Success)
	}()

	t.Logf("Given the need to execute a query.")
	{
		t.Logf("\tWhen giving a fixture")
		{
			if err := query.UpsertSet(context, db, qs1); err != nil {
				t.Fatalf("\t%s\tShould be able to create a query set : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to create a query set.", tests.Success)

			_, err := query.ExecuteQuerySet(context, db, name)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to execute query: %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to execute query", tests.Success)
		}
	}
}
