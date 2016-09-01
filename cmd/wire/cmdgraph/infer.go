package cmdgraph

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/coralproject/shelf/cmd/wire/disk"
	"github.com/coralproject/shelf/internal/wire"
	"github.com/spf13/cobra"
)

var inferLong = `Use execute to infer relationships from item documents.

Example:
	graph infer -p items.json

	graph infer -p ./items
`

// infer contains the state for this command.
var infer struct {
	path string
}

// addInferRelationships handles the inference of relationships.
func addInferRelationships() {
	cmd := &cobra.Command{
		Use:   "infer",
		Short: "infer infers relationships from item documents.",
		Long:  inferLong,
		Run:   runInferRelationships,
	}

	cmd.Flags().StringVarP(&infer.path, "path", "p", "", "Path to the JSON containing items")

	graphCmd.AddCommand(cmd)
}

// runInferRelationships is the code that implements the InferRelationships command.
func runInferRelationships(cmd *cobra.Command, args []string) {
	cmd.Printf("Inferring relationships : Path[%s]\n", infer.path)

	// Validate the path.
	if infer.path == "" {
		cmd.Help()
		return
	}

	// Get the working directory.
	pwd, err := os.Getwd()
	if err != nil {
		cmd.Println("Inferring relationships : ", err)
		return
	}

	// Join the provided path with the working directory.
	file := filepath.Join(pwd, infer.path)

	// Get the description of the file.
	stat, err := os.Stat(file)
	if err != nil {
		cmd.Println("Inferring relationships : ", err)
		return
	}

	// If a file is provided (i.e., not a directory), infer the relationships.
	if !stat.IsDir() {
		scr, err := disk.LoadItems("", file)
		if err != nil {
			cmd.Println("Inferring relationships : ", err)
			return
		}

		quadParams, err := wire.InferRelationships("", mgoDB, scr)
		if err != nil {
			cmd.Println("Inferring relationships : ", err)
			return
		}

		// Prepare the params for printing.
		data, err := json.MarshalIndent(quadParams, "", "    ")
		if err != nil {
			cmd.Println("Inferring relationships : ", err)
			return
		}

		cmd.Printf("\n%s\n\n", string(data))
		cmd.Println("\n", "Inferring relationships : Inferred")
		return
	}

	// If a directory is provided, infer relationships for all the included files.
	f := func(path string) error {
		scr, err := disk.LoadItems("", path)
		if err != nil {
			return err
		}

		quadParams, err := wire.InferRelationships("", mgoDB, scr)
		if err != nil {
			return err
		}

		// Prepare the params for printing.
		data, err := json.MarshalIndent(quadParams, "", "    ")
		if err != nil {
			cmd.Println("Inferring relationships : ", err)
			return err
		}

		cmd.Printf("\n%s\n\n", string(data))
		return nil
	}

	if err := disk.LoadDir(file, f); err != nil {
		cmd.Println("Inferring relationships : ", err)
		return
	}

	cmd.Println("\n", "Inferring relationships : Inferred")
}
