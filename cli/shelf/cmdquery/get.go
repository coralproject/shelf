package cmdquery

import (
	"encoding/json"
	"fmt"

	"github.com/coralproject/shelf/pkg/db"
	"github.com/coralproject/shelf/pkg/log"
	"github.com/coralproject/shelf/pkg/srv/query"

	"github.com/spf13/cobra"
)

var getLong = `Retrieves a query record from the system with the supplied name.

Example:
	query get -n user_advice
`

// get contains the state for this command.
var get struct {
	name string
}

// addGet handles the retrival query records, displayed in json formatted response.
func addGet() {
	cmd := &cobra.Command{
		Use:   "get [-n name]",
		Short: "Retrieves a query record",
		Long:  getLong,
		Run:   runGet,
	}

	cmd.Flags().StringVarP(&get.name, "name", "n", "", "name of the user record")

	queryCmd.AddCommand(cmd)
}

// runGet is the code that implements the get command.
func runGet(cmd *cobra.Command, args []string) {
	if get.name == "" {
		cmd.Help()
		return
	}

	db := db.NewMGO()
	defer db.CloseMGO()

	user, err := query.GetSetByName("commands", db, get.name)
	if err != nil {
		log.Error("commands", "runGet", err, "Completed")
		return
	}

	res, err := json.MarshalIndent(user, "", "\t")
	if err != nil {
		log.Error("commands", "runGet", err, "Completed")
		return
	}

	// TODO: What are you doing with doc
	fmt.Printf(`
Result of Query(%s):
	%s
`, get.name, string(res))

	return
}
