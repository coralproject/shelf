package service

import (
	"net/http"
	"time"

	"github.com/ardanlabs/kit/log"
	"github.com/ardanlabs/kit/web/app"
	"github.com/coralproject/shelf/internal/platform/auth"
	"github.com/pborman/uuid"
)

// SignServiceRequest will take a given signer, and adds a Authorization header
// with the token that is generated from the signer.
func SignServiceRequest(context interface{}, signer auth.Signer, r *http.Request) error {
	log.Dev(context, "SignServiceRequest", "Started")

	// downstreamTokenValidFor describes the time that a token is valid for before
	// expiring. It is also used to describe the time before now that the token
	// should be valid before resulting in a 2 * downstreamTokenValidFor window that
	// the token is valid for. This window should account for enough backpressure or
	// time skew to a resonable degree.
	const downstreamTokenValidFor = 5 * time.Second

	// Create the new claims object for the token that we're using to send
	// downstream. This includes a unique identifier with the expiry set in the
	// future and set as not valid before the current date.

	now := time.Now()

	// The claims currently just use the reserved keys, in the future we can add
	// any other POD into this map that can be used downstream.
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
		log.Error(context, "SignServiceRequest", err, "Can't sign the claims object")
		return err
	}

	// Include the token that we generated as the Authorization request header.
	r.Header.Set("Authorization", token)

	log.Dev(context, "SignServiceRequest", "Completed")
	return nil
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
