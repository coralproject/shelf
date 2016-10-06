package log

import (
	"net/url"
	"time"

	"github.com/ardanlabs/kit/log"
	"github.com/ardanlabs/kit/web"
)

// Midware handles the request logging.
func Midware(mongoURI *url.URL) web.Middleware {

	// Create the middleware that we can use to log requests with.
	m := func(next web.Handler) web.Handler {

		// Create the handler that will be attached in the middleware chain.
		h := func(c *web.Context) error {

			log.User(c.SessionID, "Midware", "Started : Method[%s] URL[%s] RADDR[%s]", c.Request.Method, c.Request.URL.Path, c.Request.RemoteAddr)

			if err := next(c); err != nil {
				return err
			}

			log.User(c.SessionID, "Midware", "Completed : Status[%d] Duration[%s]", c.Status, time.Since(c.Now))

			return nil
		}

		return h
	}

	return m
}
