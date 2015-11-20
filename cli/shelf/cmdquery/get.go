package cmdquery

import (
	"github.com/spf13/cobra"
)

var getLong = `Retrieves a query record from the system using the supplied name.
Example:

		user get -n user_advice

`

// get contains the state for this command.
var get struct {
	name string
}

// addGet handles the retrival users records, displayed in json formatted response.
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
}
