package cmdregex

import (
	"github.com/coralproject/xenia/pkg/regex"

	"github.com/ardanlabs/kit/db"
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
	cmd.Println("Getting Regex List")

	db := db.NewMGO()
	defer db.CloseMGO()

	names, err := regex.GetNames("", db)
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
