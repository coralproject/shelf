package service

import (
	"net/http"

	"github.com/ardanlabs/kit/web/app"
	"github.com/coralproject/shelf/internal/platform/auth"
)

// SignServiceRequest signs a request with the claims necessary to authenticate
// with downstream services.
func SignServiceRequest(context interface{}, signer auth.Signer, r *http.Request) error {
	claims := map[string]interface{}{}

	return auth.SignRequest(context, signer, claims, r)
}

// Rewrite will add service request headers to the request and add other
// standards.
func Rewrite(c *app.Context) func(*http.Request) {

	f := func(r *http.Request) {

		// Extract the signer from the application context.
		signer, ok := c.App.Ctx["signer"].(auth.Signer)
		if !ok {
			return
		}

		// Sign the service request with the signer.
		if err := SignServiceRequest(c, signer, r); err != nil {
			return
		}

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
