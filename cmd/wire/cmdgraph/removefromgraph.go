package cmdgraph

import (
	"os"
	"path/filepath"

	"github.com/coralproject/shelf/cmd/wire/disk"
	"github.com/coralproject/shelf/internal/wire"
	"github.com/spf13/cobra"
)

var removeLong = `Use execute to remove relationships from a graph.

Example:
	graph remove -p relationships.json

	graph remove -p ./relationships
`

// remove contains the state for this command.
var remove struct {
	path string
}

// addRemoveFromGraph handles the removal of relationships from a graph.
func addRemoveFromGraph() {
	cmd := &cobra.Command{
		Use:   "remove",
		Short: "remove removes relationship quads from a graph.",
		Long:  removeLong,
		Run:   runRemoveFromGraph,
	}

	cmd.Flags().StringVarP(&remove.path, "path", "p", "", "Path to the JSON containing relationships")

	graphCmd.AddCommand(cmd)
}

// runRemoveFromGraph is the code that implements the RemoveFromGraph command.
func runRemoveFromGraph(cmd *cobra.Command, args []string) {
	cmd.Printf("Removing relationships : Path[%s]\n", remove.path)

	// Validate the path.
	if remove.path == "" {
		cmd.Help()
		return
	}

	// Get the working directory.
	pwd, err := os.Getwd()
	if err != nil {
		cmd.Println("Removing relationships : ", err)
		return
	}

	// Join the provided path with the working directory.
	file := filepath.Join(pwd, remove.path)

	// Get the description of the file.
	stat, err := os.Stat(file)
	if err != nil {
		cmd.Println("Removing relationships : ", err)
		return
	}

	// If a file is provided (i.e., not a directory), remove the relationships.
	if !stat.IsDir() {
		scr, err := disk.LoadQuadParams("", file)
		if err != nil {
			cmd.Println("Removing relationships : ", err)
			return
		}

		if err := wire.RemoveFromGraph("", graphDB, scr); err != nil {
			cmd.Println("Removing relationships : ", err)
			return
		}

		cmd.Println("\n", "Removing relationships : Removed")
		return
	}

	// If a directory is provided, remove relationships for all the included files.
	f := func(path string) error {
		scr, err := disk.LoadQuadParams("", path)
		if err != nil {
			return err
		}

		return wire.RemoveFromGraph("", graphDB, scr)
	}

	if err := disk.LoadDir(file, f); err != nil {
		cmd.Println("Removing relationships : ", err)
		return
	}

	cmd.Println("\n", "Removing relationships : Removed")
}
