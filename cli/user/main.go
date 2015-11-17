// This program provides support for managing users in the coral project
// platform.
package main

import (
	"github.com/coralproject/shelf/cli/user/commands"
	"github.com/coralproject/shelf/cli/user/db"
)

func main() {
	db.InitMGO()

	commands.Run()
}
