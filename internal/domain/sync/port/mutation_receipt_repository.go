package port

import (
	"context"

	statepkg "github.com/aritxonly/deadlinerserver/internal/domain/sync/state"
)

type MutationReceiptRepository interface {
	Find(ctx context.Context, accountID int64, deviceUID, mutationID string) (*statepkg.MutationReceipt, error)
	Save(ctx context.Context, receipt *statepkg.MutationReceipt) error
}
