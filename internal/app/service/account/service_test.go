package account

import (
	"context"
	"testing"

	domainAccount "github.com/aritxonly/deadlinerserver/internal/domain/account"
)

func TestRegisterForwardsCommand(t *testing.T) {
	domainService := &fakeDomainAccountService{
		registerResult: &domainAccount.SessionBundle{},
	}
	service := NewService(domainService)

	_, err := service.Register(context.Background(), RegisterInput{
		Email:       "user@example.com",
		Password:    "secret",
		DisplayName: "User",
		DeviceUID:   "device-1",
		DeviceName:  "iPhone",
		Platform:    "ios",
	})
	if err != nil {
		t.Fatalf("Register returned error: %v", err)
	}
	if domainService.lastRegister.Email != "user@example.com" {
		t.Fatalf("expected forwarded email, got %s", domainService.lastRegister.Email)
	}
}

func TestRefreshSessionForwardsCommand(t *testing.T) {
	domainService := &fakeDomainAccountService{
		refreshResult: &domainAccount.SessionBundle{},
	}
	service := NewService(domainService)

	_, err := service.RefreshSession(context.Background(), RefreshSessionInput{
		RefreshToken: "refresh-token",
		DeviceUID:    "device-1",
	})
	if err != nil {
		t.Fatalf("RefreshSession returned error: %v", err)
	}
	if domainService.lastRefresh.RefreshToken != "refresh-token" {
		t.Fatalf("expected forwarded refresh token, got %s", domainService.lastRefresh.RefreshToken)
	}
}

type fakeDomainAccountService struct {
	lastRegister   domainAccount.RegisterCommand
	lastLogin      domainAccount.LoginCommand
	lastRefresh    domainAccount.RefreshSessionCommand
	registerResult *domainAccount.SessionBundle
	loginResult    *domainAccount.SessionBundle
	refreshResult  *domainAccount.SessionBundle
}

func (s *fakeDomainAccountService) Register(
	_ context.Context,
	cmd domainAccount.RegisterCommand,
) (*domainAccount.SessionBundle, error) {
	s.lastRegister = cmd
	return s.registerResult, nil
}

func (s *fakeDomainAccountService) Login(
	_ context.Context,
	cmd domainAccount.LoginCommand,
) (*domainAccount.SessionBundle, error) {
	s.lastLogin = cmd
	return s.loginResult, nil
}

func (s *fakeDomainAccountService) RefreshSession(
	_ context.Context,
	cmd domainAccount.RefreshSessionCommand,
) (*domainAccount.SessionBundle, error) {
	s.lastRefresh = cmd
	return s.refreshResult, nil
}
