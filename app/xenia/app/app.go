package app

import (
	"errors"
	"net/http"
	"time"

	"github.com/coralproject/shelf/pkg/db"
	"github.com/coralproject/shelf/pkg/log"

	"github.com/dimfeld/httptreemux"
	"github.com/pborman/uuid"
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

		c := Context{
			DB:             db.NewMGO(),
			ResponseWriter: w,
			Request:        r,
			Params:         p,
			SessionID:      uuid.New(),
		}
		defer c.DB.CloseMGO()

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
