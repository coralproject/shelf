package cmdscript

import (
	"github.com/coralproject/xenia/pkg/script"

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

// runList is the code that implements the lists command.
func runList(cmd *cobra.Command, args []string) {
	cmd.Println("Getting Script List")

	names, err := script.GetNames("", conn)
	if err != nil {
		cmd.Println("Getting Script List : ", err)
		return
	}

	cmd.Println("")

	for _, name := range names {
		cmd.Println(name)
	}

	cmd.Println("")
}
