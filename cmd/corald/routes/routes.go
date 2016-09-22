package routes

import (
	"net/http"
	"os"

	"github.com/ardanlabs/kit/cfg"
	"github.com/ardanlabs/kit/log"
	"github.com/ardanlabs/kit/web/app"
	"github.com/coralproject/shelf/cmd/corald/fixtures"
	"github.com/coralproject/shelf/cmd/corald/handlers"
	"github.com/coralproject/shelf/cmd/corald/midware"
)

func init() {
	// Initialize the configuration and logging systems. Plus anything
	// else the web app layer needs.
	app.Init(cfg.EnvProvider{Namespace: "XENIA"})
}

// API returns a handler for a set of routes.
func API(testing ...bool) http.Handler {
	auth, err := midware.Auth()
	if err != nil {
		log.Error("startup", "Init", err, "Initializing Auth")
		os.Exit(1)
	}

	a := app.New(auth)

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
