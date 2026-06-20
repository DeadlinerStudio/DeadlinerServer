package auth

import "context"

type AccountResolver interface {
	ResolveAccountUID(context.Context) (string, error)
}
