package service

import (
	"context"
	"fmt"
	"time"

	commandpkg "github.com/aritxonly/deadlinerserver/internal/domain/account/command"
	entitypkg "github.com/aritxonly/deadlinerserver/internal/domain/account/entity"
	portpkg "github.com/aritxonly/deadlinerserver/internal/domain/account/port"
)

type service struct {
	repo               portpkg.Repository
	passwordHasher     portpkg.PasswordHasher
	refreshTokenHasher portpkg.RefreshTokenHasher
	accessTokenCodec   portpkg.AccessTokenCodec
	tokenGenerator     portpkg.TokenGenerator
	clock              portpkg.Clock
	accessTokenTTL     time.Duration
	refreshTokenTTL    time.Duration
}

func NewService(
	repo portpkg.Repository,
	passwordHasher portpkg.PasswordHasher,
	refreshTokenHasher portpkg.RefreshTokenHasher,
	accessTokenCodec portpkg.AccessTokenCodec,
	tokenGenerator portpkg.TokenGenerator,
	clock portpkg.Clock,
	accessTokenTTL time.Duration,
	refreshTokenTTL time.Duration,
) Service {
	return &service{
		repo:               repo,
		passwordHasher:     passwordHasher,
		refreshTokenHasher: refreshTokenHasher,
		accessTokenCodec:   accessTokenCodec,
		tokenGenerator:     tokenGenerator,
		clock:              clock,
		accessTokenTTL:     accessTokenTTL,
		refreshTokenTTL:    refreshTokenTTL,
	}
}

func (s *service) Register(
	ctx context.Context,
	cmd commandpkg.RegisterCommand,
) (*entitypkg.SessionBundle, error) {
	existing, err := s.repo.FindAccountByEmail(ctx, cmd.Email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrAccountAlreadyExists
	}

	passwordHash, err := s.passwordHasher.Hash(cmd.Password)
	if err != nil {
		return nil, err
	}

	accountUID, err := s.nextID("acc")
	if err != nil {
		return nil, err
	}

	acc := &entitypkg.Account{
		AccountUID:   accountUID,
		Email:        cmd.Email,
		PasswordHash: passwordHash,
		DisplayName:  cmd.DisplayName,
	}
	if err := s.repo.SaveAccount(ctx, acc); err != nil {
		return nil, err
	}

	if err := s.repo.SaveDevice(ctx, &entitypkg.Device{
		DeviceUID:  cmd.DeviceUID,
		AccountID:  acc.ID,
		Platform:   cmd.Platform,
		DeviceName: cmd.DeviceName,
	}); err != nil {
		return nil, err
	}

	return s.issueSessionBundle(ctx, acc, cmd.DeviceUID, "")
}

func (s *service) Login(ctx context.Context, cmd commandpkg.LoginCommand) (*entitypkg.SessionBundle, error) {
	acc, err := s.repo.FindAccountByEmail(ctx, cmd.Email)
	if err != nil {
		return nil, err
	}
	if acc == nil {
		return nil, ErrInvalidCredentials
	}
	if err := s.passwordHasher.Compare(acc.PasswordHash, cmd.Password); err != nil {
		return nil, ErrInvalidCredentials
	}

	if err := s.repo.SaveDevice(ctx, &entitypkg.Device{
		DeviceUID:  cmd.DeviceUID,
		AccountID:  acc.ID,
		Platform:   cmd.Platform,
		DeviceName: cmd.DeviceName,
	}); err != nil {
		return nil, err
	}

	return s.issueSessionBundle(ctx, acc, cmd.DeviceUID, "")
}

func (s *service) RefreshSession(
	ctx context.Context,
	cmd commandpkg.RefreshSessionCommand,
) (*entitypkg.SessionBundle, error) {
	refreshTokenHash := s.refreshTokenHasher.Hash(cmd.RefreshToken)
	session, err := s.repo.FindSessionByRefreshTokenHash(ctx, refreshTokenHash)
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, ErrInvalidRefreshToken
	}
	if session.DeviceUID != "" && cmd.DeviceUID != "" && session.DeviceUID != cmd.DeviceUID {
		return nil, ErrDeviceMismatch
	}

	now := s.clock.Now().UTC()
	expiresAt, err := time.Parse(time.RFC3339, session.ExpiresAt)
	if err != nil {
		return nil, fmt.Errorf("parse session expiry: %w", err)
	}
	if now.After(expiresAt) {
		return nil, ErrExpiredRefreshToken
	}

	acc, err := s.repo.FindAccountByID(ctx, session.AccountID)
	if err != nil {
		return nil, err
	}
	if acc == nil {
		return nil, ErrInvalidRefreshToken
	}

	deviceUID := session.DeviceUID
	if cmd.DeviceUID != "" {
		deviceUID = cmd.DeviceUID
	}

	return s.issueSessionBundle(ctx, acc, deviceUID, session.SessionUID)
}

func (s *service) issueSessionBundle(
	ctx context.Context,
	acc *entitypkg.Account,
	deviceUID string,
	existingSessionUID string,
) (*entitypkg.SessionBundle, error) {
	now := s.clock.Now().UTC()
	accessExpiresAt := now.Add(s.accessTokenTTL)
	refreshExpiresAt := now.Add(s.refreshTokenTTL)

	accessToken, err := s.accessTokenCodec.Sign(entitypkg.AccessTokenClaims{
		AccountUID: acc.AccountUID,
		DeviceUID:  deviceUID,
		ExpiresAt:  accessExpiresAt,
	})
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.tokenGenerator.Generate()
	if err != nil {
		return nil, err
	}

	sessionUID := existingSessionUID
	if sessionUID == "" {
		sessionUID, err = s.nextID("sess")
		if err != nil {
			return nil, err
		}
	}

	if err := s.repo.SaveSession(ctx, &entitypkg.Session{
		SessionUID:       sessionUID,
		AccountID:        acc.ID,
		DeviceUID:        deviceUID,
		RefreshTokenHash: s.refreshTokenHasher.Hash(refreshToken),
		ExpiresAt:        refreshExpiresAt.Format(time.RFC3339),
	}); err != nil {
		return nil, err
	}

	return &entitypkg.SessionBundle{
		AccountUID:   acc.AccountUID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    accessExpiresAt.Format(time.RFC3339),
	}, nil
}

func (s *service) nextID(prefix string) (string, error) {
	token, err := s.tokenGenerator.Generate()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s_%s", prefix, token), nil
}
