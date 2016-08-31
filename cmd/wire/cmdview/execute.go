package cmdview

import (
	"encoding/json"

	"github.com/coralproject/shelf/internal/wire"
	"github.com/spf13/cobra"
)

var executeLong = `Use execute to execute a view.

Example:
	view execute -n viewname -i itemkey -c resultscollection -b bufferlimit
`

// execute contains the state for this command.
var execute struct {
	viewName          string
	itemKey           string
	resultsCollection string
	bufferLimit       int
}

// addExecute handles the execution of a view.
func addExecute() {
	cmd := &cobra.Command{
		Use:   "execute",
		Short: "Execute executes a view based on input parameters.",
		Long:  executeLong,
		Run:   runExecute,
	}

	cmd.Flags().StringVarP(&execute.viewName, "name", "n", "", "View name")
	cmd.Flags().StringVarP(&execute.itemKey, "key", "i", "", "Item key")
	cmd.Flags().StringVarP(&execute.resultsCollection, "collection", "c", "", "Results collection")
	cmd.Flags().IntVarP(&execute.bufferLimit, "buffer", "b", 0, "Buffer Limit")

	viewCmd.AddCommand(cmd)
}

// runExecute is the code that implements the execute command.
func runExecute(cmd *cobra.Command, args []string) {
	cmd.Printf("Executing View : Name[%s]\n", execute.viewName)

	// Validate the input parameters.
	if execute.viewName == "" || execute.itemKey == "" {
		cmd.Help()
		return
	}

	// Ready the view parameters.
	viewParams := wire.ViewParams{
		ViewName:          execute.viewName,
		ItemKey:           execute.itemKey,
		ResultsCollection: execute.resultsCollection,
		BufferLimit:       execute.bufferLimit,
	}

	// Execute the view.
	results, err := wire.Execute("", mgoDB, graphDB, &viewParams)
	if err != nil {
		cmd.Println("Executing View : ", err)
		return
	}

	// Prepare the results for printing.
	data, err := json.MarshalIndent(results, "", "    ")
	if err != nil {
		cmd.Println("Executing View : ", err)
		return
	}

	cmd.Printf("\n%s\n\n", string(data))
	cmd.Println("\n", "Executing View : Executed")
}
