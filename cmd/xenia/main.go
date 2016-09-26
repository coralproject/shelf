// This program provides the coral project xenia central CLI platform.
package main

import (
	"os"

	"github.com/ardanlabs/kit/cfg"
	"github.com/ardanlabs/kit/db"
	"github.com/ardanlabs/kit/db/mongo"
	"github.com/ardanlabs/kit/log"
	"github.com/coralproject/shelf/cmd/xenia/cmddb"
	"github.com/coralproject/shelf/cmd/xenia/cmdmask"
	"github.com/coralproject/shelf/cmd/xenia/cmdpattern"
	"github.com/coralproject/shelf/cmd/xenia/cmdquery"
	"github.com/coralproject/shelf/cmd/xenia/cmdregex"
	"github.com/coralproject/shelf/cmd/xenia/cmdrelationship"
	"github.com/coralproject/shelf/cmd/xenia/cmdscript"
	"github.com/coralproject/shelf/cmd/xenia/cmdview"
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
	cfgWebHost       = "WEB_HOST"
)

var xenia = &cobra.Command{
	Use:   "xenia",
	Short: "Xenia provides the central cli housing of various cli tools that interface with the API",
}

func main() {
	if err := cfg.Init(cfg.EnvProvider{Namespace: "XENIA"}); err != nil {
		xenia.Println("Unable to initialize configuration")
		os.Exit(1)
	}

	logLevel := func() int {
		ll, err := cfg.Int(cfgLoggingLevel)
		if err != nil {
			return log.NONE
		}
		return ll
	}
	log.Init(os.Stderr, logLevel, log.Ldefault)

	xenia.Println("Using log level", logLevel())

	// Pull options from the config.
	var conn *db.DB
	if _, errHost := cfg.String(cfgWebHost); errHost != nil {
		xenia.Println("Configuring MongoDB")

		cfg := mongo.Config{
			Host:     cfg.MustString(cfgMongoHost),
			AuthDB:   cfg.MustString(cfgMongoAuthDB),
			DB:       cfg.MustString(cfgMongoDB),
			User:     cfg.MustString(cfgMongoUser),
			Password: cfg.MustString(cfgMongoPassword),
		}

		err := db.RegMasterSession("startup", cfg.DB, cfg)
		if err != nil {
			xenia.Println("Unable to initialize MongoDB")
			os.Exit(1)
		}

		conn, err = db.NewMGO("", cfg.DB)
		if err != nil {
			xenia.Println("Unable to get MongoDB session")
			os.Exit(1)
		}
		defer conn.CloseMGO("")
	}

	xenia.AddCommand(
		cmddb.GetCommands(conn),
		cmdquery.GetCommands(),
		cmdscript.GetCommands(),
		cmdregex.GetCommands(),
		cmdmask.GetCommands(),
		cmdrelationship.GetCommands(),
		cmdview.GetCommands(),
		cmdpattern.GetCommands(),
	)
	xenia.Execute()
}
