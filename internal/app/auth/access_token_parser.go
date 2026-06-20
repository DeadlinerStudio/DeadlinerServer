package auth

import domainAccount "github.com/aritxonly/deadlinerserver/internal/domain/account"

type AccessTokenParser interface {
	Parse(token string) (*domainAccount.AccessTokenClaims, error)
}
