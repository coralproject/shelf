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
	"github.com/coralproject/shelf/cmd/askd/handlers"
	"github.com/coralproject/shelf/cmd/askd/midware"
	"github.com/coralproject/shelf/internal/ask/form/submission"
	"github.com/coralproject/shelf/internal/platform/auth"
)

// Environmental variables.
const (
	cfgMongoHost       = "MONGO_HOST"
	cfgMongoAuthDB     = "MONGO_AUTHDB"
	cfgMongoDB         = "MONGO_DB"
	cfgMongoUser       = "MONGO_USER"
	cfgMongoPassword   = "MONGO_PASS"
	cfgRecaptchaSecret = "RECAPTCHA_SECRET"
	cfgAuthPublicKey   = "AUTH_PUBLIC_KEY"
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
func API() http.Handler {
	if err := ensureDBIndexes(); err != nil {
		log.Error("startup", "Init", err, "Initializing DB Indexes")
		os.Exit(1)
	}

	a := app.New(midware.Mongo)

	// Load in the recaptcha secret from the config.
	if recaptcha, err := cfg.String(cfgRecaptchaSecret); err == nil {
		a.Ctx["recaptcha"] = recaptcha
		log.Dev("startup", "Init", "Recaptcha Enabled")
	} else {
		log.Dev("startup", "Init", "Recaptcha Disabled")
	}

	log.Dev("startup", "Init", "Initalizing routes")

	routes(a)

	log.Dev("startup", "Init", "Initalizing CORS")
	a.CORS()

	return a
}

func routes(a *app.App) {

	// Create a new app group which will be for internal functions that may have
	// an optional layer of auth added to it.
	internal := a.Group()

	// Now we will load in the public key from the config. If found, we'll add a
	// middleware to all internal endpoints that will ensure that we validate the
	// requests coming in.

	publicKey, err := cfg.String(cfgAuthPublicKey)
	if err != nil || publicKey == "" {
		log.User("startup", "Init", "XENIA_%s is missing, internal authentication is disabled", cfgAuthPublicKey)
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
	external := a.Group()

	external.Handle("POST", "/v1/form/:form_id/submission", handlers.FormSubmission.Create)
}

func ensureDBIndexes() error {
	// Check if mongodb is configured.
	dbName, err := cfg.String(cfgMongoDB)
	if err != nil {
		log.Dev("startup", "Init", "MongoDB Disabled")
		return nil
	}

	mgoDB, err := db.NewMGO("startup", dbName)
	if err != nil {
		return err
	}
	defer mgoDB.CloseMGO("startup")

	return submission.EnsureIndexes("startup", mgoDB)
}
