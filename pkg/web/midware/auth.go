package midware

import (
	"github.com/coralproject/shelf/pkg/log"
	"github.com/coralproject/shelf/pkg/srv/auth"
	"github.com/coralproject/shelf/pkg/web/app"
)

// Auth handles token authentication.
func Auth(h app.Handler) app.Handler {
	return func(c *app.Context) error {
		token := c.Request.Header.Get("Authorization")
		log.Dev(c.SessionID, "Auth", "Started : Token[%s]", token)

		if len(token) < 5 || token[0:5] != "Basic" {
			log.Error(c.SessionID, "Auth", app.ErrNotAuthorized, "Validating token")
			return app.ErrNotAuthorized
		}

		var err error
		if c.User, err = auth.ValidateWebToken(c.SessionID, c.DB, token[6:]); err != nil {
			log.Error(c.SessionID, "Auth", err, "Validating token")
			return app.ErrNotAuthorized
		}

		log.Dev(c.SessionID, "Auth", "Completed")
		return h(c)
	}
}
