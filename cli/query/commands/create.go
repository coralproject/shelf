package commands

import (
	"github.com/coralproject/shelf/cli/query/db"
	"github.com/coralproject/shelf/pkg/db/mongo"
	"github.com/coralproject/shelf/pkg/log"
	"github.com/spf13/cobra"
)

var createLong = `Creates a new query into to the system.
When creating a new query, you need to supply the path to the file that contains the query
to be saved.

Note: To give the query a custom name other than the filename, supply a name field in the
json document else the name of the file will be used as the name of the query.

Example:

	query create -f ./queries/user_advice.json
`

// create contains the state for this command.
var create struct {
	file string
}

// addCreate handles the creation of users.
func addCreate() {
	cmd := &cobra.Command{
		Use:   "create [-f file]",
		Short: "Creates a new query from a file",
		Long:  createLong,
		Run:   runCreate,
	}

	cmd.Flags().StringVarP(&create.file, "file", "f", "", "file path of query json file")

	rootCmd.AddCommand(cmd)
}

// runCreate is the code that implements the create command.
func runCreate(cmd *cobra.Command, args []string) {
	if create.file == "" {
		cmd.Help()
		return
	}

	q, err := db.QueryFromFile(create.file)
	if err != nil {
		log.Error("commands", "runCreate", err, "Completed")
		return
	}

	// Initialize the mongodb session.
	mongo.InitMGO()

	err2 := db.Create(q)
	if err2 != nil {
		log.Error("commands", "runCreate", err, "Completed")
		return
	}
}
