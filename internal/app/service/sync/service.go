package sync

import (
	"context"

	domainSync "github.com/aritxonly/deadlinerserver/internal/domain/sync"
)

type Service interface {
	PullChanges(context.Context, PullChangesInput) (*domainSync.PullChangesResult, error)
	PushChanges(context.Context, PushChangesInput) (*domainSync.PushChangesResult, error)
}
