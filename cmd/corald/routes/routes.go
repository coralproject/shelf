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

const (

	// cfgSpongdURL is the config key for the url to the sponged service.
	cfgSpongdURL = "SPONGED_URL"

	// cfgXeniadURL is the config key for the url to the xeniad service.
	cfgXeniadURL = "XENIAD_URL"
)

func init() {

	// Initialize the configuration and logging systems. Plus anything
	// else the web app layer needs.
	app.Init(cfg.EnvProvider{Namespace: "CORAL"})
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
	a.Handle("GET", "/v1/version", handlers.Version.List)

	spongedURL := cfg.MustURL(cfgSpongdURL).String()
	xeniadURL := cfg.MustURL(cfgXeniadURL).String()

	a.Handle("GET", "/v1/form", fixtures.Handler("forms/forms", http.StatusOK))
	a.Handle("POST", "/v1/form", fixtures.Handler("forms/form", http.StatusCreated))
	a.Handle("GET", "/v1/form/:form_id", fixtures.Handler("forms/form", http.StatusOK))
	a.Handle("PUT", "/v1/form/:form_id", fixtures.NoContent)

	a.Handle("GET", "/v1/item/:view_name/:item_key/:query_set",
		handlers.Proxy(xeniadURL,
			func(c *app.Context) string {
				return "/v1/exec/" + c.Params["query_set"]
			}))

	a.Handle("POST", "/v1/item", fixtures.Handler("items/itemid", http.StatusCreated))
	a.Handle("PUT", "/v1/item", handlers.Proxy(spongedURL, nil))
}
