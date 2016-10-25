package cmditem

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/coralproject/shelf/cmd/sponge/disk"
	"github.com/coralproject/shelf/cmd/sponge/web"
	"github.com/coralproject/shelf/internal/sponge/item"
	"github.com/spf13/cobra"
)

var upsertLong = `Use upsert to add or update an item in the system.

Example:
	item upsert -p item.json

	item upsert -p ./items
`

// upsert contains the state for this command.
var upsert struct {
	path string
}

// addUpsert handles the add or update of item records into the db.
func addUpsert() {
	cmd := &cobra.Command{
		Use:   "upsert",
		Short: "Upsert adds or updates an item from a file or directory.",
		Long:  upsertLong,
		RunE:  runUpsert,
	}

	cmd.Flags().StringVarP(&upsert.path, "path", "p", "", "Path of the item file or directory.")

	itemCmd.AddCommand(cmd)
}

// runUpsert is the code that implements the upsert command.
func runUpsert(cmd *cobra.Command, args []string) error {
	cmd.Printf("Upserting Items : Path[%s]\n", upsert.path)

	if upsert.path == "" {
		return fmt.Errorf("path must be provided")
	}

	file := upsert.path

	stat, err := os.Stat(file)
	if err != nil {
		return err
	}

	if !stat.IsDir() {
		item, err := disk.LoadItem("", file)
		if err != nil {
			return err
		}

		if err := runUpsertWeb(cmd, *item); err != nil {
			return err
		}

		cmd.Println("\n", "Upserting Items : Upserted")
		return nil
	}

	f := func(path string) error {
		item, err := disk.LoadItem("", path)
		if err != nil {
			return err
		}

		if err := runUpsertWeb(cmd, *item); err != nil {
			return err
		}

		return nil
	}

	if err := disk.LoadDir(file, f); err != nil {
		return err
	}

	cmd.Println("\n", "Upserting Items : Upserted")
	return nil
}

// runUpsertWeb issues the command talking to the web service.
func runUpsertWeb(cmd *cobra.Command, item item.Item) error {
	verb := "PUT"
	url := "/v1/item"

	data, err := json.Marshal(item)
	if err != nil {
		return err
	}

	cmd.Printf("\n%s\n\n", string(data))

	if _, err := web.Request(cmd, verb, url, bytes.NewBuffer(data)); err != nil {
		return err
	}

	return nil
}
