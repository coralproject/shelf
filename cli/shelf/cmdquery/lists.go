package cmdquery

import (
	"bytes"
	"fmt"

	"github.com/coralproject/shelf/pkg/db"
	"github.com/coralproject/shelf/pkg/log"
	"github.com/coralproject/shelf/pkg/srv/query"

	"github.com/spf13/cobra"
)

var listLong = `Retrieves a list of all available query names.

Example:
	query list
`

// addList handles the retrival query records names.
func addList() {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Retrieves a list of all available query names.",
		Long:  listLong,
		Run:   runList,
	}
	queryCmd.AddCommand(cmd)
}

// runList is the code that implements the lists command.
func runList(cmd *cobra.Command, args []string) {
	db := db.NewMGO()
	defer db.CloseMGO()

	names, err := query.GetSetNames("commands", db)
	if err != nil {
		log.Error("commands", "runGet", err, "Completed")
		return
	}

	var buf bytes.Buffer
	buf.Write([]byte("\n"))
	fmt.Fprint(&buf, fmt.Sprintf("Total Records: %d", len(names)))
	buf.Write([]byte("\n"))

	for _, name := range names {
		fmt.Fprint(&buf, name)
		buf.Write([]byte("\n"))
	}

	fmt.Println(buf.String())
	return
}
