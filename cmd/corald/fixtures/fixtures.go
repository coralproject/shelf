//go:generate go-bindata -pkg json -ignore json/json.gen.go -o json/json.gen.go json/...

package fixtures

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ardanlabs/kit/web"
	fixturesjson "github.com/coralproject/shelf/cmd/corald/fixtures/json"
)

// Load unmarshals the specified fixture into the provided
// data value.
func Load(name string, v interface{}) error {

	// Load the fixtures bytes into the byte slice.
	fixtureBytes, err := fixturesjson.Asset(fmt.Sprintf("json/%s.json", name))
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
func Error(err error) web.Handler {

	// Create this really, really simple handler.
	return web.Handler(func(c *web.Context) error {
		return err
	})
}

// NoContent simply responds with a HTTP Status Code of 204.
func NoContent(c *web.Context) error {
	c.Respond(nil, http.StatusNoContent)
	return nil
}

// Handler will serve a JSON payload as the endpoint response.
func Handler(name string, code int) web.Handler {

	// Load the fixture for response.
	var fixture interface{}
	if err := Load(name, &fixture); err != nil {
		return Error(err)
	}

	// Provide the handler which will just serve the json fixture
	// out to the client.
	return web.Handler(func(c *web.Context) error {

		// Respond with the fixture and the code.
		c.Respond(fixture, code)
		return nil
	})
}
