// Package tests provides testing support.
package tests

import (
	"bytes"
	"os"
	_ "testing"

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
	// if !testing.Verbose() {
	// 	return
	// }

	logdash.WriteTo(os.Stdout)
}

// Init is to be runned once. It initializes the necessary logs and mongodb
// connections for testing.
func Init() {
	cfg.Init("shelf")
	log.Init(&logdash, func() int { return log.DEV })

	if err := mongo.InitMGO(); err != nil {
		log.Error("Test", "Init", err, "Completed")
		DisplayLog()
		os.Exit(1)
	}
}
