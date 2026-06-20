package port

import (
	"context"

	statepkg "github.com/aritxonly/deadlinerserver/internal/domain/sync/state"
)

type SyncChangeRepository interface {
	Append(ctx context.Context, params AppendSyncChangeParams) (*statepkg.SyncChange, error)
	ListAfterChangeID(ctx context.Context, accountID, afterChangeID int64, limit int) ([]statepkg.SyncChange, error)
}

type AppendSyncChangeParams struct {
	AccountID   int64
	DeviceUID   string
	MutationID  string
	EntityKind  string
	EntityUID   string
	Action      string
	Payload     []byte
	CommittedAt string
}
