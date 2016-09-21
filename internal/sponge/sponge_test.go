package sponge_test

import (
	"testing"

	"github.com/coralproject/shelf/internal/sponge"
	"github.com/coralproject/shelf/internal/sponge/item"
	"github.com/coralproject/shelf/internal/sponge/sfix"

	"github.com/ardanlabs/kit/cfg"
	"github.com/ardanlabs/kit/db"
	"github.com/ardanlabs/kit/db/mongo"
	"github.com/ardanlabs/kit/tests"
)

// prefix is what we are looking to delete after the test.
const prefix = "SPONGE_TEST_O"

func init() {
	// Initialize the configuration and logging systems. Plus anything
	// else the web app layer needs.
	tests.Init("CORAL")

	// Initialize MongoDB using the `tests.TestSession` as the name of the
	// master session.
	cfg := mongo.Config{
		Host:     cfg.MustString("MONGO_HOST"),
		AuthDB:   cfg.MustString("MONGO_AUTHDB"),
		DB:       cfg.MustString("MONGO_DB") + "_test",
		User:     cfg.MustString("MONGO_USER"),
		Password: cfg.MustString("MONGO_PASS"),
	}
	tests.InitMongo(cfg)
}

//==============================================================================

func TestItemizeData(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	db, err := db.NewMGO(tests.Context, tests.TestSession)
	if err != nil {
		t.Fatalf("\t%s\tShould be able to get a Mongo session : %v", tests.Failed, err)
	}
	defer db.CloseMGO(tests.Context)

	d, err := sfix.Get("data.json")
	if err != nil {
		t.Fatalf("\t%s\tShould be able to load item data.json fixture: %v", tests.Failed, err)
	}

	i, err := sponge.Itemize(tests.Context, db, "test_comment", 1, d)
	if err != nil {
		t.Fatalf("\t%s\tCould not create item from data: %v", tests.Failed, err)
	}

	t.Logf("\t%s\tShould be able to insert the new item.", tests.Success)
	err = item.Upsert(tests.Context, db, &i)
	if err != nil {
		t.Fatalf("\t%s\tCould not upsert (insert) item: %v", tests.Failed, err)
	}

	t.Logf("\t%s\tShould get the same id item when itemizing again", tests.Success)
	i, err = sponge.Itemize(tests.Context, db, "test_comment", 1, d)
	if err != nil {
		t.Fatalf("\t%s\tCould not create item from data when item is present in store: %v", tests.Failed, err)
	}
	if i.ID == "" {
		t.Fatalf("\t%s\tItemize should find existing item_id when re-itemizing same data packet: %v", tests.Failed, err)
	}

}
