package cmduser

import (
	"github.com/spf13/cobra"
)

var getLong = `Retrieves a user record from the system using any of the supplied keys.
When retrieving a record from the system, atleast one of the key points must be supplied,
that is either the email(-e), the name(-n) or public id(-p).

Each flag's presence and use is based on the order of importance:

	'public_id' is first importance.
	'email' is second in importance.
	'name' is the third and last in importance.

If all are supplied, the highest flag with the most priority gets used.

Example:

	1. To get a user using it's name:

		user get -n "Alex Boulder"

	2. To get a user using it's email address:

		user get -e alex.boulder@gmail.com

	3. To get a user using it's public id number:

		user get -p 199550d7-484d-4440-801f-390d44911ade
`

// get contains the state for this command.
var get struct {
	name  string
	pid   string
	email string
}

// addGet handles the retrival users records, displayed in json formatted response.
func addGet() {
	cmd := &cobra.Command{
		Use:   "get [-n name -p public_id -e email]",
		Short: "Retrieves a user record",
		Long:  getLong,
		Run:   runGet,
	}

	cmd.Flags().StringVarP(&get.name, "name", "n", "", "name of the user record")
	cmd.Flags().StringVarP(&get.pid, "public_id", "p", "", "publicId of the user record")
	cmd.Flags().StringVarP(&get.email, "email", "e", "", "email of the user")

	userCmd.AddCommand(cmd)
}

// runGet is the code that implements the get command.
func runGet(cmd *cobra.Command, args []string) {
}
