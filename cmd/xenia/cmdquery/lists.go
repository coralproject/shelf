package cmdquery

import (
	"github.com/coralproject/shelf/cmd/xenia/web"
	"github.com/spf13/cobra"
)

var listLong = `Retrieves a list of all available Set names.

Example:
	query list
`

// addList handles the retrival Set records names.
func addList() {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Retrieves a list of all available Set names.",
		Long:  listLong,
		Run:   runList,
	}
	queryCmd.AddCommand(cmd)
}

// runList issues the command talking to the web service.
func runList(cmd *cobra.Command, args []string) {
	verb := "GET"
	url := "/v1/query"

	resp, err := web.Request(cmd, verb, url, nil)
	if err != nil {
		cmd.Println("Getting Set List : ", err)
	}

	cmd.Printf("\n%s\n\n", resp)
}
