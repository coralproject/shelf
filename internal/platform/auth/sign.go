package auth

import jwt "github.com/dgrijalva/jwt-go"

// NewSigner will return a signer that can be used to sign tokens for a given
// set of claims.
func NewSigner(privateKeyBase64Str string) (Signer, error) {
	privateKey, err := DecodePrivateKey(privateKeyBase64Str)
	if err != nil {
		return nil, err
	}

	f := func(claims map[string]interface{}) (string, error) {

		// Create the new JWT token with the RS512 signing method.
		token := jwt.NewWithClaims(jwt.SigningMethodES512, jwt.MapClaims(claims))

		// And actually sign the token.
		return token.SignedString(privateKey)
	}

	return f, nil
}

// Signer is a function that can be used to sign claims with and generate a
// signed JWT token from them.
type Signer func(claims map[string]interface{}) (string, error)
