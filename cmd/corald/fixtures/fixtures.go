//go:generate go-bindata -pkg fixtures -o assets.go json/

package fixtures

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ardanlabs/kit/web/app"
)

// Load unmarshals the specified fixture into the provided
// data value.
func Load(name string, v interface{}) error {

	// Load the fixtures bytes into the byte slice.
	fixtureBytes, err := Asset(fmt.Sprintf("json/%s.json", name))
	if err != nil {
		return err
	}

	// Unmarshal the fixture into the provided data value.
	err = json.Unmarshal(fixtureBytes, &v)
	if err != nil {
		return err
	}

	return nil
}

// Error will simply return the error to the calling request stack.
func Error(err error) app.Handler {

	// Create this really, really simple handler.
	return app.Handler(func(c *app.Context) error {
		return err
	})
}

// NoContent simply responds with a HTTP Status Code of 204.
func NoContent(c *app.Context) error {
	c.Respond(nil, http.StatusNoContent)
	return nil
}

// Handler will serve a JSON payload as the endpoint response.
func Handler(name string, code int) app.Handler {

	// Load the fixture for response.
	var fixture interface{}
	if err := Load(name, &fixture); err != nil {
		return Error(err)
	}

	// Provide the handler which will just serve the json fixture
	// out to the client.
	return app.Handler(func(c *app.Context) error {

		// Respond with the fixture and the code.
		c.Respond(fixture, code)
		return nil
	})
}
