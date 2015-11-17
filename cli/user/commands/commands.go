package commands

import (
	"encoding/json"
	"strings"

	"github.com/coralproject/shelf/cli/user/db"
	"github.com/coralproject/shelf/log"
	"github.com/spf13/cobra"
)

// rootCmd represents the parent for all cli commands.
var rootCmd = &cobra.Command{
	Use:   "user",
	Short: "user provides a shelf CLI for managing user records.",
}

// Run will process the command line arguments for the application.
func Run() {
	cmdCreate()
	cmdGet()
	cmdUpdate()
	rootCmd.Execute()
}

// cmdCreate handles the creation of users.
func cmdCreate() {
	var name string
	var pass string
	var email string

	cmdCreate := &cobra.Command{
		Use:   "create [-n name -p password -e email]",
		Short: "Creates a new user",
		Long: `Creates adds a new user to the system.
When creating a new user, the name(-n), email(-e) and password(-p) must all be supplied.

Example:

	user create -n "Alex Boulder" -e alex.boulder@gmail.com -p yefc*7fdf92
`,
		Run: func(cmd *cobra.Command, args []string) {
			if name == "" && email == "" && pass == "" {
				cmd.Help()
				return
			}

			if name == "" {
				cmd.Println("\n\tname(-n) can not be empty. please supply a name using the `-n` or `-name` flag")
				return
			}
			if pass == "" {
				cmd.Println("\n\tpassword(-p) can not be empty. Please supply a password using `-p` or `-password` flag")
				return
			}

			if email == "" {
				cmd.Println("\n\temail(-e) can not be empty. Please supply a email address using `-e` or `-email` flag")
				return
			}

			// Trying to match the complexity of email address is unecessary, as far as we
			// have a valid expectation pattern,we can skip alot of the mess.
			// TODO: should we use something more robust?
			if !strings.Contains(email, "@") {
				cmd.Println("\n\tEmail address must be a valid addresss. Please supply a correct email address.")
				return
			}

			log.Dev("CreateUser", "commands.Create", "Started : User : Create : Email %q", email)
			user, err := db.NewUser(name, email, pass)
			if err != nil {
				log.Dev("CreateUser", "commands.Create", "Completed Error : User : Create : Email %q : Error %s", email, err.Error())
				return
			}

			// Initialize the mongodb session.
			db.InitMGO()

			err2 := db.Create(user)
			if err2 != nil {
				log.Dev("CreateUser", "commands.Create", "Completed Error : User : Create : Email %q : Error %s", email, err.Error())
				return
			}

			log.Dev("CreateUser", "commands.Create", "Completed : User : Create : Email %q : Success", email)
		},
	}

	cmdCreate.Flags().StringVarP(&name, "name", "n", "", "name of user")
	cmdCreate.Flags().StringVarP(&pass, "pass", "p", "", "password for user")
	cmdCreate.Flags().StringVarP(&email, "email", "e", "", "email of user")

	rootCmd.AddCommand(cmdCreate)
}

// cmdGet handles the retrival users records, displayed in json formatted response.
func cmdGet() {
	var name string
	var pid string
	var email string

	cmdGet := &cobra.Command{
		Use:   "get [-n name -p public_id -e email]",
		Short: "Retrieves a user record",
		Long: `Retrieves a user record from the system using any of the supplied keys.
When retrieving a record from the system, atleast one of the key points must be supplied, that is either the email(-e), the name(-n) or public id(-p). Each flag is checked in the order of important, where

	'public_id' is first importance.
	'email' is second in importance.
	'name' is the third and last in importance.

If all are supplied, the highest flag with the most importance gets used.

Example:

	1. To get a user using it's name:

		user get -n "Alex Boulder"

	2. To get a user using it's email address:

		user get -e alex.boulder@gmail.com

	3. To get a user using it's public id number:

		user get -p 199550d7-484d-4440-801f-390d44911ade
`,
		Run: func(cmd *cobra.Command, args []string) {
			if name == "" && pid == "" && email == "" {
				cmd.Println("\n\tAtleast one key flag must be supplied. Use either -name(-n), -email(-e) or -public_id(-p) to retrieve the desired record.")
			}

			if email != "" {
				// Trying to match the complexity of email address is unecessary, as far as we
				// have a valid expectation pattern,we can skip alot of the mess.
				// TODO: should we use something more robust?
				if !strings.Contains(email, "@") {
					cmd.Println("\n\tEmail address must be a valid addresss. Please supply a correct email address.")
					return
				}
			}

			// Initialize the mongodb session.
			db.InitMGO()

			if pid != "" {
				log.Dev("GetUser", "commands.Get", "Started : User : Get : Pid %q", pid)
				user, err := db.GetUserByPublicID(pid)
				if err != nil {
					log.Dev("GetUser", "commands.Get", "Completed Error : User : Get : Pid %q : Error %q", pid, err.Error())
					return
				}

				jsonFormatted, err := json.MarshalIndent(user, "", "\n\n")
				if err != nil {
					log.Dev("GetUser", "commands.Get", "Completed Error : User : Get : Pid %q : Error %q", pid, err.Error())
					return
				}

				log.Dev("GetUser", "commands.Get", "Completed : User : Get : Pid %q : Success : <User \n%s\nUser>", pid, jsonFormatted)
				return
			}

			if email != "" {
				log.Dev("GetUser", "commands.Get", "Started : User : Get : Email %q", email)
				user, err := db.GetUserByEmail(email)
				if err != nil {
					log.Dev("GetUser", "commands.Get", "Completed Error : User : Get : Email %q : Error %q", email, err.Error())
					return
				}

				jsonFormatted, err := json.MarshalIndent(user, "", "\n\n")
				if err != nil {
					log.Dev("GetUser", "commands.Get", "Completed Error : User : Get : Email %q : Error %q", email, err.Error())
					return
				}

				log.Dev("GetUser", "commands.Get", "Completed : User : Get : Pid %q : Success : <User \n%s\nUser>", pid, jsonFormatted)
				return
			}

			log.Dev("GetUser", "commands.Get", "Started : User : Get : Name %q", name)
			user, err := db.GetUserByName(name)
			if err != nil {
				log.Dev("GetUser", "commands.Get", "Completed Error : User : Get : Name %q : Error %q", name, err.Error())
				return
			}

			jsonFormatted, err := json.MarshalIndent(user, "", "\n\n")
			if err != nil {
				log.Dev("GetUser", "commands.Get", "Completed Error : User : Get : Name %q : Error %q", name, err.Error())
				return
			}

			log.Dev("GetUser", "commands.Get", "Completed : User : Get : Pid %q : Success : <User \n%s\nUser>", pid, jsonFormatted)
			return
		},
	}

	cmdGet.Flags().StringVarP(&name, "name", "n", "", "name of the user record")
	cmdGet.Flags().StringVarP(&pid, "public_id", "p", "", "publicId of the user record")
	cmdGet.Flags().StringVarP(&email, "email", "e", "", "email of the user")

	rootCmd.AddCommand(cmdGet)
}

// cmdCreate handles the creation of users.
func cmdUpdate() {
	var utype string
	var email string
	var oldValue string
	var newValue string

	cmdUpdate := &cobra.Command{
		Use:   "update -t auth|email|name [-e email ....extraOptions ]",
		Short: "Updates a existing user information",
		Long: `Updates a existing user's information in the system.
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
`,
		Run: func(cmd *cobra.Command, args []string) {
			if utype == "" {
				cmd.Println("\n\tname(-n) can not be empty. please supply a name using the `-n` or `-name` flag")
			}

			// Initialize the mongodb session.
			db.InitMGO()
		},
	}

	cmdUpdate.Flags().StringVarP(&utype, "type", "t", "", "type of update")
	cmdUpdate.Flags().StringVarP(&email, "email", "e", "", "email of user")
	cmdUpdate.Flags().StringVarP(&oldValue, "old", "o", "", "old value of user")
	cmdUpdate.Flags().StringVarP(&newValue, "new", "n", "", "new value for user")

	rootCmd.AddCommand(cmdUpdate)
}

// cmdDelete handles the deletion users records.
func cmdDelete() {
	var name string
	var pid string
	var email string

	cmdDelete := &cobra.Command{
		Use:   "delete [-n name -p public_id -e email]",
		Short: "Deletes a user record",
		Long: `Deletes a user record from the system using any of the supplied keys.
When deleting a record from the system, atleast one of the key points must be supplied, that is either the email(-e), the name(-n) or public id(-p). Each flag is checked in the order of important, where

	'public_id' is first importance.
	'email' is second in importance.
	'name' is the third and last in importance.

If all are supplied, the highest flag with the most importance gets used.

Example:

	1. To delete a user using it's name:

		user get -n "Alex Boulder"

	2. To delete a user using it's email address:

		user get -e alex.boulder@gmail.com

	3. To delete a user using it's public id number:

		user get -p 199550d7-484d-4440-801f-390d44911ade
`,
		Run: func(cmd *cobra.Command, args []string) {
			if name == "" && pid == "" && email == "" {
				cmd.Println("\n\tAtleast one key flag must be supplied. Use either -name(-n), -email(-e) or -public_id(-p) to delete the desired record.")
			}

			if email != "" {
				// Trying to match the complexity of email address is unecessary, as far as we
				// have a valid expectation pattern,we can skip alot of the mess.
				// TODO: should we use something more robust?
				if !strings.Contains(email, "@") {
					cmd.Println("\n\tEmail address must be a valid addresss. Please supply a correct email address.")
					return
				}
			}

			// Initialize the mongodb session.
			db.InitMGO()

			if pid != "" {
				log.Dev("DeleteUser", "commands.Delete", "Started : User : Delete : Pid %q", pid)
				user, err := db.GetUserByPublicID(pid)
				if err != nil {
					log.Dev("DeleteUser", "commands.Delete", "Completed Error : User : Get : Pid %q : Error %q", pid, err.Error())
					return
				}

				err = db.Delete(user)
				if err != nil {
					log.Dev("DeleteUser", "commands.Delete", "Completed Error : User : Get : Pid %q : Error %q", pid, err.Error())
					return
				}

				log.Dev("DeleteUser", "commands.Delete", `Completed : User : Get : Name %q : Success`, name)
				return
			}

			if email != "" {
				log.Dev("DeleteUser", "commands.Delete", "Started : User : Get : Email %q", email)
				user, err := db.GetUserByEmail(email)
				if err != nil {
					log.Dev("DeleteUser", "commands.Delete", "Completed Error : User : Get : Email %q : Error %q", email, err.Error())
					return
				}

				err = db.Delete(user)
				if err != nil {
					log.Dev("DeleteUser", "commands.Delete", "Completed Error : User : Get : Email %q : Error %q", email, err.Error())
					return
				}

				log.Dev("DeleteUser", "commands.Delete", `Completed : User : Get : Name %q : Success`, name)
				return
			}

			log.Dev("DeleteUser", "commands.Delete", "Started : User : Get : Name %q", name)
			user, err := db.GetUserByName(name)
			if err != nil {
				log.Dev("DeleteUser", "commands.Delete", "Completed Error : User : Get : Name %q : Error %q", name, err.Error())
				return
			}

			err = db.Delete(user)
			if err != nil {
				log.Dev("DeleteUser", "commands.Delete", "Completed Error : User : Get : Name %q : Error %q", name, err.Error())
				return
			}

			log.Dev("DeleteUser", "commands.Delete", `Completed : User : Get : Name %q : Success`, name)
			return
		},
	}

	cmdDelete.Flags().StringVarP(&name, "name", "n", "", "name of the user record")
	cmdDelete.Flags().StringVarP(&pid, "public_id", "p", "", "publicId of the user record")
	cmdDelete.Flags().StringVarP(&email, "email", "e", "", "email of the user")

	rootCmd.AddCommand(cmdDelete)
}
