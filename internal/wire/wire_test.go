package wire_test

import (
	"fmt"
	"os"
	"testing"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/ardanlabs/kit/cfg"
	"github.com/ardanlabs/kit/tests"
	"github.com/cayleygraph/cayley"
	"github.com/cayleygraph/cayley/graph"
	"github.com/cayleygraph/cayley/quad"
	"github.com/coralproject/shelf/internal/platform/db"
	cayleyshelf "github.com/coralproject/shelf/internal/platform/db/cayley"
	"github.com/coralproject/shelf/internal/wire"
	"github.com/coralproject/shelf/internal/wire/wirefix"
)

const wirePrefix = "WTEST_"

func TestMain(m *testing.M) {
	os.Exit(runTest(m))
}

// runTest initializes the environment for the tests and allows for
// the proper return code if the test fails or succeeds.
func runTest(m *testing.M) int {

	// Initialize the configuration and logging systems. Plus anything
	// else the web app layer needs.
	tests.Init("XENIA")

	// Initialize MongoDB using the `tests.TestSession` as the name of the
	// master session.
	if err := db.RegMasterSession(tests.Context, tests.TestSession, cfg.MustURL("MONGO_URI").String(), 0); err != nil {
		fmt.Println("Can't register master session: " + err.Error())
		return 1
	}

	db, err := db.NewMGO(tests.Context, tests.TestSession)
	if err != nil {
		fmt.Println("MongoDB is not configured")
		return 1
	}
	defer db.CloseMGO(tests.Context)

	if err := loadTestData(tests.Context, db); err != nil {
		fmt.Println("test data is not loaded: " + err.Error())
		return 1
	}
	defer unloadTestData(tests.Context, db)

	return m.Run()
}

// loadTestData adds all the test data into the database.
func loadTestData(context interface{}, db *db.DB) error {

	// Make sure the old data is clear.
	if err := unloadTestData(context, db); err != nil {
		if !graph.IsQuadNotExist(err) {
			return err
		}
	}

	// -----------------------------------------------------------
	// Load example items, relationships, views, and patterns.

	items, rels, vs, pats, err := wirefix.Get()
	if err != nil {
		return err
	}

	if err := wirefix.Add(context, db, items, rels, vs, pats); err != nil {
		return err
	}

	// -----------------------------------------------------------
	// Build the example graph.

	mongoURI := cfg.MustURL("MONGO_URI")

	if err := cayleyshelf.InitQuadStore(mongoURI.String()); err != nil {
		return err
	}

	var quads []quad.Quad
	quads = append(quads, quad.Make(wirePrefix+"d1dfa366-d2f7-4a4a-a64f-af89d4c97d82", wirePrefix+"on", wirePrefix+"c1b2bbfe-af9f-4903-8777-bd47c4d5b20a", ""))
	quads = append(quads, quad.Make(wirePrefix+"6eaaa19f-da7a-4095-bbe3-cee7a7631dd4", wirePrefix+"on", wirePrefix+"c1b2bbfe-af9f-4903-8777-bd47c4d5b20a", ""))
	quads = append(quads, quad.Make(wirePrefix+"d16790f8-13e9-4cb4-b9ef-d82835589660", wirePrefix+"on", wirePrefix+"c1b2bbfe-af9f-4903-8777-bd47c4d5b20a", ""))
	quads = append(quads, quad.Make(wirePrefix+"80aa936a-f618-4234-a7be-df59a14cf8de", wirePrefix+"authored", wirePrefix+"d1dfa366-d2f7-4a4a-a64f-af89d4c97d82", ""))
	quads = append(quads, quad.Make(wirePrefix+"80aa936a-f618-4234-a7be-df59a14cf8de", wirePrefix+"authored", wirePrefix+"6eaaa19f-da7a-4095-bbe3-cee7a7631dd4", ""))
	quads = append(quads, quad.Make(wirePrefix+"a63af637-58af-472b-98c7-f5c00743bac6", wirePrefix+"authored", wirePrefix+"d16790f8-13e9-4cb4-b9ef-d82835589660", ""))
	quads = append(quads, quad.Make(wirePrefix+"a63af637-58af-472b-98c7-f5c00743bac6", wirePrefix+"flagged", wirePrefix+"80aa936a-f618-4234-a7be-df59a14cf8de", ""))

	tx := cayley.NewTransaction()
	for _, quad := range quads {
		tx.AddQuad(quad)
	}

	store, err := cayleyshelf.New(mongoURI.String())
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

	wirefix.Remove("context", db, wirePrefix)

	// ------------------------------------------------------------
	// Clear cayley graph.

	var quads []quad.Quad
	quads = append(quads, quad.Make(wirePrefix+"d1dfa366-d2f7-4a4a-a64f-af89d4c97d82", wirePrefix+"on", wirePrefix+"c1b2bbfe-af9f-4903-8777-bd47c4d5b20a", ""))
	quads = append(quads, quad.Make(wirePrefix+"6eaaa19f-da7a-4095-bbe3-cee7a7631dd4", wirePrefix+"on", wirePrefix+"c1b2bbfe-af9f-4903-8777-bd47c4d5b20a", ""))
	quads = append(quads, quad.Make(wirePrefix+"d16790f8-13e9-4cb4-b9ef-d82835589660", wirePrefix+"on", wirePrefix+"c1b2bbfe-af9f-4903-8777-bd47c4d5b20a", ""))
	quads = append(quads, quad.Make(wirePrefix+"80aa936a-f618-4234-a7be-df59a14cf8de", wirePrefix+"authored", wirePrefix+"d1dfa366-d2f7-4a4a-a64f-af89d4c97d82", ""))
	quads = append(quads, quad.Make(wirePrefix+"80aa936a-f618-4234-a7be-df59a14cf8de", wirePrefix+"authored", wirePrefix+"6eaaa19f-da7a-4095-bbe3-cee7a7631dd4", ""))
	quads = append(quads, quad.Make(wirePrefix+"a63af637-58af-472b-98c7-f5c00743bac6", wirePrefix+"authored", wirePrefix+"d16790f8-13e9-4cb4-b9ef-d82835589660", ""))
	quads = append(quads, quad.Make(wirePrefix+"a63af637-58af-472b-98c7-f5c00743bac6", wirePrefix+"flagged", wirePrefix+"80aa936a-f618-4234-a7be-df59a14cf8de", ""))

	tx := cayley.NewTransaction()
	for _, quad := range quads {
		tx.RemoveQuad(quad)
	}

	store, err := cayleyshelf.New(cfg.MustURL("MONGO_URI").String())
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

	store, err := cayleyshelf.New(cfg.MustURL("MONGO_URI").String())
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
				ViewName: wirePrefix + "user comments",
				ItemKey:  wirePrefix + "80aa936a-f618-4234-a7be-df59a14cf8de",
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
				t.Fatalf("\t%s\tShould be able to get 2 items in the view.", tests.Failed)
			}
			t.Logf("\t%s\tShould be able to get 2 items in the view.", tests.Success)
		}
	}
}

// TestExecuteReturnRoot tests the generation of a view, opting not to persist the view
// but returning a root item.
func TestExecuteReturnRoot(t *testing.T) {
	db, store := setup(t)
	defer teardown(t, db)

	t.Log("Given the need to generate a view and return a root item.")
	{
		t.Log("\tWhen using the view, relationship, and item fixtures.")
		{

			// Form the view parameters.
			viewParams := wire.ViewParams{
				ViewName: wirePrefix + "user comments return root",
				ItemKey:  wirePrefix + "80aa936a-f618-4234-a7be-df59a14cf8de",
			}

			// Generate the view.
			result, err := wire.Execute(tests.Context, db, store, &viewParams)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to generate the view : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to generate the view", tests.Success)

			// Check the resulting items.
			items, ok := result.Results.([]bson.M)
			if !ok || len(items) != 3 {
				t.Fatalf("\t%s\tShould be able to get 3 items in the view.", tests.Failed)
			}
			t.Logf("\t%s\tShould be able to get 3 items in the view.", tests.Success)
		}
	}
}

// TestExecuteSplitPathEmbeds tests the generation of a view from a split path, opting
// not to persist the view.
func TestExecuteSplitPathEmbeds(t *testing.T) {
	db, store := setup(t)
	defer teardown(t, db)

	t.Log("Given the need to generate a view from a split path.")
	{
		t.Log("\tWhen using the view, relationship, and item fixtures.")
		{

			// Form the view parameters.
			viewParams := wire.ViewParams{
				ViewName: wirePrefix + "split_path",
				ItemKey:  wirePrefix + "a63af637-58af-472b-98c7-f5c00743bac6",
			}

			// Generate the view.
			result, err := wire.Execute(tests.Context, db, store, &viewParams)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to generate the view : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to generate the view", tests.Success)

			// Check the resulting items.
			items, ok := result.Results.([]bson.M)
			if !ok || len(items) != 3 {
				t.Fatalf("\t%s\tShould be able to get 2 items in the view.", tests.Failed)
			}
			t.Logf("\t%s\tShould be able to get 2 items in the view.", tests.Success)

			for _, itm := range items {
				itemField, ok := itm["item_id"]
				if !ok {
					continue
				}
				itemID, ok := itemField.(string)
				if !ok {
					continue
				}
				if itemID == wirePrefix+"a63af637-58af-472b-98c7-f5c00743bac6" {
					if _, ok := itm["related"]; !ok {
						t.Fatalf("\t%s\tShould be able to get related items.", tests.Failed)
					}
				}
			}
			t.Logf("\t%s\tShould be able to get related items.", tests.Success)
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
				ViewName: wirePrefix + "thread_backwards",
				ItemKey:  wirePrefix + "80aa936a-f618-4234-a7be-df59a14cf8de",
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
				t.Fatalf("\t%s\tShould be able to get e items in the view", tests.Failed)
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
				ViewName:          wirePrefix + "thread",
				ItemKey:           wirePrefix + "c1b2bbfe-af9f-4903-8777-bd47c4d5b20a",
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
				t.Fatalf("\t%s\tShould be able to get 5 items in the view", tests.Failed)
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
				t.Fatalf("\t%s\tShould be able to get 5 items from the output collection", tests.Failed)
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
				ViewName:          wirePrefix + "thread",
				ItemKey:           wirePrefix + "c1b2bbfe-af9f-4903-8777-bd47c4d5b20a",
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
				t.Fatalf("\t%s\tShould be able to get 5 items in the view.", tests.Failed)
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
				t.Fatalf("\t%s\tShould be able to get 5 items from the output collection.", tests.Failed)
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
				ViewName: wirePrefix + "this view name does not exist",
				ItemKey:  wirePrefix + "80aa936a-f618-4234-a7be-df59a14cf8de",
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
				ViewName: wirePrefix + "comments from authors flagged by a user",
				ItemKey:  wirePrefix + "80aa936a-f618-4234-a7be-df59a14cf8de",
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
				ViewName: wirePrefix + "has invalid starting relationship",
				ItemKey:  wirePrefix + "80aa936a-f618-4234-a7be-df59a14cf8de",
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
