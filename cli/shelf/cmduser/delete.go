package cmduser

import (
	"github.com/spf13/cobra"
)

var delLong = `Deletes a user record from the system using any of the supplied keys.
When deleting a record from the system, atleast one of the key points must be supplied,
that is either the email(-e), the name(-n) or public id(-p).

Each flag's presence and use is based on the order of importance:

	'public_id' is first importance.
	'email' is second in importance.
	'name' is the third and last in importance.

If all are supplied, the highest flag with the most priority gets used.

Example:

	1. To delete a user using it's name:

		user delete -n "Alex Boulder"

	2. To delete a user using it's email address:

		user delete -e alex.boulder@gmail.com

	3. To delete a user using it's public id number:

		user delete -p 199550d7-484d-4440-801f-390d44911ade
`

// del contains the state for this command.
var del struct {
	name  string
	pid   string
	email string
}

// addDel handles the deletion of user records.
func addDel() {
	cmd := &cobra.Command{
		Use:   "delete [-n name -p public_id -e email]",
		Short: "Deletes a user record",
		Long:  delLong,
		Run:   runDel,
	}

	cmd.Flags().StringVarP(&del.name, "name", "n", "", "name of the user record")
	cmd.Flags().StringVarP(&del.pid, "public_id", "p", "", "publicId of the user record")
	cmd.Flags().StringVarP(&del.email, "email", "e", "", "email of the user")

	userCmd.AddCommand(cmd)
}

// runDel is the code that implements the delete command.
func runDel(cmd *cobra.Command, args []string) {
}
