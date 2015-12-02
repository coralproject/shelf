package cmduser

import (
	"encoding/json"

	"github.com/coralproject/shelf/pkg/db"
	"github.com/coralproject/shelf/pkg/srv/auth"

	"github.com/spf13/cobra"
)

var getLong = `Use get to retreive a user record from the system.

Example:
  ./shelf user get -p "6dcda2da-92c3-11e5-8994-feff819cdc9f"

  ./shelf user get -e "bill@ardanlabs.com"
`

// get contains the state for this command.
var get struct {
	pid   string
	email string
}

// addGet handles the retrival users records, displayed in json formatted response.
func addGet() {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Retrieves a user record by public_id or email.",
		Long:  getLong,
		Run:   runGet,
	}

	cmd.Flags().StringVarP(&get.pid, "pid", "p", "", "Public Id of the user.")
	cmd.Flags().StringVarP(&get.email, "email", "e", "", "Email of the user.")

	userCmd.AddCommand(cmd)
}

// runGet is the code that implements the get command.
func runGet(cmd *cobra.Command, args []string) {
	cmd.Printf("Getting User : Pid[%s] Email[%s]\n", get.pid, get.email)

	if get.pid == "" && get.email == "" {
		cmd.Help()
		return
	}

	db := db.NewMGO()
	defer db.CloseMGO()

	var u *auth.User
	var err error

	if get.pid != "" {
		u, err = auth.GetUserByPublicID("", db, get.pid)
	} else {
		u, err = auth.GetUserByEmail("", db, get.email)
	}

	if err != nil {
		cmd.Println("Getting User : ", err)
		return
	}

	data, err := json.MarshalIndent(&u, "", "    ")
	if err != nil {
		cmd.Println("Getting User : ", err)
		return
	}

	cmd.Printf("\n%s\n\n", string(data))
	return
}
