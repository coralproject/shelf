package cmditem

import (
	"github.com/coralproject/shelf/cmd/sponge/web"
	"github.com/spf13/cobra"
)

var deleteLong = `Removes an Item from the system using the Item id.

Example:
	item delete -i ID
`

// delete contains the state for this command.
var delete struct {
	ID string
}

// addDel handles the deletion of Item records.
func addDel() {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Removes an Item record by ID.",
		Long:  deleteLong,
		Run:   runDelete,
	}

	cmd.Flags().StringVarP(&delete.ID, "ID", "i", "", "Item ID.")

	itemCmd.AddCommand(cmd)
}

// runDelete issues the command talking to the web service.
func runDelete(cmd *cobra.Command, args []string) {
	verb := "DELETE"
	url := "/1.0/item/" + delete.ID

	if _, err := web.Request(cmd, verb, url, nil); err != nil {
		cmd.Println("Deleting Item : ", err)
	}

	cmd.Println("Deleting Item : Deleted")
}
