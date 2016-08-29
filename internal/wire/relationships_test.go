package wire_test

import (
	"testing"

	"github.com/ardanlabs/kit/tests"
	"github.com/cayleygraph/cayley"
	"github.com/cayleygraph/cayley/quad"
	"github.com/coralproject/shelf/internal/wire"
)

// setupGraph initializes an in-memory Cayley graph and logging for an individual test.
func setupGraph(t *testing.T) *cayley.Handle {
	tests.ResetLog()

	store, err := cayley.NewMemoryGraph()
	if err != nil {
		t.Fatalf("\t%s\tShould be able to create a new Cayley graph : %v", tests.Failed, err)
	}
	t.Logf("\t%s\tShould be able to create a new Cayley graph.", tests.Success)

	return store
}

// TestAddRemoveGraph tests if we can add/remove relationship quads to/from cayley.
func TestAddRemoveGraph(t *testing.T) {
	store := setupGraph(t)
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
			t.Logf("\t%s\tShould be able to verify the empty graph.", tests.Success)
		}
	}
}

// TestGraphParamFail tests if we can handle invalid quad parameters.
func TestGraphParamFail(t *testing.T) {
	store := setupGraph(t)
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
