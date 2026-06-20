package sync

import (
	"context"

	domainSync "github.com/aritxonly/deadlinerserver/internal/domain/sync"
)

func (s *service) PullChanges(ctx context.Context, input PullChangesInput) (*domainSync.PullChangesResult, error) {
	accountUID, err := s.accountResolver.ResolveAccountUID(ctx)
	if err != nil {
		return nil, err
	}

	return s.domainService.PullChanges(ctx, domainSync.PullChangesCommand{
		AccountUID:    accountUID,
		DeviceUID:     input.DeviceUID,
		Cursor:        input.Cursor,
		Limit:         s.normalizePullLimit(input.Limit),
		IncludeDelete: input.IncludeDeleted,
	})
}
