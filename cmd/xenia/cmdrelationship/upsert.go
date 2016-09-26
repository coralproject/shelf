package cmdrelationship

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/coralproject/shelf/cmd/xenia/disk"
	"github.com/coralproject/shelf/cmd/xenia/web"
	"github.com/coralproject/shelf/internal/wire/relationship"
	"github.com/spf13/cobra"
)

var upsertLong = `Use upsert to add or update a relationship in the system.

Example:
	relationship upsert -p relationship.json

	relationship upsert -p ./relationships
`

// upsert contains the state for this command.
var upsert struct {
	path string
}

// addUpsert handles the add or update of relationship records into the db.
func addUpsert() {
	cmd := &cobra.Command{
		Use:   "upsert",
		Short: "Upsert adds or updates a relationship from a file or directory.",
		Long:  upsertLong,
		Run:   runUpsert,
	}

	cmd.Flags().StringVarP(&upsert.path, "path", "p", "", "Path of relationship file or directory.")

	relationshipCmd.AddCommand(cmd)
}

// runUpsert is the code that implements the upsert command.
func runUpsert(cmd *cobra.Command, args []string) {
	cmd.Printf("Upserting Relationship : Path[%s]\n", upsert.path)

	if upsert.path == "" {
		cmd.Help()
		return
	}

	pwd, err := os.Getwd()
	if err != nil {
		cmd.Println("Upserting Relationship : ", err)
		return
	}

	file := filepath.Join(pwd, upsert.path)

	stat, err := os.Stat(file)
	if err != nil {
		cmd.Println("Upserting Relationship : ", err)
		return
	}

	if !stat.IsDir() {
		rel, err := disk.LoadRelationship("", file)
		if err != nil {
			cmd.Println("Upserting Relationship : ", err)
			return
		}

		if err := runUpsertWeb(cmd, rel); err != nil {
			cmd.Println("Upserting Relationship : ", err)
			return
		}

		cmd.Println("\n", "Upserting Relationship : Upserted")
		return
	}

	f := func(path string) error {
		rel, err := disk.LoadRelationship("", path)
		if err != nil {
			return err
		}

		return runUpsertWeb(cmd, rel)
	}

	if err := disk.LoadDir(file, f); err != nil {
		cmd.Println("Upserting Relationship : ", err)
		return
	}

	cmd.Println("\n", "Upserting Relationship : Upserted")
}

// runUpsertWeb issues the command talking to the web service.
func runUpsertWeb(cmd *cobra.Command, rel relationship.Relationship) error {
	verb := "PUT"
	url := "/1.0/relationship"

	data, err := json.Marshal(rel)
	if err != nil {
		return err
	}

	cmd.Printf("\n%s\n\n", string(data))

	if _, err := web.Request(cmd, verb, url, bytes.NewBuffer(data)); err != nil {
		return err
	}

	return nil
}
