package service

import "errors"

var (
	ErrAccountAlreadyExists = errors.New("account already exists")
	ErrInvalidCredentials   = errors.New("invalid credentials")
	ErrInvalidRefreshToken  = errors.New("invalid refresh token")
	ErrExpiredRefreshToken  = errors.New("refresh token expired")
	ErrDeviceMismatch       = errors.New("device mismatch")
)
