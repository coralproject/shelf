package cmdmask

import (
	"github.com/coralproject/shelf/cmd/xenia/web"
	"github.com/spf13/cobra"
)

var getLong = `Retrieves a Mask record from the system with the supplied collection and/or field.

Example:
	mask get -c collection -f field
`

// get contains the state for this command.
var get struct {
	collection string
	field      string
}

// addGet handles the retrival Script records, displayed in json formatted response.
func addGet() {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Retrieves a Mask record by collection/field.",
		Long:  getLong,
		Run:   runGet,
	}

	cmd.Flags().StringVarP(&get.collection, "collection", "c", "", "Name of the Collection.")
	cmd.Flags().StringVarP(&get.field, "field", "f", "", "Name of the Field.")

	maskCmd.AddCommand(cmd)
}

// runGet issues the command talking to the web service.
func runGet(cmd *cobra.Command, args []string) {
	verb := "GET"
	url := "/1.0/mask"

	if get.collection != "" {
		url += "/" + get.collection
	}

	if get.field != "" {
		url += "/" + get.field
	}

	resp, err := web.Request(cmd, verb, url, nil)
	if err != nil {
		cmd.Println("Getting Mask : ", err)
	}

	cmd.Printf("\n%s\n\n", resp)
}
