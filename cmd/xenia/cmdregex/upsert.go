package cmdregex

import (
	"os"
	"path/filepath"

	"github.com/coralproject/xenia/cmd/xenia/disk"
	"github.com/coralproject/xenia/pkg/regex"

	"github.com/ardanlabs/kit/db"
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

	db := db.NewMGO()
	defer db.CloseMGO()

	if !stat.IsDir() {
		rgx, err := disk.LoadRegex("", file)
		if err != nil {
			cmd.Println("Upserting Regex : ", err)
			return
		}

		if err := regex.Upsert("", db, rgx); err != nil {
			cmd.Println("Upserting Regex : ", err)
			return
		}

		return
	}

	f := func(path string) error {
		rgx, err := disk.LoadRegex("", path)
		if err != nil {
			return err
		}

		return regex.Upsert("", db, rgx)
	}

	if err := disk.LoadDir(file, f); err != nil {
		cmd.Println("Upserting Regex : ", err)
		return
	}

	cmd.Println("Upserting Regex : Upserted")
}
