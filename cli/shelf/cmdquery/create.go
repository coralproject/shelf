package cmdquery

import (
	"github.com/coralproject/shelf/pkg/log"
	"github.com/coralproject/shelf/pkg/srv/mongo"
	"github.com/coralproject/shelf/pkg/srv/query"
	"github.com/spf13/cobra"
	"gopkg.in/mgo.v2"
)

var createLong = `Creates a new query into to the system.
When creating a new query, you need to supply the path to the file that contains the query
to be saved.

Note: To give the query a custom name other than the filename, supply a name field in the
json document else the name of the file will be used as the name of the query.

Example:

	1. To load a single file

		query create -f ./queries/user_advice.json

	2. To load a directory of rules
	 By default: if no path is giving, two things will occur:
		- It will check environment variables for a "SHELF_RULES_DIR"
	 	- It will default to a directory called "rules"

		query create -d ./{path_to_dir}
`

// create contains the state for this command.
var create struct {
	file string
	dir  string
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
	cmd.Flags().StringVarP(&create.dir, "dir", "d", "rules", "dir contain json files")

	queryCmd.AddCommand(cmd)
}

// runCreate is the code that implements the create command.
func runCreate(cmd *cobra.Command, args []string) {
	if create.file == "" && create.dir == "" {
		cmd.Help()
		return
	}

	session := mongo.GetSession()
	defer session.Close()

	// If the file option is not an empty string, then
	// try to load file path.
	if create.file != "" {
		if err := loadNew(create.file, session); err != nil {
			log.Error("commands", "runCreate", err, "Completed")
			return
		}
	}

	// Attempt to load directory path.
	if err := loadDir(create.dir, session, loadNew); err != nil {
		log.Error("commands", "runCreate", err, "Completed")
		return
	}

}

// loadNew loads a given file path and attempts to save into the query db.
func loadNew(file string, ses *mgo.Session) error {
	q, err := setFromFile("commands", file)
	if err != nil {
		return err
	}

	return query.CreateSet("commands", ses, q)
}
