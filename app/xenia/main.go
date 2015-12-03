// Xenia is a web service for handling query related calls.
// Use gin during development : // go get github.com/codegangsta/gin
// Run this command in the xenia folder: gin -p 5000 -a 4000 -i run
package main

import (
	"net/http"
	"os"
	"os/signal"

	"github.com/coralproject/shelf/app/xenia/routes"
	"github.com/coralproject/shelf/pkg/cfg"
	"github.com/coralproject/shelf/pkg/db/mongo"
	"github.com/coralproject/shelf/pkg/log"
)

func init() {
	logLevel := func() int {
		ll, err := cfg.Int("LOGGING_LEVEL")
		if err != nil {
			return log.USER
		}
		return ll
	}

	log.Init(os.Stderr, logLevel)

	if err := cfg.Init("SHELF"); err != nil {
		log.Error("startup", "init", err, "Initializing config")
		os.Exit(1)
	}

	err := mongo.InitMGO()
	if err != nil {
		log.Error("startup", "init", err, "Initializing MongoDB")
		os.Exit(1)
	}
}

func main() {
	log.Dev("startup", "main", "Start")

	// Check for a configured host value.
	host, err := cfg.String("XENIA_HOST")
	if err != nil {
		host = ":4000"
	}

	// Create this goroutine to run the web server.
	go func() {
		log.Dev("listener", "main", "Listening on: %s", host)
		http.ListenAndServe(host, routes.API())
	}()

	// Listen for an interrupt signal from the OS.
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	<-sigChan

	log.Dev("shutdown", "main", "Complete")
}
