package cmduser

import (
	"github.com/spf13/cobra"
)

var createLong = `Creates adds a new user to the system.
When creating a new user, the name(-n), email(-e) and password(-p) must all be supplied.

Example:

	user create -n "Alex Boulder" -e alex.boulder@gmail.com -p yefc*7fdf92
`

// create contains the state for this command.
var create struct {
	name  string
	pass  string
	email string
}

// addCreate handles the creation of users.
func addCreate() {
	cmd := &cobra.Command{
		Use:   "create [-n name -p password -e email]",
		Short: "Creates a new user",
		Long:  createLong,
		Run:   runCreate,
	}

	cmd.Flags().StringVarP(&create.name, "name", "n", "", "name of user")
	cmd.Flags().StringVarP(&create.pass, "pass", "p", "", "password for user")
	cmd.Flags().StringVarP(&create.email, "email", "e", "", "email of user")

	userCmd.AddCommand(cmd)
}

// runCreate is the code that implements the create command.
func runCreate(cmd *cobra.Command, args []string) {
}
