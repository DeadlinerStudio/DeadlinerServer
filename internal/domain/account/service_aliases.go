package account

import (
	"time"

	portpkg "github.com/aritxonly/deadlinerserver/internal/domain/account/port"
	servicepkg "github.com/aritxonly/deadlinerserver/internal/domain/account/service"
)

type Service = servicepkg.Service

var (
	ErrAccountAlreadyExists = servicepkg.ErrAccountAlreadyExists
	ErrInvalidCredentials   = servicepkg.ErrInvalidCredentials
	ErrInvalidRefreshToken  = servicepkg.ErrInvalidRefreshToken
	ErrExpiredRefreshToken  = servicepkg.ErrExpiredRefreshToken
	ErrDeviceMismatch       = servicepkg.ErrDeviceMismatch
)

func NewService(
	repo portpkg.Repository,
	passwordHasher portpkg.PasswordHasher,
	refreshTokenHasher portpkg.RefreshTokenHasher,
	accessTokenCodec portpkg.AccessTokenCodec,
	tokenGenerator portpkg.TokenGenerator,
	clock portpkg.Clock,
	accessTokenTTL time.Duration,
	refreshTokenTTL time.Duration,
) Service {
	return servicepkg.NewService(
		repo,
		passwordHasher,
		refreshTokenHasher,
		accessTokenCodec,
		tokenGenerator,
		clock,
		accessTokenTTL,
		refreshTokenTTL,
	)
}
