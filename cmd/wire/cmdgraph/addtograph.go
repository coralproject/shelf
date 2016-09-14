package cmdgraph

import (
	"os"
	"path/filepath"

	"github.com/coralproject/shelf/cmd/wire/disk"
	"github.com/coralproject/shelf/internal/wire"
	"github.com/spf13/cobra"
)

var addLong = `Use execute to add relationships, inferred from an item, to a graph.

Example:
	graph add -p item.json

	graph add -p ./items
`

// add contains the state for this command.
var add struct {
	path string
}

// addAddToGraph handles the addition of relationship quads to a graph.
func addAddToGraph() {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "add adds relationships to a graph.",
		Long:  addLong,
		Run:   runAddToGraph,
	}

	cmd.Flags().StringVarP(&add.path, "path", "p", "", "Path to the JSON containing an item")

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

	// Get the working directory.
	pwd, err := os.Getwd()
	if err != nil {
		cmd.Println("Adding relationships : ", err)
		return
	}

	// Join the provided path with the working directory.
	file := filepath.Join(pwd, add.path)

	// Get the description of the file.
	stat, err := os.Stat(file)
	if err != nil {
		cmd.Println("Adding relationships : ", err)
		return
	}

	// If a file is provided (i.e., not a directory), add the relationships.
	if !stat.IsDir() {
		scr, err := disk.LoadItem("", file)
		if err != nil {
			cmd.Println("Adding relationships : ", err)
			return
		}

		if err := wire.AddToGraph("", mgoDB, graphDB, scr); err != nil {
			cmd.Println("Adding relationships : ", err)
			return
		}

		cmd.Println("\n", "Adding relationships : Added")
		return
	}

	// If a directory is provided, add relationships for all the included files.
	f := func(path string) error {
		scr, err := disk.LoadItem("", path)
		if err != nil {
			return err
		}

		return wire.AddToGraph("", mgoDB, graphDB, scr)
	}

	if err := disk.LoadDir(file, f); err != nil {
		cmd.Println("Adding relationships : ", err)
		return
	}

	cmd.Println("\n", "Adding relationships : Added")
}
