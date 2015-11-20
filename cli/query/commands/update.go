package commands

import (
	"github.com/coralproject/shelf/cli/query/db"
	"github.com/coralproject/shelf/pkg/log"
	"github.com/coralproject/shelf/pkg/mongo"
	"github.com/spf13/cobra"
)

var updateLong = `Updates a query in the system using the giving file and  name.

Note: Regardless of name in the json file,the name of the record remains intact as it
was created in the system

Example:

	query update -n user_advice -f ./queries/user_advice.json
`

// update contains the state for this command.
var update struct {
	file string
	name string
}

// addUpd handles the update of query record.
func addUpd() {
	cmd := &cobra.Command{
		Use:   "update [-name -f file]",
		Short: "Updates a query in the system from a file",
		Long:  updateLong,
		Run:   runUpdate,
	}

	cmd.Flags().StringVarP(&update.name, "name", "n", "", "name of query record")
	cmd.Flags().StringVarP(&update.file, "file", "f", "", "file path of query json file")

	rootCmd.AddCommand(cmd)
}

// runUpdate is the code that implements the create command.
func runUpdate(cmd *cobra.Command, args []string) {
	if update.name == "" {
		cmd.Help()
		return
	}

	q, err := db.GetByName(update.name)
	if err != nil {
		log.Error("commands", "runCreate", err, "Completed")
		return
	}

	if err := q.LoadFile(update.file); err != nil {
		log.Error("commands", "runCreate", err, "Completed")
		return
	}

	// Safeguard the name of the record.
	q.Name = update.name

	// Initialize the mongodb session.
	mongo.InitMGO()

	err2 := db.Update(q)
	if err2 != nil {
		log.Error("commands", "runCreate", err, "Completed")
		return
	}
}
