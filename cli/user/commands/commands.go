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
	var email string

	cmdCreate := &cobra.Command{
		Use:   "create [-n name -p password -e email]",
		Short: "Creates a new user",
		Long:  `Creates adds a new user to the system.`,
		Run: func(cmd *cobra.Command, args []string) {

		},
	}

	cmdCreate.Flags().StringVarP(&name, "name", "n", "", "name of user")
	cmdCreate.Flags().StringVarP(&pass, "pass", "p", "", "password for user")
	cmdCreate.Flags().StringVarP(&email, "email", "e", "", "email of user")

	rootCmd.AddCommand(cmdCreate)
}

// cmdCreate handles the creation of users.
func cmdUpdate() {
	var email string
	var name string
	var pass string

	cmdUpdate := &cobra.Command{
		Use:   "create [-e email -n name -p password]",
		Short: "Updates a existing user.",
		Long:  `Updates a existing user's information in the system.`,
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: Handle the create
			fmt.Println(name, pass)
		},
	}

	cmdUpdate.Flags().StringVarP(&email, "email", "e", "", "email of user")
	cmdUpdate.Flags().StringVarP(&name, "name", "n", "", "name of user")
	cmdUpdate.Flags().StringVarP(&pass, "password", "n", "", "password of user")

	rootCmd.AddCommand(cmdUpdate)
}
