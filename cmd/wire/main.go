// This program provides a set of commands for relationship/view functionality.
package main

import (
	"os"

	"github.com/ardanlabs/kit/cfg"
	"github.com/ardanlabs/kit/db"
	"github.com/ardanlabs/kit/db/mongo"
	"github.com/ardanlabs/kit/log"
	"github.com/cayleygraph/cayley"
	_ "github.com/cayleygraph/cayley/graph/mongo"
	"github.com/coralproject/shelf/cmd/wire/cmdgraph"
	"github.com/coralproject/shelf/cmd/wire/cmdview"
	"github.com/spf13/cobra"
)

// Config environmental variables.
const (
	cfgLoggingLevel  = "LOGGING_LEVEL"
	cfgMongoHost     = "MONGO_HOST"
	cfgMongoAuthDB   = "MONGO_AUTHDB"
	cfgMongoDB       = "MONGO_DB"
	cfgMongoUser     = "MONGO_USER"
	cfgMongoPassword = "MONGO_PASS"
)

// wire includes information about the wire cobra command.
var wire = &cobra.Command{
	Use:   "wire",
	Short: "Wire provides the central cli housing of various cli tools that interface with the internal wire API",
}

func main() {

	// Initialize the configuration
	if err := cfg.Init(cfg.EnvProvider{Namespace: "WIRE"}); err != nil {
		wire.Println("Unable to initialize configuration")
		os.Exit(1)
	}

	// Initialize the logging
	logLevel := func() int {
		ll, err := cfg.Int(cfgLoggingLevel)
		if err != nil {
			return log.NONE
		}
		return ll
	}

	log.Init(os.Stderr, logLevel, log.Ldefault)
	wire.Println("Using log level", logLevel())

	// Pull options from the config.
	var mgoDB *db.DB
	var graphDB *cayley.Handle

	// Configure MongoDB.
	wire.Println("Configuring MongoDB")

	mongoCfg := mongo.Config{
		Host:     cfg.MustString(cfgMongoHost),
		AuthDB:   cfg.MustString(cfgMongoAuthDB),
		DB:       cfg.MustString(cfgMongoDB),
		User:     cfg.MustString(cfgMongoUser),
		Password: cfg.MustString(cfgMongoPassword),
	}

	err := db.RegMasterSession("startup", mongoCfg.DB, mongoCfg)
	if err != nil {
		wire.Println("Unable to initialize MongoDB")
		os.Exit(1)
	}

	mgoDB, err = db.NewMGO("", mongoCfg.DB)
	if err != nil {
		wire.Println("Unable to get MongoDB session")
		os.Exit(1)
	}
	defer mgoDB.CloseMGO("")

	// Configure Cayley.
	wire.Println("Configuring Cayley")

	opts := map[string]interface{}{
		"database_name": cfg.MustString(cfgMongoDB),
		"username":      cfg.MustString(cfgMongoUser),
		"password":      cfg.MustString(cfgMongoPassword),
	}
	graphDB, err = cayley.NewGraph("mongo", cfg.MustString(cfgMongoHost), opts)
	if err != nil {
		wire.Println("Unable to get Cayley handle")
		os.Exit(1)
	}

	// Add the graph and view commands to the CLI tool.
	wire.AddCommand(
		cmdview.GetCommands(mgoDB, graphDB),
		cmdgraph.GetCommands(mgoDB, graphDB),
	)

	// Execute the command.
	wire.Execute()
}
