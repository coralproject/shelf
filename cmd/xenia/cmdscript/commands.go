package cmdscript

import "github.com/spf13/cobra"

// scriptCmd represents the parent for all script cli commands.
var scriptCmd = &cobra.Command{
	Use:   "script",
	Short: "script provides a xenia CLI for managing scripts.",
}

// GetCommands returns the query commands.
func GetCommands() *cobra.Command {
	addUpsert()
	addGet()
	addDel()
	addList()
	return scriptCmd
}
