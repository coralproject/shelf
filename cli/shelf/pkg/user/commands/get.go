package commands

import (
	"fmt"
	"strings"

	"github.com/coralproject/shelf/cli/shelf/pkg/user/db"
	"github.com/coralproject/shelf/pkg/db/mongo"
	"github.com/coralproject/shelf/pkg/log"
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

	rootCmd.AddCommand(cmd)
}

// runGet is the code that implements the get command.
func runGet(cmd *cobra.Command, args []string) {
	if get.name == "" && get.pid == "" && get.email == "" {
		// cmd.Println("\n\tError: Atleast one key flag must be supplied. Use either -name(-n), -email(-e) or -public_id(-p) to retrieve the desired record.")
		cmd.Help()
		return
	}

	if get.email != "" {
		// Trying to match the complexity of email address is unecessary, as far as we
		// have a valid expectation pattern,we can skip alot of the mess.
		// TODO: should we use something more robust?
		if !strings.Contains(get.email, "@") {
			cmd.Println("\n\tError: Email address must be a valid addresss. Please supply a correct email address.")
			return
		}
	}

	// Initialize the mongodb session.
	mongo.InitMGO()

	if get.pid != "" {
		log.Dev("commands", "runGet", "Pid[%s]", get.pid)
		user, err := db.GetUserByPublicID(get.pid)
		if err != nil {
			log.Error("commands", "runGet", err, "Completed")
			return
		}

		fmt.Printf(`
Record for User(%s):
	 Name: %s
	 Email: %s
	 Token: %s
	 PublicID: %s
	 PrivateID: %s
	 Record Creation Date: %s
	 Modified At: %s
`, get.pid, user.Name, user.Email, user.Token, user.PublicID, user.PrivateID, user.CreatedAt.String(), user.ModifiedAt.String())

		return
	}

	if get.email != "" {
		log.Dev("commands", "runGet", "Email[%s]", get.email)
		user, err := db.GetUserByEmail(get.email)
		if err != nil {
			log.Error("commands", "runGet", err, "Completed")
			return
		}

		fmt.Printf(`
Record for User(%s):
	 Name: %s
	 Email: %s
	 Token: %s
	 PublicID: %s
	 PrivateID: %s
	 Record Creation Date: %s
	 Modified At: %s
`, get.email, user.Name, user.Email, user.Token, user.PublicID, user.PrivateID, user.CreatedAt.String(), user.ModifiedAt.String())

		return
	}

	log.Dev("commands", "runGet", "Name[%s]", get.name)
	user, err := db.GetUserByName(get.name)
	if err != nil {
		log.Error("commands", "runGet", err, "Completed")
		return
	}

	// _, err = json.MarshalIndent(user, "", "\n")
	// if err != nil {
	// 	log.Error("GetUser", "runGet", err, "Completed")
	// 	return
	// }

	// TODO: What are you doing with doc
	fmt.Printf(`
Record for User(%s):
	 Name: %s
	 Email: %s
	 Token: %s
	 PublicID: %s
	 PrivateID: %s
	 Record Creation Date: %s
	 Modified At: %s
`, get.name, user.Name, user.Email, user.Token, user.PublicID, user.PrivateID, user.CreatedAt.String(), user.ModifiedAt.String())

	return
}
