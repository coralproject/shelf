package midware

import (
	"github.com/anvilresearch/go-anvil"
	"github.com/ardanlabs/kit/cfg"
	"github.com/ardanlabs/kit/log"
	"github.com/ardanlabs/kit/web/app"
)

const cfgAnvilHost = "ANVIL_HOST"

// Auth handles token authentication.
func Auth(h app.Handler) app.Handler {

	// Check if authentication is turned off.
	if _, err := cfg.String(cfgAnvilHost); err != nil {
		return func(c *app.Context) error {
			log.Dev(c.SessionID, "Auth", "******> Authentication Off")
			return h(c)
		}
	}

	// Turn authentication on.
	return func(c *app.Context) error {
		a := c.App.Ctx["anvil"].(*anvil.Anvil)

		claims, err := a.ValidateFromRequest(c.Request)
		if err != nil {
			log.Error(c.SessionID, "Auth", err, "Validating token")
			return app.ErrNotAuthorized
		}

		c.Ctx["claims"] = claims

		log.Dev(c.SessionID, "Auth", "Completed")
		return h(c)
	}
}
