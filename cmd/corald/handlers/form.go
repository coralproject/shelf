package handlers

import (
	"net/http"

	"github.com/ardanlabs/kit/web/app"
	"github.com/coralproject/shelf/cmd/corald/fixtures"
)

// formHandle maintains the set of handlers for the form api.
type formHandle struct{}

// Form fronts the access to the form service functionality.
var Form formHandle

//==============================================================================

// List returns all the existing forms in the system.
// 200 Success, 404 Not Found, 500 Internal
func (formHandle) List(c *app.Context) error {

	// Create the projected data value to return.
	var forms []struct {
		ID string `json:"id"`
	}

	// Load the array of forms.
	if err := fixtures.Load("forms", &forms); err != nil {
		return app.ErrNotFound
	}

	c.Respond(forms, http.StatusOK)
	return nil
}
