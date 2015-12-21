// Package app provides a thin layer of support for writing web services. It
// integrates with the ardanlabs kit repo to provide support for logging,
// configuration, database, routing and application context. The base things
// you need to write a web service is provided.
package app

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/ardanlabs/kit/cfg"
	"github.com/ardanlabs/kit/db"
	"github.com/ardanlabs/kit/db/mongo"
	"github.com/ardanlabs/kit/log"

	"github.com/dimfeld/httptreemux"
	"github.com/pborman/uuid"
)

// Web config environmental variables.
const (
	cfgLoggingLevel = "LOGGING_LEVEL"
	cfgHost         = "HOST"
)

// Mongo config environmental variables.
const (
	cfgMongoHost     = "MONGO_HOST"
	cfgMongoAuthDB   = "MONGO_AUTHDB"
	cfgMongoDB       = "MONGO_DB"
	cfgMongoUser     = "MONGO_USER"
	cfgMongoPassword = "MONGO_PASS"
)

var (
	// ErrNotAuthorized occurs when the call is not authorized.
	ErrNotAuthorized = errors.New("Not authorized")

	// ErrNotFound is abstracting the mgo not found error.
	ErrNotFound = errors.New("Entity Not found")

	// ErrInvalidID occurs when an ID is not in a valid form.
	ErrInvalidID = errors.New("ID is not in it's proper form")

	// ErrValidation occurs when there are validation errors.
	ErrValidation = errors.New("Validation errors occurred")
)

type (
	// A Handler is a type that handles an http request within our own little mini
	// framework. The fun part is that our context is fully controlled and
	// configured by us so we can extend the functionality of the Context whenever
	// we want.
	Handler func(*Context) error

	// A Middleware is a type that wraps a handler to remove boilerplate or other
	// concerns not direct to any given Handler.
	Middleware func(Handler) Handler
)

// app maintains some framework state.
var app struct {
	useMongo    bool
	userHeaders map[string]string // Extra headers for each response.
}

//==============================================================================

// App is the entrypoint into our application and what configures our context
// object for each of our http handlers. Feel free to add any configuration
// data/logic on this App struct
type App struct {
	*httptreemux.TreeMux
	mw []Middleware
}

// New create an App value that handle a set of routes for the application.
// You can provide any number of middleware and they'll be used to wrap every
// request handler.
func New(mw ...Middleware) *App {
	return &App{
		TreeMux: httptreemux.New(),
		mw:      mw,
	}
}

// Handle is our mechanism for mounting Handlers for a given HTTP verb and path
// pair, this makes for really easy, convenient routing.
func (a *App) Handle(verb, path string, handler Handler, mw ...Middleware) {

	// The function to execute for each request.
	h := func(w http.ResponseWriter, r *http.Request, p map[string]string) {
		start := time.Now()

		var dbConn *db.DB
		if app.useMongo {
			dbConn = db.NewMGO()
		}

		c := Context{
			DB:             dbConn,
			ResponseWriter: w,
			Request:        r,
			Params:         p,
			SessionID:      uuid.New(),
		}

		if app.useMongo {
			defer c.DB.CloseMGO()
		}

		log.User(c.SessionID, "Request", "Started : Method[%s] URL[%s] RADDR[%s]", c.Request.Method, c.Request.URL.Path, c.Request.RemoteAddr)

		// Wrap the handler in all associated middleware.
		wrap := func(h Handler) Handler {
			// Wrap up the application-wide first...
			for i := len(a.mw) - 1; i >= 0; i-- {
				h = a.mw[i](h)
			}

			// Then wrap with our route specific ones.
			for i := len(mw) - 1; i >= 0; i-- {
				h = mw[i](h)
			}

			return h
		}

		// Call the wrapped handler and handle any possible error.
		if err := wrap(handler)(&c); err != nil {
			c.Error(err)
		}

		log.User(c.SessionID, "Request", "Completed : Status[%d] Duration[%s]", c.Status, time.Since(start))
	}

	// Add this handler for the specified verb and route.
	a.TreeMux.Handle(verb, path, h)
}

// CORS providing support for Cross-Origin Resource Sharing.
// https://metajack.im/2010/01/19/crossdomain-ajax-for-xmpp-http-binding-made-easy/
func (a *App) CORS() {
	h := func(w http.ResponseWriter, r *http.Request, p map[string]string) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Max-Age", "86400")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		//fmt.Fprintf(w, "")
	}

	a.TreeMux.OptionsHandler = h

	// This header is now required for all requests.
	app.userHeaders["Access-Control-Allow-Origin"] = "*"
}

//==============================================================================

// Init is called to initialize the application.
func Init(configKey string) {

	// Init the configuration system.
	if err := cfg.Init(configKey); err != nil {
		fmt.Println("Error initalizing configuration system", err)
		os.Exit(1)
	}

	// Init the log system.
	logLevel := func() int {
		ll, err := cfg.Int(cfgLoggingLevel)
		if err != nil {
			return log.USER
		}
		return ll
	}
	log.Init(os.Stderr, logLevel)

	// Log all the configuration options
	log.User("startup", "Init", "\n\nConfig Settings: %s\n%s\n", configKey, cfg.Log())

	// Init MongoDB if configured.
	if _, err := cfg.String(cfgMongoHost); err == nil {
		app.useMongo = true

		cfg := mongo.Config{
			Host:     cfg.MustString(cfgMongoHost),
			AuthDB:   cfg.MustString(cfgMongoAuthDB),
			DB:       cfg.MustString(cfgMongoDB),
			User:     cfg.MustString(cfgMongoUser),
			Password: cfg.MustString(cfgMongoPassword),
		}

		if err := mongo.Init(cfg); err != nil {
			log.Error("startup", "Init", err, "Initializing MongoDB")
			os.Exit(1)
		}
	}

	// Load user defined custom headers. HEADERS should be key:value,key:value
	if hs, err := cfg.String("HEADERS"); err == nil {
		app.userHeaders = make(map[string]string)
		hdrs := strings.Split(hs, ",")
		for _, hdr := range hdrs {
			if kv := strings.Split(hdr, ":"); len(kv) == 2 {
				log.User("startup", "Init", "User Headers : %s:%s", kv[0], kv[1])
				app.userHeaders[kv[0]] = kv[1]
			}
		}
	}
}

// Run is called to start the web service.
func Run(defaultHost string, routes http.Handler) {
	log.Dev("startup", "Run", "Start : defaultHost[%s]", defaultHost)

	// Check for a configured host value.
	host, err := cfg.String(cfgHost)
	if err != nil {
		host = defaultHost
	}

	// Create this goroutine to run the web server.
	go func() {
		log.Dev("listener", "Run", "Listening on: %s", host)
		http.ListenAndServe(host, routes)
	}()

	// Listen for an interrupt signal from the OS.
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	<-sigChan

	log.Dev("shutdown", "Run", "Complete")
}
