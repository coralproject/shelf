package cmdquery

import "github.com/spf13/cobra"

// queryCmd represents the parent for all query cli commands.
var queryCmd = &cobra.Command{
	Use:   "query",
	Short: "query provides a xenia CLI for managing and executing queries.",
}

// GetCommands returns the query commands.
func GetCommands() *cobra.Command {
	addUpsert()
	addGet()
	addDel()
	addExec()
	addList()
	return queryCmd
}
