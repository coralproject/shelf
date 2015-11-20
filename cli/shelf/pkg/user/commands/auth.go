package commands

import (
	"errors"
	"strings"

	"github.com/coralproject/shelf/cli/shelf/pkg/user/db"
	"github.com/coralproject/shelf/pkg/db/mongo"
	"github.com/coralproject/shelf/pkg/log"
	"github.com/spf13/cobra"
)

// authLong provides a detailed description on the auth subcommand.
var authLong = `Authenticates the given user crendentails.
When authenticating a credentail using the CLI, its required to first specify the authetication type.
The type(-t : token|password) determines the expecting authentication credentails expected.
The credentails are passed into the '-key' and '-pass' flags.

 Examples:

  1. To authenticate using the user's Public Id and Token,set the type to 'token':

		 user auth -t token -k {User PublicID} -p {User Token}

  2. To authenticate using the user's Email and Password, set the type to 'pass':

		 user auth -t password -k shou.lou@gmail.com -p Shen5A43*2f3e

`

// auth contains the state for this command.
var auth struct {
	utype string
	key   string
	pass  string
}

// addAuth handles the authentication of user credentails.
func addAuth() {
	cmd := &cobra.Command{
		Use:   "auth",
		Short: "Authenticates user credentails",
		Long:  authLong,
		Run:   runAuth,
	}

	cmd.Flags().StringVarP(&auth.utype, "type", "t", "", "sets authentication type")
	cmd.Flags().StringVarP(&auth.key, "key", "k", "", "sets the key(email|publicId) of the user")
	cmd.Flags().StringVarP(&auth.pass, "pass", "p", "", "sets the Pass(password|token) of the user")

	rootCmd.AddCommand(cmd)
}

// runAuth provides the operation logic for the auth command.
func runAuth(cmd *cobra.Command, args []string) {
	if auth.utype == "" && auth.key == "" && auth.pass == "" {
		cmd.Help()
		return
	}

	if auth.utype == "" {
		cmd.Println("\n\tError: type(-t) can not be empty. Please supply a name using the `-t` or `-type` flag")
		cmd.Help()
		return
	}

	if auth.utype == "token" {
		if auth.key == "" {
			cmd.Println("\n\tError: key(-k) can not be empty. Please supply the key(Public Id) using the `-k` or `-key` flag")
			return
		}

		if auth.pass == "" {
			cmd.Println("\n\tError: pass(-p) can not be empty. Please supply the pass(Token) using the `-p` or `-pass` flag")
			return
		}
	}

	if auth.utype == "password" {
		if auth.key == "" {
			cmd.Println("\n\tError: key(-k) can not be empty. Please supply the key(Email) using the `-k` or `-key` flag")
			return
		}

		if auth.key != "" {
			// Trying to match the complexity of email address is unecessary, as far as we
			// have a valid expectation pattern,we can skip alot of the mess.
			// TODO: should we use something more robust?
			if !strings.Contains(auth.key, "@") {
				cmd.Println("\n\tError: Email address must be a valid addresss. Please supply a correct email address.")
				return
			}
		}

		if auth.pass == "" {
			cmd.Println("\n\tError: pass(-p) can not be empty. Please supply the pass(Password) using the `-p` or `-pass` flag")
			return
		}
	}

	// Initialize the mongodb session.
	mongo.InitMGO()

	switch auth.utype {
	case "auth":
		authToken()
	case "password":
		authPassword()
	}
}

// authToken checks the token against the database.
func authToken() {
	log.Dev("commands", "runAuth", "Public Id[%s]", auth.key)

	user, err := db.GetUserByPublicID(auth.key)
	if err != nil {
		log.Error("commands", "runAuth", err, "Completed")
		return
	}

	if !user.IsTokenValid(auth.pass) {
		log.Error("commands", "runAuth", errors.New("Invalid Token"), "Completed")
		return
	}

	log.Dev("commands", "runAuth", "Auth: 200 Ok!")
}

// authPassword checks the password against the database.
func authPassword() {
	log.Dev("commands", "runAuth", "Email[%s]", auth.key)

	user, err := db.GetUserByEmail(auth.key)
	if err != nil {
		log.Error("commands", "runAuth", err, "Completed")
		return
	}

	if !user.IsPasswordValid(auth.pass) {
		log.Error("commands", "runAuth", errors.New("Invalid Password"), "Completed")
		return
	}

	log.Dev("commands", "runAuth", "Auth: 200 Ok!")
}
