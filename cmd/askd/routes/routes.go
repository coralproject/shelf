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
)

// Environmental variables.
const (
	cfgMongoHost       = "MONGO_HOST"
	cfgMongoAuthDB     = "MONGO_AUTHDB"
	cfgMongoDB         = "MONGO_DB"
	cfgMongoUser       = "MONGO_USER"
	cfgMongoPassword   = "MONGO_PASS"
	cfgAnvilHost       = "ANVIL_HOST"
	cfgRecaptchaSecret = "RECAPTCHA_SECRET"
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

	// If authentication is on then configure Anvil.
	/*

		// Anvil is temporarily disabled pending auth strategy

		var anv *anvil.Anvil
		if url, err := cfg.String(cfgAnvilHost); err == nil {

			log.Dev("startup", "Init", "Initalizing Anvil")
			anv, err = anvil.New(url)
			if err != nil {
				log.Error("startup", "Init", err, "Initializing Anvil: %s", url)
				os.Exit(1)
			}
		}
	*/

	a := app.New(midware.Mongo, midware.Auth)
	//		a.Ctx["anvil"] = anv

	// Load in the recaptcha secret from the config.
	if recaptcha, err := cfg.String(cfgRecaptchaSecret); err == nil {
		a.Ctx["recaptcha"] = recaptcha
		log.Dev("startup", "Init", "Recaptcha Enabled")
	} else {
		log.Dev("startup", "Init", "Recaptcha Disabled")
	}

	log.Dev("startup", "Init", "Initalizing routes")

	//oldRoutes(a) // FIXME: remove on next API release
	routes(a)

	log.Dev("startup", "Init", "Initalizing CORS")
	a.CORS()

	return a
}

// oldRoutes manages the handling of the API endpoints for the old style ported
// from Ask.
//
// FIXME: remove on next API release
func oldRoutes(a *app.App) {
	// forms
	a.Handle("POST", "/api/form", handlers.Form.Upsert)
	a.Handle("PUT", "/api/form", handlers.Form.Upsert)
	a.Handle("PUT", "/api/form/:id/status/:status", handlers.Form.UpdateStatus)
	a.Handle("GET", "/api/forms", handlers.Form.List)
	a.Handle("GET", "/api/form/:id", handlers.Form.Retrieve)
	a.Handle("DELETE", "/api/form/:id", handlers.Form.Delete)

	// form submissions
	a.Handle("POST", "/api/form_submission/:id", handlers.FormSubmission.Create)
	a.Handle("PUT", "/api/form_submission/:id/status/:status", handlers.FormSubmission.UpdateStatus)
	a.Handle("GET", "/api/form_submissions/:form_id", handlers.FormSubmission.Search)
	a.Handle("GET", "/api/form_submission/:id", handlers.FormSubmission.Retrieve)
	a.Handle("PUT", "/api/form_submission/:id/:answer_id", handlers.FormSubmission.UpdateAnswer)
	a.Handle("PUT", "/api/form_submission/:id/flag/:flag", handlers.FormSubmission.AddFlag)
	a.Handle("DELETE", "/api/form_submission/:id/flag/:flag", handlers.FormSubmission.RemoveFlag)
	a.Handle("DELETE", "/api/form_submission/:id", handlers.FormSubmission.Delete)

	// form galleries
	a.Handle("GET", "/api/form_gallery/:id", handlers.FormGallery.Retrieve)
	a.Handle("GET", "/api/form_galleries/:form_id", handlers.FormGallery.RetrieveForForm)
	a.Handle("GET", "/api/form_galleries/form/:form_id", handlers.FormGallery.RetrieveForForm)
	a.Handle("PUT", "/api/form_gallery/:id/add/:submission_id/:answer_id", handlers.FormGallery.AddAnswer)
	a.Handle("PUT", "/api/form_gallery/:id", handlers.FormGallery.Update)
	a.Handle("DELETE", "/api/form_gallery/:id/remove/:submission_id/:answer_id", handlers.FormGallery.RemoveAnswer)
}

func routes(a *app.App) {
	// global
	a.Handle("GET", "/v1/version", handlers.Version.List)

	// forms
	a.Handle("POST", "/v1/form", handlers.Form.Upsert)
	a.Handle("GET", "/v1/form", handlers.Form.List)
	a.Handle("PUT", "/v1/form/:id", handlers.Form.Upsert)
	a.Handle("PUT", "/v1/form/:id/status/:status", handlers.Form.UpdateStatus)
	a.Handle("GET", "/v1/form/:id", handlers.Form.Retrieve)
	a.Handle("DELETE", "/v1/form/:id", handlers.Form.Delete)

	// form form submissions
	a.Handle("POST", "/v1/form/:form_id/submission", handlers.FormSubmission.Create)
	a.Handle("GET", "/v1/form/:form_id/submission", handlers.FormSubmission.Search)
	a.Handle("GET", "/v1/form/:form_id/submission/:id", handlers.FormSubmission.Retrieve)
	a.Handle("PUT", "/v1/form/:form_id/submission/:id/status/:status", handlers.FormSubmission.UpdateStatus)
	a.Handle("POST", "/v1/form/:form_id/submission/:id/flag/:flag", handlers.FormSubmission.AddFlag)
	a.Handle("DELETE", "/v1/form/:form_id/submission/:id/flag/:flag", handlers.FormSubmission.RemoveFlag)
	a.Handle("PUT", "/v1/form/:form_id/submission/:id/answer/:answer_id", handlers.FormSubmission.UpdateAnswer)
	a.Handle("DELETE", "/v1/form/:form_id/submission/:id", handlers.FormSubmission.Delete)

	// temporal route to get CSV file - TO DO : move into a different service
	a.Handle("GET", "/v1/form/:form_id/submission/export", handlers.FormSubmission.Download)

	// form form galleries
	a.Handle("GET", "/v1/form/:form_id/gallery", handlers.FormGallery.RetrieveForForm)

	// form galleries
	a.Handle("GET", "/v1/form_gallery/:id", handlers.FormGallery.Retrieve)
	a.Handle("PUT", "/v1/form_gallery/:id", handlers.FormGallery.Update)
	a.Handle("POST", "/v1/form_gallery/:id/submission/:submission_id/:answer_id", handlers.FormGallery.AddAnswer)
	a.Handle("DELETE", "/v1/form_gallery/:id/submission/:submission_id/:answer_id", handlers.FormGallery.RemoveAnswer)
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
