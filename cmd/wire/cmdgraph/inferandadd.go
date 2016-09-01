package cmdgraph

import (
	"os"
	"path/filepath"

	"github.com/coralproject/shelf/cmd/wire/disk"
	"github.com/coralproject/shelf/internal/wire"
	"github.com/spf13/cobra"
)

var inferAddLong = `Use inferadd to infer relationships from item documents and add them to a graph.

Example:
	graph inferadd -p items.json

	graph inferadd -p ./items
`

// inferAdd contains the state for this command.
var inferAdd struct {
	path string
}

// addInferAddRelationships handles the inference of relationships.
func addInferAddRelationships() {
	cmd := &cobra.Command{
		Use:   "inferadd",
		Short: "inferadd infers relationships and adds them to a graph.",
		Long:  inferAddLong,
		Run:   runInferAddRelationships,
	}

	cmd.Flags().StringVarP(&inferAdd.path, "path", "p", "", "Path to the JSON containing items")

	graphCmd.AddCommand(cmd)
}

// runInferAddRelationships is the code that implements the InferRelationships
// and AddToGraph commands.
func runInferAddRelationships(cmd *cobra.Command, args []string) {
	cmd.Printf("Inferring and adding relationships : Path[%s]\n", inferAdd.path)

	// Validate the path.
	if inferAdd.path == "" {
		cmd.Help()
		return
	}

	// Get the working directory.
	pwd, err := os.Getwd()
	if err != nil {
		cmd.Println("Inferring and adding relationships : ", err)
		return
	}

	// Join the provided path with the working directory.
	file := filepath.Join(pwd, inferAdd.path)

	// Get the description of the file.
	stat, err := os.Stat(file)
	if err != nil {
		cmd.Println("Inferring and adding relationships : ", err)
		return
	}

	// If a file is provided (i.e., not a directory), infer and add the relationships.
	if !stat.IsDir() {
		scr, err := disk.LoadItems("", file)
		if err != nil {
			cmd.Println("Inferring and adding relationships : ", err)
			return
		}

		quadParams, err := wire.InferRelationships("", mgoDB, scr)
		if err != nil {
			cmd.Println("Inferring and adding relationships : ", err)
			return
		}

		if err := wire.AddToGraph("", graphDB, quadParams); err != nil {
			cmd.Println("Inferring and adding relationships : ", err)
			return
		}

		cmd.Println("\n", "Inferring and adding relationships : Added")
		return
	}

	// If a directory is provided, infer and add relationships for all the included files.
	f := func(path string) error {
		scr, err := disk.LoadItems("", path)
		if err != nil {
			return err
		}

		quadParams, err := wire.InferRelationships("", mgoDB, scr)
		if err != nil {
			cmd.Println("Inferring and adding relationships : ", err)
			return err
		}

		return wire.AddToGraph("", graphDB, quadParams)
	}

	if err := disk.LoadDir(file, f); err != nil {
		cmd.Println("Inferring and adding relationships : ", err)
		return
	}

	cmd.Println("\n", "Inferring and adding relationships : Added")
}
