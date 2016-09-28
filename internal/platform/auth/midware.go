package auth

import (
	"fmt"

	"github.com/ardanlabs/kit/log"
	"github.com/ardanlabs/kit/web/app"
	jwt "github.com/dgrijalva/jwt-go"
)

// Midware handles token authentication for external authentication
// sources.
func Midware(publicKeyBase64Str string) (app.Middleware, error) {
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

			if tokenString == "" {
				log.Error(c.SessionID, "auth : Midware", ErrInvalidToken, "No token on request")
				return app.ErrNotAuthorized
			}

			token, err := jwt.ParseWithClaims(tokenString, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {

				// Don't forget to validate the alg is what you expect.
				if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
					return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
				}

				// Return with the public key that was provided in the config.
				return publicKey, nil
			})

			// Return with the error if there was an issue parsing the token.
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
