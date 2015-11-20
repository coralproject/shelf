package commands

import "github.com/spf13/cobra"

// rootCmd represents the parent for all cli commands.
var rootCmd = &cobra.Command{
	Use:   "query",
	Short: "query provides a shelf CLI for managing and executing queries.",
}

// GetCommand returns the user root commander.
func GetCommand() *cobra.Command {
	addCreate()
	addGet()
	addDel()
	addExecute()
	addUpd()
	return rootCmd
}

// Run will process the command line arguments for the application.
func Run() {
	GetCommand().Execute()
}
