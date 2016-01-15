package routes

import (
	"net/http"

	"github.com/coralproject/xenia/app/xenia/handlers"

	"github.com/ardanlabs/kit/web/app"
	"github.com/ardanlabs/kit/web/midware"
)

func init() {
	app.Init("XENIA")
}

//==============================================================================

// API returns a handler for a set of routes.
func API() http.Handler {
	a := app.New(midware.Auth)

	a.Handle("GET", "/1.0/query", handlers.Query.List)
	a.Handle("PUT", "/1.0/query", handlers.Query.Upsert)
	a.Handle("GET", "/1.0/query/:name", handlers.Query.Retrieve)
	a.Handle("DELETE", "/1.0/query/:name", handlers.Query.Delete)

	a.Handle("POST", "/1.0/exec", handlers.Query.ExecuteCustom)
	a.Handle("GET", "/1.0/exec/:name", handlers.Query.Execute)

	a.CORS()

	return a
}
