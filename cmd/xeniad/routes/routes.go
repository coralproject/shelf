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
	"github.com/coralproject/shelf/cmd/xeniad/handlers"
	"github.com/coralproject/shelf/cmd/xeniad/midware"
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
	app.Init(cfg.EnvProvider{Namespace: "XENIA"})

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
func API(testing ...bool) http.Handler {
	a := app.New()

	publicKey, err := cfg.String(cfgAuthPublicKey)
	if err != nil {
		log.User("startup", "Init", "XENIA_%s is missing, internal authentication is disabled", cfgAuthPublicKey)
	}

	// If the public key is provided then add the auth middleware or fail using
	// the provided public key.
	if publicKey != "" {
		log.Dev("startup", "Init", "Initializing Auth")

		authm, err := auth.Midware(publicKey, auth.MidwareOpts{})
		if err != nil {
			log.Error("startup", "Init", err, "Initializing Auth")
			os.Exit(1)
		}

		// Apply the authentication middleware on top of the application as the
		// first middleware.
		a.Use(authm)
	}

	a.Use(midware.Mongo)

	log.Dev("startup", "Init", "Initalizing routes")
	routes(a)

	log.Dev("startup", "Init", "Initalizing CORS")
	a.CORS()

	return a
}

// routes manages the handling of the API endpoints.
func routes(a *app.App) {

	a.Handle("GET", "/v1/version", handlers.Version.List)

	a.Handle("GET", "/v1/script", handlers.Script.List)
	a.Handle("PUT", "/v1/script", handlers.Script.Upsert)
	a.Handle("GET", "/v1/script/:name", handlers.Script.Retrieve)
	a.Handle("DELETE", "/v1/script/:name", handlers.Script.Delete)

	a.Handle("GET", "/v1/query", handlers.Query.List)
	a.Handle("PUT", "/v1/query", handlers.Query.Upsert)
	a.Handle("GET", "/v1/query/:name", handlers.Query.Retrieve)
	a.Handle("DELETE", "/v1/query/:name", handlers.Query.Delete)

	a.Handle("PUT", "/v1/index/:name", handlers.Query.EnsureIndexes)

	a.Handle("GET", "/v1/regex", handlers.Regex.List)
	a.Handle("PUT", "/v1/regex", handlers.Regex.Upsert)
	a.Handle("GET", "/v1/regex/:name", handlers.Regex.Retrieve)
	a.Handle("DELETE", "/v1/regex/:name", handlers.Regex.Delete)

	a.Handle("GET", "/v1/mask", handlers.Mask.List)
	a.Handle("PUT", "/v1/mask", handlers.Mask.Upsert)
	a.Handle("GET", "/v1/mask/:collection/:field", handlers.Mask.Retrieve)
	a.Handle("GET", "/v1/mask/:collection", handlers.Mask.Retrieve)
	a.Handle("DELETE", "/v1/mask/:collection/:field", handlers.Mask.Delete)

	a.Handle("POST", "/v1/exec", handlers.Exec.Custom)
	a.Handle("GET", "/v1/exec/:name", handlers.Exec.Name)
	a.Handle("GET", "/v1/exec/:name/view/:view/:item", handlers.Exec.NameOnView, midware.Cayley)
	a.Handle("POST", "/v1/exec/view/:view/:item", handlers.Exec.CustomOnView, midware.Cayley)

	a.Handle("GET", "/v1/relationship", handlers.Relationship.List)
	a.Handle("PUT", "/v1/relationship", handlers.Relationship.Upsert)
	a.Handle("GET", "/v1/relationship/:predicate", handlers.Relationship.Retrieve)
	a.Handle("DELETE", "/v1/relationship/:predicate", handlers.Relationship.Delete)

	a.Handle("GET", "/v1/view", handlers.View.List)
	a.Handle("PUT", "/v1/view", handlers.View.Upsert)
	a.Handle("GET", "/v1/view/:name", handlers.View.Retrieve)
	a.Handle("DELETE", "/v1/view/:name", handlers.View.Delete)

	a.Handle("GET", "/v1/pattern", handlers.Pattern.List)
	a.Handle("PUT", "/v1/pattern", handlers.Pattern.Upsert)
	a.Handle("GET", "/v1/pattern/:type", handlers.Pattern.Retrieve)
	a.Handle("DELETE", "/v1/pattern/:type", handlers.Pattern.Delete)

}
