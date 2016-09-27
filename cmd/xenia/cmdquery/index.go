package cmdquery

import (
	"bytes"
	"encoding/json"

	"github.com/coralproject/shelf/cmd/xenia/web"
	"github.com/coralproject/shelf/internal/xenia/query"
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

	cmd.Flags().StringVarP(&index.name, "name", "n", "", "Name of the Set.")

	queryCmd.AddCommand(cmd)
}

// runIndex issues the command talking to the web service.
func runIndex(cmd *cobra.Command, args []string) {
	cmd.Printf("Ensure Indexes : Name[%s]\n", index.name)

	set, err := runGetSet(cmd, index.name)
	if err != nil {
		cmd.Println("Ensure Indexes : ", err)
		return
	}

	verb := "PUT"
	url := "/v1/index/" + index.name

	data, err := json.Marshal(set)
	if err != nil {
		cmd.Println("Ensure Indexes : ", err)
		return
	}

	cmd.Printf("\n%s\n\n", string(data))

	if _, err := web.Request(cmd, verb, url, bytes.NewBuffer(data)); err != nil {
		cmd.Println("Ensure Indexes : ", err)
		return
	}

	cmd.Println("\n", "Ensure Indexes : Ensured")
	return
}

// runGetSet get a query set by name.
func runGetSet(cmd *cobra.Command, name string) (query.Set, error) {
	verb := "GET"
	url := "/v1/query/" + name

	resp, err := web.Request(cmd, verb, url, nil)
	if err != nil {
		return query.Set{}, err
	}

	var set query.Set
	if err = json.Unmarshal([]byte(resp), &set); err != nil {
		return query.Set{}, err
	}

	return set, nil
}
