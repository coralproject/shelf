package cmdgraph

import (
	"github.com/ardanlabs/kit/db"
	"github.com/cayleygraph/cayley"
	"github.com/spf13/cobra"
)

// graphCmd represents the parent for all graph cli commands.
var graphCmd = &cobra.Command{
	Use:   "graph",
	Short: "graph provides a CLI for add/removing relationships to/from a graph.",
}

var (
	// mgoDB holds the session for the DB access.
	mgoDB *db.DB

	// graphDB holds the graph handle for graph access.
	graphDB *cayley.Handle
)

// GetCommands returns the graph commands.
func GetCommands(conn *db.DB, store *cayley.Handle) *cobra.Command {
	mgoDB = conn
	graphDB = store

	addAddToGraph()
	addRemoveFromGraph()
	return graphCmd
}
