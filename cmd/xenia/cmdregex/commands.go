package cmdregex

import "github.com/spf13/cobra"

// regexCmd represents the parent for all regex cli commands.
var regexCmd = &cobra.Command{
	Use:   "regex",
	Short: "regex provides a xenia CLI for managing regexs.",
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
	return regexCmd
}
