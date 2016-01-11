package cmdquery

import (
	"os"
	"path/filepath"

	"github.com/coralproject/xenia/pkg/query"

	"github.com/ardanlabs/kit/db"

	"github.com/spf13/cobra"
)

var upsertLong = `Use upsert to add or update a query in the system.
Adding can be done per file or per directory.

Example:
	query upsert -p user_advice.json

	query upsert -p ./queries
`

// upsert contains the state for this command.
var upsert struct {
	path string
}

// addUpsert handles the add or update of query records into the db.
func addUpsert() {
	cmd := &cobra.Command{
		Use:   "upsert",
		Short: "Upsert adds or updates a query from a file or directory.",
		Long:  upsertLong,
		Run:   runUpsert,
	}

	cmd.Flags().StringVarP(&upsert.path, "path", "p", "", "Path of query file or directory.")

	queryCmd.AddCommand(cmd)
}

// runUpsert is the code that implements the upsert command.
func runUpsert(cmd *cobra.Command, args []string) {
	cmd.Printf("Upserting Query : Path[%s]\n", upsert.path)

	if upsert.path == "" {
		cmd.Help()
		return
	}

	pwd, err := os.Getwd()
	if err != nil {
		cmd.Println("Upserting Query : ", err)
		return
	}

	file := filepath.Join(pwd, upsert.path)

	stat, err := os.Stat(file)
	if err != nil {
		cmd.Println("Upserting Query : ", err)
		return
	}

	db := db.NewMGO()
	defer db.CloseMGO()

	if !stat.IsDir() {
		q, err := setFromFile("", file)
		if err != nil {
			cmd.Println("Upserting Query : ", err)
			return
		}

		if err := query.Upsert("", db, q); err != nil {
			cmd.Println("Upserting Query : ", err)
			return
		}

		return
	}

	err2 := loadDir(file, func(path string) error {
		q, err := setFromFile("", path)
		if err != nil {
			return err
		}

		return query.Upsert("", db, q)
	})

	if err2 != nil {
		cmd.Println("Upserting Query : ", err)
		return
	}

	cmd.Println("Upserting Query : Upserted")
}
