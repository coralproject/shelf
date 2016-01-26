package cmdscript

import "github.com/spf13/cobra"

// scriptCmd represents the parent for all script cli commands.
var scriptCmd = &cobra.Command{
	Use:   "script",
	Short: "script provides a xenia CLI for managing scripts.",
}

// mgoSession holds the master session for the DB access.
var mgoSession string

// GetCommands returns the query commands.
func GetCommands(mgoSes string) *cobra.Command {
	mgoSession = mgoSes

	addUpsert()
	addGet()
	addDel()
	addList()
	return scriptCmd
}
