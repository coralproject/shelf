package routes

import (
	"net/http"
	"os"
	"time"

	anvil "github.com/anvilresearch/go-anvil"

	"github.com/ardanlabs/kit/cfg"
	"github.com/ardanlabs/kit/db"
	"github.com/ardanlabs/kit/db/mongo"
	"github.com/ardanlabs/kit/log"
	"github.com/ardanlabs/kit/web/app"

	"github.com/coralproject/xenia/cmd/sponged/handlers"
	"github.com/coralproject/xenia/cmd/sponged/midware"

	"github.com/coralproject/xenia/internal/item"
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

	// Initialize the item / relationship system
	err := item.Initialize()
	if err != nil {
		log.Dev("startup", "Main", "Error Initializing Items: %s", err)

	}

}

//==============================================================================

// API returns a handler for a set of routes.
func API(testing ...bool) http.Handler {

	// If authentication is on then configure Anvil.
	var anv *anvil.Anvil
	if url, err := cfg.String(cfgAnvilHost); err == nil {

		log.Dev("startup", "Init", "Initalizing Anvil")
		anv, err = anvil.New(url)
		if err != nil {
			log.Error("startup", "Init", err, "Initializing Anvil: %s", url)
			os.Exit(1)
		}
	}

	a := app.New(midware.Mongo, midware.Auth)
	a.Ctx["anvil"] = anv

	log.Dev("startup", "Init", "Initalizing routes")
	routes(a)

	log.Dev("startup", "Init", "Initalizing CORS")
	a.CORS()

	return a
}

// routes manages the handling of the API endpoints.
func routes(a *app.App) {

	a.Handle("POST", "/1.0/item/type/:type", handlers.Item.Upsert)
	a.Handle("PUT", "/1.0/item/type/:type", handlers.Item.Upsert)

	a.Handle("GET", "/1.0/item/:id", handlers.Item.Get)
	a.Handle("GET", "/1.0/item/:id/rels", handlers.Item.GetRels)

	a.Handle("GET", "/1.0/types", handlers.Item.Types)
	a.Handle("GET", "/1.0/version", handlers.Version.List)

}
