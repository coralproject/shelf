package service

import (
	"net/http"

	"github.com/ardanlabs/kit/web/app"
)

// Rewrite will add service request headers to the request and add other
// standards.
func Rewrite(c *app.Context) func(*http.Request) {

	f := func(r *http.Request) {
		// TODO: Add authentication headers here
	}

	return f
}

// RewritePath will rewrite the path given a PathRewriter and return the request
// director function.
func RewritePath(c *app.Context, targetPath string) func(*http.Request) {

	f := func(r *http.Request) {

		// Rewrite the request for the services which will add authentication
		// headers and/or other default headers.
		Rewrite(c)(r)

		// Update the target path.
		r.URL.Path = targetPath
	}

	return f
}
