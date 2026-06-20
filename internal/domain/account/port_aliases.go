package account

import portpkg "github.com/aritxonly/deadlinerserver/internal/domain/account/port"

type Repository = portpkg.Repository
type PasswordHasher = portpkg.PasswordHasher
type RefreshTokenHasher = portpkg.RefreshTokenHasher
type AccessTokenCodec = portpkg.AccessTokenCodec
type TokenGenerator = portpkg.TokenGenerator
type Clock = portpkg.Clock
