package cmdrelationship

import (
	"github.com/coralproject/shelf/cmd/xenia/web"
	"github.com/spf13/cobra"
)

var deleteLong = `Removes a Relationship from the system using the Relationship predicate.

Example:
	relationship delete -p predicate
`

// delete contains the state for this command.
var delete struct {
	predicate string
}

// addDel handles the deletion of Relationship records.
func addDel() {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Removes a Relationship record by predicate.",
		Long:  deleteLong,
		Run:   runDelete,
	}

	cmd.Flags().StringVarP(&delete.predicate, "predicate", "p", "", "Preciate of the Relationship.")

	relationshipCmd.AddCommand(cmd)
}

// runDelete issues the command talking to the web service.
func runDelete(cmd *cobra.Command, args []string) {
	verb := "DELETE"
	url := "/1.0/relationship/" + delete.predicate

	if _, err := web.Request(cmd, verb, url, nil); err != nil {
		cmd.Println("Deleting Relationship : ", err)
	}

	cmd.Println("Deleting Relationship : Deleted")
}
