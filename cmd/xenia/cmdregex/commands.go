package cmdregex

import "github.com/spf13/cobra"

// regexCmd represents the parent for all regex cli commands.
var regexCmd = &cobra.Command{
	Use:   "regex",
	Short: "regex provides a xenia CLI for managing regexs.",
}

// GetCommands returns the query commands.
func GetCommands() *cobra.Command {
	addUpsert()
	addGet()
	addDel()
	addList()
	return regexCmd
}
