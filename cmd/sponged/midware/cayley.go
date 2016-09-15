package midware

import (
	"github.com/ardanlabs/kit/cfg"
	"github.com/ardanlabs/kit/log"
	"github.com/ardanlabs/kit/web/app"
	"github.com/cayleygraph/cayley"

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
		opts := map[string]interface{}{
			"database_name": cfg.MustString(cfgMongoDB),
			"username":      cfg.MustString(cfgMongoUser),
			"password":      cfg.MustString(cfgMongoPassword),
		}
		store, err := cayley.NewGraph("mongo", mongoHost, opts)
		if err != nil {
			return app.ErrDBNotConfigured
		}

		log.Dev(c.SessionID, "Cayley", "******> Capture Cayley Session")
		c.Ctx["Graph"] = store
		defer func() {
			log.Dev(c.SessionID, "Cayley", "******> Release Cayley Session")
			store.Close()
		}()

		return h(c)
	}
}
