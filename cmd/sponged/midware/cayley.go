package midware

import (
	"github.com/ardanlabs/kit/cfg"
	"github.com/ardanlabs/kit/log"
	"github.com/ardanlabs/kit/web/app"
	"github.com/cayleygraph/cayley"

	// mongo is needed to utilize mongoDB as the backend store for cayley.
	_ "github.com/cayleygraph/cayley/graph/mongo"
)

// cfgMongoDB config environmental variables.
const cfgMongoHost = "MONGO_HOST"

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
		store, err := cayley.NewGraph("mongo", mongoHost, nil)
		if err != nil {
			return app.ErrDBNotConfigured
		}

		log.Dev(c.SessionID, "Cayley", "******> Capture Cayley Session")
		c.Ctx["Cayley"] = store
		defer func() {
			log.Dev(c.SessionID, "Cayley", "******> Release Cayley Session")
			store.Close()
		}()

		return h(c)
	}
}
