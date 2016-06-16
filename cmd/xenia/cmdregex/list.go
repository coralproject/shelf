package cmdregex

import (
	"github.com/coralproject/xenia/cmd/xenia/web"
	"github.com/coralproject/xenia/internal/regex"

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

// runList is the code that implements the lists command.
func runList(cmd *cobra.Command, args []string) {
	if conn == nil {
		runListWeb(cmd)
		return
	}

	runListDB(cmd)
}

// runListWeb issues the command talking to the web service.
func runListWeb(cmd *cobra.Command) {
	verb := "GET"
	url := "/1.0/regex"

	resp, err := web.Request(cmd, verb, url, nil)
	if err != nil {
		cmd.Println("Getting Regex List : ", err)
	}

	cmd.Printf("\n%s\n\n", resp)
}

// runListDB issues the command talking to the DB.
func runListDB(cmd *cobra.Command) {
	cmd.Println("Getting Regex List")

	names, err := regex.GetNames("", conn)
	if err != nil {
		cmd.Println("Getting Regex List : ", err)
		return
	}

	cmd.Println("")

	for _, name := range names {
		cmd.Println(name)
	}

	cmd.Println("")
}
