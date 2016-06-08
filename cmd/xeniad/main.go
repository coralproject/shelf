// Xenia is a web service for handling query related calls.
package main

import (
	"runtime"

	"github.com/coralproject/xenia/cmd/xeniad/handlers"
	"github.com/coralproject/xenia/cmd/xeniad/routes"

	"github.com/ardanlabs/kit/log"
	"github.com/ardanlabs/kit/web/app"
)

// These are set by the makefile with:
// go build -ldflags "-X main.GitVersion=8e830ff -X main.GitRevision=123123123123 -X main.BuildDate=2016-01-25"
var (
	GitRevision = "<unknown>"
	GitVersion  = "<unknown>"
	BuildDate   = "<unknown>"
	IntVersion  = "201606081030"
)

func main() {
	log.User("startup", "Init", "Revision     : %q", GitRevision)
	log.User("startup", "Init", "Version      : %q", GitVersion)
	log.User("startup", "Init", "Build Date   : %q", BuildDate)
	log.User("startup", "Init", "Int Version  : %q", IntVersion)
	log.User("startup", "Init", "Go Version   : %q", runtime.Version())
	log.User("startup", "Init", "Go Compiler  : %q", runtime.Compiler)
	log.User("startup", "Init", "Go ARCH      : %q", runtime.GOARCH)
	log.User("startup", "Init", "Go OS        : %q", runtime.GOOS)

	handlers.Version.GitRevision = GitRevision
	handlers.Version.GitVersion = GitVersion
	handlers.Version.BuildDate = BuildDate
	handlers.Version.IntVersion = IntVersion

	app.Run(":4000", routes.API())
}
