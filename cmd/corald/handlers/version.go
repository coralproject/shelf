// Package handlers contains the handler logic for processing requests.
package handlers

import (
	"net/http"

	"github.com/ardanlabs/kit/web"
)

// verHandle maintains the set of handlers for the ver api.
type verHandle struct {
	GitRevision string
	GitVersion  string
	BuildDate   string
	IntVersion  string
}

// Version fronts the access to the ver service functionality.
var Version verHandle

//==============================================================================

// List returns the version information.
// 200 Success, 404 Not Found, 500 Internal
func (verHandle) List(c *web.Context) error {
	c.Respond(Version, http.StatusOK)
	return nil
}
