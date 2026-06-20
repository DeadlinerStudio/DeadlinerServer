package kitex

import (
	"context"
	"errors"
	"fmt"

	accountmapper "github.com/aritxonly/deadlinerserver/internal/app/transport/kitex/accountmapper"
	v1 "github.com/aritxonly/deadlinerserver/kitex_gen/deadliner/v1"
)

func (h *Handler) Register(ctx context.Context, req *v1.RegisterRequest) (*v1.RegisterResponse, error) {
	if h.accountService == nil {
		return nil, errors.New("account service is not configured")
	}

	result, err := h.accountService.Register(ctx, accountmapper.ToRegisterInput(req))
	if err != nil {
		return nil, fmt.Errorf("register: %w", err)
	}

	return accountmapper.ToRegisterResponse(result), nil
}

func (h *Handler) Login(ctx context.Context, req *v1.LoginRequest) (*v1.LoginResponse, error) {
	if h.accountService == nil {
		return nil, errors.New("account service is not configured")
	}

	result, err := h.accountService.Login(ctx, accountmapper.ToLoginInput(req))
	if err != nil {
		return nil, fmt.Errorf("login: %w", err)
	}

	return accountmapper.ToLoginResponse(result), nil
}

func (h *Handler) RefreshSession(
	ctx context.Context,
	req *v1.RefreshSessionRequest,
) (*v1.RefreshSessionResponse, error) {
	if h.accountService == nil {
		return nil, errors.New("account service is not configured")
	}

	result, err := h.accountService.RefreshSession(ctx, accountmapper.ToRefreshSessionInput(req))
	if err != nil {
		return nil, fmt.Errorf("refresh session: %w", err)
	}

	return accountmapper.ToRefreshSessionResponse(result), nil
}
