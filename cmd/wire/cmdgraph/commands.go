package cmdgraph

import (
	"github.com/cayleygraph/cayley"
	"github.com/spf13/cobra"
)

// graphCmd represents the parent for all graph cli commands.
var graphCmd = &cobra.Command{
	Use:   "graph",
	Short: "graph provides a CLI for add/removing relationships to/from a graph.",
}

// graphDB holds the graph handle for graph access.
var graphDB *cayley.Handle

// GetCommands returns the graph commands.
func GetCommands(store *cayley.Handle) *cobra.Command {
	graphDB = store

	addAddToGraph()
	addRemoveFromGraph()
	return graphCmd
}
