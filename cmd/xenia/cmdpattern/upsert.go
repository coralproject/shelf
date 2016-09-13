package cmdpattern

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/coralproject/shelf/cmd/xenia/disk"
	"github.com/coralproject/shelf/cmd/xenia/web"
	"github.com/coralproject/shelf/internal/wire/pattern"
	"github.com/spf13/cobra"
)

var upsertLong = `Use upsert to add or update a pattern in the system.

Example:
	pattern upsert -p pattern.json

	pattern upsert -p ./patterns
`

// upsert contains the state for this command.
var upsert struct {
	path string
}

// addUpsert handles the add or update of pattern records into the db.
func addUpsert() {
	cmd := &cobra.Command{
		Use:   "upsert",
		Short: "Upsert adds or updates a pattern from a file or directory.",
		Long:  upsertLong,
		Run:   runUpsert,
	}

	cmd.Flags().StringVarP(&upsert.path, "path", "p", "", "Path of pattern file or directory.")

	patternCmd.AddCommand(cmd)
}

// runUpsert is the code that implements the upsert command.
func runUpsert(cmd *cobra.Command, args []string) {
	cmd.Printf("Upserting Pattern : Path[%s]\n", upsert.path)

	if upsert.path == "" {
		cmd.Help()
		return
	}

	pwd, err := os.Getwd()
	if err != nil {
		cmd.Println("Upserting Pattern : ", err)
		return
	}

	file := filepath.Join(pwd, upsert.path)

	stat, err := os.Stat(file)
	if err != nil {
		cmd.Println("Upserting Pattern : ", err)
		return
	}

	if !stat.IsDir() {
		p, err := disk.LoadPattern("", file)
		if err != nil {
			cmd.Println("Upserting Pattern : ", err)
			return
		}

		if err := runUpsertWeb(cmd, p); err != nil {
			cmd.Println("Upserting Pattern : ", err)
			return
		}

		cmd.Println("\n", "Upserting Pattern : Upserted")
		return
	}

	f := func(path string) error {
		p, err := disk.LoadPattern("", path)
		if err != nil {
			return err
		}

		return runUpsertWeb(cmd, p)
	}

	if err := disk.LoadDir(file, f); err != nil {
		cmd.Println("Upserting Pattern : ", err)
		return
	}

	cmd.Println("\n", "Upserting Pattern : Upserted")
}

// runUpsertWeb issues the command talking to the web service.
func runUpsertWeb(cmd *cobra.Command, p pattern.Pattern) error {
	verb := "PUT"
	url := "/1.0/pattern"

	data, err := json.Marshal(p)
	if err != nil {
		return err
	}

	cmd.Printf("\n%s\n\n", string(data))

	if _, err := web.Request(cmd, verb, url, bytes.NewBuffer(data)); err != nil {
		return err
	}

	return nil
}
