package routes

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/coralproject/xenia/app/xenia/handlers"

	"github.com/ardanlabs/kit/cfg"
	"github.com/ardanlabs/kit/db"
	"github.com/ardanlabs/kit/db/mongo"
	"github.com/ardanlabs/kit/log"
	"github.com/ardanlabs/kit/web/app"
	"github.com/ardanlabs/kit/web/midware"
)

// Mongo config environmental variables.
const (
	cfgMongoHost     = "MONGO_HOST"
	cfgMongoAuthDB   = "MONGO_AUTHDB"
	cfgMongoDB       = "MONGO_DB"
	cfgMongoUser     = "MONGO_USER"
	cfgMongoPassword = "MONGO_PASS"
)

func init() {
	// Initialize the configuration and logging systems. Plus anything
	// else the web app layer needs.
	app.Init("XENIA")

	// Initialize MongoDB.
	if _, err := cfg.String(cfgMongoHost); err == nil {
		cfg := mongo.Config{
			Host:     cfg.MustString(cfgMongoHost),
			AuthDB:   cfg.MustString(cfgMongoAuthDB),
			DB:       cfg.MustString(cfgMongoDB),
			User:     cfg.MustString(cfgMongoUser),
			Password: cfg.MustString(cfgMongoPassword),
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
	a := app.New(midware.Mongo, midware.Auth)

	log.Dev("startup", "Init", "Initalizing routes")
	routes(a)

	log.Dev("startup", "Init", "Initalizing CORS")
	a.CORS()

	log.Dev("startup", "Init", "Initalizing website")
	website(a)

	return a
}

// routes manages the handling of the API endpoints.
func routes(a *app.App) {
	a.Handle("GET", "/1.0/script", handlers.Script.List)
	a.Handle("PUT", "/1.0/script", handlers.Script.Upsert)
	a.Handle("GET", "/1.0/script/:name", handlers.Script.Retrieve)
	a.Handle("DELETE", "/1.0/script/:name", handlers.Script.Delete)

	a.Handle("GET", "/1.0/query", handlers.Query.List)
	a.Handle("PUT", "/1.0/query", handlers.Query.Upsert)
	a.Handle("GET", "/1.0/query/:name", handlers.Query.Retrieve)
	a.Handle("DELETE", "/1.0/query/:name", handlers.Query.Delete)

	a.Handle("GET", "/1.0/regex", handlers.Regex.List)
	a.Handle("PUT", "/1.0/regex", handlers.Regex.Upsert)
	a.Handle("GET", "/1.0/regex/:name", handlers.Regex.Retrieve)
	a.Handle("DELETE", "/1.0/regex/:name", handlers.Regex.Delete)

	a.Handle("POST", "/1.0/exec", handlers.Exec.Custom)
	a.Handle("GET", "/1.0/exec/:name", handlers.Exec.Name)
}

// website manages the serving of web files for the project.
func website(a *app.App) {
	fs := http.FileServer(http.Dir("static"))
	h1 := func(rw http.ResponseWriter, r *http.Request, p map[string]string) {
		fs.ServeHTTP(rw, r)
	}

	a.TreeMux.Handle("GET", "/dist/*path", h1)
	a.TreeMux.Handle("GET", "/img/*path", h1)
	a.TreeMux.Handle("GET", "/", h1)

	h2 := func(rw http.ResponseWriter, r *http.Request, p map[string]string) {
		data, _ := ioutil.ReadFile("static/index.html")
		io.WriteString(rw, string(data))
	}

	file, err := os.Open("static/routes.json")
	if err != nil {
		log.Error("startup", "Init", err, "Initializing website")
		os.Exit(1)
	}

	defer file.Close()

	var routes []struct {
		URL string
	}

	if err := json.NewDecoder(file).Decode(&routes); err != nil {
		log.Error("startup", "Init", err, "Initializing website")
		os.Exit(1)
	}

	for _, route := range routes {
		a.TreeMux.Handle("GET", route.URL, h2)
	}
}
