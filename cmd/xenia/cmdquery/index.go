package cmdquery

import (
	"bytes"
	"encoding/json"

	"github.com/coralproject/xenia/cmd/xenia/web"
	"github.com/coralproject/xenia/internal/query"

	"github.com/spf13/cobra"
)

var indexLong = `Use index to add or update a Set in the system.
Adding can be done per file or per directory.

Example:
	query index -n user_advice
`

// index contains the state for this command.
var index struct {
	name string
}

// addIndex handles the add or update of Set records into the db.
func addIndex() {
	cmd := &cobra.Command{
		Use:   "index",
		Short: "Index adds or updates a Set from a file or directory.",
		Long:  indexLong,
		Run:   runIndex,
	}

	cmd.Flags().StringVarP(&get.name, "name", "n", "", "Name of the Set.")

	queryCmd.AddCommand(cmd)
}

// runIndex is the code that implements the index command.
func runIndex(cmd *cobra.Command, args []string) {
	cmd.Printf("Ensure Indexes : Name[%s]\n", index.name)

	if index.name == "" {
		cmd.Help()
		return
	}

	set, err := query.GetByName("", conn, index.name)
	if err != nil {
		cmd.Println("Ensure Indexes : ", err)
		return
	}

	if err := query.EnsureIndexes("", conn, set); err != nil {
		cmd.Println("Ensure Indexes : ", err)
		return
	}

	cmd.Println("\n", "Ensure Indexes : Ensured")
}

// runIndexWeb issues the command talking to the web service.
func runIndexWeb(cmd *cobra.Command, set *query.Set) error {
	verb := "PUT"
	url := "/1.0/index/" + set.Name

	data, err := json.Marshal(set)
	if err != nil {
		return err
	}

	cmd.Printf("\n%s\n\n", string(data))

	if _, err := web.Request(cmd, verb, url, bytes.NewBuffer(data)); err != nil {
		return err
	}

	return nil
}
