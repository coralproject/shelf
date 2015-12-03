package routes

import (
	"net/http"

	"github.com/coralproject/shelf/app/xenia/app"
	"github.com/coralproject/shelf/app/xenia/handlers"
	"github.com/coralproject/shelf/app/xenia/midware"
)

// API returns a handler for a set of routes.
func API() http.Handler {
	a := app.New(midware.Auth)

	// Initialize the routes for the API.
	a.Handle("GET", "/1.0/query/names", handlers.Query.List)
	// a.Handle("GET", "/1.0/:name", handlers.Query.Retrieve)
	// a.Handle("GET", "/1.0/run", handlers.Query.Execute)
	// a.Handle("POST", "/1.0/run/custom", handlers.Query.ExecuteCustom)

	return a
}
