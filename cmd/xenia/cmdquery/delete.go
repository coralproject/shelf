package cmdquery

import (
	"github.com/coralproject/shelf/cmd/xenia/web"
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
		RunE:  runDelete,
	}

	cmd.Flags().StringVarP(&delete.name, "name", "n", "", "Name of the Set record.")

	queryCmd.AddCommand(cmd)
}

// runDelete issues the command talking to the web service.
func runDelete(cmd *cobra.Command, args []string) error {
	verb := "DELETE"
	url := "/v1/query/" + delete.name

	if _, err := web.Request(cmd, verb, url, nil); err != nil {
		return err
	}

	cmd.Println("Deleting Set : Deleted")
	return nil
}
