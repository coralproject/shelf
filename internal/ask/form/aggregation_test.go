package form_test

import (
	"testing"

	"github.com/ardanlabs/kit/tests"
	"github.com/coralproject/shelf/internal/ask/form/aggfix"
	"github.com/coralproject/shelf/internal/ask/form/formfix"
	"github.com/coralproject/shelf/internal/platform/db"
)

// aggPrefix is what we are looking to delete after the test.
const aggPrefix = "ATEST_"

func setupAgg(t *testing.T, fixture string) *db.DB {
	tests.ResetLog()

	fms, subs, err := aggfix.Get()
	if err != nil {
		t.Fatalf("%s\tShould be able retrieve form and submission fixture : %s", tests.Failed, err)
	}

	db, err := db.NewMGO(tests.Context, tests.TestSession)
	if err != nil {
		t.Fatalf("Should be able to get a Mongo session : %v", err)
	}

	if err := aggfix.Add(tests.Context, db, fms, subs); err != nil {
		t.Fatalf("Should be able to add forms and submissions to the database : %v", err)
	}

	return db
}

func teardownAgg(t *testing.T, db *db.DB) {
	if err := formfix.Remove(tests.Context, db, aggPrefix); err != nil {
		t.Fatalf("%s\tShould be able to remove the forms and submissions : %v", tests.Failed, err)
	}
	t.Logf("%s\tShould be able to remove the forms and submissions.", tests.Success)

	db.CloseMGO(tests.Context)
	tests.DisplayLog()
}

func TestAggregation(t *testing.T) {
	setupAgg(t, "form")
}
