package error

import (
	"fmt"
	"net/http"
	"runtime"
	"runtime/debug"

	"github.com/ardanlabs/kit/log"
	"github.com/ardanlabs/kit/web"
)

// Midware for catching errors.
func Midware(next web.Handler) web.Handler {

	// Create the handler that will be attached in the middleware chain.
	h := func(c *web.Context) error {

		log.Dev(c.SessionID, "error : Midware", "Started")

		// In the event of a panic, we want to capture it here so we can send an
		// error down the stack.
		defer func() {
			if err := recover(); err != nil {

				// Respond with the error.
				c.RespondError("internal server error", http.StatusInternalServerError)

				// Log out that we caught the error.
				switch err := err.(type) {
				case error:
					log.Error(c.SessionID, "error : Midware", err, "Panic Caught")
				default:
					log.Error(c.SessionID, "error : Midware", fmt.Errorf("%v", err), "Panic Caught")
				}

				// Print out the stack.
				log.Dev(c.SessionID, "error : Midware", "Panic Stacktrace:\n%s", debug.Stack())

				_, filePath, line, _ := runtime.Caller(4)
				log.Dev(c.SessionID, "error : Midware", "Panic Traced to %s:%d", filePath, line)
			}
		}()

		if err := next(c); err != nil {

			// Respond with the error.
			c.Error(err)

			// Log out that we caught the error.
			log.Dev(c.SessionID, "error : Midware", "Completed : Error Caught : %s", err.Error())

			return nil
		}

		log.Dev(c.SessionID, "error : Midware", "Completed")

		return nil
	}

	return h
}
