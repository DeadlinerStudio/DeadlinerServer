package sync

import (
	"context"

	domainSync "github.com/aritxonly/deadlinerserver/internal/domain/sync"
)

func (s *service) PushChanges(ctx context.Context, input PushChangesInput) (*domainSync.PushChangesResult, error) {
	accountUID, err := s.accountResolver.ResolveAccountUID(ctx)
	if err != nil {
		return nil, err
	}

	return s.domainService.PushChanges(ctx, domainSync.PushChangesCommand{
		AccountUID: accountUID,
		DeviceUID:  input.DeviceUID,
		BaseCursor: input.BaseCursor,
		Mutations:  input.Mutations,
	})
}
