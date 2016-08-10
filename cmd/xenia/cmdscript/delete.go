package cmdscript

import (
	"github.com/coralproject/xenia/cmd/xenia/web"
	"github.com/coralproject/xenia/internal/xenia/script"
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
	if conn == nil {
		runDeleteWeb(cmd)
		return
	}

	runDeleteDB(cmd)
}

// runDeleteWeb issues the command talking to the web service.
func runDeleteWeb(cmd *cobra.Command) {
	verb := "DELETE"
	url := "/1.0/script/" + get.name

	if _, err := web.Request(cmd, verb, url, nil); err != nil {
		cmd.Println("Deleting Script : ", err)
	}

	cmd.Println("Deleting Script : Deleted")
}

// runDeleteDB issues the command talking to the DB.
func runDeleteDB(cmd *cobra.Command) {
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
