package cmduser

import (
	"github.com/coralproject/shelf/pkg/db"
	"github.com/coralproject/shelf/pkg/srv/auth"

	"github.com/spf13/cobra"
)

var createLong = `Use create to add a new user to the system. The user email
must be unique for every user.

Example:
  ./shelf user create -n "Bill Kennedy" -e "bill@ardanlabs.com" -p "yefc*7fdf92"
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
		Use:   "create",
		Short: "Add a new user to the system.",
		Long:  createLong,
		Run:   runCreate,
	}

	cmd.Flags().StringVarP(&create.name, "name", "n", "", "Full name of the user")
	cmd.Flags().StringVarP(&create.email, "email", "e", "", "Email for the user")
	cmd.Flags().StringVarP(&create.pass, "pass", "p", "", "Password for the user")

	userCmd.AddCommand(cmd)
}

// runCreate is the code that implements the create command.
func runCreate(cmd *cobra.Command, args []string) {
	cmd.Printf("Creating User : Name[%s] Email[%s] Pass[%s]\n", create.name, create.email, create.pass)

	u, err := auth.NewUser(auth.NUser{
		Status:   auth.StatusActive,
		FullName: create.name,
		Email:    create.email,
		Password: create.pass,
	})
	if err != nil {
		cmd.Println("Creating User : ", err)
		cmd.Help()
		return
	}

	db := db.NewMGO()
	defer db.CloseMGO()

	if err := auth.CreateUser("", db, u); err != nil {
		cmd.Println("Creating User : ", err)
		return
	}

	cmd.Println("Creating User : Created")
}
