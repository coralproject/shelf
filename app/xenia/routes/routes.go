package routes

import (
	"net/http"

	"github.com/coralproject/xenia/app/xenia/handlers"

	"github.com/ardanlabs/kit/cfg"
	"github.com/ardanlabs/kit/web/app"
	"github.com/ardanlabs/kit/web/midware"
)

func init() {
	set := app.Settings{
		ConfigKey: "XENIA",
		UseMongo:  true,
	}

	app.Init(set)
}

//==============================================================================

// API returns a handler for a set of routes.
func API() http.Handler {
	var a *app.App

	// Check is authentication is turned off.
	auth, err := cfg.Bool("AUTH")
	if err != nil {

		// AUTH variable is not set.
		a = app.New(midware.Auth)
	} else {
		if auth {

			// AUTH variable is set to true
			a = app.New(midware.Auth)
		} else {

			// AUTH variable is set to false
			a = app.New()
		}
	}

	// Initialize the routes for the API.
	a.Handle("GET", "/1.0/query", handlers.Query.List)
	a.Handle("GET", "/1.0/query/:name", handlers.Query.Retrieve)
	a.Handle("GET", "/1.0/query/:name/exec", handlers.Query.Execute)

	return a
}
