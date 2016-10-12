package routes

import (
	"net/http"
	"os"

	"github.com/ardanlabs/kit/cfg"
	"github.com/ardanlabs/kit/log"
	"github.com/ardanlabs/kit/web"
	"github.com/coralproject/shelf/cmd/corald/fixtures"
	"github.com/coralproject/shelf/cmd/corald/handlers"
	"github.com/coralproject/shelf/internal/platform/app"
	"github.com/coralproject/shelf/internal/platform/auth"
	authm "github.com/coralproject/shelf/internal/platform/midware/auth"
	errorm "github.com/coralproject/shelf/internal/platform/midware/error"
	logm "github.com/coralproject/shelf/internal/platform/midware/log"
)

const (

	// Namespace is the key that is the prefix for configuration in the
	// environment.
	Namespace = "CORAL"

	// cfgSpongdURL is the config key for the url to the sponged service.
	cfgSpongdURL = "SPONGED_URL"

	// cfgXeniadURL is the config key for the url to the xeniad service.
	cfgXeniadURL = "XENIAD_URL"

	// cfgAuthPublicKey is the key for the public key used for verifying the
	// inbound requests.
	cfgAuthPublicKey = "AUTH_PUBLIC_KEY"

	// cfgPlatformPrivateKey is the private key used to sign new requests to the
	// downstream service layer.
	cfgPlatformPrivateKey = "PLATFORM_PRIVATE_KEY"

	// cfgEnableCORS is set the key to the state for CORS on the service.
	cfgEnableCORS = "ENABLE_CORS"
)

func init() {

	// Initialize the configuration and logging systems. Plus anything
	// else the web app layer needs.
	app.Init(cfg.EnvProvider{Namespace: Namespace})
}

// API returns a handler for a set of routes.
func API() http.Handler {
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

	platformPrivateKey, err := cfg.String(cfgPlatformPrivateKey)
	if err != nil || platformPrivateKey == "" {
		log.User("startup", "Init", "%s is missing, downstream platform authentication is disabled", cfgPlatformPrivateKey)
	}

	// If the platformPrivateKey is provided, then we should generate the token
	// signing function to be used when composing requests down to the platform.
	if platformPrivateKey != "" {
		log.Dev("startup", "Init", "Initializing Downstream Platform Auth")

		signer, err := auth.NewSigner(platformPrivateKey)
		if err != nil {
			log.Error("startup", "Init", err, "Initializing Downstream Platform Auth")
			os.Exit(1)
		}

		// Requests can now be signed with the given signer function which we will
		// save on the application wide context. In the event that a function
		// requires a call down to a downstream platform, we will include a signed
		// header using the signer function here.
		w.Ctx["signer"] = signer
	}

	if cors, err := cfg.Bool(cfgEnableCORS); err == nil && cors {
		log.Dev("startup", "Init", "Initializing CORS : CORS Enabled")
		w.Use(w.CORS())
	} else {
		log.Dev("startup", "Init", "CORS Disabled")
	}

	log.Dev("startup", "Init", "Initalizing routes")
	routes(w)

	return w
}

// routes manages the handling of the API endpoints.
func routes(w *web.Web) {
	w.Handle("GET", "/v1/version", handlers.Version.List)

	spongedURL := cfg.MustURL(cfgSpongdURL).String()
	xeniadURL := cfg.MustURL(cfgXeniadURL).String()

	// CRU- for forms
	w.Handle("GET", "/v1/form", fixtures.Handler("forms/forms", http.StatusOK))
	w.Handle("POST", "/v1/form", fixtures.Handler("forms/form", http.StatusCreated))
	w.Handle("GET", "/v1/form/:form_id", fixtures.Handler("forms/form", http.StatusOK))
	w.Handle("PUT", "/v1/form/:form_id", fixtures.NoContent)

	// Execute the :query_set on the view :view_name on this :item_key.
	w.Handle("GET", "/v1/exec/:query_set/view/:view_name/:item_key",
		handlers.Proxy(xeniadURL,
			func(c *web.Context) string {
				return "/v1/exec/" + c.Params["query_set"] + "/view/" + c.Params["view_name"] + "/" + c.Params["item_key"]
			}))

	// Get all the items from the view :view_name on this :item_key.
	w.Handle("POST", "/v1/exec/view/:view_name/:item_key",
		handlers.Proxy(xeniadURL,
			func(c *web.Context) string {
				return "/v1/exec/view/" + c.Params["view_name"] + "/" + c.Params["item_key"]
			}))

	// Execute xenia queries directly.
	w.Handle("GET", "/v1/exec/:query_set",
		handlers.Proxy(xeniadURL, func(c *web.Context) string { return "/v1/exec/" + c.Params["query_set"] }))

	// Send a new query to xenia. ********* TEMPORAL *********
	w.Handle("PUT", "/v1/query",
		handlers.Proxy(xeniadURL, func(c *web.Context) string { return "/v1/query" }))

	// Execute a custom xenia query. ********* TEMPORAL *********
	w.Handle("POST", "/v1/exec",
		handlers.Proxy(xeniadURL, func(c *web.Context) string { return "/v1/exec" }))

	// Create or removes Actions.
	w.Handle("POST", "/v1/action/:action/on/item/:item_key", fixtures.Handler("items/actionWFlag", http.StatusOK))    //handlers.Action.Create)
	w.Handle("DELETE", "/v1/action/:action/on/item/:item_key", fixtures.Handler("items/actionWoFlag", http.StatusOK)) //handlers.Action.Remove)

	// Save or Update Items
	w.Handle("PUT", "/v1/item",
		handlers.Proxy(spongedURL, func(c *web.Context) string { return "/v1/item" }))

	w.Handle("POST", "/v1/item",
		handlers.Proxy(spongedURL, func(c *web.Context) string { return "/v1/item" }))
}
