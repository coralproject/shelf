// This program provides the coral project shelf central CLI
// platform.
package main

import (
	"bytes"
	"os"

	"github.com/coralproject/shelf/cli/shelf/cmddb"
	"github.com/coralproject/shelf/cli/shelf/cmdquery"
	"github.com/coralproject/shelf/cli/shelf/cmduser"
	"github.com/coralproject/shelf/pkg/cfg"
	"github.com/coralproject/shelf/pkg/log"
	"github.com/coralproject/shelf/pkg/srv/mongo"

	"github.com/spf13/cobra"
)

var shelf = &cobra.Command{
	Use:   "shelf",
	Short: "Shelf provides the central cli housing of various cli tools that interface with the API",
}

var logdash bytes.Buffer

func displayLog() {
	defer logdash.Reset()
	logdash.WriteTo(os.Stdout)
}

func verbosity() {
	verbose, err := cfg.Int("LOGGING")
	if err != nil {
		//default is to logg out
		displayLog()
		return
	}

	if verbose == 1 {
		displayLog()
	}
}

func main() {
	defer verbosity()

	log.Init(&logdash, func() int { return log.DEV })

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
