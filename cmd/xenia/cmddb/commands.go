package cmddb

import (
	"github.com/coralproject/shelf/internal/platform/db"
	"github.com/spf13/cobra"
)

// dbCmd represents the parent for all database cli commands.
var dbCmd = &cobra.Command{
	Use:   "db",
	Short: "db will create a kit database and validate everything exists.",
}

// conn holds the session for the DB access.
var conn *db.DB

// GetCommands returns the db commands.
func GetCommands(db *db.DB) *cobra.Command {
	conn = db

	addCreate()
	return dbCmd
}
