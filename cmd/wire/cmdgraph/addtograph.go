package cmdgraph

import (
	"os"
	"path/filepath"

	"github.com/coralproject/shelf/cmd/wire/disk"
	"github.com/coralproject/shelf/internal/wire"
	"github.com/spf13/cobra"
)

var addLong = `Use execute to add relationships to a graph.

Example:
	graph add -p relationships.json

	graph add -p ./relationships
`

// add contains the state for this command.
var add struct {
	path string
}

// addAddToGraph handles the addition of relationship quads to a graph.
func addAddToGraph() {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "add adds relationship quads to a graph.",
		Long:  addLong,
		Run:   runAddToGraph,
	}

	cmd.Flags().StringVarP(&add.path, "path", "p", "", "Path to the JSON containing relationships")

	graphCmd.AddCommand(cmd)
}

// runAddToGraph is the code that implements the AddToGraph command.
func runAddToGraph(cmd *cobra.Command, args []string) {
	cmd.Printf("Adding relationships : Path[%s]\n", add.path)

	// Validate the path.
	if add.path == "" {
		cmd.Help()
		return
	}

	pwd, err := os.Getwd()
	if err != nil {
		cmd.Println("Adding relationships : ", err)
		return
	}

	file := filepath.Join(pwd, add.path)

	stat, err := os.Stat(file)
	if err != nil {
		cmd.Println("Adding relationships : ", err)
		return
	}

	if !stat.IsDir() {
		scr, err := disk.LoadQuadParams("", file)
		if err != nil {
			cmd.Println("Adding relationships : ", err)
			return
		}

		if err := wire.AddToGraph("", graphDB, scr); err != nil {
			cmd.Println("Adding relationships : ", err)
			return
		}

		cmd.Println("\n", "Adding relationships : Added")
		return
	}

	f := func(path string) error {
		scr, err := disk.LoadQuadParams("", path)
		if err != nil {
			return err
		}

		return wire.AddToGraph("", graphDB, scr)
	}

	if err := disk.LoadDir(file, f); err != nil {
		cmd.Println("Adding relationships : ", err)
		return
	}

	cmd.Println("\n", "Adding relationships : Added")
}
