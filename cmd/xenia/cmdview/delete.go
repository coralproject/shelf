package cmdview

import (
	"github.com/coralproject/shelf/cmd/xenia/web"
	"github.com/spf13/cobra"
)

var deleteLong = `Removes a View from the system using the View name.

Example:
	view delete -n name
`

// delete contains the state for this command.
var delete struct {
	name string
}

// addDel handles the deletion of View records.
func addDel() {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Removes a View record by name.",
		Long:  deleteLong,
		Run:   runDelete,
	}

	cmd.Flags().StringVarP(&delete.name, "name", "n", "", "View name.")

	viewCmd.AddCommand(cmd)
}

// runDelete issues the command talking to the web service.
func runDelete(cmd *cobra.Command, args []string) {
	verb := "DELETE"
	url := "/v1/view/" + delete.name

	if _, err := web.Request(cmd, verb, url, nil); err != nil {
		cmd.Println("Deleting View : ", err)
	}

	cmd.Println("Deleting View : Deleted")
}
