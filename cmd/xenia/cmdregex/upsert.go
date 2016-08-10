package cmdregex

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/coralproject/xenia/cmd/xenia/disk"
	"github.com/coralproject/xenia/cmd/xenia/web"
	"github.com/coralproject/xenia/internal/xenia/regex"
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
		Run:   runUpsert,
	}

	cmd.Flags().StringVarP(&upsert.path, "path", "p", "", "Path of Regex file or directory.")

	regexCmd.AddCommand(cmd)
}

// runUpsert is the code that implements the upsert command.
func runUpsert(cmd *cobra.Command, args []string) {
	cmd.Printf("Upserting Regex : Path[%s]\n", upsert.path)

	if upsert.path == "" {
		cmd.Help()
		return
	}

	pwd, err := os.Getwd()
	if err != nil {
		cmd.Println("Upserting Regex : ", err)
		return
	}

	file := filepath.Join(pwd, upsert.path)

	stat, err := os.Stat(file)
	if err != nil {
		cmd.Println("Upserting Regex : ", err)
		return
	}

	if !stat.IsDir() {
		rgx, err := disk.LoadRegex("", file)
		if err != nil {
			cmd.Println("Upserting Regex : ", err)
			return
		}

		if conn != nil {
			cmd.Printf("\n%+v\n", rgx)
			if err := regex.Upsert("", conn, rgx); err != nil {
				cmd.Println("Upserting Regex : ", err)
				return
			}
		} else {
			if err := runUpsertWeb(cmd, rgx); err != nil {
				cmd.Println("Upserting Regex : ", err)
				return
			}
		}

		cmd.Println("\n", "Upserting Regex : Upserted")
		return
	}

	f := func(path string) error {
		rgx, err := disk.LoadRegex("", path)
		if err != nil {
			return err
		}

		if conn != nil {
			return regex.Upsert("", conn, rgx)
		}

		return runUpsertWeb(cmd, rgx)
	}

	if err := disk.LoadDir(file, f); err != nil {
		cmd.Println("Upserting Regex : ", err)
		return
	}

	cmd.Println("\n", "Upserting Regex : Upserted")
}

// runUpsertWeb issues the command talking to the web service.
func runUpsertWeb(cmd *cobra.Command, rgx regex.Regex) error {
	verb := "PUT"
	url := "/1.0/regex"

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
