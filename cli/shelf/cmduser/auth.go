package cmduser

import (
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
var auths struct {
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

	cmd.Flags().StringVarP(&auths.utype, "type", "t", "", "sets authentication type")
	cmd.Flags().StringVarP(&auths.key, "key", "k", "", "sets the key(email|publicId) of the user")
	cmd.Flags().StringVarP(&auths.pass, "pass", "p", "", "sets the Pass(password|token) of the user")

	userCmd.AddCommand(cmd)
}

// runAuth provides the operation logic for the auth command.
func runAuth(cmd *cobra.Command, args []string) {
}

// authPassword checks the password against the database.
func authPassword() {
}
