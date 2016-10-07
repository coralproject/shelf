package routes

import (
	"net/http"
	"os"
	"time"

	"github.com/ardanlabs/kit/cfg"
	"github.com/ardanlabs/kit/log"
	"github.com/ardanlabs/kit/web"
	"github.com/coralproject/shelf/cmd/sponged/handlers"
	"github.com/coralproject/shelf/internal/platform/app"
	"github.com/coralproject/shelf/internal/platform/db"
	authm "github.com/coralproject/shelf/internal/platform/midware/auth"
	"github.com/coralproject/shelf/internal/platform/midware/cayley"
	errorm "github.com/coralproject/shelf/internal/platform/midware/error"
	logm "github.com/coralproject/shelf/internal/platform/midware/log"
	"github.com/coralproject/shelf/internal/platform/midware/mongo"
)

const (

	// Namespace is the key that is the prefix for configuration in the
	// environment.
	Namespace = "SPONGE"

	// cfgMongoURI is the key for the URI to the MongoDB service.
	cfgMongoURI = "MONGO_URI"

	// cfgAuthPublicKey is the key for the public key used for verifying the
	// inbound requests.
	cfgAuthPublicKey = "AUTH_PUBLIC_KEY"
)

func init() {

	// Initialize the configuration and logging systems. Plus anything
	// else the web app layer needs.
	app.Init(cfg.EnvProvider{Namespace: Namespace})
}

//==============================================================================

// API returns a handler for a set of routes.
func API() http.Handler {
	mongoURI := cfg.MustURL(cfgMongoURI)

	// The web framework middleware for Mongo is using the name of the
	// database as the name of the master session by convention. So use
	// cfg.DB as the second argument when creating the master session.
	if err := db.RegMasterSession("startup", mongoURI.Path, mongoURI.String(), 25*time.Second); err != nil {
		log.Error("startup", "Init", err, "Initializing MongoDB")
		os.Exit(1)
	}

	w := web.New(logm.Midware, errorm.Midware)

	publicKey, err := cfg.String(cfgAuthPublicKey)
	if err != nil || publicKey == "" {
		log.User("startup", "Init", "%s is missing, internal authentication is disabled", cfgAuthPublicKey)
	}

	// If the public key is provided then add the auth middleware or fail using
	// the provided public key.
	if publicKey != "" {
		log.Dev("startup", "Init", "Initializing Auth")

		authm, err := authm.Midware(publicKey, authm.MidwareOpts{})
		if err != nil {
			log.Error("startup", "Init", err, "Initializing Auth")
			os.Exit(1)
		}

		// Apply the authentication middleware on top of the application as the
		// first middleware.
		w.Use(authm)
	}

	// Add the Mongo and Cayley middlewares possibly after the auth middleware.
	w.Use(mongo.Midware(mongoURI), cayley.Midware(mongoURI))

	log.Dev("startup", "Init", "Initalizing CORS")
	w.Use(w.CORS())

	log.Dev("startup", "Init", "Initalizing routes")
	routes(w)

	return w
}

// routes manages the handling of the API endpoints.
func routes(w *web.Web) {
	w.Handle("GET", "/v1/version", handlers.Version.List)

	w.Handle("GET", "/v1/item/:id", handlers.Item.Retrieve)
	w.Handle("PUT", "/v1/item", handlers.Item.Import)
	w.Handle("POST", "/v1/item", handlers.Item.Import)
	w.Handle("DELETE", "/v1/item/:id", handlers.Item.Remove)

	w.Handle("GET", "/v1/pattern", handlers.Pattern.List)
	w.Handle("PUT", "/v1/pattern", handlers.Pattern.Upsert)
	w.Handle("GET", "/v1/pattern/:type", handlers.Pattern.Retrieve)
	w.Handle("DELETE", "/v1/pattern/:type", handlers.Pattern.Delete)

	w.Handle("POST", "/v1/data/:type", handlers.Data.Import)
}
