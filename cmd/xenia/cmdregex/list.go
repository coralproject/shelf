package cmdregex

import (
	"github.com/coralproject/xenia/cmd/xenia/web"
	"github.com/spf13/cobra"
)

var listLong = `Retrieves a list of all available Regex names.

Example:
	regex list
`

// addList handles the retrival Regex records names.
func addList() {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Retrieves a list of all available Regex names.",
		Long:  listLong,
		Run:   runList,
	}
	regexCmd.AddCommand(cmd)
}

// runList issues the command talking to the web service.
func runList(cmd *cobra.Command, args []string) {
	verb := "GET"
	url := "/1.0/regex"

	resp, err := web.Request(cmd, verb, url, nil)
	if err != nil {
		cmd.Println("Getting Regex List : ", err)
	}

	cmd.Printf("\n%s\n\n", resp)
}
