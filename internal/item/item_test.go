package item_test

import (
	"testing"

	"github.com/coralproject/xenia/internal/item"
	"github.com/coralproject/xenia/internal/item/ifix"

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

	// get db connection
	//  should this be moved to a shared testing package?
	db, err := db.NewMGO(tests.Context, tests.TestSession)
	if err != nil {
		t.Fatalf("\t%s\tShould be able to get a Mongo session : %v", tests.Failed, err)
	}
	defer db.CloseMGO(tests.Context)

	// register the item types
	err = ifix.RegisterTypes("types.json")
	if err != nil {
		t.Fatalf("\t%s\tCould not register the types from types.json. %s", err)
	}

	// create the indexes
	err = item.EnsureTypeIndexes(tests.Context, db, item.Types)
	if err != nil {
		t.Fatalf("\t%s\tUnable to create indices", tests.Failed, err)

	}

	//Todo, how can we test to ensure inexes are created?

}

//==============================================================================

func TestRels(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	// get db connection
	//  should this be moved to a shared testing package?
	db, err := db.NewMGO(tests.Context, tests.TestSession)
	if err != nil {
		t.Fatalf("\t%s\tShould be able to get a Mongo session : %v", tests.Failed, err)
	}
	defer db.CloseMGO(tests.Context)

	// register the item types
	err = ifix.RegisterTypes("types.json")
	if err != nil {
		t.Fatalf("\t%s\tCould not register the types from types.json. %s", err)
	}

	err = ifix.InsertItemsFromDataFile(tests.Context, db, "coral_asset_data.json", "coral_asset")
	if err != nil {
		t.Fatalf("\t%s\tCould not load asset items from coral_asset_data.json : %v", tests.Failed, err)
	}

	err = ifix.InsertItemsFromDataFile(tests.Context, db, "coral_user_data.json", "coral_user")
	if err != nil {
		t.Fatalf("\t%s\tCould not load asset items from coral_asset_data.json : %v", tests.Failed, err)
	}

	dataSets, err := ifix.Get("coral_comment_data.json")
	if err != nil {
		t.Fatalf("\t%s\tShould be able to load item data.json fixture", tests.Failed, err)
	}

	t.Logf("\t%s\tCreate, items, get relationships and save.", tests.Success)
	for _, d := range *dataSets {

		_, err := item.Create(tests.Context, db, "coral_comment", 1, d)
		if err != nil {
			t.Fatalf("\t%s\tCould not create item from data: %v", tests.Failed, err)
		}

		// relationships are automatically determened as part of item Create()

	}
}

func TestCreateAndUpsertItem(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	// get db connection
	//  should this be moved to a shared testing package?
	db, err := db.NewMGO(tests.Context, tests.TestSession)
	if err != nil {
		t.Fatalf("\t%s\tShould be able to get a Mongo session : %v", tests.Failed, err)
	}
	defer db.CloseMGO(tests.Context)

	// register the item types
	err = ifix.RegisterTypes("types.json")
	if err != nil {
		t.Fatalf("\t%s\tShould be able to register the types from types.json. %s", err)
	}

	dataSets, err := ifix.Get("coral_comment_data.json")
	if err != nil {
		t.Fatalf("\t%s\tShould be able to load item data.json fixture", tests.Failed, err)
	}

	t.Logf("\t%s\tCreate, save and update items.", tests.Success)
	for _, d := range *dataSets {

		i, err := item.Create(tests.Context, db, "coral_comment", 1, d)
		if err != nil {
			t.Fatalf("\t%s\tCould not create item from data: %v", tests.Failed, err)
		}

		rels, err := item.GetRels(tests.Context, db, &i)
		if err != nil {
			t.Fatalf("\t%s\tFailed to get an item's relationships: %v", tests.Failed, err)
		}
		i.Rels = *rels

		_, err = item.Create(tests.Context, db, "an_unregistered_type", 1, d)
		if err == nil {
			t.Fatalf("\t%s\tShould not be able to create with unregistered type: %v", tests.Failed, err)
		}

		t.Logf("\t%s\tShould be able to insert items.", tests.Success)
		err = item.Upsert(tests.Context, db, &i)
		if err != nil {
			t.Fatalf("\t%s\tCould not upsert (insert) item: %v", tests.Failed, err)
		}

		t.Logf("\t%s\tShould be able to update items and get updated items from the db.", tests.Success)
		// bump the version
		i.Version = 2
		// upsert (update) the item
		err = item.Upsert(tests.Context, db, &i)
		if err != nil {
			t.Fatalf("\t%s\tCould not upsert (update) item: %v", tests.Failed, err)
		}

		// get the item by id
		i2, err := item.GetById(tests.Context, db, i.Id)
		if err != nil {
			t.Fatalf("\t%s\tCould not GetById item: %v", tests.Failed, err)
		}

		// verify that the version change has been saved
		if i.Version != i2.Version {
			t.Fatalf("\t%s\tDid not see update of version", tests.Failed)

		}

	}

}
