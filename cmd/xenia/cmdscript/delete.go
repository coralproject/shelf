package cmdscript

import (
	"github.com/coralproject/shelf/cmd/xenia/web"
	"github.com/spf13/cobra"
)

var deleteLong = `Removes a Script from the system using the Script name.

Example:
	script delete -n user_advice
`

// delete contains the state for this command.
var delete struct {
	name string
}

// addDel handles the retrival Script records, displayed in json formatted response.
func addDel() {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Removes a Script record by name.",
		Long:  deleteLong,
		RunE:  runDelete,
	}

	cmd.Flags().StringVarP(&delete.name, "name", "n", "", "Name of the Script record.")

	scriptCmd.AddCommand(cmd)
}

// runDelete issues the command talking to the web service.
func runDelete(cmd *cobra.Command, args []string) error {
	verb := "DELETE"
	url := "/v1/script/" + get.name

	if _, err := web.Request(cmd, verb, url, nil); err != nil {
		return err
	}

	cmd.Println("Deleting Script : Deleted")
	return nil
}
