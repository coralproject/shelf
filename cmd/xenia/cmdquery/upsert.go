package cmdquery

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/coralproject/xenia/cmd/xenia/disk"
	"github.com/coralproject/xenia/cmd/xenia/web"
	"github.com/coralproject/xenia/internal/query"

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

	file, err := filepath.Abs(upsert.path)
	if err != nil {
		cmd.Println("Upserting Set : ", err)
		return
	}

	stat, err := os.Stat(file)
	if err != nil {
		cmd.Println("Upserting Set : ", err)
		return
	}

	if !stat.IsDir() {
		set, err := disk.LoadSet("", file)
		if err != nil {
			cmd.Println("Upserting Set : ", err)
			return
		}

		if conn != nil {
			cmd.Printf("\n%+v\n", set)
			if err := query.Upsert("", conn, set); err != nil {
				cmd.Println("Upserting Set : ", err)
				return
			}
		} else {
			if err := runUpsertWeb(cmd, set); err != nil {
				cmd.Println("Upserting Set : ", err)
				return
			}
		}

		cmd.Println("\n", "Upserting Set : Upserted")
		return
	}

	f := func(path string) error {
		set, err := disk.LoadSet("", path)
		if err != nil {
			return err
		}

		if conn != nil {
			return query.Upsert("", conn, set)
		}

		return runUpsertWeb(cmd, set)
	}

	if err := disk.LoadDir(file, f); err != nil {
		cmd.Println("Upserting Set : ", err)
		return
	}

	cmd.Println("\n", "Upserting Set : Upserted")
}

// runUpsertWeb issues the command talking to the web service.
func runUpsertWeb(cmd *cobra.Command, set *query.Set) error {
	verb := "PUT"
	url := "/1.0/query"

	data, err := json.Marshal(set)
	if err != nil {
		return err
	}

	cmd.Printf("\n%s\n\n", string(data))

	if _, err := web.Request(cmd, verb, url, bytes.NewBuffer(data)); err != nil {
		return err
	}

	return nil
}
