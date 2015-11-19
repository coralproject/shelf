package commands

import "github.com/spf13/cobra"

// rootCmd represents the parent for all cli commands.
var rootCmd = &cobra.Command{
	Use:   "user",
	Short: "user provides a shelf CLI for managing user records.",
}

// Run will process the command line arguments for the application.
func Run() {
	addAuth()
	addCreate()
	addGet()
	addUpd()
	addDel()
	rootCmd.Execute()
}
