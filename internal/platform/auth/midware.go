package auth

import (
	"errors"
	"fmt"

	"github.com/ardanlabs/kit/log"
	"github.com/ardanlabs/kit/web/app"
	jwt "github.com/dgrijalva/jwt-go"
)

// ErrInvalidToken is returned when the token provided is not valid.
var ErrInvalidToken = errors.New("invalid token")

// MidwareOpts describes the options for configuring the Midware.
type MidwareOpts struct {

	// AllowQueryString is true when we want to allow accessing the tokenString
	// from the query string as a fallback.
	AllowQueryString bool
}

// Midware handles token authentication for external authentication
// sources.
func Midware(publicKeyBase64Str string, config MidwareOpts) (app.Middleware, error) {
	publicKey, err := DecodePublicKey(publicKeyBase64Str)
	if err != nil {
		log.Error("startup", "auth : Midware", err, "Can not decode the public key base64 encoding")
		return nil, err
	}

	// Create the middleware to actually return.
	m := func(h app.Handler) app.Handler {

		// Create the handler that we should return as a part of the middleware
		// chain.
		f := func(c *app.Context) error {
			log.Dev(c.SessionID, "auth : Midware", "Started")

			// Extract the token from the Authorization header provided on the request.
			tokenString := c.Request.Header.Get("Authorization")

			// In the event that the request does not have a header key for the
			// Authorization header, and we are allowed to check the query string, then
			// we need to try and access it from the URL query parameters.
			if tokenString == "" && config.AllowQueryString {
				tokenString = c.Request.URL.Query().Get("access_token")
			}

			if tokenString == "" {
				log.Error(c.SessionID, "auth : Midware", ErrInvalidToken, "No token on request")
				return app.ErrNotAuthorized
			}

			// This describes the key validation function to provide the certificate
			// to validate the signature on the passed in JWT.
			keyValidation := func(token *jwt.Token) (interface{}, error) {

				// Don't forget to validate the alg is what you expect.
				if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
					return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
				}

				// Return with the public key that was provided in the config.
				return publicKey, nil
			}

			// Here we actually parse/verify the signature on the JWT and extract the
			// claims.
			token, err := jwt.ParseWithClaims(tokenString, &jwt.MapClaims{}, keyValidation)
			if err != nil {
				log.Error(c.SessionID, "auth : Midware", err, "Token could not be parsed")
				return app.ErrNotAuthorized
			}

			// Return with an error if the token is not valid.
			if !token.Valid {
				log.Error(c.SessionID, "auth : Midware", ErrInvalidToken, "Token not valid")
				return app.ErrNotAuthorized
			}

			// Ensure that the claims that are inside the token are indeed the MapClaims
			// that we expect.
			claims, ok := token.Claims.(*jwt.MapClaims)
			if !ok {
				log.Error(c.SessionID, "auth : Midware", ErrInvalidToken, "Claims not valid")
				return app.ErrNotAuthorized
			}

			// Validate that all the parameters we expect are correct, noteably, the
			// expiry date, and not before claims should be verified.
			if err := claims.Valid(); err != nil {
				log.Error(c.SessionID, "auth : Midware", err, "Claims not valid")
				return app.ErrNotAuthorized
			}

			// Add the claims to the context.
			c.Ctx["claims"] = claims

			log.Dev(c.SessionID, "auth : Midware", "Completed : Valid")
			return h(c)
		}

		return f
	}

	return m, nil
}
