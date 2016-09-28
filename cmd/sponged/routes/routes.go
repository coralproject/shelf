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
	"github.com/coralproject/shelf/cmd/sponged/handlers"
	"github.com/coralproject/shelf/cmd/sponged/midware"
	"github.com/coralproject/shelf/internal/platform/auth"
)

// Environmental variables.
const (
	cfgMongoHost     = "MONGO_HOST"
	cfgMongoAuthDB   = "MONGO_AUTHDB"
	cfgMongoDB       = "MONGO_DB"
	cfgMongoUser     = "MONGO_USER"
	cfgMongoPassword = "MONGO_PASS"
	cfgAuthPublicKey = "AUTH_PUBLIC_KEY"
)

func init() {
	// Initialize the configuration and logging systems. Plus anything
	// else the web app layer needs.
	app.Init(cfg.EnvProvider{Namespace: "SPONGE"})

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

//==============================================================================

// API returns a handler for a set of routes.
func API() http.Handler {
	a := app.New()

	publicKey, err := cfg.String(cfgAuthPublicKey)
	if err != nil {
		log.User("startup", "Init", "SPONGE_%s is missing, internal authentication is disabled", cfgAuthPublicKey)
	}

	// If the public key is provided then add the auth middleware or fail using
	// the provided public key.
	if publicKey != "" {
		authm, err := auth.Midware(publicKey)
		if err != nil {
			log.Error("startup", "Init", err, "Initializing Auth")
			os.Exit(1)
		}

		// Apply the authentication middleware on top of the application as the
		// first middleware.
		a.Use(authm)
	}

	// Add the Mongo and Cayley middlewares possibly after the auth middleware.
	a.Use(midware.Mongo, midware.Cayley)

	log.Dev("startup", "Init", "Initalizing routes")
	routes(a)

	log.Dev("startup", "Init", "Initalizing CORS")
	a.CORS()

	return a
}

// routes manages the handling of the API endpoints.
func routes(a *app.App) {
	a.Handle("GET", "/1.0/version", handlers.Version.List)

	a.Handle("GET", "/1.0/item/:id", handlers.Item.Retrieve)
	a.Handle("PUT", "/1.0/item", handlers.Item.Upsert)
	a.Handle("DELETE", "/1.0/item/:id", handlers.Item.Delete)

	a.Handle("POST", "/1.0/data/:type", handlers.Data.Upsert)
}
