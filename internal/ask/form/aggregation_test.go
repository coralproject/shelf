package form_test

import (
	"testing"

	"github.com/ardanlabs/kit/tests"
	"github.com/coralproject/shelf/internal/ask/form"
	"github.com/coralproject/shelf/internal/ask/form/aggfix"
	"github.com/coralproject/shelf/internal/ask/form/submission"
	"github.com/coralproject/shelf/internal/platform/db"
)

// aggPrefix is what we are looking to delete after the test.
const aggPrefix = "ATEST_"

func setupAgg(t *testing.T) (*form.Form, []submission.Submission, *db.DB) {
	tests.ResetLog()

	fm, subs, err := aggfix.Get()
	if err != nil {
		t.Fatalf("%s\tShould be able retrieve form and submission fixture : %s", tests.Failed, err)
	}

	db, err := db.NewMGO(tests.Context, tests.TestSession)
	if err != nil {
		t.Fatalf("Should be able to get a Mongo session : %v", err)
	}

	if err := aggfix.Add(tests.Context, db, fm, subs); err != nil {
		t.Fatalf("Should be able to add forms and submissions to the database : %v", err)
	}

	return fm, subs, db
}

func teardownAgg(t *testing.T, db *db.DB) {
	if err := aggfix.Remove(tests.Context, db, aggPrefix); err != nil {
		t.Fatalf("%s\tShould be able to remove the forms and submissions : %v", tests.Failed, err)
	}
	t.Logf("%s\tShould be able to remove the forms and submissions.", tests.Success)

	db.CloseMGO(tests.Context)
	tests.DisplayLog()
}

func TestAggregation(t *testing.T) {
	fm, _, db := setupAgg(t)
	defer teardownAgg(t, db)

	t.Log("Given the need to aggregate submissions.")
	{
		t.Log("\tWhen starting from a form and submission fixtures")
		{
			//----------------------------------------------------------------------
			// Aggregate the submissions.

			aggs, err := form.AggregateFormSubmissions(tests.Context, db, fm.ID.Hex())
			if err != nil {
				t.Fatalf("\t%s\tShould be able to aggregate submissions : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to aggregate submissions.", tests.Success)

			//----------------------------------------------------------------------
			// Check the aggregations.

			if len(aggs) != 11 {
				t.Fatalf("\t%s\tShould be able to get 11 aggregations : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to get 11 aggregations.", tests.Success)
		}
	}
}

func TestTextAggregate(t *testing.T) {
	fm, subs, db := setupAgg(t)
	defer teardownAgg(t, db)

	t.Log("Given the need to aggregate text.")
	{
		t.Log("\tWhen starting from a form and submission fixtures")
		{
			//----------------------------------------------------------------------
			// Aggregate the submissions.

			aggs, err := form.TextAggregate(tests.Context, db, fm.ID.Hex(), subs)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to aggregate submissions : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to aggregate submissions.", tests.Success)

			//----------------------------------------------------------------------
			// Check the aggregations.

			if len(aggs) != 10 {
				t.Fatalf("\t%s\tShould be able to get 11 aggregations : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to get 11 aggregations.", tests.Success)
		}
	}
}
