package auth

import "errors"

// ErrInvalidToken is returned when the token provided is not valid.
var ErrInvalidToken = errors.New("invalid token")
