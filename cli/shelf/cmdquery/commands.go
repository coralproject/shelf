package cmdquery

import "github.com/spf13/cobra"

// queryCmd represents the parent for all query cli commands.
var queryCmd = &cobra.Command{
	Use:   "query",
	Short: "query provides a shelf CLI for managing and executing queries.",
}

// GetCommands returns the query commands.
func GetCommands() *cobra.Command {
	addCreate()
	addGet()
	addDel()
	addExecute()
	addUpd()
	return queryCmd
}
