// This program provides the coral project shelf central CLI
// platform.
package main

import (
	"os"

	"github.com/coralproject/shelf/cli/shelf/cmddb"
	"github.com/coralproject/shelf/cli/shelf/cmdquery"
	"github.com/coralproject/shelf/cli/shelf/cmduser"
	"github.com/coralproject/shelf/pkg/cfg"
	"github.com/coralproject/shelf/pkg/srv/mongo"

	"github.com/spf13/cobra"
)

var shelf = &cobra.Command{
	Use:   "shelf",
	Short: "Shelf provides the central cli housing shelf's various cli tools",
}

func main() {
	if err := cfg.Init("SHELF"); err != nil {
		shelf.Println("Unable to initialize configuration")
		os.Exit(1)
	}

	err := mongo.InitMGO()
	if err != nil {
		shelf.Println("Unable to initialize MongoDB")
		os.Exit(1)
	}

	shelf.AddCommand(cmduser.GetCommands(), cmdquery.GetCommands(), cmddb.GetCommands())
	shelf.Execute()
}
