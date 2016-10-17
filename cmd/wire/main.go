// This program provides a set of commands for relationship/view functionality.
package main

import (
	"os"

	"github.com/ardanlabs/kit/cfg"
	"github.com/ardanlabs/kit/log"
	"github.com/cayleygraph/cayley"
	"github.com/coralproject/shelf/cmd/wire/cmdview"
	"github.com/coralproject/shelf/internal/platform/db"
	cayleyshelf "github.com/coralproject/shelf/internal/platform/db/cayley"
	"github.com/spf13/cobra"
)

const (

	// cfgLoggingLevel is the key for the logging level.
	cfgLoggingLevel = "LOGGING_LEVEL"

	// cfgMongoURI is the key for the URI to the MongoDB service.
	cfgMongoURI = "MONGO_URI"
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

	mongoURI := cfg.MustURL(cfgMongoURI)

	err := db.RegMasterSession("startup", mongoURI.Path, mongoURI.String(), 0)
	if err != nil {
		wire.Println("Unable to initialize MongoDB")
		os.Exit(1)
	}

	mgoDB, err = db.NewMGO("", mongoURI.Path)
	if err != nil {
		wire.Println("Unable to get MongoDB session")
		os.Exit(1)
	}
	defer mgoDB.CloseMGO("")

	// Configure Cayley.
	wire.Println("Configuring Cayley")

	graphDB, err = cayleyshelf.New(mongoURI.String(), nil)
	if err != nil {
		wire.Println("Unable to get Cayley handle")
		os.Exit(1)
	}

	// Add the graph and view commands to the CLI tool.
	wire.AddCommand(
		cmdview.GetCommands(mgoDB, graphDB),
	)

	// Execute the command.
	wire.Execute()
}
