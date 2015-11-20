package commands

import "github.com/spf13/cobra"

// rootCmd represents the parent for all cli commands.
var rootCmd = &cobra.Command{
	Use:   "query",
	Short: "query provides a shelf CLI for managing and executing queries.",
}

// Run will process the command line arguments for the application.
func Run() {
	addCreate()
	addGet()
	addDel()
	addExecute()
	addUpd()
	rootCmd.Execute()
}
