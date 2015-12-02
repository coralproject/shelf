package cmduser

import (
	"github.com/spf13/cobra"
)

var updLong = `Use create to add a new user to the system. The user email
must be unique for every user.

Example:
  ./shelf user update -n "Bill Kennedy" -e "bill@ardanlabs.com" -p "yefc*7fdf92"
`

// upd contains the state for this command.
var upd struct {
	utype    string
	email    string
	oldValue string
	newValue string
}

// addUpd handles the update of user records.
func addUpd() {
	cmd := &cobra.Command{
		Use:   "update -t auth|email|name [-e email ....extraOptions ]",
		Short: "Updates a existing user information",
		Long:  updLong,
		Run:   runUpdate,
	}

	cmd.Flags().StringVarP(&upd.utype, "type", "t", "", "type of update")
	cmd.Flags().StringVarP(&upd.email, "email", "e", "", "email of user")
	cmd.Flags().StringVarP(&upd.oldValue, "old", "o", "", "old value of user")
	cmd.Flags().StringVarP(&upd.newValue, "new", "n", "", "new value for user")

	userCmd.AddCommand(cmd)
}

// runUpdate is the code that implements the update command.
func runUpdate(cmd *cobra.Command, args []string) {
}
