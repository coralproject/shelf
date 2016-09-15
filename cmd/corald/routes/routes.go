package routes

import (
	"net/http"

	"github.com/ardanlabs/kit/log"
	"github.com/ardanlabs/kit/web/app"
	"github.com/coralproject/shelf/cmd/corald/handlers"
)

// API returns a handler for a set of routes.
func API() http.Handler {
	// Create a new App.
	a := app.New()

	log.Dev("startup", "Init", "Initalizing routes")

	// Add the routes to the API.
	setupRoutes(a)

	log.Dev("startup", "Init", "Initalizing CORS")

	// Enable CORS on the endpoints.
	a.CORS()

	return a
}

// setupRoutes adds all the routes that the corald command will serve.
func setupRoutes(a *app.App) {

	//----------------------------------------------------------------------
	// Implemented handlers.
	//

	a.Handle("GET", "/version", handlers.Version.List)

	//----------------------------------------------------------------------
	// Fixture handlers.
	//

	a.Handle("GET", "/1.0/form", handlers.Fixture("forms", http.StatusOK))
	a.Handle("POST", "/1.0/form", handlers.Fixture("form", http.StatusCreated))
	a.Handle("GET", "/1.0/form/:form_id", handlers.Fixture("form", http.StatusOK))
	a.Handle("PUT", "/1.0/form/:form_id", handlers.NoContent)
}
