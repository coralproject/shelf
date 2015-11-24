// Package tests provides testing support.
package tests

import (
	"bytes"
	"os"
	"testing"

	"github.com/coralproject/shelf/pkg/cfg"
	"github.com/coralproject/shelf/pkg/log"
	"github.com/coralproject/shelf/pkg/srv/mongo"
)

// Success is a unicode codepoint for a check mark.
var Success = "\u2713"

// Failed is a unicode codepoint for a check X mark.
var Failed = "\u2717"

// logdash is the central buffer where all logs are stored.
var logdash bytes.Buffer

// ResetLog resets the contents of logdash.
func ResetLog() {
	logdash.Reset()
}

// DisplayLog writes the logdash data to standand out, if testing in verbose mode
// was turned on.
func DisplayLog() {
	if !testing.Verbose() {
		return
	}

	logdash.WriteTo(os.Stdout)
}

// Init is to be runned once. It initializes the necessary logs and mongodb
// connections for testing.
func Init() {
	initConfig()
	cfg.Init("shelf")

	// TODO: read current log mode from configuration.
	log.Init(&logdash, func() int { return log.DEV })
	err := mongo.InitMGO()
	if err != nil {
		log.Error("Test", "Init", err, "Test : Init : Error")
		DisplayLog()
		os.Exit(1)
	}
}

// initConfig initializes tests configuration for mongo.
func initConfig() {
	os.Setenv("SHELF_MONGO_HOST", "ds053894.mongolab.com:53894")
	os.Setenv("SHELF_MONGO_AUTHDB", "monsterbox")
	os.Setenv("SHELF_MONGO_DB", "monsterbox")
	os.Setenv("SHELF_MONGO_USER", "monsterbox")
	os.Setenv("SHELF_MONGO_PASS", "box")
}
