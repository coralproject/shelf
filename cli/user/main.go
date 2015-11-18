// This program provides support for managing users in the coral project
// platform.
package main

import (
	"os"

	"github.com/coralproject/shelf/cli/user/commands"
	"github.com/coralproject/shelf/pkg/log"
)

func main() {
	log.Init(os.Stderr, func() int { return log.DEV })

	commands.Run()
}
