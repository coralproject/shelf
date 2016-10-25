package cmdscript

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/coralproject/shelf/cmd/xenia/disk"
	"github.com/coralproject/shelf/cmd/xenia/web"
	"github.com/coralproject/shelf/internal/xenia/script"
	"github.com/spf13/cobra"
)

var upsertLong = `Use upsert to add or update a script in the system.
Adding can be done per file or per directory.

Example:
	script upsert -p pre_script.json

	script upsert -p ./pre_scripts
`

// upsert contains the state for this command.
var upsert struct {
	path string
}

// addUpsert handles the add or update of script records into the db.
func addUpsert() {
	cmd := &cobra.Command{
		Use:   "upsert",
		Short: "Upsert adds or updates a script from a file or directory.",
		Long:  upsertLong,
		RunE:  runUpsert,
	}

	cmd.Flags().StringVarP(&upsert.path, "path", "p", "", "Path of script file or directory.")

	scriptCmd.AddCommand(cmd)
}

// runUpsert is the code that implements the upsert command.
func runUpsert(cmd *cobra.Command, args []string) error {
	cmd.Printf("Upserting Script : Path[%s]\n", upsert.path)

	if upsert.path == "" {
		return fmt.Errorf("path must be provided")
	}

	file := upsert.path

	stat, err := os.Stat(file)
	if err != nil {
		return err
	}

	if !stat.IsDir() {
		scr, err := disk.LoadScript("", file)
		if err != nil {
			return err
		}

		if err := runUpsertWeb(cmd, scr); err != nil {
			return err
		}

		cmd.Println("\n", "Upserting Script : Upserted")
		return nil
	}

	f := func(path string) error {
		scr, err := disk.LoadScript("", path)
		if err != nil {
			return err
		}

		return runUpsertWeb(cmd, scr)
	}

	if err := disk.LoadDir(file, f); err != nil {
		return err
	}

	cmd.Println("\n", "Upserting Script : Upserted")
	return nil
}

// runUpsertWeb issues the command talking to the web service.
func runUpsertWeb(cmd *cobra.Command, scr script.Script) error {
	verb := "PUT"
	url := "/v1/script"

	data, err := json.Marshal(scr)
	if err != nil {
		return err
	}

	cmd.Printf("\n%s\n\n", string(data))

	if _, err := web.Request(cmd, verb, url, bytes.NewBuffer(data)); err != nil {
		return err
	}

	return nil
}
