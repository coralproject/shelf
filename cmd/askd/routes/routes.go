package routes

import (
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/ardanlabs/kit/cfg"
	"github.com/ardanlabs/kit/log"
	"github.com/ardanlabs/kit/web"
	"github.com/coralproject/shelf/cmd/askd/handlers"
	"github.com/coralproject/shelf/internal/ask/form/submission"
	"github.com/coralproject/shelf/internal/platform/app"
	"github.com/coralproject/shelf/internal/platform/db"
	"github.com/coralproject/shelf/internal/platform/midware/auth"
	errorm "github.com/coralproject/shelf/internal/platform/midware/error"
	logm "github.com/coralproject/shelf/internal/platform/midware/log"
	"github.com/coralproject/shelf/internal/platform/midware/mongo"
)

const (

	// Namespace is the key that is the prefix for configuration in the
	// environment.
	Namespace = "ASK"

	// cfgMongoURI is the key for the URI to the MongoDB service.
	cfgMongoURI = "MONGO_URI"

	// cfgRecaptchaSecret is the key for the secret used for the recaptcha
	// service.
	cfgRecaptchaSecret = "RECAPTCHA_SECRET"

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

	// Ensure that the database indexes are setup on the underlying MongoDB
	// database.
	if err := ensureDBIndexes(mongoURI); err != nil {
		log.Error("startup", "Init", err, "Initializing DB Indexes")
		os.Exit(1)
	}

	w := web.New(logm.Midware, errorm.Midware, mongo.Midware(mongoURI))

	// Load in the recaptcha secret from the config.
	if recaptcha, err := cfg.String(cfgRecaptchaSecret); err == nil {
		w.Ctx["recaptcha"] = recaptcha
		log.Dev("startup", "Init", "Recaptcha Enabled")
	} else {
		log.Dev("startup", "Init", "Recaptcha Disabled")
	}

	log.Dev("startup", "Init", "Initalizing routes")

	routes(w)

	log.Dev("startup", "Init", "Initalizing CORS")
	w.Use(w.CORS())

	return w
}

func routes(w *web.Web) {

	// Create a new app group which will be for internal functions that may have
	// an optional layer of auth added to it.
	internal := w.Group()

	// Now we will load in the public key from the config. If found, we'll add a
	// middleware to all internal endpoints that will ensure that we validate the
	// requests coming in.

	publicKey, err := cfg.String(cfgAuthPublicKey)
	if err != nil || publicKey == "" {
		log.User("startup", "Init", "%s is missing, internal authentication is disabled", cfgAuthPublicKey)
	}

	// If the public key is provided then add the auth middleware or fail using
	// the provided public key.
	if publicKey != "" {
		log.Dev("startup", "Init", "Initializing Auth")

		// We are allowing the query string to act as the access token provider
		// because this service has endpoints that are accessed directly currently
		// and we need someway to authenticate to these endpoints.
		authmOpts := auth.MidwareOpts{
			AllowQueryString: true,
		}

		authm, err := auth.Midware(publicKey, authmOpts)
		if err != nil {
			log.Error("startup", "Init", err, "Initializing Auth")
			os.Exit(1)
		}

		// Apply the authentication middleware on top of the application as the
		// first middleware.
		internal.Use(authm)
	}

	// global
	internal.Handle("GET", "/v1/version", handlers.Version.List)

	// forms
	internal.Handle("POST", "/v1/form", handlers.Form.Upsert)
	internal.Handle("GET", "/v1/form", handlers.Form.List)
	internal.Handle("PUT", "/v1/form/:id", handlers.Form.Upsert)
	internal.Handle("PUT", "/v1/form/:id/status/:status", handlers.Form.UpdateStatus)
	internal.Handle("GET", "/v1/form/:id", handlers.Form.Retrieve)
	internal.Handle("DELETE", "/v1/form/:id", handlers.Form.Delete)

	// form submissions
	internal.Handle("GET", "/v1/form/:form_id/submission", handlers.FormSubmission.Search)
	internal.Handle("GET", "/v1/form/:form_id/submission/:id", handlers.FormSubmission.Retrieve)
	internal.Handle("PUT", "/v1/form/:form_id/submission/:id/status/:status", handlers.FormSubmission.UpdateStatus)
	internal.Handle("POST", "/v1/form/:form_id/submission/:id/flag/:flag", handlers.FormSubmission.AddFlag)
	internal.Handle("DELETE", "/v1/form/:form_id/submission/:id/flag/:flag", handlers.FormSubmission.RemoveFlag)
	internal.Handle("PUT", "/v1/form/:form_id/submission/:id/answer/:answer_id", handlers.FormSubmission.UpdateAnswer)
	internal.Handle("DELETE", "/v1/form/:form_id/submission/:id", handlers.FormSubmission.Delete)

	// temporal route to get CSV file - TO DO : move into a different service
	internal.Handle("GET", "/v1/form/:form_id/submission/export", handlers.FormSubmission.Download)

	// form form galleries
	internal.Handle("GET", "/v1/form/:form_id/gallery", handlers.FormGallery.RetrieveForForm)

	// form galleries
	internal.Handle("GET", "/v1/form_gallery/:id", handlers.FormGallery.Retrieve)
	internal.Handle("PUT", "/v1/form_gallery/:id", handlers.FormGallery.Update)
	internal.Handle("POST", "/v1/form_gallery/:id/submission/:submission_id/:answer_id", handlers.FormGallery.AddAnswer)
	internal.Handle("DELETE", "/v1/form_gallery/:id/submission/:submission_id/:answer_id", handlers.FormGallery.RemoveAnswer)

	// Create a new app group which will be for external functions that will need
	// to be publically exposed.
	external := w.Group()

	external.Handle("POST", "/v1/form/:form_id/submission", handlers.FormSubmission.Create)
}

func ensureDBIndexes(mongoURI *url.URL) error {
	mgoDB, err := db.NewMGO("startup", mongoURI.Path)
	if err != nil {
		return err
	}
	defer mgoDB.CloseMGO("startup")

	return submission.EnsureIndexes("startup", mgoDB)
}
