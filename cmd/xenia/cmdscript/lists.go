package cmdscript

import (
	"github.com/coralproject/shelf/cmd/xenia/web"
	"github.com/spf13/cobra"
)

var listLong = `Retrieves a list of all available Script names.

Example:
	script list
`

// addList handles the retrival Script records names.
func addList() {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Retrieves a list of all available Script names.",
		Long:  listLong,
		Run:   runList,
	}
	scriptCmd.AddCommand(cmd)
}

// runList issues the command talking to the web service.
func runList(cmd *cobra.Command, args []string) {
	verb := "GET"
	url := "/v1/script"

	resp, err := web.Request(cmd, verb, url, nil)
	if err != nil {
		cmd.Println("Getting Script List : ", err)
	}

	cmd.Printf("\n%s\n\n", resp)
}
