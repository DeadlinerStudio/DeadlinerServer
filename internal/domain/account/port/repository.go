package port

import (
	"context"

	entitypkg "github.com/aritxonly/deadlinerserver/internal/domain/account/entity"
)

type Repository interface {
	FindAccountByEmail(context.Context, string) (*entitypkg.Account, error)
	FindAccountByUID(context.Context, string) (*entitypkg.Account, error)
	FindAccountByID(context.Context, int64) (*entitypkg.Account, error)
	FindSessionByRefreshTokenHash(context.Context, string) (*entitypkg.Session, error)
	SaveAccount(context.Context, *entitypkg.Account) error
	SaveDevice(context.Context, *entitypkg.Device) error
	SaveSession(context.Context, *entitypkg.Session) error
}
