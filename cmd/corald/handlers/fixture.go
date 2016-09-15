//go:generate go-bindata -pkg fixtures -o fixtures/fixtures.go fixtures/

package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ardanlabs/kit/web/app"
	"github.com/coralproject/shelf/cmd/corald/handlers/fixtures"
)

// Error will simply return the error to the calling request stack.
func Error(err error) app.Handler {

	// Create this really, really simple handler.
	return app.Handler(func(c *app.Context) error {
		return err
	})
}

// Fixture will serve a JSON payload as the endpoint response.
func Fixture(name string, code int) app.Handler {

	// Load the fixtures bytes into the byte slice.
	fixtureBytes, err := fixtures.Asset(fmt.Sprintf("fixtures/%s.json", name))
	if err != nil {
		return Error(err)
	}

	// Decode the payload into a blank interface{}, this could be an array or a
	// object so we can serve the data.
	var fixture interface{}
	err = json.Unmarshal(fixtureBytes, &fixture)
	if err != nil {
		return Error(err)
	}

	// Provide the handler which will just serve the json fixture out to the
	// client.
	return app.Handler(func(c *app.Context) error {

		// Respond with the fixture and the code.
		c.Respond(fixture, code)
		return nil
	})
}

// NoContent simply responds with a HTTP Status Code of 204.
func NoContent(c *app.Context) error {
	c.Respond(nil, http.StatusNoContent)
	return nil
}
