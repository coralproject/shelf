// This program provides a sample web service that implements a
// RESTFul CRUD API against a MongoDB database.
package main

import (
	"net/http"
	"os"
	"os/signal"

	"github.com/coralproject/shelf/app/xenia/routes"
	"github.com/coralproject/shelf/pkg/db/mongo"
	"github.com/coralproject/shelf/pkg/log"
)

func init() {
	// TODO: Need to read configuration.
	log.Init(os.Stdout, func() int { return log.DEV })
	mongo.InitMGO()
}

func main() {
	log.Dev("startup", "main", "Start")

	// Check the environment for a configured port value.
	// TODO: Need to read configuration.
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	// Create this goroutine to run the web server.
	go func() {
		log.Dev("listener", "main-func", "Listening on: http://localhost:%d", port)
		http.ListenAndServe(":"+port, routes.API())
	}()

	// Listen for an interrupt signal from the OS.
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	<-sigChan

	log.Dev("shutdown", "main", "Complete")
}
