package port

import entitypkg "github.com/aritxonly/deadlinerserver/internal/domain/account/entity"

type AccessTokenCodec interface {
	Sign(claims entitypkg.AccessTokenClaims) (string, error)
	Parse(token string) (*entitypkg.AccessTokenClaims, error)
}
