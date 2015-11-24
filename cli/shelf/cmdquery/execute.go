package cmdquery

import (
	"log"

	"github.com/spf13/cobra"
)

var executeLong = `Executes a giving query using the name and a map of parameters.
If the map is provided, it will be converted to a map object else ignored.

Example:

	query execute -n user_advice

	query execute -n user_advice -p {"name":"john"}

`

// execute contains the state for this command.
var execute struct {
	name   string
	params string
}

// addExecute handles the execution of queries.
func addExecute() {
	cmd := &cobra.Command{
		Use:   "execute [-n name -p {parameters...}]",
		Short: "executes a query using its name and a map of parameters",
		Long:  executeLong,
		Run:   runExecute,
	}

	cmd.Flags().StringVarP(&execute.name, "name", "n", "", "name of query in db")
	cmd.Flags().StringVarP(&execute.params, "params", "p", "", "parameter map for query")

	queryCmd.AddCommand(cmd)
}

// runExecute is the code that implements the execute command.
func runExecute(cmd *cobra.Command, args []string) {
	log.Fatal("commands", "runExecute", "The query execute functionality is still pending")
}
