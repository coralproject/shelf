package cmdquery

import (
	"github.com/coralproject/shelf/pkg/db"
	"github.com/coralproject/shelf/pkg/log"
	"github.com/coralproject/shelf/pkg/srv/query"

	"github.com/spf13/cobra"
)

var deleteLong = `Removes a query record from the system using the supplied name.
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
		Use:   "delete [-n name]",
		Short: "Removes a query record",
		Long:  deleteLong,
		Run:   runDelete,
	}

	cmd.Flags().StringVarP(&delete.name, "name", "n", "", "name of the user record")

	queryCmd.AddCommand(cmd)
}

// runDelete is the code that implements the delete command.
func runDelete(cmd *cobra.Command, args []string) {
	if delete.name == "" {
		cmd.Help()
		return
	}

	db := db.NewMGO()
	defer db.CloseMGO()

	err := query.DeleteSet("commands", db, delete.name)
	if err != nil {
		log.Error("commands", "runGet", err, "Completed")
		return
	}

	return
}
