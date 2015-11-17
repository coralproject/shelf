package crypto

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/scrypt"
)

// Entity represents an interface for authenticating varied entities.
type Entity interface {
	Pwd() ([]byte, error)
	Salt() ([]byte, error)
}

// BcryptHash generates a hash using the bcrypt encoding standard from a provided
// []byte of the password string.
func BcryptHash(pwd []byte) (string, error) {
	crypted, err := bcrypt.GenerateFromPassword(pwd, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(crypted), nil
}

// CompareBcryptHash compares a given password to a bcrypted hash key. It returns
// a non-nil error if its not a match.
func CompareBcryptHash(hash []byte, pwd []byte) error {
	return bcrypt.CompareHashAndPassword(hash, pwd)
}

// SignedSHAHash generates a signed hash using the SHA256 encoding standard from a given
// password and salt.
func SignedSHAHash(pwd []byte, salt []byte) ([]byte, error) {
	key, err := scrypt.Key([]byte(pwd), []byte(salt), 16384, 8, 1, 32)
	if err != nil {
		return nil, err
	}

	// Append salt into pwd.
	hm := hmac.New(sha256.New, key)
	hm.Write(pwd)
	return hm.Sum(nil), nil
}

// SHAHash generates a hash using the SHA256 encoding standard from a given hash
// string.
func SHAHash(pwd string) string {
	pwdBytes := []byte(pwd)
	shm := sha256.New()
	shm.Write(pwdBytes)
	return hex.EncodeToString(shm.Sum(nil))
}

// Base64Token takes a token and returns a base64
// version for safe transimission.
func Base64Token(token []byte) string {
	return base64.StdEncoding.EncodeToString(token)
}

// TokenForEntity generates a unique signed SHA256 token for a given entity.
func TokenForEntity(e Entity) ([]byte, error) {
	var err error
	var pwd, salt []byte

	pwd, err = e.Pwd()
	if err != nil {
		return nil, err
	}

	salt, err = e.Salt()
	if err != nil {
		return nil, err
	}

	return SignedSHAHash(pwd, salt)
}

// IsTokenValidForEntity validates if a given token is correct for an entity.
// It returns a non-nil error if the given token is not a match.
func IsTokenValidForEntity(e Entity, token string) error {
	decoded, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return err
	}

	// Generate the unique token for the entity.
	eToken, err := TokenForEntity(e)
	if err != nil {
		return err
	}

	if hmac.Equal(decoded, eToken) == false {
		return fmt.Errorf("Invalid Token")
	}

	return nil
}
