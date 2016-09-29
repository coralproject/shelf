package auth

import jwt "github.com/dgrijalva/jwt-go"

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
