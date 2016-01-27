package cmdscript

import (
	"github.com/ardanlabs/kit/db"
	"github.com/spf13/cobra"
)

// scriptCmd represents the parent for all script cli commands.
var scriptCmd = &cobra.Command{
	Use:   "script",
	Short: "script provides a xenia CLI for managing scripts.",
}

// conn holds the session for the DB access.
var conn *db.DB

// GetCommands returns the script commands.
func GetCommands(db *db.DB) *cobra.Command {
	conn = db

	addUpsert()
	addGet()
	addDel()
	addList()
	return scriptCmd
}
