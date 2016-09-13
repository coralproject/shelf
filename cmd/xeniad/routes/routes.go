package routes

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/anvilresearch/go-anvil"
	"github.com/ardanlabs/kit/cfg"
	"github.com/ardanlabs/kit/db"
	"github.com/ardanlabs/kit/db/mongo"
	"github.com/ardanlabs/kit/log"
	"github.com/ardanlabs/kit/web/app"
	"github.com/coralproject/shelf/cmd/xeniad/handlers"
	"github.com/coralproject/shelf/cmd/xeniad/midware"
)

// Environmental variables.
const (
	cfgMongoHost     = "MONGO_HOST"
	cfgMongoAuthDB   = "MONGO_AUTHDB"
	cfgMongoDB       = "MONGO_DB"
	cfgMongoUser     = "MONGO_USER"
	cfgMongoPassword = "MONGO_PASS"
	cfgAnvilHost     = "ANVIL_HOST"
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

	// If authentication is on then configure Anvil.
	var anv *anvil.Anvil
	if url, err := cfg.String(cfgAnvilHost); err == nil {

		log.Dev("startup", "Init", "Initalizing Anvil")
		anv, err = anvil.New(url)
		if err != nil {
			log.Error("startup", "Init", err, "Initializing Anvil: %s", url)
			os.Exit(1)
		}
	}

	a := app.New(midware.Mongo, midware.Auth)
	a.Ctx["anvil"] = anv

	log.Dev("startup", "Init", "Initalizing routes")
	routes(a)

	log.Dev("startup", "Init", "Initalizing CORS")
	a.CORS()

	// It has been decided the website is no longer required.
	// if testing == nil {
	// 	log.Dev("startup", "Init", "Initalizing website")
	// 	website(a)
	// }

	return a
}

// routes manages the handling of the API endpoints.
func routes(a *app.App) {
	a.Handle("GET", "/1.0/version", handlers.Version.List)

	a.Handle("GET", "/1.0/script", handlers.Script.List)
	a.Handle("PUT", "/1.0/script", handlers.Script.Upsert)
	a.Handle("GET", "/1.0/script/:name", handlers.Script.Retrieve)
	a.Handle("DELETE", "/1.0/script/:name", handlers.Script.Delete)

	a.Handle("GET", "/1.0/query", handlers.Query.List)
	a.Handle("PUT", "/1.0/query", handlers.Query.Upsert)
	a.Handle("GET", "/1.0/query/:name", handlers.Query.Retrieve)
	a.Handle("DELETE", "/1.0/query/:name", handlers.Query.Delete)

	a.Handle("PUT", "/1.0/index/:name", handlers.Query.EnsureIndexes)

	a.Handle("GET", "/1.0/regex", handlers.Regex.List)
	a.Handle("PUT", "/1.0/regex", handlers.Regex.Upsert)
	a.Handle("GET", "/1.0/regex/:name", handlers.Regex.Retrieve)
	a.Handle("DELETE", "/1.0/regex/:name", handlers.Regex.Delete)

	a.Handle("GET", "/1.0/mask", handlers.Mask.List)
	a.Handle("PUT", "/1.0/mask", handlers.Mask.Upsert)
	a.Handle("GET", "/1.0/mask/:collection/:field", handlers.Mask.Retrieve)
	a.Handle("GET", "/1.0/mask/:collection", handlers.Mask.Retrieve)
	a.Handle("DELETE", "/1.0/mask/:collection/:field", handlers.Mask.Delete)

	a.Handle("POST", "/1.0/exec", handlers.Exec.Custom)
	a.Handle("GET", "/1.0/exec/:name", handlers.Exec.Name)

	a.Handle("GET", "/1.0/relationship", handlers.Relationship.List)
	a.Handle("PUT", "/1.0/relationship", handlers.Relationship.Upsert)
	a.Handle("GET", "/1.0/relationship/:predicate", handlers.Relationship.Retrieve)
	a.Handle("DELETE", "/1.0/relationship/:predicate", handlers.Relationship.Delete)

	a.Handle("GET", "/1.0/view", handlers.View.List)
	a.Handle("PUT", "/1.0/view", handlers.View.Upsert)
	a.Handle("GET", "/1.0/view/:name", handlers.View.Retrieve)
	a.Handle("DELETE", "/1.0/view/:name", handlers.View.Delete)

	a.Handle("GET", "/1.0/pattern", handlers.Pattern.List)
	a.Handle("PUT", "/1.0/pattern", handlers.Pattern.Upsert)
	a.Handle("GET", "/1.0/pattern/:type", handlers.Pattern.Retrieve)
	a.Handle("DELETE", "/1.0/pattern/:type", handlers.Pattern.Delete)

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
