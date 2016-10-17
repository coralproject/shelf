// Package tests implements users tests for the API layer.
package sponge_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/ardanlabs/kit/cfg"
	"github.com/ardanlabs/kit/tests"
	"github.com/cayleygraph/cayley"
	"github.com/cayleygraph/cayley/quad"
	"github.com/coralproject/shelf/internal/platform/db"
	cayleyshelf "github.com/coralproject/shelf/internal/platform/db/cayley"
	"github.com/coralproject/shelf/internal/sponge"
	"github.com/coralproject/shelf/internal/sponge/item/itemfix"
	"github.com/coralproject/shelf/internal/wire/pattern/patternfix"
)

const (
	// itemPrefix is the base name for items.
	itemPrefix = "ITEST_"

	// patternPrefix is the base name for patterns.
	patternPrefix = "PTEST_"
)

func init() {

	// Initialize the configuration and logging systems. Plus anything
	// else the web app layer needs.
	tests.Init("SPONGE")
}

// TestMain helps to clean up the test data.
func TestMain(m *testing.M) {
	os.Exit(runTest(m))
}

// runTest initializes the environment for the tests and allows for
// the proper return code if the test fails or succeeds.
func runTest(m *testing.M) int {

	// Initialize MongoDB using the `tests.TestSession` as the name of the
	// master session.
	if err := db.RegMasterSession(tests.Context, tests.TestSession, cfg.MustURL("MONGO_URI").String(), 0); err != nil {
		fmt.Println("Can't register master session: " + err.Error())
		return 1
	}

	db, err := db.NewMGO(tests.Context, tests.TestSession)
	if err != nil {
		fmt.Println("Can't get mongo session: " + err.Error())
		return 1
	}
	defer db.CloseMGO(tests.Context)

	loadPatterns("context", db)
	defer patternfix.Remove("context", db, patternPrefix)

	defer itemfix.Remove("context", db, itemPrefix)

	return m.Run()
}

// loadPatterns adds patterns to run tests.
func loadPatterns(context interface{}, db *db.DB) error {
	ps, _, err := patternfix.Get()
	if err != nil {
		return err
	}

	if err := patternfix.Add(context, db, ps[0:2]); err != nil {
		return err
	}

	return nil
}

// setup initializes for each individual test.
func setup(t *testing.T) (*db.DB, *cayley.Handle) {
	tests.ResetLog()

	db, err := db.NewMGO(tests.Context, tests.TestSession)
	if err != nil {
		t.Fatalf("%s\tShould be able to get a Mongo session : %v", tests.Failed, err)
	}

	store, err := cayleyshelf.New(cfg.MustURL("MONGO_URI").String(), nil)
	if err != nil {
		t.Fatalf("%s\tShould be able to get a cayley handle : %v", tests.Failed, err)
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

// TestImportRemoveItem tests the insert and update of an item.
func TestImportRemoveItem(t *testing.T) {
	db, store := setup(t)
	defer teardown(t, db, store)

	t.Log("Given the need to import an item.")
	{
		//----------------------------------------------------------------------
		// Get the fixture.

		items, err := itemfix.Get()
		if err != nil {
			t.Fatalf("\t%s\tShould be able to retrieve the fixture : %v", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to retrieve the fixture.", tests.Success)

		//----------------------------------------------------------------------
		// Import the Item.

		if err := sponge.Import(tests.Context, db, store, &items[0]); err != nil {
			t.Fatalf("\t%s\tShould be able to import an item : %s", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to import an item", tests.Success)

		//----------------------------------------------------------------------
		// Check the inferred relationship.

		p := cayley.StartPath(store, quad.String("ITEST_80aa936a-f618-4234-a7be-df59a14cf8de")).Out(quad.String("authored"))
		it, _ := p.BuildIterator().Optimize()
		defer it.Close()
		for it.Next() {
			token := it.Result()
			value := store.NameOf(token)
			if quad.NativeOf(value) != "ITEST_d1dfa366-d2f7-4a4a-a64f-af89d4c97d82" {
				t.Fatalf("\t%s\tShould be able to get the inferred relationships from the graph", tests.Failed)
			}
		}
		if err := it.Err(); err != nil {
			t.Fatalf("\t%s\tShould be able to get the inferred relationships from the graph : %s", tests.Failed, err)
		}
		it.Close()
		t.Logf("\t%s\tShould be able to get the inferred relationships from the graph.", tests.Success)

		//----------------------------------------------------------------------
		// Import the Item again to test for duplicate imports.

		if err := sponge.Import(tests.Context, db, store, &items[0]); err != nil {
			t.Fatalf("\t%s\tShould be able to import a duplicate item : %s", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to import a duplicate item", tests.Success)

		//----------------------------------------------------------------------
		// Remove the item.

		if err := sponge.Remove(tests.Context, db, store, items[0].ID); err != nil {
			t.Fatalf("\t%s\tShould be able to remove the item : %s", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to remove the item", tests.Success)

		//----------------------------------------------------------------------
		// Check the inferred relationships.

		p = cayley.StartPath(store, quad.String("ITEST_80aa936a-f618-4234-a7be-df59a14cf8de")).Out(quad.String("authored"))
		it, _ = p.BuildIterator().Optimize()
		defer it.Close()

		var count int
		for it.Next() {
			count++
		}
		if err := it.Err(); err != nil {
			t.Fatalf("\t%s\tShould be able to confirm removed relationships : %s", tests.Failed, err)
		}

		if count > 0 {
			t.Fatalf("\t%s\tShould be able to confirm removed relationships.", tests.Failed)
		}
		t.Logf("\t%s\tShould be able to confirm removed relationships.", tests.Success)
	}
}
