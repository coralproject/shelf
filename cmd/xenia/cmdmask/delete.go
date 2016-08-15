package cmdmask

import (
	"github.com/coralproject/shelf/cmd/xenia/web"
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

// runDelete issues the command talking to the web service.
func runDelete(cmd *cobra.Command, args []string) {
	verb := "DELETE"
	url := "/1.0/mask/" + delete.collection + "/" + delete.field

	if _, err := web.Request(cmd, verb, url, nil); err != nil {
		cmd.Println("Deleting Mask : ", err)
	}

	cmd.Println("Deleting Mask : Deleted")
}
