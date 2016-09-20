package routes

import (
	"net/http"

	"github.com/ardanlabs/kit/cfg"
	"github.com/ardanlabs/kit/log"
	"github.com/ardanlabs/kit/web/app"
	"github.com/coralproject/shelf/cmd/corald/fixtures"
	"github.com/coralproject/shelf/cmd/corald/handlers"
)

func init() {
	// Initialize the configuration and logging systems. Plus anything
	// else the web app layer needs.
	app.Init(cfg.EnvProvider{Namespace: "XENIA"})
}

// API returns a handler for a set of routes.
func API(testing ...bool) http.Handler {

	// TODO: If authentication is on then configure
	// it and provide proper middleware.

	a := app.New()

	log.Dev("startup", "Init", "Initalizing routes")
	routes(a)

	log.Dev("startup", "Init", "Initalizing CORS")
	a.CORS()

	return a
}

// routes manages the handling of the API endpoints.
func routes(a *app.App) {
	a.Handle("GET", "/1.0/version", handlers.Version.List)

	// TODO: For now these are sample routes.
	a.Handle("GET", "/1.0/form", handlers.Form.List)
	a.Handle("POST", "/1.0/form", fixtures.Handler("form", http.StatusCreated))
	a.Handle("GET", "/1.0/form/:form_id", fixtures.Handler("form", http.StatusOK))
	a.Handle("PUT", "/1.0/form/:form_id", fixtures.NoContent)
}
