package cmdmask

import (
	"github.com/coralproject/xenia/cmd/xenia/web"
	"github.com/coralproject/xenia/internal/mask"

	"github.com/spf13/cobra"
)

var deleteLong = `Removes a Mask from the system using the Mask collection/field name.

Example:
	mask delete -c * -f test
`

// delete contains the state for this command.
var delete struct {
	collection string
	field      string
}

// addDel handles the retrival Script records, displayed in json formatted response.
func addDel() {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Removes a Mask record by collection/field.",
		Long:  deleteLong,
		Run:   runDelete,
	}

	cmd.Flags().StringVarP(&get.collection, "collection", "c", "", "Name of the Collection.")
	cmd.Flags().StringVarP(&get.field, "field", "f", "", "Name of the Field.")

	maskCmd.AddCommand(cmd)
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
	url := "/1.0/mask/" + delete.collection + "/" + delete.field

	if _, err := web.Request(cmd, verb, url, nil); err != nil {
		cmd.Println("Deleting Mask : ", err)
	}

	cmd.Println("Deleting Mask : Deleted")
}

// runDeleteDB issues the command talking to the DB.
func runDeleteDB(cmd *cobra.Command) {
	cmd.Printf("Deleting Mask : Collection[%s] Field[%s]\n", delete.collection, delete.field)

	if delete.collection == "" || delete.field == "" {
		cmd.Help()
		return
	}

	if err := mask.Delete("", conn, delete.collection, delete.field); err != nil {
		cmd.Println("Deleting Mask : ", err)
		return
	}

	cmd.Println("Deleting Mask : Deleted")
}
