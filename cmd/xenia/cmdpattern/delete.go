package cmdpattern

import (
	"github.com/coralproject/shelf/cmd/xenia/web"
	"github.com/spf13/cobra"
)

var deleteLong = `Removes a Pattern from the system using the Pattern type.

Example:
	pattern delete -t ptype
`

// delete contains the state for this command.
var delete struct {
	ptype string
}

// addDel handles the deletion of Pattern records.
func addDel() {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Removes a Pattern record by type.",
		Long:  deleteLong,
		Run:   runDelete,
	}

	cmd.Flags().StringVarP(&delete.ptype, "type", "t", "", "Type of the Pattern.")

	patternCmd.AddCommand(cmd)
}

// runDelete issues the command talking to the web service.
func runDelete(cmd *cobra.Command, args []string) {
	verb := "DELETE"
	url := "/1.0/pattern/" + delete.ptype

	if _, err := web.Request(cmd, verb, url, nil); err != nil {
		cmd.Println("Deleting Pattern : ", err)
	}

	cmd.Println("Deleting Pattern : Deleted")
}
