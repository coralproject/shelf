package cmdquery

import (
	"github.com/coralproject/xenia/pkg/query"

	"github.com/ardanlabs/kit/db"
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

// addDel handles the retrival set records, displayed in json formatted response.
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
	cmd.Printf("Deleting Set : Name[%s]\n", delete.name)

	if delete.name == "" {
		cmd.Help()
		return
	}

	db := db.NewMGO()
	defer db.CloseMGO()

	err := query.Delete("", db, delete.name)
	if err != nil {
		cmd.Println("Deleting Set : ", err)
		return
	}

	cmd.Println("Deleting Set : Deleted")
}
