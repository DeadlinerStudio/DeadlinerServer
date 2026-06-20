package service

import (
	"context"

	commandpkg "github.com/aritxonly/deadlinerserver/internal/domain/account/command"
	entitypkg "github.com/aritxonly/deadlinerserver/internal/domain/account/entity"
)

type Service interface {
	Register(context.Context, commandpkg.RegisterCommand) (*entitypkg.SessionBundle, error)
	Login(context.Context, commandpkg.LoginCommand) (*entitypkg.SessionBundle, error)
	RefreshSession(context.Context, commandpkg.RefreshSessionCommand) (*entitypkg.SessionBundle, error)
}
