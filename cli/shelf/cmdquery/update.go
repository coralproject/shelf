package cmdquery

import (
	"github.com/coralproject/shelf/pkg/log"
	"github.com/coralproject/shelf/pkg/srv/mongo"
	"github.com/coralproject/shelf/pkg/srv/query"
	"github.com/spf13/cobra"
)

var updateLong = `Updates a query in the system using the giving file and  name.

Note: Regardless of name in the json file,the name of the record remains intact as it
was created in the system

Example:

	query update -f ./queries/user_advice.json
`

// update contains the state for this command.
var update struct {
	file string
	// name string
}

// addUpd handles the update of query record.
func addUpd() {
	cmd := &cobra.Command{
		Use:   "update [-name -f file]",
		Short: "Updates a query in the system from a file",
		Long:  updateLong,
		Run:   runUpdate,
	}

	// cmd.Flags().StringVarP(&update.name, "name", "n", "", "name of query record")
	cmd.Flags().StringVarP(&update.file, "file", "f", "", "file path of query json file")

	queryCmd.AddCommand(cmd)
}

// runUpdate is the code that implements the create command.
func runUpdate(cmd *cobra.Command, args []string) {
	if update.file == "" {
		cmd.Help()
		return
	}

	// Initialize the mongodb session.
	mongo.InitMGO()

	session := mongo.GetSession()
	defer session.Close()

	// // check if the record exists with the giving name
	// _, err := query.Get("commands", session, update.name)
	// if err != nil {
	// 	log.Error("commands", "runUpdate", err, "Completed")
	// 	return
	// }

	q, err := setFromFile("commands", update.file)
	if err != nil {
		log.Error("commands", "runUpdate", err, "Completed")
		return
	}

	// persist the name of the record
	// q.Name = update.name

	err2 := query.UpdateSet("commands", session, *(q))
	if err2 != nil {
		log.Error("commands", "runUpdate", err, "Completed")
		return
	}
}
