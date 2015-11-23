package cmduser

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"github.com/astaxie/beego/validation"
	"github.com/coralproject/shelf/pkg/db/auth"
	"github.com/coralproject/shelf/pkg/db/mongo"
	"github.com/coralproject/shelf/pkg/log"
	"github.com/spf13/cobra"
)

var createLong = `Creates adds a new user to the system.
When creating a new user, the name(-n), email(-e) and password(-p) must all be supplied.

Example:

	user create -n "Alex Boulder" -e alex.boulder@gmail.com -p yefc*7fdf92
`

// create contains the state for this command.
var create struct {
	name   string
	pass   string
	email  string
	postal string
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
	cmd.Flags().StringVarP(&create.postal, "postal", "po", "", "postal code for user")
	cmd.Flags().StringVarP(&create.email, "email", "e", "", "email of user")

	userCmd.AddCommand(cmd)
}

// runCreate is the code that implements the create command.
func runCreate(cmd *cobra.Command, args []string) {
	if create.name == "" && create.email == "" && create.pass == "" {
		cmd.Help()
		return
	}

	if create.name == "" {
		cmd.Println("\n\tError: name(-n) can not be empty. please supply a name using the `-n` or `-name` flag")
		return
	}

	if create.pass == "" {
		cmd.Println("\n\tError: password(-p) can not be empty. Please supply a password using `-p` or `-password` flag")
		return
	}

	if create.email == "" {
		cmd.Println("\n\tError: email(-e) can not be empty. Please supply a email address using `-e` or `-email` flag")
		return
	}

	// Trying to match the complexity of email address is unecessary, as far as we
	// have a valid expectation pattern,we can skip alot of the mess.
	// TODO: should we use something more robust?
	if !strings.Contains(create.email, "@") {
		cmd.Println("\n\tError: Email address must be a valid addresss. Please supply a correct email address.")
		return
	}

	log.Dev("commands", "runCreate", "Email[%s]", create.email)

	newUser := &auth.NewUser{
		FullName:        create.name,
		Email:           create.email,
		Password:        create.pass,
		PasswordConfirm: create.pass,
		PostalCode:      create.postal,
	}

	validator := new(validation.Validation)
	newUser.Validate("commands", validator)

	if validator.HasErrors() {
		errbuf := bytes.NewBuffer([]byte{})
		for _, errs := range validator.Errors {
			errbuf.Write([]byte(fmt.Sprintf("%s: %s\n", errs.Key, errs.Message)))
		}

		log.Error("commands", "runCreate", errors.New("Invalid User Credentails"), "Completed")
		fmt.Println(errbuf.String())
	}

	user, err := newUser.Create("commands")
	if err != nil {
		log.Error("commands", "runCreate", err, "Completed")
		return
	}

	// Initialize the mongodb session.
	mongo.InitMGO()

	session := mongo.GetSession()
	defer session.Close()

	err2 := auth.Create("commands", session, user)
	if err2 != nil {
		log.Error("commands", "runCreate", err, "Completed")
		return
	}
}
