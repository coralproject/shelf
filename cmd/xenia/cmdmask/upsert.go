package cmdmask

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/coralproject/shelf/cmd/xenia/disk"
	"github.com/coralproject/shelf/cmd/xenia/web"
	"github.com/coralproject/shelf/internal/xenia/mask"
	"github.com/spf13/cobra"
)

var upsertLong = `Use upsert to add or update a mask in the system.
Adding can be done per file or per directory.

Example:
	mask upsert -p mask.json

	mask upsert -p ./masks
`

// upsert contains the state for this command.
var upsert struct {
	path string
}

// addUpsert handles the add or update of mask records into the db.
func addUpsert() {
	cmd := &cobra.Command{
		Use:   "upsert",
		Short: "Upsert adds or updates a mask from a file or directory.",
		Long:  upsertLong,
		RunE:  runUpsert,
	}

	cmd.Flags().StringVarP(&upsert.path, "path", "p", "", "Path of mask file or directory.")

	maskCmd.AddCommand(cmd)
}

// runUpsert is the code that implements the upsert command.
func runUpsert(cmd *cobra.Command, args []string) error {
	cmd.Printf("Upserting Mask : Path[%s]\n", upsert.path)

	if upsert.path == "" {
		return fmt.Errorf("path must be provided")
	}

	file := upsert.path

	stat, err := os.Stat(file)
	if err != nil {
		return err
	}

	if !stat.IsDir() {
		msk, err := disk.LoadMask("", file)
		if err != nil {
			return err
		}

		if err := runUpsertWeb(cmd, msk); err != nil {
			return err
		}

		cmd.Println("\n", "Upserting Mask : Upserted")
		return nil
	}

	f := func(path string) error {
		msk, err := disk.LoadMask("", path)
		if err != nil {
			return err
		}

		return runUpsertWeb(cmd, msk)
	}

	if err := disk.LoadDir(file, f); err != nil {
		return err
	}

	cmd.Println("\n", "Upserting Mask : Upserted")
	return nil
}

// runUpsertWeb issues the command talking to the web service.
func runUpsertWeb(cmd *cobra.Command, msk mask.Mask) error {
	verb := "PUT"
	url := "/v1/mask"

	data, err := json.Marshal(msk)
	if err != nil {
		return err
	}

	cmd.Printf("\n%s\n\n", string(data))

	if _, err := web.Request(cmd, verb, url, bytes.NewBuffer(data)); err != nil {
		return err
	}

	return nil
}
