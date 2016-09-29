package routes

import (
	"net/http"
	"os"

	"github.com/ardanlabs/kit/cfg"
	"github.com/ardanlabs/kit/log"
	"github.com/ardanlabs/kit/web/app"
	"github.com/coralproject/shelf/cmd/corald/fixtures"
	"github.com/coralproject/shelf/cmd/corald/handlers"
	"github.com/coralproject/shelf/internal/platform/auth"
)

const (

	// cfgSpongdURL is the config key for the url to the sponged service.
	cfgSpongdURL = "SPONGED_URL"

	// cfgXeniadURL is the config key for the url to the xeniad service.
	cfgXeniadURL = "XENIAD_URL"

	// cfgAuthPublicKey is the config key for the public key used to verify
	// tokens.
	cfgAuthPublicKey = "AUTH_PUBLIC_KEY"

	// cfgPlatformPrivateKey is the private key used to sign new requests to the
	// downstream service layer.
	cfgPlatformPrivateKey = "PLATFORM_PRIVATE_KEY"
)

func init() {

	// Initialize the configuration and logging systems. Plus anything
	// else the web app layer needs.
	app.Init(cfg.EnvProvider{Namespace: "CORAL"})
}

// API returns a handler for a set of routes.
func API(testing ...bool) http.Handler {
	a := app.New()

	publicKey, err := cfg.String(cfgAuthPublicKey)
	if err != nil {
		log.User("startup", "Init", "CORAL_%s is missing, internal authentication is disabled", cfgAuthPublicKey)
	}

	// If the public key is provided then add the auth middleware or fail using
	// the provided public key.
	if publicKey != "" {
		log.Dev("startup", "Init", "Initializing Auth")

		authm, err := auth.Midware(publicKey)
		if err != nil {
			log.Error("startup", "Init", err, "Initializing Auth")
			os.Exit(1)
		}

		// Apply the authentication middleware on top of the application as the
		// first middleware.
		a.Use(authm)
	}

	platformPrivateKey, err := cfg.String(cfgPlatformPrivateKey)
	if err != nil {
		log.User("startup", "Init", "CORAL_%s is missing, downstream platform authentication is disabled", cfgPlatformPrivateKey)
	}

	// If the platformPrivateKey is provided, then we should generate the token
	// signing function to be used when composing requests down to the platform.
	if platformPrivateKey != "" {
		log.Dev("startup", "Init", "Initializing Downstream Platform Auth")

		signer, err := auth.NewSigner(platformPrivateKey)
		if err != nil {
			log.Error("startup", "Init", err, "Initializing Auth")
			os.Exit(1)
		}

		// Requests can now be signed with the given signer function which we will
		// save on the application wide context. In the event that a function
		// requires a call down to a downstream platform, we will include a signed
		// header using the signer function here.
		a.Ctx["signer"] = signer
	}

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

	// Execute the :query_set on the view :view_name on this :item_key.
	a.Handle("GET", "/v1/exec/:query_set/view/:view_name/:item_key",
		handlers.Proxy(xeniadURL,
			func(c *app.Context) string {
				return "/v1/exec/" + c.Params["query_set"]
			}))

	// Execute xenia queries directly.
	a.Handle("GET", "/v1/exec/:query_set", handlers.Proxy(xeniadURL, nil))

	// Send a new query to xenia.
	a.Handle("PUT", "/v1/query", handlers.Proxy(xeniadURL, nil))

	a.Handle("PUT", "/v1/item",
		handlers.Proxy(spongedURL,
			func(c *app.Context) string {
				return "/v1/item"
			}))

	a.Handle("POST", "/v1/item",
		handlers.Proxy(spongedURL,
			func(c *app.Context) string {
				return "/v1/item"
			}))
}
