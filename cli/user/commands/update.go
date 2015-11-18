package commands

import (
	"github.com/coralproject/shelf/cli/user/db"
	"github.com/spf13/cobra"
)

var updLong = `Updates a existing user's information in the system.
	When updating a user record, a update type is required, to determine what values
	are to be expected including the records email which is used as the record search point.

	Examples:

		1. To update the 'name' of a giving record. Simple set the update type to "name",
		and supply the email and new name.

			user update -t name -e shou.lou@gmail -n "Shou Lou FengZhu"

		2. To update the 'email' of a giving record. Simple set the update type to "email",
		and supply the current email and new email.

			user update -t email -e shou.lou@gmail -n shou.lou.fengzhu@gmail.com

		3. To update the 'password' of a giving record. Simple set the update type to "auth",
		and supply the current email of the record, the current password of the record and
		the new password

			user update -t auth -e shou.lou@gmail -o oldPassword -n newPassword
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
		Run:   runUpd,
	}

	cmd.Flags().StringVarP(&upd.utype, "type", "t", "", "type of update")
	cmd.Flags().StringVarP(&upd.email, "email", "e", "", "email of user")
	cmd.Flags().StringVarP(&upd.oldValue, "old", "o", "", "old value of user")
	cmd.Flags().StringVarP(&upd.newValue, "new", "n", "", "new value for user")

	rootCmd.AddCommand(cmd)
}

// runUpd is the code that implements the update command.
func runUpd(cmd *cobra.Command, args []string) {
	if upd.utype == "" {
		cmd.Println("\n\tname(-n) can not be empty. please supply a name using the `-n` or `-name` flag")
	}

	// Initialize the mongodb session.
	db.InitMGO()
}
