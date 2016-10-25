package cmdscript

import (
	"github.com/coralproject/shelf/cmd/xenia/web"
	"github.com/spf13/cobra"
)

var getLong = `Retrieves a Script record from the system with the supplied name.

Example:
	script get -n pre_script
`

// get contains the state for this command.
var get struct {
	name string
}

// addGet handles the retrival Script records, displayed in json formatted response.
func addGet() {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Retrieves a Script record by name.",
		Long:  getLong,
		RunE:  runGet,
	}

	cmd.Flags().StringVarP(&get.name, "name", "n", "", "Name of the Script.")

	scriptCmd.AddCommand(cmd)
}

// runGet issues the command talking to the web service.
func runGet(cmd *cobra.Command, args []string) error {
	verb := "GET"
	url := "/v1/script/" + get.name

	resp, err := web.Request(cmd, verb, url, nil)
	if err != nil {
		return err
	}

	cmd.Printf("\n%s\n\n", resp)
	return nil
}
