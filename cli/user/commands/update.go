package commands

import (
	"strings"

	"github.com/coralproject/shelf/cli/user/db"
	"github.com/coralproject/shelf/pkg/log"
	"github.com/coralproject/shelf/pkg/mongo"
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
		Run:   runUpdate,
	}

	cmd.Flags().StringVarP(&upd.utype, "type", "t", "", "type of update")
	cmd.Flags().StringVarP(&upd.email, "email", "e", "", "email of user")
	cmd.Flags().StringVarP(&upd.oldValue, "old", "o", "", "old value of user")
	cmd.Flags().StringVarP(&upd.newValue, "new", "n", "", "new value for user")

	rootCmd.AddCommand(cmd)
}

// runUpdate is the code that implements the update command.
func runUpdate(cmd *cobra.Command, args []string) {
	if upd.utype == "" && upd.email == "" {
		cmd.Help()
		return
	}

	if upd.utype == "" {
		cmd.Println("\n\tError: type(-t) can not be empty. Please supply a name using the `-t` or `-type` flag")
		return
	}

	if upd.email == "" {
		cmd.Println("\n\tError: email(-e) can not be empty. Please supply a email using the `-e` or `-email` flag")
		return
	}

	if upd.newValue == "" {
		cmd.Println("\n\tError: newValue(-n) can not be empty. Please supply the new value using the `-n` or `-new` flag")
		return
	}

	if upd.utype == "auth" {
		if upd.oldValue == "" {
			cmd.Println("\n\tError: oldValue(-o) can not be empty. Please supply the old value using the `-o` or `-old` flag")
			return
		}
	}

	if upd.email != "" {
		// Trying to match the complexity of email address is unecessary, as far as we
		// have a valid expectation pattern,we can skip alot of the mess.
		// TODO: should we use something more robust?
		if !strings.Contains(upd.email, "@") {
			cmd.Println("\n\tError: Email address must be a valid addresss. Please supply a correct email address.")
			return
		}
	}

	// Initialize the mongodb session.
	mongo.InitMGO()

	var updateEmail = func() {
		log.Dev("commands", "runUpdate", "Email[%s]", upd.email)

		user, err := db.GetUserByEmail(upd.email)
		if err != nil {
			log.Error("commands", "runUpdate", err, "Completed")
			return
		}

		if err := db.UpdateEmail(user, upd.newValue); err != nil {
			log.Error("commands", "runUpdate", err, "Completed")
		}

	}

	var updateName = func() {
		log.Dev("commands", "runUpdate", "Name[%s]", upd.newValue)

		user, err := db.GetUserByEmail(upd.email)
		if err != nil {
			log.Error("commands", "runUpdate", err, "Completed")
			return
		}

		if err := db.UpdateName(user, upd.newValue); err != nil {
			log.Error("commands", "runUpdate", err, "Completed")
		}

	}

	var updatePassword = func() {
		log.Dev("commands", "runUpdate", "Password[%s]", upd.newValue)

		user, err := db.GetUserByEmail(upd.email)
		if err != nil {
			log.Error("commands", "runUpdate", err, "Completed")
			return
		}

		if err := db.UpdatePassword(user, upd.oldValue, upd.newValue); err != nil {
			log.Error("commands", "runUpdate", err, "Completed")
		}
	}

	switch upd.utype {
	case "auth":
		updatePassword()
	case "name":
		updateName()
	case "email":
		updateEmail()
	}
}
