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
	To create a single query:
	query create -p user_advice.json

	To create all the queries in a directory:
	query create -p ./queries

	To create all the queries in the env directory:
	query create
`

// create contains the state for this command.
var create struct {
	path string
}

// addCreate handles the creation of query records into the db.
func addCreate() {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates add a new query from a file or directory.",
		Long:  createLong,
		Run:   runCreate,
	}

	cmd.Flags().StringVarP(&create.path, "path", "p", "", "Path of query file or directory.")

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
