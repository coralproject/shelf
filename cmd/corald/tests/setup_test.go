// Package tests implements users tests for the API layer.
package tests

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/cayleygraph/cayley"

	"github.com/ardanlabs/kit/cfg"
	"github.com/ardanlabs/kit/tests"
	"github.com/coralproject/shelf/cmd/corald/routes"
	"github.com/coralproject/shelf/internal/platform/db"
	"github.com/coralproject/shelf/internal/sponge"
	"github.com/coralproject/shelf/internal/sponge/item"
	"github.com/coralproject/shelf/internal/sponge/item/itemfix"
	"github.com/coralproject/shelf/tstdata"
)

var a http.Handler

//==============================================================================

// TestMain helps to clean up the test data.
func TestMain(m *testing.M) {
	os.Exit(runTest(m))
}

// runTest initializes the environment for the tests and allows for
// the proper return code if the test fails or succeeds.
func runTest(m *testing.M) int {

	// Create stub server for Sponged.
	server := setup()
	cfg.SetString("SPONGED_URL", server)

	mongoURI := cfg.MustURL("MONGO_URI")

	// Initialize MongoDB using the `tests.TestSession` as the name of the
	// master session.
	if err := db.RegMasterSession(tests.Context, tests.TestSession, mongoURI.String(), 0); err != nil {
		fmt.Println("Can't register master session: " + err.Error())
		return 1
	}

	a = routes.API()

	// Snatch the mongo session so we can create some test data.
	db, err := db.NewMGO(tests.Context, tests.TestSession)
	if err != nil {
		fmt.Println("Unable to get Mongo session")
		return 1
	}
	defer db.CloseMGO(tests.Context)

	if err = db.NewCayley(tests.Context, tests.TestSession); err != nil {
		fmt.Println("Unable to get Cayley support")
	}

	store, err := db.GraphHandle(tests.Context)
	if err != nil {
		fmt.Println("Unable to get Cayley handle")
		return 1
	}
	defer store.Close()

	if err := tstdata.Generate(db); err != nil {
		fmt.Println("Could not generate test data.")
		return 1
	}
	defer tstdata.Drop(db)

	if err := loadItems("context", db, store); err != nil {
		fmt.Println("Could not import items")
		return 1
	}
	defer unloadItems("context", db, store)

	return m.Run()
}

// loadItems adds items to run tests.
func loadItems(context interface{}, db *db.DB, store *cayley.Handle) error {
	items, err := itemfix.Get()
	if err != nil {
		return err
	}

	for _, itm := range items {
		if err := sponge.Import(context, db, store, &itm); err != nil {
			return err
		}
	}

	return nil
}

// unloadItems removes items from the items collection and the graph.
func unloadItems(context interface{}, db *db.DB, store *cayley.Handle) error {
	items, err := itemfix.Get()
	if err != nil {
		return err
	}

	for _, itm := range items {
		if err := sponge.Remove(context, db, store, itm.ID); err != nil {
			return err
		}
	}

	return nil
}

// Setup sets the stubbed Sponged service.
// It will check that the data it gets is the appropiate.
func setup() string {

	// Initialization of stub server for Sponged.
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var err error
		var itm item.Item

		// check that the row is what we want it to be
		switch r.RequestURI {
		// TO DO : CHECK DATA THAT ARRIVES TO EACH ENDPOINT ON SPONGED TO BE THE APPROPIATE
		case "/v1/item/ITEST_d16790f8-13e9-4cb4-b9ef-d82835589660":
			itm = item.Item{ID: "ITEST_d16790f8-13e9-4cb4-b9ef-d82835589660", Type: "comment", Data: map[string]interface{}{"body": "Something."}}
			itm.Data["flagged_by"] = []string{"ITEST_80aa936a-f618-4234-a7be-df59a14cf8de"}
			b, _ := json.Marshal([]item.Item{itm})
			w.Write(b)
		case "/v1/item/ITEST_d1dfa366-d2f7-4a4a-a64f-af89d4c97d82":
			itm = item.Item{ID: "ITEST_d1dfa366-d2f7-4a4a-a64f-af89d4c97d82", Type: "comment", Data: map[string]interface{}{"body": "Something."}}
			itm.Data["flagged_by"] = []string{"ITEST_a63af637-58af-472b-98c7-f5c00743bac6"}
			b, _ := json.Marshal([]item.Item{itm})
			w.Write(b)
		case "/v1/item/wrongitem":
			b, _ := json.Marshal([]item.Item{})
			w.Write(b)

		case "/v1/item":
			// Decode the item.
			err = json.NewDecoder(r.Body).Decode(&itm)
		default:
			err = errors.New("Bad Request")
		}

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
		if err == nil {
			w.WriteHeader(http.StatusOK)
		}
		w.Header().Set("Content-Type", "application/json")

		fmt.Fprintln(w, err)

	}))

	return server.URL
}
