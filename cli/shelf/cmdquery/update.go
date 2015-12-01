package cmdquery

import (
	"github.com/coralproject/shelf/pkg/log"
	"github.com/coralproject/shelf/pkg/srv/mongo"
	"github.com/coralproject/shelf/pkg/srv/query"
	"github.com/spf13/cobra"
	"gopkg.in/mgo.v2"
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
	dir  string
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
	cmd.Flags().StringVarP(&update.dir, "dir", "d", "", "dir to load updated json files from")

	queryCmd.AddCommand(cmd)
}

// runUpdate is the code that implements the create command.
func runUpdate(cmd *cobra.Command, args []string) {
	if update.file == "" && update.dir == "" {
		cmd.Help()
		return
	}

	session := mongo.GetSession()
	defer session.Close()

	// If the file option is not an empty string, then
	// try to load file path.
	if update.file != "" {
		if err := loadUpdate(update.file, session); err != nil {
			log.Error("commands", "runUpdate", err, "Completed")
			return
		}
	}

	// Attempt to load directory path.
	if err := loadDir(update.dir, session, loadUpdate); err != nil {
		log.Error("commands", "runUpdate", err, "Completed")
		return
	}
}

// loadUpdate loads a given file path and attempts to update the query db.
func loadUpdate(file string, ses *mgo.Session) error {
	q, err := setFromFile("commands", file)
	if err != nil {
		return err
	}

	return query.UpdateSet("commands", ses, q)
}
