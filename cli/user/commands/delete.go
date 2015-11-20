package commands

import (
	"strings"

	"github.com/coralproject/shelf/cli/user/db"
	"github.com/coralproject/shelf/pkg/db/mongo"
	"github.com/coralproject/shelf/pkg/log"
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

	rootCmd.AddCommand(cmd)
}

// runDel is the code that implements the delete command.
func runDel(cmd *cobra.Command, args []string) {
	if del.name == "" && del.pid == "" && del.email == "" {
		// cmd.Println("\n\tError: Atleast one key flag must be supplied. Use either -name(-n), -email(-e) or -public_id(-p) to delete the desired record.")
		cmd.Help()
		return
	}

	if del.email != "" {
		// Trying to match the complexity of email address is unecessary, as far as we
		// have a valid expectation pattern,we can skip alot of the mess.
		// TODO: should we use something more robust?
		if !strings.Contains(del.email, "@") {
			cmd.Println("\n\tError: Email address must be a valid addresss. Please supply a correct email address.")
			return
		}
	}

	// Initialize the mongodb session.
	mongo.InitMGO()

	if del.pid != "" {
		log.Dev("commands", "runDel", "Pid[%s]", del.pid)
		user, err := db.GetUserByPublicID(del.pid)
		if err != nil {
			log.Error("commands", "runDel", err, "Completed")
			return
		}

		err = db.Delete(user)
		if err != nil {
			log.Error("commands", "runDel", err, "Completed")
			return
		}

		return
	}

	if del.email != "" {
		log.Dev("commands", "runDel", "Email[%s]", del.email)
		user, err := db.GetUserByEmail(del.email)
		if err != nil {
			log.Error("commands", "runDel", err, "Completed")
			return
		}

		err = db.Delete(user)
		if err != nil {
			log.Error("commands", "runDel", err, "Completed")
			return
		}

		return
	}

	log.Dev("commands", "runDel", "Name[%s]", del.name)
	user, err := db.GetUserByName(del.name)
	if err != nil {
		log.Error("commands", "runDel", err, "Completed")
		return
	}

	err = db.Delete(user)
	if err != nil {
		log.Error("commands", "runDel", err, "Completed")
		return
	}

	return
}
