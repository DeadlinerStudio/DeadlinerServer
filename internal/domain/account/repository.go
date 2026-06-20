package account

import "context"

type Repository interface {
	FindAccountByEmail(context.Context, string) (*Account, error)
	FindAccountByUID(context.Context, string) (*Account, error)
	SaveAccount(context.Context, *Account) error
	SaveDevice(context.Context, *Device) error
	SaveSession(context.Context, *Session) error
}

