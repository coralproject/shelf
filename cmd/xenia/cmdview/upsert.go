package cmdview

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/coralproject/shelf/cmd/xenia/disk"
	"github.com/coralproject/shelf/cmd/xenia/web"
	"github.com/coralproject/shelf/internal/wire/view"
	"github.com/spf13/cobra"
)

var upsertLong = `Use upsert to add or update a view in the system.

Example:
	view upsert -p view.json

	view upsert -p ./views
`

// upsert contains the state for this command.
var upsert struct {
	path string
}

// addUpsert handles the add or update of view records into the db.
func addUpsert() {
	cmd := &cobra.Command{
		Use:   "upsert",
		Short: "Upsert adds or updates a view from a file or directory.",
		Long:  upsertLong,
		Run:   runUpsert,
	}

	cmd.Flags().StringVarP(&upsert.path, "path", "p", "", "Path of view file or directory.")

	viewCmd.AddCommand(cmd)
}

// runUpsert is the code that implements the upsert command.
func runUpsert(cmd *cobra.Command, args []string) {
	cmd.Printf("Upserting View : Path[%s]\n", upsert.path)

	if upsert.path == "" {
		cmd.Help()
		return
	}

	pwd, err := os.Getwd()
	if err != nil {
		cmd.Println("Upserting View : ", err)
		return
	}

	file := filepath.Join(pwd, upsert.path)

	stat, err := os.Stat(file)
	if err != nil {
		cmd.Println("Upserting View : ", err)
		return
	}

	if !stat.IsDir() {
		v, err := disk.LoadView("", file)
		if err != nil {
			cmd.Println("Upserting View : ", err)
			return
		}

		if err := runUpsertWeb(cmd, v); err != nil {
			cmd.Println("Upserting View : ", err)
			return
		}

		cmd.Println("\n", "Upserting View : Upserted")
		return
	}

	f := func(path string) error {
		v, err := disk.LoadView("", path)
		if err != nil {
			return err
		}

		return runUpsertWeb(cmd, v)
	}

	if err := disk.LoadDir(file, f); err != nil {
		cmd.Println("Upserting View : ", err)
		return
	}

	cmd.Println("\n", "Upserting View : Upserted")
}

// runUpsertWeb issues the command talking to the web service.
func runUpsertWeb(cmd *cobra.Command, v view.View) error {
	verb := "PUT"
	url := "/v1/view"

	data, err := json.Marshal(v)
	if err != nil {
		return err
	}

	cmd.Printf("\n%s\n\n", string(data))

	if _, err := web.Request(cmd, verb, url, bytes.NewBuffer(data)); err != nil {
		return err
	}

	return nil
}
