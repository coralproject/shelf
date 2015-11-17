package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

// rootCmd represents the parent for all cli commands.
var rootCmd = &cobra.Command{Use: "user"}

// Run will process the command line arguments for the application.
func Run() {
	cmdCreate()
	rootCmd.Execute()
}

// cmdCreate handles the creation of users.
func cmdCreate() {
	var name string
	var pass string

	cmdCreate := &cobra.Command{
		Use:   "create [-n name -p password]",
		Short: "Creates a new user",
		Long:  `Creates adds a new user to the system.`,
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: Handle the create
			fmt.Println(name, pass)
		},
	}

	cmdCreate.Flags().StringVarP(&name, "name", "n", "", "name of user")
	cmdCreate.Flags().StringVarP(&pass, "pass", "p", "", "password for user")

	rootCmd.AddCommand(cmdCreate)
}
