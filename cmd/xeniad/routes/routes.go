package routes

import (
	"net/http"
	"os"
	"time"

	"github.com/ardanlabs/kit/cfg"
	"github.com/ardanlabs/kit/log"
	"github.com/ardanlabs/kit/web"
	"github.com/coralproject/shelf/cmd/xeniad/handlers"
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
	Namespace = "XENIA"

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

	w.Use(mongo.Midware(mongoURI))

	log.Dev("startup", "Init", "Initalizing CORS")
	w.Use(w.CORS())

	log.Dev("startup", "Init", "Initalizing routes")
	routes(w)

	return w
}

// routes manages the handling of the API endpoints.
func routes(w *web.Web) {
	w.Handle("GET", "/v1/version", handlers.Version.List)

	w.Handle("GET", "/v1/script", handlers.Script.List)
	w.Handle("PUT", "/v1/script", handlers.Script.Upsert)
	w.Handle("GET", "/v1/script/:name", handlers.Script.Retrieve)
	w.Handle("DELETE", "/v1/script/:name", handlers.Script.Delete)

	w.Handle("GET", "/v1/query", handlers.Query.List)
	w.Handle("PUT", "/v1/query", handlers.Query.Upsert)
	w.Handle("GET", "/v1/query/:name", handlers.Query.Retrieve)
	w.Handle("DELETE", "/v1/query/:name", handlers.Query.Delete)

	w.Handle("PUT", "/v1/index/:name", handlers.Query.EnsureIndexes)

	w.Handle("GET", "/v1/regex", handlers.Regex.List)
	w.Handle("PUT", "/v1/regex", handlers.Regex.Upsert)
	w.Handle("GET", "/v1/regex/:name", handlers.Regex.Retrieve)
	w.Handle("DELETE", "/v1/regex/:name", handlers.Regex.Delete)

	w.Handle("GET", "/v1/mask", handlers.Mask.List)
	w.Handle("PUT", "/v1/mask", handlers.Mask.Upsert)
	w.Handle("GET", "/v1/mask/:collection/:field", handlers.Mask.Retrieve)
	w.Handle("GET", "/v1/mask/:collection", handlers.Mask.Retrieve)
	w.Handle("DELETE", "/v1/mask/:collection/:field", handlers.Mask.Delete)

	w.Handle("POST", "/v1/exec", handlers.Exec.Custom)
	w.Handle("GET", "/v1/exec/:name", handlers.Exec.Name)

	// Create the Cayley middleware which will only be binded to specific
	// endpoints.
	cayleym := cayley.Midware(cfg.MustURL(cfgMongoURI))

	// These endpoints require Cayley, we will add the middleware onto the routes.
	w.Handle("GET", "/v1/exec/:name/view/:view/:item", handlers.Exec.NameOnView, cayleym)
	w.Handle("POST", "/v1/exec/view/:view/:item", handlers.Exec.CustomOnView, cayleym)

	w.Handle("GET", "/v1/relationship", handlers.Relationship.List)
	w.Handle("PUT", "/v1/relationship", handlers.Relationship.Upsert)
	w.Handle("GET", "/v1/relationship/:predicate", handlers.Relationship.Retrieve)
	w.Handle("DELETE", "/v1/relationship/:predicate", handlers.Relationship.Delete)

	w.Handle("GET", "/v1/view", handlers.View.List)
	w.Handle("PUT", "/v1/view", handlers.View.Upsert)
	w.Handle("GET", "/v1/view/:name", handlers.View.Retrieve)
	w.Handle("DELETE", "/v1/view/:name", handlers.View.Delete)
}
