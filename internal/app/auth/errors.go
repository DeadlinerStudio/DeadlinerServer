package auth

import "errors"

var ErrUnauthenticated = errors.New("account uid metadata is required")
