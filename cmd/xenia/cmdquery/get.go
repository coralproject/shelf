package cmdquery

import (
	"encoding/json"

	"github.com/coralproject/xenia/pkg/query"

	"github.com/spf13/cobra"
)

var getLong = `Retrieves a Set record from the system with the supplied name.

Example:
	query get -n user_advice
`

// get contains the state for this command.
var get struct {
	name string
}

// addGet handles the retrival Set records, displayed in json formatted response.
func addGet() {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Retrieves a Set record by name.",
		Long:  getLong,
		Run:   runGet,
	}

	cmd.Flags().StringVarP(&get.name, "name", "n", "", "Name of the Set.")

	queryCmd.AddCommand(cmd)
}

// runGet is the code that implements the get command.
func runGet(cmd *cobra.Command, args []string) {
	cmd.Printf("Getting Set : Name[%s]\n", get.name)

	if get.name == "" {
		cmd.Help()
		return
	}

	set, err := query.GetByName("", conn, get.name)
	if err != nil {
		cmd.Println("Getting Set : ", err)
		return
	}

	data, err := json.MarshalIndent(&set, "", "    ")
	if err != nil {
		cmd.Println("Getting Set : ", err)
		return
	}

	cmd.Printf("\n%s\n\n", string(data))
	return
}
