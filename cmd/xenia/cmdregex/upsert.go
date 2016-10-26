package cmdregex

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/coralproject/shelf/cmd/xenia/disk"
	"github.com/coralproject/shelf/cmd/xenia/web"
	"github.com/coralproject/shelf/internal/xenia/regex"
	"github.com/spf13/cobra"
)

var upsertLong = `Use upsert to add or update a regex in the system.
Adding can be done per file or per directory.

Example:
	regex upsert -p alpha.json

	regex upsert -p ./regexs
`

// upsert contains the state for this command.
var upsert struct {
	path string
}

// addUpsert handles the add or update of Regex records into the db.
func addUpsert() {
	cmd := &cobra.Command{
		Use:   "upsert",
		Short: "Upsert adds or updates a Regex from a file or directory.",
		Long:  upsertLong,
		RunE:  runUpsert,
	}

	cmd.Flags().StringVarP(&upsert.path, "path", "p", "", "Path of Regex file or directory.")

	regexCmd.AddCommand(cmd)
}

// runUpsert is the code that implements the upsert command.
func runUpsert(cmd *cobra.Command, args []string) error {
	cmd.Printf("Upserting Regex : Path[%s]\n", upsert.path)

	if upsert.path == "" {
		return fmt.Errorf("path must be provided")
	}

	file := upsert.path

	stat, err := os.Stat(file)
	if err != nil {
		return err
	}

	if !stat.IsDir() {
		rgx, err := disk.LoadRegex("", file)
		if err != nil {
			return err
		}

		if err := runUpsertWeb(cmd, rgx); err != nil {
			return err
		}

		cmd.Println("\n", "Upserting Regex : Upserted")
		return nil
	}

	f := func(path string) error {
		rgx, err := disk.LoadRegex("", path)
		if err != nil {
			return err
		}

		return runUpsertWeb(cmd, rgx)
	}

	if err := disk.LoadDir(file, f); err != nil {
		return err
	}

	cmd.Println("\n", "Upserting Regex : Upserted")
	return nil
}

// runUpsertWeb issues the command talking to the web service.
func runUpsertWeb(cmd *cobra.Command, rgx regex.Regex) error {
	verb := "PUT"
	url := "/v1/regex"

	data, err := json.Marshal(rgx)
	if err != nil {
		return err
	}

	cmd.Printf("\n%s\n\n", string(data))

	if _, err := web.Request(cmd, verb, url, bytes.NewBuffer(data)); err != nil {
		return err
	}

	return nil
}
