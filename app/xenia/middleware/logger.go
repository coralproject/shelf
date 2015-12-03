package middleware

import (
	"time"

	"github.com/coralproject/shelf/app/xenia/app"
	"github.com/coralproject/shelf/pkg/log"
)

// Logger writes a record of the start and end of the request.
func Logger(h app.Handler) app.Handler {
	return func(c *app.Context) error {
		start := time.Now()

		log.User(c.SessionID, "Request", "Started : Method[%s] URL[%s] RADDR[%s]", c.Request.Method, c.Request.URL.Path, c.Request.RemoteAddr)

		err := h(c)

		log.User(c.SessionID, "Request", "Completed : Status[%d] Duration[%s]", c.Status, time.Since(start))

		return err
	}
}
