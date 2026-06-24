package adminconfig

import "context"

type Service interface {
	GetSnapshot(context.Context) (*Snapshot, error)
	Update(context.Context, UpdateInput) (*Snapshot, error)
}
