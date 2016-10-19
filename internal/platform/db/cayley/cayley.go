// Package cayley provides support for the cayley Graph DB with a Mongo backend.
package cayley

import (
	"net/url"
	"strings"

	mgo "gopkg.in/mgo.v2"

	"github.com/cayleygraph/cayley"
	"github.com/cayleygraph/cayley/graph"

	// Every instance of our cayley instance will be using mongo to connect.
	_ "github.com/cayleygraph/cayley/graph/mongo"
)

//==============================================================================

func parseMongoURL(cfg *url.URL) map[string]interface{} {
	opts := make(map[string]interface{})

	// Load the database name from the path, but the path will contain the
	// leading slash as well.
	opts["database_name"] = strings.TrimPrefix(cfg.Path, "/")

	if cfg.User != nil {
		if password, ok := cfg.User.Password(); ok {
			opts["password"] = password
		}

		opts["username"] = cfg.User.Username()
	}

	return opts
}

//==============================================================================

// New creates a new cayley handle.
func New(mongoURL string, ses *mgo.Session) (*cayley.Handle, error) {

	// If the session is not provided, create the Cayley Handle from
	// the mongo URL.
	if ses == nil {

		// Parse the provied mongoURL.
		cfg, err := url.Parse(mongoURL)
		if err != nil {
			return nil, err
		}

		// Form the Cayley connection options.
		opts := parseMongoURL(cfg)

		// Create the cayley handle that maintains a connection to the
		// Cayley graph database in Mongo.
		store, err := cayley.NewGraph("mongo", cfg.Host, opts)
		if err != nil {
			return store, err
		}

		return store, nil
	}

	// If the session is not nil, create the Cayley handle with the
	// provided mongo session.
	opts := make(map[string]interface{})
	opts["session"] = ses

	// Create the cayley handle that maintains a connection to the
	// Cayley graph database in Mongo.
	store, err := cayley.NewGraph("mongo", "", opts)
	if err != nil {
		return store, err
	}

	return store, nil
}

// InitQuadStore initializes the quadstore.
func InitQuadStore(mongoURL string) error {
	cfg, err := url.Parse(mongoURL)
	if err != nil {
		return err
	}

	// Form the Cayley connection options.
	opts := parseMongoURL(cfg)

	if err := graph.InitQuadStore("mongo", cfg.Host, opts); err != nil {
		return err
	}

	return nil
}
