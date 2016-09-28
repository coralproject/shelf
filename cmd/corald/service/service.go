package service

import (
	"net/http"
	"time"

	"github.com/ardanlabs/kit/web/app"
	"github.com/coralproject/shelf/internal/platform/auth"
	"github.com/pborman/uuid"
)

const downstreamTokenValidFor = 5 * time.Second

// Rewrite will add service request headers to the request and add other
// standards.
func Rewrite(c *app.Context) func(*http.Request) {

	f := func(r *http.Request) {

		// Extract the signer from the application context.
		signer, ok := c.App.Ctx["signer"].(auth.Signer)
		if !ok {
			return
		}

		// Create the new claims object for the token that we're using to send
		// downstream. This includes a unique identifier with the expiry set in the
		// future and set as not valid before the current date.

		now := time.Now()

		claims := map[string]interface{}{
			"jti": uuid.New(),
			"exp": int64(now.Add(downstreamTokenValidFor).Unix()),
			"iat": int64(now.Unix()),
			"nbf": int64(now.Add(-downstreamTokenValidFor).Unix()),
		}

		// Use the signer that was available on the application context to sign the
		// claims to get a token.
		token, err := signer(claims)
		if err != nil {
			return
		}

		// Include the token that we generated as the Authorization request header.
		r.Header.Set("Authorization", token)
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
