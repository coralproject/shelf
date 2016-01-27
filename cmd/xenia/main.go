// This program provides the coral project xenia central CLI platform.
package main

import (
	"os"

	"github.com/coralproject/xenia/cmd/xenia/cmdquery"
	"github.com/coralproject/xenia/cmd/xenia/cmdregex"
	"github.com/coralproject/xenia/cmd/xenia/cmdscript"

	"github.com/ardanlabs/kit/cfg"
	"github.com/ardanlabs/kit/cmd/kit/cmdauth"
	"github.com/ardanlabs/kit/cmd/kit/cmddb"
	"github.com/ardanlabs/kit/db"
	"github.com/ardanlabs/kit/db/mongo"
	"github.com/ardanlabs/kit/log"

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
	log.Init(os.Stderr, logLevel)

	cfg := mongo.Config{
		Host:     cfg.MustString(cfgMongoHost),
		AuthDB:   cfg.MustString(cfgMongoAuthDB),
		DB:       cfg.MustString(cfgMongoDB),
		User:     cfg.MustString(cfgMongoUser),
		Password: cfg.MustString(cfgMongoPassword),
	}

	if err := db.RegMasterSession("startup", cfg.DB, cfg); err != nil {
		xenia.Println("Unable to initialize MongoDB")
		os.Exit(1)
	}

	db, err := db.NewMGO("", cfg.DB)
	if err != nil {
		xenia.Println("Unable to get MongoDB session")
		os.Exit(1)
	}
	defer db.CloseMGO("")

	xenia.AddCommand(
		cmdauth.GetCommands(db),
		cmddb.GetCommands(db),
		cmdquery.GetCommands(db),
		cmdscript.GetCommands(db),
		cmdregex.GetCommands(db),
	)
	xenia.Execute()
}
