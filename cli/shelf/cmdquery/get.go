package cmdquery

import (
	"encoding/json"

	"github.com/coralproject/shelf/pkg/query"

	"github.com/ardanlabs/kit/db"

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
		Use:   "get",
		Short: "Retrieves a query record by name.",
		Long:  getLong,
		Run:   runGet,
	}

	cmd.Flags().StringVarP(&get.name, "name", "n", "", "Name of the query.")

	queryCmd.AddCommand(cmd)
}

// runGet is the code that implements the get command.
func runGet(cmd *cobra.Command, args []string) {
	cmd.Printf("Getting Query : Name[%s]\n", get.name)

	if get.name == "" {
		cmd.Help()
		return
	}

	db := db.NewMGO()
	defer db.CloseMGO()

	set, err := query.GetSetByName("", db, get.name)
	if err != nil {
		cmd.Println("Getting Query : ", err)
		return
	}

	data, err := json.MarshalIndent(&set, "", "    ")
	if err != nil {
		cmd.Println("Getting Query : ", err)
		return
	}

	cmd.Printf("\n%s\n\n", string(data))
	return
}
