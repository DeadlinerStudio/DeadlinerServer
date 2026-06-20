package port

import (
	"context"

	documentpkg "github.com/aritxonly/deadlinerserver/internal/domain/sync/document"
	statepkg "github.com/aritxonly/deadlinerserver/internal/domain/sync/state"
)

type DeadlineRepository interface {
	FindByUID(ctx context.Context, accountID int64, uid string) (*statepkg.DeadlineChange, error)
	Save(ctx context.Context, params SaveDeadlineParams) error
	ListAfterChangeID(
		ctx context.Context,
		accountID int64,
		afterChangeID int64,
		limit int,
		includeDeleted bool,
	) ([]statepkg.DeadlineChange, error)
}

type SaveDeadlineParams struct {
	AccountID          int64
	Deleted            bool
	Document           documentpkg.DeadlineDocument
	ServerVersion      statepkg.ServerVersion
	ClientVersion      *statepkg.LogicalVersion
	UpdatedByDeviceUID string
}
