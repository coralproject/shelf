package wire_test

import (
	"reflect"
	"testing"

	"github.com/ardanlabs/kit/db"
	"github.com/ardanlabs/kit/tests"
	"github.com/cayleygraph/cayley"
	"github.com/cayleygraph/cayley/quad"
	"github.com/coralproject/shelf/internal/wire"
	"github.com/coralproject/shelf/internal/wire/pattern/patternfix"
	"gopkg.in/mgo.v2/bson"
)

// setupGraph initializes an in-memory Cayley graph and logging for an individual test.
func setupGraph(t *testing.T) (*db.DB, *cayley.Handle, []bson.M) {
	tests.ResetLog()

	_, items, err := patternfix.Get()
	if err != nil {
		t.Fatalf("%s\tShould load item records from the fixture file : %v", tests.Failed, err)
	}
	t.Logf("%s\tShould load item records from the fixture file.", tests.Success)

	db, err := db.NewMGO(tests.Context, tests.TestSession)
	if err != nil {
		t.Fatalf("%s\tShould be able to get a Mongo session : %v", tests.Failed, err)
	}

	store, err := cayley.NewMemoryGraph()
	if err != nil {
		t.Fatalf("\t%s\tShould be able to create a new Cayley graph : %v", tests.Failed, err)
	}
	t.Logf("\t%s\tShould be able to create a new Cayley graph.", tests.Success)

	return db, store, items
}

// TestAddRemoveGraph tests if we can add/remove relationship quads to/from cayley.
func TestAddRemoveGraph(t *testing.T) {
	_, store, _ := setupGraph(t)
	defer tests.DisplayLog()

	t.Log("Given the need to add/remove relationship quads from the Cayley graph.")
	{
		t.Log("\tWhen starting from an empty graph")
		{
			//----------------------------------------------------------------------
			// Create some example parameters to import into the graph.

			params1 := wire.QuadParams{
				Subject:   "frodo",
				Predicate: "carries",
				Object:    "the ring",
			}
			params2 := wire.QuadParams{
				Subject:   "orcs",
				Predicate: "chase",
				Object:    "frodo",
			}
			params := []wire.QuadParams{params1, params2}

			//----------------------------------------------------------------------
			// Add the relationships to the graph.

			if err := wire.AddToGraph(tests.Context, store, params); err != nil {
				t.Fatalf("\t%s\tShould be able to add relationships to the graph : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to add relationships to the graph.", tests.Success)

			//----------------------------------------------------------------------
			// Get the relationship quads from the graph.

			p := cayley.StartPath(store, quad.String("frodo")).Out(quad.String("carries"))
			it, _ := p.BuildIterator().Optimize()
			defer it.Close()
			for it.Next() {
				token := it.Result()
				value := store.NameOf(token)
				if quad.NativeOf(value) != "the ring" {
					t.Fatalf("\t%s\tShould be able to get the relationships from the graph", tests.Failed)
				}
			}
			if err := it.Err(); err != nil {
				t.Fatalf("\t%s\tShould be able to get the relationships from the graph : %s", tests.Failed, err)
			}
			it.Close()

			p = cayley.StartPath(store, quad.String("orcs")).Out(quad.String("chase"))
			it, _ = p.BuildIterator().Optimize()
			defer it.Close()
			for it.Next() {
				token := it.Result()
				value := store.NameOf(token)
				if quad.NativeOf(value) != "frodo" {
					t.Fatalf("\t%s\tShould be able to get the relationships from the graph", tests.Failed)
				}
			}
			if err := it.Err(); err != nil {
				t.Fatalf("\t%s\tShould be able to get the relationships from the graph : %s", tests.Failed, err)
			}
			it.Close()
			t.Logf("\t%s\tShould be able to get relationships from the graph.", tests.Success)

			//----------------------------------------------------------------------
			// Remove the relationships from the graph.

			if err := wire.RemoveFromGraph(tests.Context, store, params); err != nil {
				t.Fatalf("\t%s\tShould be able to remove relationships from the graph : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to remove relationships from the graph.", tests.Success)

			//----------------------------------------------------------------------
			// Try to get the relationships.

			var count int
			p = cayley.StartPath(store, quad.String("frodo")).Out(quad.String("carries"))
			it, _ = p.BuildIterator().Optimize()
			defer it.Close()
			for it.Next() {
				count++
			}
			if err := it.Err(); err != nil {
				t.Fatalf("\t%s\tShould be able to verify the empty graph : %s", tests.Failed, err)
			}
			it.Close()

			p = cayley.StartPath(store, quad.String("orcs")).Out(quad.String("chase"))
			it, _ = p.BuildIterator().Optimize()
			defer it.Close()
			for it.Next() {
				count++
			}
			if err := it.Err(); err != nil {
				t.Fatalf("\t%s\tShould be able to verify the empty graph : %s", tests.Failed, err)
			}
			it.Close()

			if count != 0 {
				t.Fatalf("\t%s\tShould be able to verify the empty graph", tests.Failed)
			}
			t.Logf("\t%s\tShould be able to verify the empty graph.", tests.Success)
		}
	}
}

// TestGraphParamFail tests if we can handle invalid quad parameters.
func TestGraphParamFail(t *testing.T) {
	_, store, _ := setupGraph(t)
	defer tests.DisplayLog()

	t.Log("Given the need to add/remove relationship quads from the Cayley graph.")
	{
		t.Log("\tWhen starting from an empty graph")
		{
			//----------------------------------------------------------------------
			// Create some example parameters to import into the graph.

			params1 := wire.QuadParams{
				Subject:   "",
				Predicate: "",
				Object:    "the ring",
			}
			params2 := wire.QuadParams{
				Subject:   "orcs",
				Predicate: "chase",
				Object:    "frodo",
			}
			params := []wire.QuadParams{params1, params2}

			//----------------------------------------------------------------------
			// Try to add the invalid relationship to the graph.

			if err := wire.AddToGraph(tests.Context, store, params); err == nil {
				t.Fatalf("\t%s\tShould be able to catch invalid quad parameters on add : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to catch invalid quad parameters on add.", tests.Success)

			//----------------------------------------------------------------------
			// Try to remove the invalid relationship to the graph.

			if err := wire.RemoveFromGraph(tests.Context, store, params); err == nil {
				t.Fatalf("\t%s\tShould be able to catch invalid quad parameters on remove : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to catch invalid quad parameters on remove.", tests.Success)
		}
	}
}

// TestInferRelationships tests if we can infer relationships based on patterns.
func TestInferRelationships(t *testing.T) {
	db, _, items := setupGraph(t)
	defer tests.DisplayLog()

	t.Log("Given the need to infer relationships from items.")
	{
		t.Log("\tWhen starting from a slice of input item documents.")
		{
			//----------------------------------------------------------------------
			// Infer Relationships.

			quadParams, err := wire.InferRelationships(tests.Context, db, items)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to infer relationships : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to infer relationships.", tests.Success)

			//----------------------------------------------------------------------
			// Check the inferred relationships.

			if len(quadParams) != 10 {
				t.Fatalf("\t%s\tShould infer 10 relationships", tests.Failed)
			}
			t.Logf("\t%s\tShould infer 10 relationships.", tests.Success)

			relCounts := make(map[string]int)
			for _, params := range quadParams {
				relCounts[params.Predicate]++
			}

			expectedCounts := map[string]int{
				"authored":    3,
				"on":          3,
				"parented_by": 1,
				"has_role":    2,
				"part_of":     1,
			}

			if eq := reflect.DeepEqual(relCounts, expectedCounts); !eq {
				t.Fatalf("\t%s\tShould have expected numbers of relationship predicates", tests.Failed)
			}
			t.Logf("\t%s\tShould have expected numbers of relationship predicates.", tests.Success)
		}
	}
}
