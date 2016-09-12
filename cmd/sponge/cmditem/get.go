package cmditem

import (
	"github.com/coralproject/shelf/cmd/sponge/web"
	"github.com/spf13/cobra"
)

var getLong = `Retrieves item records from the system having one of the supplied IDs.

Example:
	item get -i ids
`

// get contains the state for this command.
var get struct {
	IDs string
}

// addGet handles the retrival of item records, displayed in json formatted response.
func addGet() {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Retrieves all item records matching the supplied IDs.",
		Long:  getLong,
		Run:   runGet,
	}

	cmd.Flags().StringVarP(&get.IDs, "IDs", "i", "", "Item IDs.")

	itemCmd.AddCommand(cmd)
}

// runGet issues the command talking to the web service.
func runGet(cmd *cobra.Command, args []string) {
	verb := "GET"
	url := "/1.0/item"

	if get.IDs == "" {
		cmd.Help()
		return
	}

	url += "/" + get.IDs
	resp, err := web.Request(cmd, verb, url, nil)
	if err != nil {
		cmd.Println("Getting Items : ", err)
	}

	cmd.Printf("\n%s\n\n", resp)
}
