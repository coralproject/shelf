package auth

import (
	"crypto/ecdsa"
	"encoding/base64"

	jwt "github.com/dgrijalva/jwt-go"
)

// DecodePublicKey pulls the public key out of the string passed by first
// decoding from base64 and parsing the PEM encoding.
func DecodePublicKey(publicKeyBase64Str string) (*ecdsa.PublicKey, error) {

	// Our public key has been encoded from a PEM encoded public ECDSA key into this
	// publicKeyBase64Str. We need to decode the base64 string in order to get the
	// PEM encoded certificate back out.
	publicKeyPEM, err := base64.StdEncoding.DecodeString(publicKeyBase64Str)
	if err != nil {
		return nil, err
	}

	// Now that we have our PEM encoded public ECDSA key, we can parse it using the
	// methods built into the jwt librairy into something we can use to verify the
	// incomming JWT's.
	publicKey, err := jwt.ParseECPublicKeyFromPEM(publicKeyPEM)
	if err != nil {
		return nil, err
	}

	return publicKey, nil
}

// DecodePrivateKey pulls the private key out of the string passed by first
// decoding from base64 and parsing the PEM encoding.
func DecodePrivateKey(privateKeyBase64Str string) (*ecdsa.PrivateKey, error) {

	// Our private key has been encoded from a PEM encoded private ECDSA key into this
	// privateKeyBase64Str. We need to decode the base64 string in order to get the
	// PEM encoded certificate back out.
	privateKeyPEM, err := base64.StdEncoding.DecodeString(privateKeyBase64Str)
	if err != nil {
		return nil, err
	}

	// Now that we have our PEM encoded private ECDSA key, we can parse it using the
	// methods built into the jwt librairy into something we can use to verify the
	// incomming JWT's.
	privateKey, err := jwt.ParseECPrivateKeyFromPEM(privateKeyPEM)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}
