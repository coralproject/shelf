package log

import (
	"time"

	"github.com/ardanlabs/kit/log"
	"github.com/ardanlabs/kit/web"
)

// Midware handles the request logging.
func Midware(next web.Handler) web.Handler {

	// Create the handler that will be attached in the middleware chain.
	h := func(c *web.Context) error {

		log.User(c.SessionID, "log : Midware", "Started : Method[%s] URL[%s] RADDR[%s]", c.Request.Method, c.Request.URL.Path, c.Request.RemoteAddr)

		if err := next(c); err != nil {
			return err
		}

		log.User(c.SessionID, "log : Midware", "Completed : Status[%d] Duration[%s]", c.Status, time.Since(c.Now))

		return nil
	}

	return h
}
