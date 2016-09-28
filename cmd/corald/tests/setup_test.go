// Package tests implements users tests for the API layer.
package tests

import (
	"os"
	"testing"

	"github.com/ardanlabs/kit/web/app"
	"github.com/coralproject/shelf/cmd/corald/routes"
)

var a *app.App

func init() {
	// The call to API will force the init() function to initialize
	// cfg, log and mongodb.
	a = routes.API(true).(*app.App)
}

//==============================================================================

// TestMain helps to clean up the test data.
func TestMain(m *testing.M) {
	os.Exit(runTest(m))
}

// runTest initializes the environment for the tests and allows for
// the proper return code if the test fails or succeeds.
func runTest(m *testing.M) int {
	return m.Run()
}
