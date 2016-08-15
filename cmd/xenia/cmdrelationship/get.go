package cmdrelationship

import (
	"github.com/coralproject/xenia/cmd/xenia/web"
	"github.com/spf13/cobra"
)

var getLong = `Retrieves relationships record from the system with the optional supplied predicate.

Example:
	relationship get
	
	relationship get -p predicate
`

// get contains the state for this command.
var get struct {
	predicate string
}

// addGet handles the retrival of relationship records, displayed in json formatted response.
func addGet() {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Retrieves all relationship records, or those matching an optional predicate.",
		Long:  getLong,
		Run:   runGet,
	}

	cmd.Flags().StringVarP(&get.predicate, "predicate", "p", "", "Relationship predicate.")

	relationshipCmd.AddCommand(cmd)
}

// runGet issues the command talking to the web service.
func runGet(cmd *cobra.Command, args []string) {
	verb := "GET"
	url := "/1.0/relationship"

	if get.predicate != "" {
		url += "/" + get.predicate
	}

	resp, err := web.Request(cmd, verb, url, nil)
	if err != nil {
		cmd.Println("Getting Relationship : ", err)
	}

	cmd.Printf("\n%s\n\n", resp)
}
