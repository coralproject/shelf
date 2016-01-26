package cmdquery

import (
	"os"
	"path/filepath"

	"github.com/coralproject/xenia/cmd/xenia/disk"
	"github.com/coralproject/xenia/pkg/query"

	"github.com/ardanlabs/kit/db"
	"github.com/spf13/cobra"
)

var upsertLong = `Use upsert to add or update a Set in the system.
Adding can be done per file or per directory.

Example:
	query upsert -p user_advice.json

	query upsert -p ./sets
`

// upsert contains the state for this command.
var upsert struct {
	path string
}

// addUpsert handles the add or update of Set records into the db.
func addUpsert() {
	cmd := &cobra.Command{
		Use:   "upsert",
		Short: "Upsert adds or updates a Set from a file or directory.",
		Long:  upsertLong,
		Run:   runUpsert,
	}

	cmd.Flags().StringVarP(&upsert.path, "path", "p", "", "Path of Set file or directory.")

	queryCmd.AddCommand(cmd)
}

// runUpsert is the code that implements the upsert command.
func runUpsert(cmd *cobra.Command, args []string) {
	cmd.Printf("Upserting Set : Path[%s]\n", upsert.path)

	if upsert.path == "" {
		cmd.Help()
		return
	}

	pwd, err := os.Getwd()
	if err != nil {
		cmd.Println("Upserting Set : ", err)
		return
	}

	file := filepath.Join(pwd, upsert.path)

	stat, err := os.Stat(file)
	if err != nil {
		cmd.Println("Upserting Set : ", err)
		return
	}

	db, err := db.NewMGO("", mgoSession)
	if err != nil {
		cmd.Println("Upserting Set : ", err)
		return
	}
	defer db.CloseMGO("")

	if !stat.IsDir() {
		set, err := disk.LoadSet("", file)
		if err != nil {
			cmd.Println("Upserting Set : ", err)
			return
		}

		if err := query.Upsert("", db, set); err != nil {
			cmd.Println("Upserting Set : ", err)
			return
		}

		return
	}

	f := func(path string) error {
		set, err := disk.LoadSet("", path)
		if err != nil {
			return err
		}

		return query.Upsert("", db, set)
	}

	if err := disk.LoadDir(file, f); err != nil {
		cmd.Println("Upserting Set : ", err)
		return
	}

	cmd.Println("Upserting Set : Upserted")
}
