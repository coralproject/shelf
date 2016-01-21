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

	err := mongo.Init(cfg)
	if err != nil {
		xenia.Println("Unable to initialize MongoDB")
		os.Exit(1)
	}

	xenia.AddCommand(
		cmdauth.GetCommands(),
		cmddb.GetCommands(),
		cmdquery.GetCommands(),
		cmdscript.GetCommands(),
		cmdregex.GetCommands(),
	)
	xenia.Execute()
}
