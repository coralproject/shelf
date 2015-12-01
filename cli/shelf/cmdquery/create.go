package cmdquery

import (
	"os"
	"path/filepath"

	"github.com/coralproject/shelf/pkg/cfg"
	"github.com/coralproject/shelf/pkg/log"
	"github.com/coralproject/shelf/pkg/mongo"
	"github.com/coralproject/shelf/pkg/srv/query"
	"github.com/spf13/cobra"
)

var createLong = `Creates a new query into to the system by using a supplied
file/dir name, else falls back to using a path set in the environment
variable "SHELF_SCRIPT_DIR".

Example:

1. To load a single file

	query create -p user_advice.json

2. To load a directory of query scripts

	- It will check environment variables for a "SHELF_SCRIPT_DIR"

	query create -p ./{dir_name}
`

// create contains the state for this command.
var create struct {
	path string
}

// addCreate handles the creation of query records into the db.
func addCreate() {
	cmd := &cobra.Command{
		Use:   "create [-p filename/dirname]",
		Short: "Creates a new query from a file",
		Long:  createLong,
		Run:   runCreate,
	}

	cmd.Flags().StringVarP(&create.path, "path", "p", "", "path of query file or directory")

	queryCmd.AddCommand(cmd)
}

// runCreate is the code that implements the create command.
func runCreate(cmd *cobra.Command, args []string) {
	if create.path == "" {
		dir, err := cfg.String(envKey)
		if err != nil {
			create.path = defDir
		} else {
			create.path = filepath.Join(dir, defDir)
		}
	}

	session := mongo.GetSession()
	defer session.Close()

	pwd, err := os.Getwd()
	if err != nil {
		log.Error("commands", "runCreate", err, "Completed")
		return
	}

	file := filepath.Join(pwd, create.path)

	stat, err := os.Stat(file)
	if err != nil {
		log.Error("commands", "runCreate", err, "Completed")
		return
	}

	if !stat.IsDir() {
		q, err := setFromFile("commands", file)
		if err != nil {
			log.Error("commands", "runCreate", err, "Completed")
			return
		}

		if err := query.CreateSet("commands", session, q); err != nil {
			log.Error("commands", "runCreate", err, "Completed")
			return
		}

		return
	}

	err2 := loadDir(file, func(path string) error {
		q, err := setFromFile("commands", path)
		if err != nil {
			return err
		}

		return query.CreateSet("commands", session, q)
	})

	if err2 != nil {
		log.Error("commands", "runCreate", err, "Completed")
		return
	}

}
