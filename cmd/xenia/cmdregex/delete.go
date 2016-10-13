package cmdregex

import (
	"github.com/coralproject/shelf/cmd/xenia/web"
	"github.com/spf13/cobra"
)

var deleteLong = `Removes a Regex from the system using the regex name.

Example:
	regex delete -n user_advice
`

// delete contains the state for this command.
var delete struct {
	name string
}

// addDel handles the removal of a regex document.
func addDel() {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Removes a Regex record by name.",
		Long:  deleteLong,
		Run:   runDelete,
	}

	cmd.Flags().StringVarP(&delete.name, "name", "n", "", "Name of the Regex record.")

	regexCmd.AddCommand(cmd)
}

// runDelete issues the command talking to the web service.
func runDelete(cmd *cobra.Command, args []string) {
	verb := "DELETE"
	url := "/v1/regex/" + get.name

	if _, err := web.Request(cmd, verb, url, nil); err != nil {
		cmd.Println("Deleting Regex : ", err)
	}

	cmd.Println("Deleting Regex : Deleted")
}
