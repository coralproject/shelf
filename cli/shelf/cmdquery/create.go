package cmdquery

import (
	"os"
	"path/filepath"

	"github.com/coralproject/shelf/pkg/cfg"
	"github.com/coralproject/shelf/pkg/db"
	"github.com/coralproject/shelf/pkg/log"
	"github.com/coralproject/shelf/pkg/srv/query"

	"github.com/spf13/cobra"
)

var createLong = `Use create to add a new query to the system.
Adding can be done per file or per directory.

Note: Create will check for a $SHELF_PATH environment variable of which it
appends a './queries' to, when no dirPath or fileName is given.

Example:
1. To load a single file
	query create -p user_advice.json

2. To load a directory
	query create -p ./queries

3. To load using the environment variable path
	query create
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

	db := db.NewMGO()
	defer db.CloseMGO()

	if !stat.IsDir() {
		q, err := setFromFile("commands", file)
		if err != nil {
			log.Error("commands", "runCreate", err, "Completed")
			return
		}

		if err := query.CreateSet("commands", db, q); err != nil {
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

		return query.CreateSet("commands", db, q)
	})

	if err2 != nil {
		log.Error("commands", "runCreate", err, "Completed")
		return
	}

}
