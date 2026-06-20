package account

import (
	"context"

	domainAccount "github.com/aritxonly/deadlinerserver/internal/domain/account"
)

type Service interface {
	Register(context.Context, RegisterInput) (*domainAccount.SessionBundle, error)
	Login(context.Context, LoginInput) (*domainAccount.SessionBundle, error)
	RefreshSession(context.Context, RefreshSessionInput) (*domainAccount.SessionBundle, error)
}
