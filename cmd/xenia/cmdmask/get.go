package cmdmask

import (
	"encoding/json"

	"github.com/coralproject/xenia/cmd/xenia/web"
	"github.com/coralproject/xenia/internal/mask"
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

// runGet is the code that implements the get command.
func runGet(cmd *cobra.Command, args []string) {
	if conn == nil {
		runGetWeb(cmd)
		return
	}

	runGetDB(cmd)
}

// runListWeb issues the command talking to the web service.
func runGetWeb(cmd *cobra.Command) {
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

// runGetDB issues the command talking to the DB.
func runGetDB(cmd *cobra.Command) {
	cmd.Printf("Getting Mask : Collection[%s] Field[%s]\n", get.collection, get.field)

	if get.collection == "" {
		get.collection = "*"
	}

	if get.field != "" {
		msk, err := mask.GetByName("", conn, get.collection, get.field)
		if err != nil {
			cmd.Println("Getting Mask : ", err)
			return
		}

		data, err := json.MarshalIndent(msk, "", "    ")
		if err != nil {
			cmd.Println("Getting Mask : ", err)
			return
		}

		cmd.Printf("\n%s\n\n", string(data))
		return
	}

	masks, err := mask.GetByCollection("", conn, get.collection)
	if err != nil {
		cmd.Println("Getting Mask : ", err)
		return
	}

	data, err := json.MarshalIndent(masks, "", "    ")
	if err != nil {
		cmd.Println("Getting Mask : ", err)
		return
	}

	cmd.Printf("\n%s\n\n", string(data))
	return
}
