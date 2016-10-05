package midware

import (
	"github.com/ardanlabs/kit/cfg"
	"github.com/ardanlabs/kit/db"
	kitcayley "github.com/ardanlabs/kit/db/cayley"
	"github.com/ardanlabs/kit/log"
	"github.com/ardanlabs/kit/web/app"

	// mongo is needed to utilize mongoDB as the backend store for cayley.
	_ "github.com/cayleygraph/cayley/graph/mongo"
)

const (
	cfgMongoHost     = "MONGO_HOST"
	cfgMongoUser     = "MONGO_USER"
	cfgMongoPassword = "MONGO_PASS"
)

// Cayley handles session management.
func Cayley(h app.Handler) app.Handler {

	// Check if mongodb is configured.
	mongoHost, err := cfg.String(cfgMongoHost)
	if err != nil {
		return func(c *app.Context) error {
			log.Dev(c.SessionID, "Cayley", "******> Cayley Not Configured")
			return h(c)
		}
	}

	// Wrap the handlers inside a session copy/close.
	return func(c *app.Context) error {
		db := c.Ctx["DB"].(*db.DB)
		cayleyCfg := kitcayley.Config{
			Host:     mongoHost,
			DB:       cfg.MustString(cfgMongoDB),
			User:     cfg.MustString(cfgMongoUser),
			Password: cfg.MustString(cfgMongoPassword),
		}
		if err := db.OpenCayley(c.SessionID, cayleyCfg); err != nil {
			return app.ErrDBNotConfigured
		}

		log.Dev(c.SessionID, "Cayley", "Capture Cayley Session")
		defer func() {
			log.Dev(c.SessionID, "Cayley", "Release Cayley Session")
			db.CloseCayley(c.SessionID)
		}()

		return h(c)
	}
}
