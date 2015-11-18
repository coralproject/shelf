// Package tests provides support for tests.
package tests

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"os"
	"testing"

	"github.com/coralproject/shelf/pkg/log"
	"github.com/coralproject/shelf/pkg/mongo"
)

// Succeed is the Unicode codepoint for a check mark.
const Succeed = "\u2713"

// Failed is the Unicode codepoint for an X mark.
const Failed = "\u2717"

// logdest implements io.Writer and is the log package destination.
var logdest bytes.Buffer

// ResetLog can be called at the beginning of a test or example.
func ResetLog() { logdest.Reset() }

// DisplayLog can be called at the end of a test or example.
// It only prints the log contents if the -test.v flag is set.
func DisplayLog() {
	if !testing.Verbose() {
		return
	}
	logdest.WriteTo(os.Stdout)
}

func init() {
	// TODO: Need to read configuration.
	log.Init(&logdest, func() int { return log.DEV })
	mongo.InitMGO()
}

// NewRequest used to setup a request for mocking API calls with httptreemux.
func NewRequest(method, path string, body io.Reader) *http.Request {
	r, _ := http.NewRequest(method, path, body)
	u, _ := url.Parse(path)
	r.URL = u
	r.RequestURI = path
	return r
}
