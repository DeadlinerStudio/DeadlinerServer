package port

import (
	"context"

	documentpkg "github.com/aritxonly/deadlinerserver/internal/domain/sync/document"
	statepkg "github.com/aritxonly/deadlinerserver/internal/domain/sync/state"
)

type HabitRepository interface {
	FindByDDLUID(ctx context.Context, accountID int64, ddlUID string) (*statepkg.HabitChange, error)
	Save(ctx context.Context, params SaveHabitParams) error
	ListAfterChangeID(
		ctx context.Context,
		accountID int64,
		afterChangeID int64,
		limit int,
		includeDeleted bool,
	) ([]statepkg.HabitChange, error)
}

type SaveHabitParams struct {
	AccountID          int64
	Deleted            bool
	Document           documentpkg.HabitDocument
	ServerVersion      statepkg.ServerVersion
	ClientVersion      *statepkg.LogicalVersion
	UpdatedByDeviceUID string
}
