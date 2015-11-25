package cmdquery

import (
	"github.com/coralproject/shelf/pkg/log"
	"github.com/coralproject/shelf/pkg/srv/mongo"
	"github.com/coralproject/shelf/pkg/srv/query"
	"github.com/spf13/cobra"
)

var deleteLong = `Removes a query record from the system using the supplied name.
Example:

		user delete -n user_advice

`

// delete contains the state for this command.
var delete struct {
	name string
}

// addDel handles the retrival users records, displayed in json formatted response.
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

	// Initialize the mongodb session.
	mongo.InitMGO()

	session := mongo.GetSession()
	defer session.Close()

	err := query.DeleteSet("commands", session, delete.name)
	if err != nil {
		log.Error("commands", "runGet", err, "Completed")
		return
	}

	return
}
