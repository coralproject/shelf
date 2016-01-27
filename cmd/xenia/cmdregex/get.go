package cmdregex

import (
	"encoding/json"

	"github.com/coralproject/xenia/pkg/regex"

	"github.com/spf13/cobra"
)

var getLong = `Retrieves a Regex record from the system with the supplied name.

Example:
	regex get -n pre_script
`

// get contains the state for this command.
var get struct {
	name string
}

// addGet handles the retrival Regex records, displayed in json formatted response.
func addGet() {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Retrieves a Regex record by name.",
		Long:  getLong,
		Run:   runGet,
	}

	cmd.Flags().StringVarP(&get.name, "name", "n", "", "Name of the Regex.")

	regexCmd.AddCommand(cmd)
}

// runGet is the code that implements the get command.
func runGet(cmd *cobra.Command, args []string) {
	cmd.Printf("Getting Regex : Name[%s]\n", get.name)

	if get.name == "" {
		cmd.Help()
		return
	}

	set, err := regex.GetByName("", conn, get.name)
	if err != nil {
		cmd.Println("Getting Regex : ", err)
		return
	}

	data, err := json.MarshalIndent(&set, "", "    ")
	if err != nil {
		cmd.Println("Getting Regex : ", err)
		return
	}

	cmd.Printf("\n%s\n\n", string(data))
	return
}
