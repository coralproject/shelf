package mongo

import (
	"net/url"

	"github.com/ardanlabs/kit/log"
	"github.com/ardanlabs/kit/web"
	kitdb "github.com/coralproject/shelf/internal/platform/db"
)

// Midware handles databse session management and manages a MongoDB session.
func Midware(mongoURI *url.URL) web.Middleware {

	// Create the middleware that we can use to create MongoDB sessions with.
	m := func(next web.Handler) web.Handler {

		// Create the handler that will be attached in the middleware chain.
		h := func(c *web.Context) error {

			// Pull in the mongo session from the master session so we can load it
			// onto the request context. It is keyed by the path on the uri.
			db, err := kitdb.NewMGO(c.SessionID, mongoURI.Path)
			if err != nil {
				log.Error(c.SessionID, "mongo : Midware", err, "Method[%s] URL[%s] RADDR[%s]", c.Request.Method, c.Request.URL.Path, c.Request.RemoteAddr)
				return web.ErrDBNotConfigured
			}

			// Load the mongo database onto the request context.
			c.Ctx["DB"] = db

			log.Dev(c.SessionID, "mongo : Midware", "Capture Mongo Session")

			// Close the MongoDB session when the handler returns.
			defer func() {
				log.Dev(c.SessionID, "mongo : Midware", "Release Mongo Session")
				db.CloseMGO(c.SessionID)
			}()

			return next(c)
		}

		return h
	}

	return m
}
