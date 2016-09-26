package cmdpattern

import (
	"github.com/coralproject/shelf/cmd/xenia/web"
	"github.com/spf13/cobra"
)

var getLong = `Retrieves pattern records from the system with the optional supplied type.

Example:
	pattern get
	
	pattern get -t type
`

// get contains the state for this command.
var get struct {
	ptype string
}

// addGet handles the retrival of pattern records, displayed in json formatted response.
func addGet() {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Retrieves all pattern records, or those matching an optional type.",
		Long:  getLong,
		Run:   runGet,
	}

	cmd.Flags().StringVarP(&get.ptype, "type", "t", "", "Pattern type.")

	patternCmd.AddCommand(cmd)
}

// runGet issues the command talking to the web service.
func runGet(cmd *cobra.Command, args []string) {
	verb := "GET"
	url := "/1.0/pattern"

	if get.ptype != "" {
		url += "/" + get.ptype
	}

	resp, err := web.Request(cmd, verb, url, nil)
	if err != nil {
		cmd.Println("Getting Pattern : ", err)
	}

	cmd.Printf("\n%s\n\n", resp)
}
