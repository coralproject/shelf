package middleware

import (
	"time"

	"github.com/coralproject/shelf/log"
	"github.com/coralproject/shelf/xenia/app"
)

// RequestLogger writes some information about the request to the logs in
// the format: SESSIONID : (200) GET /foo -> IP ADDR (latency)
func RequestLogger(h app.Handler) app.Handler {
	return func(c *app.Context) error {
		log.Dev(c.SessionID, "RequestLogger", "Started")

		start := time.Now()
		err := h(c)

		log.User(c.SessionID, "Request", "(%d) %s %s -> %s (%s)",
			c.Status, c.Request.Method, c.Request.URL.Path,
			c.Request.RemoteAddr, time.Since(start),
		)

		log.Dev(c.SessionID, "RequestLogger", "Completed")
		return err
	}
}
