package error

import (
	"net/http"

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
				log.Dev(c.SessionID, "error : Midware", "Completed : Panic Caught : %v", err)
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
