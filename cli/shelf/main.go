// This program provides the coral project shelf central CLI
// platform.
package main

import (
	"os"

	"github.com/coralproject/shelf/cli/shelf/cmdquery"
	"github.com/coralproject/shelf/cli/shelf/cmduser"
	"github.com/coralproject/shelf/pkg/log"
	"github.com/spf13/cobra"
)

var shelf = &cobra.Command{
	Use:   "shelf",
	Short: "Shelf provides the central cli housing shelf's various cli tools",
}

func main() {
	log.Init(os.Stderr, func() int { return log.DEV })

	shelf.AddCommand(cmduser.GetCommands(), cmdquery.GetCommands())
	shelf.Execute()
}
