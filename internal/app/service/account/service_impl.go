package account

import (
	"context"

	domainAccount "github.com/aritxonly/deadlinerserver/internal/domain/account"
)

type service struct {
	domainService domainAccount.Service
}

func NewService(domainService domainAccount.Service) Service {
	return &service{domainService: domainService}
}

func (s *service) Register(ctx context.Context, input RegisterInput) (*domainAccount.SessionBundle, error) {
	return s.domainService.Register(ctx, domainAccount.RegisterCommand{
		Email:       input.Email,
		Password:    input.Password,
		DisplayName: input.DisplayName,
		DeviceUID:   input.DeviceUID,
		DeviceName:  input.DeviceName,
		Platform:    input.Platform,
	})
}

func (s *service) Login(ctx context.Context, input LoginInput) (*domainAccount.SessionBundle, error) {
	return s.domainService.Login(ctx, domainAccount.LoginCommand{
		Email:      input.Email,
		Password:   input.Password,
		DeviceUID:  input.DeviceUID,
		DeviceName: input.DeviceName,
		Platform:   input.Platform,
	})
}

func (s *service) RefreshSession(
	ctx context.Context,
	input RefreshSessionInput,
) (*domainAccount.SessionBundle, error) {
	return s.domainService.RefreshSession(ctx, domainAccount.RefreshSessionCommand{
		RefreshToken: input.RefreshToken,
		DeviceUID:    input.DeviceUID,
	})
}
