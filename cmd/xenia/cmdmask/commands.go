package cmdmask

import (
	"github.com/ardanlabs/kit/db"
	"github.com/spf13/cobra"
)

// maskCmd represents the parent for all mask cli commands.
var maskCmd = &cobra.Command{
	Use:   "mask",
	Short: "mask provides a xenia CLI for managing and executing masks.",
}

// Capture the database connection.
var conn *db.DB

// GetCommands returns the query commands.
func GetCommands(db *db.DB) *cobra.Command {
	conn = db

	addUpsert()
	addGet()
	addDel()
	return maskCmd
}
