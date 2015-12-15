// This program provides the coral project xenia central CLI platform.
package main

import (
	"os"

	"github.com/coralproject/xenia/cmd/xenia/cmdquery"

	"github.com/ardanlabs/kit/cfg"
	"github.com/ardanlabs/kit/cmd/kit/cmdauth"
	"github.com/ardanlabs/kit/cmd/kit/cmddb"
	"github.com/ardanlabs/kit/db/mongo"
	"github.com/ardanlabs/kit/log"

	"github.com/spf13/cobra"
)

var xenia = &cobra.Command{
	Use:   "xenia",
	Short: "Xenia provides the central cli housing of various cli tools that interface with the API",
}

func main() {
	logLevel := func() int {
		ll, err := cfg.Int("LOGGING_LEVEL")
		if err != nil {
			return log.NONE
		}
		return ll
	}

	log.Init(os.Stderr, logLevel)

	if err := cfg.Init("XENIA"); err != nil {
		xenia.Println("Unable to initialize configuration")
		os.Exit(1)
	}

	err := mongo.Init()
	if err != nil {
		xenia.Println("Unable to initialize MongoDB")
		os.Exit(1)
	}

	xenia.AddCommand(cmdauth.GetCommands(), cmddb.GetCommands(), cmdquery.GetCommands())
	xenia.Execute()
}
