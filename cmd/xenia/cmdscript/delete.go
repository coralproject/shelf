package cmdscript

import (
	"github.com/coralproject/xenia/cmd/xenia/web"
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
		Run:   runDelete,
	}

	cmd.Flags().StringVarP(&delete.name, "name", "n", "", "Name of the Script record.")

	scriptCmd.AddCommand(cmd)
}

// runDelete issues the command talking to the web service.
func runDelete(cmd *cobra.Command, args []string) {
	verb := "DELETE"
	url := "/1.0/script/" + get.name

	if _, err := web.Request(cmd, verb, url, nil); err != nil {
		cmd.Println("Deleting Script : ", err)
	}

	cmd.Println("Deleting Script : Deleted")
}
