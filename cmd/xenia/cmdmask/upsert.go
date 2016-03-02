package cmdmask

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/coralproject/xenia/cmd/xenia/disk"
	"github.com/coralproject/xenia/cmd/xenia/web"
	"github.com/coralproject/xenia/pkg/mask"

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
		Run:   runUpsert,
	}

	cmd.Flags().StringVarP(&upsert.path, "path", "p", "", "Path of mask file or directory.")

	maskCmd.AddCommand(cmd)
}

// runUpsert is the code that implements the upsert command.
func runUpsert(cmd *cobra.Command, args []string) {
	cmd.Printf("Upserting Mask : Path[%s]\n", upsert.path)

	if upsert.path == "" {
		cmd.Help()
		return
	}

	pwd, err := os.Getwd()
	if err != nil {
		cmd.Println("Upserting Mask : ", err)
		return
	}

	file := filepath.Join(pwd, upsert.path)

	stat, err := os.Stat(file)
	if err != nil {
		cmd.Println("Upserting Mask : ", err)
		return
	}

	if !stat.IsDir() {
		msk, err := disk.LoadMask("", file)
		if err != nil {
			cmd.Println("Upserting Mask : ", err)
			return
		}

		if conn != nil {
			cmd.Printf("\n%+v\n", msk)
			if err := mask.Upsert("", conn, msk); err != nil {
				cmd.Println("Upserting Mask : ", err)
				return
			}
		} else {
			if err := runUpsertWeb(cmd, msk); err != nil {
				cmd.Println("Upserting Mask : ", err)
				return
			}
		}

		cmd.Println("\n", "Upserting Mask : Upserted")
		return
	}

	f := func(path string) error {
		msk, err := disk.LoadMask("", path)
		if err != nil {
			return err
		}

		if conn != nil {
			return mask.Upsert("", conn, msk)
		}

		return runUpsertWeb(cmd, msk)
	}

	if err := disk.LoadDir(file, f); err != nil {
		cmd.Println("Upserting Mask : ", err)
		return
	}

	cmd.Println("\n", "Upserting Mask : Upserted")
}

// runUpsertWeb issues the command talking to the web service.
func runUpsertWeb(cmd *cobra.Command, msk mask.Mask) error {
	verb := "PUT"
	url := "/1.0/mask"

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
