package auth

import (
	"net/http"
	"time"

	"github.com/ardanlabs/kit/log"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/pborman/uuid"
)

// Signer is a function that can be used to sign claims with and generate a
// signed JWT token from them.
type Signer func(claims map[string]interface{}) (string, error)

// NewSigner will return a signer that can be used to sign tokens for a given
// set of claims.
func NewSigner(privateKeyBase64Str string) (Signer, error) {
	privateKey, err := DecodePrivateKey(privateKeyBase64Str)
	if err != nil {
		return nil, err
	}

	// This is the signer function that provides the support for signing tokens.
	f := func(claims map[string]interface{}) (string, error) {

		// Create the new JWT token with the ES384 signing method.
		token := jwt.NewWithClaims(jwt.SigningMethodES384, jwt.MapClaims(claims))

		// Actually sign the token.
		return token.SignedString(privateKey)
	}

	return f, nil
}

// SignRequest will take a given signer, and adds a Authorization header
// with the token that is generated from the signer.
func SignRequest(context interface{}, signer Signer, claims map[string]interface{}, r *http.Request) error {
	log.Dev(context, "SignRequest", "Started")

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

	// Ensure that the claims object is not nil, if it is, then we should create a
	// fresh one.
	if claims == nil {
		claims = map[string]interface{}{}
	}

	// The claims currently just use the reserved keys, in the future we can add
	// any other POD into this map that can be used downstream.
	claims["jti"] = uuid.New()
	claims["exp"] = int64(now.Add(downstreamTokenValidFor).Unix())
	claims["iat"] = int64(now.Unix())
	claims["nbf"] = int64(now.Add(-downstreamTokenValidFor).Unix())

	// Use the signer that was available on the application context to sign the
	// claims to get a token.
	token, err := signer(claims)
	if err != nil {
		log.Error(context, "SignRequest", err, "Can't sign the claims object")
		return err
	}

	// Include the token that we generated as the Authorization request header.
	r.Header.Set("Authorization", token)

	log.Dev(context, "SignRequest", "Completed")
	return nil
}
