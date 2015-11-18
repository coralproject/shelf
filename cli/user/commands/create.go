package commands

import (
	"strings"

	"github.com/coralproject/shelf/cli/user/db"
	"github.com/coralproject/shelf/log"
	"github.com/coralproject/shelf/mongo"
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

	rootCmd.AddCommand(cmd)
}

// runCreate is the code that implements the create command.
func runCreate(cmd *cobra.Command, args []string) {
	if create.name == "" && create.email == "" && create.pass == "" {
		cmd.Help()
		return
	}

	if create.name == "" {
		cmd.Println("\n\tname(-n) can not be empty. please supply a name using the `-n` or `-name` flag")
		return
	}

	if create.pass == "" {
		cmd.Println("\n\tpassword(-p) can not be empty. Please supply a password using `-p` or `-password` flag")
		return
	}

	if create.email == "" {
		cmd.Println("\n\temail(-e) can not be empty. Please supply a email address using `-e` or `-email` flag")
		return
	}

	// Trying to match the complexity of email address is unecessary, as far as we
	// have a valid expectation pattern,we can skip alot of the mess.
	// TODO: should we use something more robust?
	if !strings.Contains(create.email, "@") {
		cmd.Println("\n\tEmail address must be a valid addresss. Please supply a correct email address.")
		return
	}

	log.Dev("commands", "runCreate", "Email[%s]", create.email)
	user, err := db.NewUser(create.name, create.email, create.pass)
	if err != nil {
		log.Error("commands", "runCreate", err, "Completed")
		return
	}

	// Initialize the mongodb session.
	mongo.InitMGO()

	err2 := db.Create(user)
	if err2 != nil {
		log.Error("commands", "runCreate", err, "Completed")
		return
	}
}
