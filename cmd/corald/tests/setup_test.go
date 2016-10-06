// Package tests implements users tests for the API layer.
package tests

import (
	"net/http"
	"os"
	"testing"

	"github.com/coralproject/shelf/cmd/corald/routes"
)

var a http.Handler

//==============================================================================

// TestMain helps to clean up the test data.
func TestMain(m *testing.M) {
	os.Exit(runTest(m))
}

// runTest initializes the environment for the tests and allows for
// the proper return code if the test fails or succeeds.
func runTest(m *testing.M) int {
	a = routes.API()

	return m.Run()
}
