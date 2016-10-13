package app

import (
	"fmt"
	"os"

	"github.com/ardanlabs/kit/cfg"
	"github.com/ardanlabs/kit/log"
)

// cfgLoggingLevel is the key for the logging level.
const cfgLoggingLevel = "LOGGING_LEVEL"

// Init sets up the configuration and logging systems.
func Init(p cfg.Provider) {
	if err := cfg.Init(p); err != nil {
		fmt.Println("Error initalizing configuration system", err)
		os.Exit(1)
	}

	// Init the log system.
	logLevel := func() int {
		ll, err := cfg.Int(cfgLoggingLevel)
		if err != nil {
			return log.USER
		}
		return ll
	}
	log.Init(os.Stderr, logLevel, log.Ldefault)
}
