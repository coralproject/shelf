package cmdquery

import (
	"github.com/coralproject/xenia/cmd/xenia/web"
	"github.com/coralproject/xenia/internal/query"

	"github.com/spf13/cobra"
)

var deleteLong = `Removes a set from the system using the set name.

Example:
	query delete -n user_advice
`

// delete contains the state for this command.
var delete struct {
	name string
}

// addDel handles the removal of a set document.
func addDel() {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Removes a Set record by name.",
		Long:  deleteLong,
		Run:   runDelete,
	}

	cmd.Flags().StringVarP(&delete.name, "name", "n", "", "Name of the Set record.")

	queryCmd.AddCommand(cmd)
}

// runDelete is the code that implements the delete command.
func runDelete(cmd *cobra.Command, args []string) {
	if conn == nil {
		runDeleteWeb(cmd)
		return
	}

	runDeleteDB(cmd)
}

// runDeleteWeb issues the command talking to the web service.
func runDeleteWeb(cmd *cobra.Command) {
	verb := "DELETE"
	url := "/1.0/query/" + get.name

	if _, err := web.Request(cmd, verb, url, nil); err != nil {
		cmd.Println("Deleting Set : ", err)
	}

	cmd.Println("Deleting Set : Deleted")
}

// runDeleteDB issues the command talking to the DB.
func runDeleteDB(cmd *cobra.Command) {
	cmd.Printf("Deleting Set : Name[%s]\n", delete.name)

	if delete.name == "" {
		cmd.Help()
		return
	}

	if err := query.Delete("", conn, delete.name); err != nil {
		cmd.Println("Deleting Set : ", err)
		return
	}

	cmd.Println("Deleting Set : Deleted")
}
