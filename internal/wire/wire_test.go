package wire_test

import (
	"fmt"
	"os"
	"testing"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/ardanlabs/kit/cfg"
	"github.com/ardanlabs/kit/db"
	"github.com/ardanlabs/kit/db/mongo"
	"github.com/ardanlabs/kit/tests"
	"github.com/cayleygraph/cayley"
	"github.com/cayleygraph/cayley/graph"
	_ "github.com/cayleygraph/cayley/graph/mongo"
	"github.com/cayleygraph/cayley/quad"
	"github.com/coralproject/shelf/internal/sponge/item/itemfix"
	"github.com/coralproject/shelf/internal/wire"
	"github.com/coralproject/shelf/internal/wire/pattern/patternfix"
	"github.com/coralproject/shelf/internal/wire/relationship/relationshipfix"
	"github.com/coralproject/shelf/internal/wire/view/viewfix"
)

var mgoCfg mongo.Config

const (
	relPrefix     = "RTEST_"
	itemPrefix    = "ITEST_"
	viewPrefix    = "VTEST_"
	patternPrefix = "PTEST_"
)

func init() {
	// Initialize the configuration and logging systems. Plus anything
	// else the web app layer needs.
	tests.Init("XENIA")

	// Initialize MongoDB using the `tests.TestSession` as the name of the
	// master session.
	mgoCfg = mongo.Config{
		Host:     cfg.MustString("MONGO_HOST"),
		AuthDB:   cfg.MustString("MONGO_AUTHDB"),
		DB:       cfg.MustString("MONGO_DB"),
		User:     cfg.MustString("MONGO_USER"),
		Password: cfg.MustString("MONGO_PASS"),
	}
	tests.InitMongo(mgoCfg)
}

// TestMain helps to clean up the test data.
func TestMain(m *testing.M) {
	os.Exit(runTest(m))
}

// runTest initializes the environment for the tests and allows for
// the proper return code if the test fails or succeeds.
func runTest(m *testing.M) int {
	db, err := db.NewMGO(tests.Context, tests.TestSession)
	if err != nil {
		fmt.Println("MongoDB is not configured")
		return 1
	}
	defer db.CloseMGO(tests.Context)

	if err := loadTestData(tests.Context, db); err != nil {
		fmt.Println("test data is not loaded")
		return 1
	}
	defer unloadTestData(tests.Context, db)

	return m.Run()
}

// loadTestData adds all the test data into the database.
func loadTestData(context interface{}, db *db.DB) error {

	// -----------------------------------------------------------
	// Load example Items.

	items, err := itemfix.Get()
	if err != nil {
		return err
	}

	if err := itemfix.Add(context, db, items); err != nil {
		return err
	}

	// -----------------------------------------------------------
	// Load example patterns.

	patterns, _, err := patternfix.Get()
	if err != nil {
		return err
	}

	if err := patternfix.Add(context, db, patterns); err != nil {
		return err
	}

	// -----------------------------------------------------------
	// Load example Relationships.

	rels, err := relationshipfix.Get()
	if err != nil {
		return err
	}

	if err := relationshipfix.Add(context, db, rels); err != nil {
		return err
	}

	// -----------------------------------------------------------
	// Load example Views.

	views, err := viewfix.Get()
	if err != nil {
		return err
	}

	if err := viewfix.Add(context, db, views); err != nil {
		return err
	}

	// -----------------------------------------------------------
	// Build the example graph.

	opts := make(map[string]interface{})
	opts["database_name"] = cfg.MustString("MONGO_DB")
	opts["username"] = cfg.MustString("MONGO_USER")
	opts["password"] = cfg.MustString("MONGO_PASS")
	if err := graph.InitQuadStore("mongo", cfg.MustString("MONGO_HOST"), opts); err != nil {
		return err
	}

	var quads []quad.Quad
	quads = append(quads, quad.Make("ITEST_d1dfa366-d2f7-4a4a-a64f-af89d4c97d82", relPrefix+"on", "ITEST_c1b2bbfe-af9f-4903-8777-bd47c4d5b20a", ""))
	quads = append(quads, quad.Make("ITEST_6eaaa19f-da7a-4095-bbe3-cee7a7631dd4", relPrefix+"on", "ITEST_c1b2bbfe-af9f-4903-8777-bd47c4d5b20a", ""))
	quads = append(quads, quad.Make("ITEST_d16790f8-13e9-4cb4-b9ef-d82835589660", relPrefix+"on", "ITEST_c1b2bbfe-af9f-4903-8777-bd47c4d5b20a", ""))
	quads = append(quads, quad.Make("ITEST_80aa936a-f618-4234-a7be-df59a14cf8de", relPrefix+"authored", "ITEST_d1dfa366-d2f7-4a4a-a64f-af89d4c97d82", ""))
	quads = append(quads, quad.Make("ITEST_80aa936a-f618-4234-a7be-df59a14cf8de", relPrefix+"authored", "ITEST_6eaaa19f-da7a-4095-bbe3-cee7a7631dd4", ""))
	quads = append(quads, quad.Make("ITEST_a63af637-58af-472b-98c7-f5c00743bac6", relPrefix+"authored", "ITEST_d16790f8-13e9-4cb4-b9ef-d82835589660", ""))

	tx := cayley.NewTransaction()
	for _, quad := range quads {
		tx.AddQuad(quad)
	}

	store, err := cayley.NewGraph("mongo", cfg.MustString("MONGO_HOST"), opts)
	if err != nil {
		return err
	}

	if err := store.ApplyTransaction(tx); err != nil {
		return err
	}

	return nil
}

// unloadTestData removes all the test data from the database.
func unloadTestData(context interface{}, db *db.DB) error {

	// ------------------------------------------------------------
	// Clear items, relationships, and views.

	itemfix.Remove("context", db, itemPrefix)
	relationshipfix.Remove("context", db, relPrefix)
	viewfix.Remove("context", db, viewPrefix)
	patternfix.Remove("context", db, patternPrefix)

	// ------------------------------------------------------------
	// Clear cayley graph.

	opts := make(map[string]interface{})
	opts["database_name"] = cfg.MustString("MONGO_DB")
	opts["username"] = cfg.MustString("MONGO_USER")
	opts["password"] = cfg.MustString("MONGO_PASS")

	var quads []quad.Quad
	quads = append(quads, quad.Make("ITEST_d1dfa366-d2f7-4a4a-a64f-af89d4c97d82", relPrefix+"on", "ITEST_c1b2bbfe-af9f-4903-8777-bd47c4d5b20a", ""))
	quads = append(quads, quad.Make("ITEST_6eaaa19f-da7a-4095-bbe3-cee7a7631dd4", relPrefix+"on", "ITEST_c1b2bbfe-af9f-4903-8777-bd47c4d5b20a", ""))
	quads = append(quads, quad.Make("ITEST_d16790f8-13e9-4cb4-b9ef-d82835589660", relPrefix+"on", "ITEST_c1b2bbfe-af9f-4903-8777-bd47c4d5b20a", ""))
	quads = append(quads, quad.Make("ITEST_80aa936a-f618-4234-a7be-df59a14cf8de", relPrefix+"authored", "ITEST_d1dfa366-d2f7-4a4a-a64f-af89d4c97d82", ""))
	quads = append(quads, quad.Make("ITEST_80aa936a-f618-4234-a7be-df59a14cf8de", relPrefix+"authored", "ITEST_6eaaa19f-da7a-4095-bbe3-cee7a7631dd4", ""))
	quads = append(quads, quad.Make("ITEST_a63af637-58af-472b-98c7-f5c00743bac6", relPrefix+"authored", "ITEST_d16790f8-13e9-4cb4-b9ef-d82835589660", ""))

	tx := cayley.NewTransaction()
	for _, quad := range quads {
		tx.RemoveQuad(quad)
	}

	store, err := cayley.NewGraph("mongo", cfg.MustString("MONGO_HOST"), opts)
	if err != nil {
		return err
	}

	if err := store.ApplyTransaction(tx); err != nil {
		return err
	}

	return nil

}

// setup initializes for each indivdual test.
func setup(t *testing.T) (*db.DB, *cayley.Handle) {
	tests.ResetLog()

	db, err := db.NewMGO(tests.Context, tests.TestSession)
	if err != nil {
		t.Fatalf("%s\tShould be able to get a Mongo session : %v", tests.Failed, err)
	}

	opts := make(map[string]interface{})
	opts["database_name"] = cfg.MustString("MONGO_DB")
	opts["username"] = cfg.MustString("MONGO_USER")
	opts["password"] = cfg.MustString("MONGO_PASS")
	store, err := cayley.NewGraph("mongo", cfg.MustString("MONGO_HOST"), opts)
	if err != nil {
		t.Fatalf("\t%s\tShould be able to get a Cayley handle : %v", tests.Failed, err)
	}

	return db, store
}

// teardown deinitializes for each indivdual test.
func teardown(t *testing.T, db *db.DB) {
	db.CloseMGO(tests.Context)
	tests.DisplayLog()
}

//==============================================================================

// TestExecuteView tests the generation of a view, opting not to persist the view.
func TestExecuteView(t *testing.T) {
	db, store := setup(t)
	defer teardown(t, db)

	t.Log("Given the need to generate a view.")
	{
		t.Log("\tWhen using the view, relationship, and item fixtures.")
		{

			// Form the view parameters.
			viewParams := wire.ViewParams{
				ViewName: viewPrefix + "user comments",
				ItemKey:  "ITEST_80aa936a-f618-4234-a7be-df59a14cf8de",
			}

			// Generate the view.
			result, err := wire.Execute(tests.Context, db, store, &viewParams)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to generate the view : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to generate the view", tests.Success)

			// Check the resulting items.
			items, ok := result.Results.([]bson.M)
			if !ok || len(items) != 2 {
				t.Fatalf("\t%s\tShould be able to get 2 items in the view : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to get 2 items in the view.", tests.Success)
		}
	}
}

// TestExecuteBackwardsView tests the generation of a view with multiple
// out direction relationships, opting not to persist the view.
func TestExecuteBackwardsView(t *testing.T) {
	db, store := setup(t)
	defer teardown(t, db)

	t.Log("Given the need to generate a view with multiple backwards direction relationships.")
	{
		t.Log("\tWhen using the view, relationship, and item fixtures.")
		{

			// Form the view parameters.
			viewParams := wire.ViewParams{
				ViewName: viewPrefix + "thread_backwards",
				ItemKey:  "ITEST_80aa936a-f618-4234-a7be-df59a14cf8de",
			}

			// Generate the view.
			result, err := wire.Execute(tests.Context, db, store, &viewParams)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to generate the backwards view : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to generate the backwards view", tests.Success)

			// Check the resulting items.
			items, ok := result.Results.([]bson.M)
			if !ok || len(items) != 3 {
				t.Fatalf("\t%s\tShould be able to get e items in the view : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to get 3 items in the view.", tests.Success)
		}
	}
}

// TestPersistView tests the generation of a view, opting to persist the view.
func TestPersistView(t *testing.T) {
	db, store := setup(t)
	defer teardown(t, db)

	t.Log("Given the need to generate and persist a view.")
	{
		t.Log("\tWhen using the view, relationship, and item fixtures.")
		{

			// Form the view parameters.
			viewParams := wire.ViewParams{
				ViewName:          viewPrefix + "thread",
				ItemKey:           "ITEST_c1b2bbfe-af9f-4903-8777-bd47c4d5b20a",
				ResultsCollection: "testcollection",
			}

			// Generate the view.
			result, err := wire.Execute(tests.Context, db, store, &viewParams)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to generate the view : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to generate the view", tests.Success)

			// Check the result message.
			msg, ok := result.Results.(bson.M)
			if !ok || msg["number_of_results"] != 5 {
				t.Fatalf("\t%s\tShould be able to get 5 items in the view : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to get 5 items in the view.", tests.Success)

			// Verify that the output collection exists.
			var viewItems []bson.M
			f := func(c *mgo.Collection) error {
				return c.Find(nil).All(&viewItems)
			}

			if err := db.ExecuteMGO(tests.Context, "testcollection", f); err != nil {
				t.Fatalf("\t%s\tShould be able to query the output collection : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to query the output collection.", tests.Success)

			if len(viewItems) != 5 {
				t.Fatalf("\t%s\tShould be able to get 5 items from the output collection : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to get 5 items from the output collection.", tests.Success)

			// Delete the persisted collection to clean up.
			f = func(c *mgo.Collection) error {
				return c.DropCollection()
			}

			if err := db.ExecuteMGO(tests.Context, "testcollection", f); err != nil {
				t.Fatalf("\t%s\tShould be able to drop the output collection : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to drop the output collection.", tests.Success)
		}
	}
}

// TestPersistViewWithBuffer tests the buffered saving of a view.
func TestPersistViewWithBuffer(t *testing.T) {
	db, store := setup(t)
	defer teardown(t, db)

	t.Log("Given the need to perform a buffered save of a view.")
	{
		t.Log("\tWhen using the view, relationship, and item fixtures.")
		{

			// Form the view parameters.
			viewParams := wire.ViewParams{
				ViewName:          viewPrefix + "thread",
				ItemKey:           "ITEST_c1b2bbfe-af9f-4903-8777-bd47c4d5b20a",
				ResultsCollection: "testcollection",
				BufferLimit:       2,
			}

			// Generate the view.
			result, err := wire.Execute(tests.Context, db, store, &viewParams)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to generate the view : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to generate the view", tests.Success)

			// Check the result message.
			msg, ok := result.Results.(bson.M)
			if !ok || msg["number_of_results"] != 5 {
				t.Fatalf("\t%s\tShould be able to get 5 items in the view : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to get 5 items in the view.", tests.Success)

			// Verify that the output collection exists.
			var viewItems []bson.M
			f := func(c *mgo.Collection) error {
				return c.Find(nil).All(&viewItems)
			}

			if err := db.ExecuteMGO(tests.Context, "testcollection", f); err != nil {
				t.Fatalf("\t%s\tShould be able to query the output collection : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to query the output collection.", tests.Success)

			if len(viewItems) != 5 {
				t.Fatalf("\t%s\tShould be able to get 5 items from the output collection : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to get 5 items from the output collection.", tests.Success)

			// Delete the persisted collection to clean up.
			f = func(c *mgo.Collection) error {
				return c.DropCollection()
			}

			if err := db.ExecuteMGO(tests.Context, "testcollection", f); err != nil {
				t.Fatalf("\t%s\tShould be able to drop the output collection : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to drop the output collection.", tests.Success)
		}
	}
}

// TestExecuteNameFail tests that the correct result is returned when
// an invalid view name is provided.
func TestExecuteNameFail(t *testing.T) {
	db, store := setup(t)
	defer teardown(t, db)

	t.Log("Given the need to catch an invalid view name.")
	{
		t.Log("\tWhen using the view, relationship, and item fixtures.")
		{

			// Form the view parameters.
			viewParams := wire.ViewParams{
				ViewName: viewPrefix + "this view name does not exist",
				ItemKey:  "ITEST_80aa936a-f618-4234-a7be-df59a14cf8de",
			}

			// Generate the view.
			result, err := wire.Execute(tests.Context, db, store, &viewParams)
			if err == nil {
				t.Fatalf("\t%s\tShould return an error for an invalid view name: %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould return an error for an invalid view name", tests.Success)

			// Check the resulting items.
			errDoc, ok := result.Results.(bson.M)
			if !ok || len(errDoc["error"].(string)) == 0 {
				t.Fatalf("\t%s\tShould return a single error document : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould return a single error document.", tests.Success)
		}
	}
}

// TestExecuteTypeFail tests that the correct result is returned when
// an invalid start type is defined in view metadata.
func TestExecuteTypeFail(t *testing.T) {
	db, store := setup(t)
	defer teardown(t, db)

	t.Log("Given the need to catch an invalid start type.")
	{
		t.Log("\tWhen using the view, relationship, and item fixtures.")
		{

			// Form the view parameters.
			viewParams := wire.ViewParams{
				ViewName: viewPrefix + "comments from authors flagged by a user",
				ItemKey:  "ITEST_80aa936a-f618-4234-a7be-df59a14cf8de",
			}

			// Generate the view.
			result, err := wire.Execute(tests.Context, db, store, &viewParams)
			if err == nil {
				t.Fatalf("\t%s\tShould return an error for an invalid start type: %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould return an error for an invalid start type", tests.Success)

			// Check the resulting items.
			errDoc, ok := result.Results.(bson.M)
			if !ok || len(errDoc["error"].(string)) == 0 {
				t.Fatalf("\t%s\tShould return a single error document : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould return a single error document.", tests.Success)
		}
	}
}

// TestExecuteRelationshipFail tests that the correct result is returned when
// an invalid relationship is defined in view metadata.
func TestExecuteRelationshipFail(t *testing.T) {
	db, store := setup(t)
	defer teardown(t, db)

	t.Log("Given the need to catch an invalid relationship.")
	{
		t.Log("\tWhen using the view, relationship, and item fixtures.")
		{

			// Form the view parameters.
			viewParams := wire.ViewParams{
				ViewName: viewPrefix + "has invalid starting relationship",
				ItemKey:  "ITEST_80aa936a-f618-4234-a7be-df59a14cf8de",
			}

			// Generate the view.
			result, err := wire.Execute(tests.Context, db, store, &viewParams)
			if err == nil {
				t.Fatalf("\t%s\tShould return an error for an invalid relationship: %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould return an error for an invalid relationship", tests.Success)

			// Check the resulting items.
			errDoc, ok := result.Results.(bson.M)
			if !ok || len(errDoc["error"].(string)) == 0 {
				t.Fatalf("\t%s\tShould return a single error document : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould return a single error document.", tests.Success)
		}
	}
}
