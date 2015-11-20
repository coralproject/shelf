package commands

import (
	"bytes"
	"encoding/json"

	"github.com/coralproject/shelf/cli/query/db"
	"github.com/coralproject/shelf/pkg/log"
	"github.com/coralproject/shelf/pkg/mongo"
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

	rootCmd.AddCommand(cmd)
}

// runExecute is the code that implements the execute command.
func runExecute(cmd *cobra.Command, args []string) {
	if execute.name == "" {
		cmd.Help()
		return
	}

	q, err := db.GetByName(execute.name)
	if err != nil {
		log.Error("commands", "runCreate", err, "Completed")
		return
	}

	// Convert paramter into map
	params := map[string]string{}

	if execute.params != "" && execute.params != "{}" {
		err = json.NewDecoder(bytes.NewBufferString(execute.params)).Decode(&params)
		if err != nil {
			log.Error("commands", "runCreate", err, "Completed")
			return
		}
	}

	// Initialize the mongodb session.
	mongo.InitMGO()

	result, err := db.Execute(q, params)
	if err != nil {
		log.Error("commands", "runCreate", err, "Completed")
		return
	}

	_ = result

}
