package cmdview

import (
	"github.com/coralproject/shelf/cmd/xenia/web"
	"github.com/spf13/cobra"
)

var getLong = `Retrieves view records from the system with the optional supplied name.

Example:
	view get

	view get -n name
`

// get contains the state for this command.
var get struct {
	name string
}

// addGet handles the retrival of view records, displayed in json formatted response.
func addGet() {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Retrieves all view records, or those matching an optional name.",
		Long:  getLong,
		RunE:  runGet,
	}

	cmd.Flags().StringVarP(&get.name, "name", "n", "", "View name.")

	viewCmd.AddCommand(cmd)
}

// runGet issues the command talking to the web service.
func runGet(cmd *cobra.Command, args []string) error {
	verb := "GET"
	url := "/v1/view"

	if get.name != "" {
		url += "/" + get.name
	}

	resp, err := web.Request(cmd, verb, url, nil)
	if err != nil {
		return err
	}

	cmd.Printf("\n%s\n\n", resp)
	return nil
}
