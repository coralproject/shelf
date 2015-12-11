package cmdquery

import (
	"github.com/coralproject/shelf/pkg/query"

	"github.com/ardanlabs/kit/db"

	"github.com/spf13/cobra"
)

var deleteLong = `Removes a query from the system using the query name.

Example:
	query delete -n user_advice
`

// delete contains the state for this command.
var delete struct {
	name string
}

// addDel handles the retrival query records, displayed in json formatted response.
func addDel() {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Removes a query record by name.",
		Long:  deleteLong,
		Run:   runDelete,
	}

	cmd.Flags().StringVarP(&delete.name, "name", "n", "", "Name of the user record.")

	queryCmd.AddCommand(cmd)
}

// runDelete is the code that implements the delete command.
func runDelete(cmd *cobra.Command, args []string) {
	cmd.Printf("Deleting Query : Path[%s]\n", upsert.path)

	if delete.name == "" {
		cmd.Help()
		return
	}

	db := db.NewMGO()
	defer db.CloseMGO()

	err := query.DeleteSet("", db, delete.name)
	if err != nil {
		cmd.Println("Deleting Query : ", err)
		return
	}

	cmd.Println("Deleting Query : Deleted")
}
