package routes

import (
	"net/http"

	"github.com/ardanlabs/kit/cfg"
	"github.com/ardanlabs/kit/log"
	"github.com/ardanlabs/kit/web/app"
	"github.com/coralproject/shelf/cmd/corald/fixtures"
	"github.com/coralproject/shelf/cmd/corald/handlers"
)

// Environmental variables.
const (
	cfgSpongeHost = "MONGO_HOST"
	cfgXeniaHost  = "MONGO_AUTHD"
	cfgAnvilHost  = "ANVIL_HOST"
)

func init() {
	// Initialize the configuration and logging systems. Plus anything
	// else the web app layer needs.
	app.Init(cfg.EnvProvider{Namespace: "CORAL"})
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
	a.Handle("GET", "/v1/version", handlers.Version.List)

	// TODO: For now these are sample routes.
	a.Handle("GET", "/v1/form", handlers.Form.List)
	a.Handle("POST", "/v1/form", fixtures.Handler("form", http.StatusCreated))
	a.Handle("GET", "/v1/form/:form_id", fixtures.Handler("form", http.StatusOK))
	a.Handle("PUT", "/v1/form/:form_id", fixtures.NoContent)

	// Items
	a.Handle("GET", "/v1/item/:view_name/:item_key/:query_set", handlers.Item.Retrieve)
	a.Handle("POST", "/v1/item", fixtures.Handler("itemid", http.StatusCreated))
	a.Handle("PUT", "/v1/item/:item_id", fixtures.NoContent)
}
