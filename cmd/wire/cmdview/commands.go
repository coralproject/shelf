package cmdview

import (
	"github.com/ardanlabs/kit/db"
	"github.com/cayleygraph/cayley"
	"github.com/spf13/cobra"
)

// viewCmd represents the parent for all view cli commands.
var viewCmd = &cobra.Command{
	Use:   "view",
	Short: "view provides a CLI for executing a view.",
}

var (
	// mgoDB holds the session for the DB access.
	mgoDB *db.DB

	// graphDB holds the graph handle for graph access.
	graphDB *cayley.Handle
)

// GetCommands returns the view commands.
func GetCommands(conn *db.DB, store *cayley.Handle) *cobra.Command {
	mgoDB = conn
	graphDB = store

	addExecute()
	return viewCmd
}
