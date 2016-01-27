package cmdscript

import (
	"github.com/coralproject/xenia/pkg/script"

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

// runDelete is the code that implements the delete command.
func runDelete(cmd *cobra.Command, args []string) {
	cmd.Printf("Deleting Script : Name[%s]\n", delete.name)

	if delete.name == "" {
		cmd.Help()
		return
	}

	if err := script.Delete("", conn, delete.name); err != nil {
		cmd.Println("Deleting Script : ", err)
		return
	}

	cmd.Println("Deleting Script : Deleted")
}
