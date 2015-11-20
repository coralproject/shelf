// This program provides the coral project shelf central CLI
// platform.
package main

import (
	"os"

	query "github.com/coralproject/shelf/cli/shelf/pkg/query/commands"
	user "github.com/coralproject/shelf/cli/shelf/pkg/user/commands"
	"github.com/coralproject/shelf/pkg/log"
	"github.com/spf13/cobra"
)

var shelf = &cobra.Command{
	Use:   "shelf",
	Short: "Shelf provides the central cli housing shelf's various cli tools",
}

func main() {
	log.Init(os.Stderr, func() int { return log.DEV })

	shelf.AddCommand(user.GetCommand(), query.GetCommand())
	shelf.Execute()
}
