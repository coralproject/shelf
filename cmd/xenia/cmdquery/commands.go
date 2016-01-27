package cmdquery

import (
	"github.com/ardanlabs/kit/db"
	"github.com/spf13/cobra"
)

// queryCmd represents the parent for all query cli commands.
var queryCmd = &cobra.Command{
	Use:   "query",
	Short: "query provides a xenia CLI for managing and executing queries.",
}

// conn holds the session for the DB access.
var conn *db.DB

// GetCommands returns the query commands.
func GetCommands(db *db.DB) *cobra.Command {
	conn = db

	addUpsert()
	addGet()
	addDel()
	addExec()
	addList()
	return queryCmd
}
