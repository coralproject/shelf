package cmdquery

import "github.com/spf13/cobra"

// queryCmd represents the parent for all query cli commands.
var queryCmd = &cobra.Command{
	Use:   "query",
	Short: "query provides a xenia CLI for managing and executing queries.",
}

// mgoSession holds the master session for the DB access.
var mgoSession string

// GetCommands returns the query commands.
func GetCommands(mgoSes string) *cobra.Command {
	mgoSession = mgoSes

	addUpsert()
	addGet()
	addDel()
	addExec()
	addList()
	return queryCmd
}
