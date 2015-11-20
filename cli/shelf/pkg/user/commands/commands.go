package commands

import "github.com/spf13/cobra"

// rootCmd represents the parent for all cli commands.
var rootCmd = &cobra.Command{
	Use:   "user",
	Short: "user provides a shelf CLI for managing user records.",
}

// GetCommand returns the user root commander.
func GetCommand() *cobra.Command {
	addAuth()
	addCreate()
	addGet()
	addUpd()
	addDel()
	return rootCmd
}

// Run will process the command line arguments for the application.
func Run() {
	GetCommand().Execute()
}
