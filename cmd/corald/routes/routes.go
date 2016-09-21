package routes

import (
	"net/http"
	"os"
	"time"

	"github.com/ardanlabs/kit/cfg"
	"github.com/ardanlabs/kit/db"
	"github.com/ardanlabs/kit/db/mongo"
	"github.com/ardanlabs/kit/log"
	"github.com/ardanlabs/kit/web/app"
	"github.com/coralproject/shelf/cmd/corald/fixtures"
	"github.com/coralproject/shelf/cmd/corald/handlers"
	"github.com/coralproject/shelf/cmd/corald/midware"
)

// Environmental variables.
const (
	cfgMongoHost     = "MONGO_HOST"
	cfgMongoAuthDB   = "MONGO_AUTHDB"
	cfgMongoDB       = "MONGO_DB"
	cfgMongoUser     = "MONGO_USER"
	cfgMongoPassword = "MONGO_PASS"
	cfgAnvilHost     = "ANVIL_HOST"
)

func init() {
	// Initialize the configuration and logging systems. Plus anything
	// else the web app layer needs.
	app.Init(cfg.EnvProvider{Namespace: "CORAL"})

	// Initialize MongoDB.
	if _, err := cfg.String(cfgMongoHost); err == nil {
		cfg := mongo.Config{
			Host:     cfg.MustString(cfgMongoHost),
			AuthDB:   cfg.MustString(cfgMongoAuthDB),
			DB:       cfg.MustString(cfgMongoDB),
			User:     cfg.MustString(cfgMongoUser),
			Password: cfg.MustString(cfgMongoPassword),
			Timeout:  25 * time.Second,
		}

		// The web framework middleware for Mongo is using the name of the
		// database as the name of the master session by convention. So use
		// cfg.DB as the second argument when creating the master session.
		if err := db.RegMasterSession("startup", cfg.DB, cfg); err != nil {
			log.Error("startup", "Init", err, "Initializing MongoDB")
			os.Exit(1)
		}
	}

}

// API returns a handler for a set of routes.
func API(testing ...bool) http.Handler {

	// TODO: If authentication is on then configure
	// it and provide proper middleware.

	a := app.New(midware.Mongo)

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
	a.Handle("GET", "/v1/item", handlers.Item.List)
	a.Handle("GET", "/v1/item/type/:type", handlers.Item.FilterByType)
	a.Handle("POST", "/v1/item", fixtures.Handler("item", http.StatusCreated))
	a.Handle("GET", "/v1/item/:item_id", fixtures.Handler("item", http.StatusOK))
	a.Handle("PUT", "/v1/item/:item_id", fixtures.NoContent)
}
