package midware

import (
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/ardanlabs/kit/cfg"
	"github.com/ardanlabs/kit/log"
	"github.com/ardanlabs/kit/web/app"
	"github.com/dgrijalva/jwt-go"
)

var (
	// ErrInvalidToken is returned when the token provided is not valid.
	ErrInvalidToken = errors.New("invalid token")

	// ErrInvalidClaims is returned when the claims inside a valid token are not
	// valid.
	ErrInvalidClaims = errors.New("invalid claims")
)

// cfgAuthPublicKey is the key for which the actual base64 + PEM encoded public
// RSA key is stored.
const cfgAuthPublicKey = "AUTH_PUBLIC_KEY"

// authOff is used when authentication is turned off by not providing a public
// key in the environment.
func authOff(h app.Handler) app.Handler {
	f := func(c *app.Context) error {

		// Log out the process for verbosity.
		log.Dev(c.SessionID, "Auth", "Started")
		log.Dev(c.SessionID, "Auth", "Authentication Off")
		log.Dev(c.SessionID, "Auth", "Completed")
		return h(c)
	}

	return f
}

// Auth handles token authentication.
func Auth() (app.Middleware, error) {

	// Load in the public key to validate the JWT tokens.
	publicKeyBase64Str, err := cfg.String(cfgAuthPublicKey)
	if err != nil {
		return authOff, nil
	}

	// Our public key has been encoded from a PEM encoded public RSA key into this
	// publicKeyBase64Str. We need to decode the base64 string in order to get the
	// PEM encoded certificate back out.
	publicKeyPEM, err := base64.StdEncoding.DecodeString(publicKeyBase64Str)
	if err != nil {
		log.Error("startup", "Auth", err, "Can not setup Auth middleware")
		return nil, err
	}

	// Now that we have our PEM encoded public RSA key, we can parse it using the
	// methods built into the jwt librairy into something we can use to verify the
	// incomming JWT's.
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyPEM)
	if err != nil {
		log.Error("startup", "Auth", err, "Can not setup Auth middleware")
		return nil, err
	}

	log.Dev("startup", "Auth", "Initalizing Auth")

	// Create the middleware to actually return.
	m := func(h app.Handler) app.Handler {

		// Create the handler that we should return as a part of the middleware
		// chain.
		f := func(c *app.Context) error {
			log.Dev(c.SessionID, "Auth", "Started")

			// Extract the token from the Authorization header provided on the request.
			tokenString := c.Request.Header.Get("Authorization")

			if tokenString == "" {
				log.Error(c.SessionID, "Auth", ErrInvalidToken, "No token on request")
				return ErrInvalidToken
			}

			token, err := jwt.ParseWithClaims(tokenString, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {

				// Don't forget to validate the alg is what you expect.
				if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
					return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
				}

				// Return with the public key that was provided in the config.
				return publicKey, nil
			})

			// Return with the error if there was an issue parsing the token.
			if err != nil {
				log.Error(c.SessionID, "Auth", err, "Token could not be parsed")
				return ErrInvalidToken
			}

			// Return with an error if the token is not valid.
			if !token.Valid {
				log.Error(c.SessionID, "Auth", ErrInvalidToken, "Token not valid")
				return ErrInvalidToken
			}

			// Ensure that the claims that are inside the token are indeed the MapClaims
			// that we expect.
			claims, ok := token.Claims.(*jwt.MapClaims)
			if !ok {
				log.Error(c.SessionID, "Auth", ErrInvalidClaims, "Claims not valid")
				return ErrInvalidClaims
			}

			// Validate that all the parameters we expect are correct, noteably, the
			// expiry date, and not before claims should be verified.
			if err := claims.Valid(); err != nil {
				log.Error(c.SessionID, "Auth", err, "Claims not valid")
				return ErrInvalidToken
			}

			// Add the claims to the context.
			c.Ctx["claims"] = claims

			log.Dev(c.SessionID, "Auth", "Completed : Valid")
			return h(c)
		}

		return f
	}

	return m, nil
}
