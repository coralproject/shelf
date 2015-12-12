// Xenia is a web service for handling query related calls.
// Use gin during development : // go get github.com/codegangsta/gin
// Run this command in the xenia folder: gin -p 5000 -a 4000 -i run
//
// To Run in DEBUG mode use this Header on all request in Dev
// Authorization: Basic NmQ3MmU2ZGQtOTNkMC00NDEzLTliNGMtODU0N
package main

import (
	"github.com/coralproject/shelf/app/xenia/routes"

	"github.com/ardanlabs/kit/web/app"
)

func main() {
	app.Run("XENIA_HOST", ":4000", routes.API())
}
