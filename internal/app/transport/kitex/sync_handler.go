package kitex

import (
	"context"
	"errors"
	"fmt"

	syncmapper "github.com/aritxonly/deadlinerserver/internal/app/transport/kitex/syncmapper"
	v1 "github.com/aritxonly/deadlinerserver/kitex_gen/deadliner/v1"
)

func (h *Handler) PullChanges(ctx context.Context, req *v1.PullChangesRequest) (*v1.PullChangesResponse, error) {
	if h.syncService == nil {
		return nil, errors.New("sync service is not configured")
	}

	result, err := h.syncService.PullChanges(ctx, syncmapper.ToPullChangesInput(req))
	if err != nil {
		return nil, fmt.Errorf("pull changes: %w", err)
	}

	return syncmapper.ToPullChangesResponse(result), nil
}

func (h *Handler) PushChanges(ctx context.Context, req *v1.PushChangesRequest) (*v1.PushChangesResponse, error) {
	if h.syncService == nil {
		return nil, errors.New("sync service is not configured")
	}

	input, err := syncmapper.ToPushChangesInput(req)
	if err != nil {
		return nil, fmt.Errorf("map push changes request: %w", err)
	}

	result, err := h.syncService.PushChanges(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("push changes: %w", err)
	}

	return syncmapper.ToPushChangesResponse(result), nil
}

var _ v1.DeadlinerService = (*Handler)(nil)
