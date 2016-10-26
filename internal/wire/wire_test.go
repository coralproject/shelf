package wire_test

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
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

	if err := db.NewCayley(tests.Context, tests.TestSession); err != nil {
		return err
	}

	store, err := db.GraphHandle(tests.Context)
	if err != nil {
		return err
	}
	defer store.Close()

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

	if err := db.NewCayley(tests.Context, tests.TestSession); err != nil {
		return err
	}

	store, err := db.GraphHandle(tests.Context)
	if err != nil {
		return err
	}
	defer store.Close()

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

	if err := db.NewCayley(tests.Context, tests.TestSession); err != nil {
		t.Fatalf("%s\tShould be able to get Cayley support : %v", tests.Failed, err)
	}

	store, err := db.GraphHandle(tests.Context)
	if err != nil {
		t.Fatalf("\t%s\tShould be able to get a Cayley handle : %v", tests.Failed, err)
	}

	return db, store
}

// teardown deinitializes for each indivdual test.
func teardown(t *testing.T, db *db.DB, store *cayley.Handle) {
	db.CloseMGO(tests.Context)
	store.Close()
	tests.DisplayLog()
}

//==============================================================================

// TestExecuteViews tests the execution of different views.
func TestExecuteViews(t *testing.T) {
	db, store := setup(t)
	defer teardown(t, db, store)

	// Build our table of the different test sets.
	execViews := []struct {
		typ   string
		views []execView
	}{
		{typ: "Positive", views: getPosViews()},
		{typ: "Negative", views: getNegViews()},
	}

	// Iterate over all the different test view.
	for _, ev := range execViews {

		t.Logf("Given the need to execute %s view tests.", ev.typ)
		{
			for _, vw := range ev.views {

				// Setup a sub-test for each test view.
				tf := func(t *testing.T) {
					t.Logf("\tWhen using the view named %s", vw.viewName)
					{
						// Form the view parameters.
						viewParams := wire.ViewParams{
							ViewName:          wirePrefix + vw.viewName,
							ItemKey:           wirePrefix + vw.itemKey,
							ResultsCollection: vw.collection,
							BufferLimit:       vw.bufferLimit,
						}

						// Generate the view.
						result, err := wire.Execute(tests.Context, db, store, &viewParams)
						if err != nil && !vw.fail {
							t.Fatalf("\t%s\tShould be able to execute the view", tests.Failed)
						}
						t.Logf("\t%s\tShould be able to execute the view.", tests.Success)
						if err != nil && vw.fail {
							errDoc, ok := result.Results.(bson.M)
							if !ok || len(errDoc["error"].(string)) == 0 {
								t.Fatalf("\t%s\tShould return a single error document : %s", tests.Failed, err)
							}
							t.Logf("\t%s\tShould return a single error document.", tests.Success)
							return
						}

						// Process the results in mongo if the view is persisted.
						var viewItems []bson.M
						if len(vw.collection) > 0 {

							// Check the result message.
							msg, ok := result.Results.(bson.M)
							if !ok || msg["number_of_results"] != vw.number {
								t.Fatalf("\t%s\tShould be able to get %d items in the view", tests.Failed, vw.number)
							}
							t.Logf("\t%s\tShould be able to get %d items in the view.", tests.Success, vw.number)

							// Query the output collection.
							f := func(c *mgo.Collection) error {
								return c.Find(nil).All(&viewItems)
							}

							if err := db.ExecuteMGO(tests.Context, "testcollection", f); err != nil {
								t.Fatalf("\t%s\tShould be able to query the output collection : %s", tests.Failed, err)
							}
							t.Logf("\t%s\tShould be able to query the output collection.", tests.Success)

							// Verify that we get the same number of items back.
							if len(viewItems) != vw.number {
								t.Fatalf("\t%s\tShould be able to get %d items from the output collection", tests.Failed, vw.number)
							}
							t.Logf("\t%s\tShould be able to get %d items from the output collection.", tests.Success, vw.number)

							// Delete the persisted collection to clean up.
							f = func(c *mgo.Collection) error {
								return c.DropCollection()
							}

							if err := db.ExecuteMGO(tests.Context, vw.collection, f); err != nil {
								t.Fatalf("\t%s\tShould be able to drop the output collection : %s", tests.Failed, err)
							}
							t.Logf("\t%s\tShould be able to drop the output collection.", tests.Success)

						}

						// Otherwise, get the items directly from the result.
						var ok bool
						if vw.collection == "" {
							viewItems, ok = result.Results.([]bson.M)
							if !ok || len(viewItems) != vw.number {
								t.Fatalf("\t%s\tShould be able to get %d items in the view", tests.Failed, vw.number)
							}
							t.Logf("\t%s\tShould be able to get %d items in the view.", tests.Success, vw.number)
						}

						// Check the content of the items returned.
						commonData := wire.Result{
							Results: viewItems,
						}
						data, err := json.Marshal(commonData)
						if err != nil {
							t.Errorf("\t%s\tShould be able to marshal the result : %s", tests.Failed, err)
							return
						}
						t.Logf("\t%s\tShould be able to marshal the result.", tests.Success)

						for _, rslt := range vw.results {

							// We just need to find the string inside the result.
							if !strings.Contains(string(data), rslt) {
								t.Log("\t\tRsl:", string(data))
								for _, rslt := range vw.results {
									t.Log("\t\tExp:", rslt)
								}
								t.Errorf("\t%s\tShould have the correct result.", tests.Failed)
								return
							}
						}
						t.Logf("\t%s\tShould have the correct result", tests.Success)
						return
					}
				}

				t.Run(vw.viewName, tf)
			}
		}
	}
}

//==============================================================================

// execView represents the table of view execution tests.
type execView struct {
	fail        bool
	viewName    string
	itemKey     string
	number      int
	collection  string
	bufferLimit int
	results     []string
}
