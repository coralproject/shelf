package cayley

import (
	"net/url"

	"github.com/ardanlabs/kit/log"
	"github.com/ardanlabs/kit/web"
	"github.com/coralproject/shelf/internal/platform/db"
)

// Midware handles the Cayley session management.
func Midware(mongoURI *url.URL) web.Middleware {

	// Create the middleware that we can use to create Cayley sessions with.
	m := func(next web.Handler) web.Handler {

		// Create the handler that will be attached in the middleware chain.
		h := func(c *web.Context) error {

			// Load the mongo db from the request context.
			db := c.Ctx["DB"].(*db.DB)

			// Create the new cayley session based on the mongo connection credentials
			// which will add it to the db object itself.
			if err := db.OpenCayley(c.SessionID, mongoURI.String()); err != nil {
				return web.ErrDBNotConfigured
			}

			log.Dev(c.SessionID, "Midware", "Capture Cayley Session")

			// Close the Cayley session when the handler returns.
			defer func() {
				log.Dev(c.SessionID, "Midware", "Release Cayley Session")
				db.CloseCayley(c.SessionID)
			}()

			return next(c)
		}

		return h
	}

	return m
}
