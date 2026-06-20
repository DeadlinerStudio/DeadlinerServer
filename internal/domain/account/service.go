package account

import "context"

type Service interface {
	Register(context.Context, RegisterCommand) (*SessionBundle, error)
	Login(context.Context, LoginCommand) (*SessionBundle, error)
	RefreshSession(context.Context, RefreshSessionCommand) (*SessionBundle, error)
}
