package data_test

import (
	"testing"

	"github.com/coralproject/shelf/internal/sponge/data"
	"github.com/coralproject/shelf/internal/sponge/data/dfix"
	"github.com/coralproject/shelf/internal/sponge/item"

	"github.com/ardanlabs/kit/cfg"
	"github.com/ardanlabs/kit/db"
	"github.com/ardanlabs/kit/db/mongo"
	"github.com/ardanlabs/kit/tests"
)

// prefix is what we are looking to delete after the test.
const prefix = "ITEM_TEST_O"

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

func TestEnsureTypeIndexes(t *testing.T) {

	db, err := db.NewMGO(tests.Context, tests.TestSession)
	if err != nil {
		t.Fatalf("\t%s\tShould be able to get a Mongo session: %v", tests.Failed, err)
	}
	defer db.CloseMGO(tests.Context)

	err = dfix.RegisterTypes("types.json")
	if err != nil {
		t.Fatalf("\t%s\tCould not register the types from types.json: %v", tests.Failed, err)
	}

	err = data.EnsureTypeIndexes(tests.Context, db, data.Types)
	if err != nil {
		t.Fatalf("\t%s\tUnable to create indexes: %v", tests.Failed, err)

	}

}

func TestItemizeData(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	db, err := db.NewMGO(tests.Context, tests.TestSession)
	if err != nil {
		t.Fatalf("\t%s\tShould be able to get a Mongo session : %v", tests.Failed, err)
	}
	defer db.CloseMGO(tests.Context)

	err = dfix.RegisterTypes("types.json")
	if err != nil {
		t.Fatalf("\t%s\tShould be able to register the types from types.json. %s", tests.Failed, err)
	}

	d, err := dfix.Get("data.json")
	if err != nil {
		t.Fatalf("\t%s\tShould be able to load item data.json fixture: %v", tests.Failed, err)
	}

	i, err := data.Itemize(tests.Context, db, "test_comment", 1, d)
	if err != nil {
		t.Fatalf("\t%s\tCould not create item from data: %v", tests.Failed, err)
	}

	t.Logf("\t%s\tShould be able to insert the new item.", tests.Success)
	err = item.Upsert(tests.Context, db, &i)
	if err != nil {
		t.Fatalf("\t%s\tCould not upsert (insert) item: %v", tests.Failed, err)
	}

	t.Logf("\t%s\tShould get the id from inserted item when itemizing again", tests.Success)
	i, err = data.Itemize(tests.Context, db, "test_comment", 1, d)
	if err != nil {
		t.Fatalf("\t%s\tCould not create item from data when item is present in store: %v", tests.Failed, err)
	}
	if i.ID == "" {
		t.Fatalf("\t%s\tItemize should find existing item_id when re-itemizing same data packet: %v", tests.Failed, err)
	}

}
