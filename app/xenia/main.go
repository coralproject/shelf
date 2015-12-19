// Xenia is a web service for handling query related calls.
// Use gin during development : // go get github.com/codegangsta/gin
// Run this command in the xenia folder: gin -p 5000 -a 4000 -i run
package main

import (
	"github.com/coralproject/xenia/app/xenia/routes"

	"github.com/ardanlabs/kit/web/app"
)

// The routes package initializes xenia. It has been placed here to help
// with the initialization of tests.

func main() {
	app.Run(":4000", routes.API())
}
