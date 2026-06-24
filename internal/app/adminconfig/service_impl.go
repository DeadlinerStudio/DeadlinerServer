package adminconfig

import (
	"context"
	"strings"

	"github.com/aritxonly/deadlinerserver/internal/config"
)

type service struct {
	configPath string
	secretPath string
}

func NewService(configPath string, secretPath string) Service {
	return &service{
		configPath: configPath,
		secretPath: secretPath,
	}
}

func (s *service) GetSnapshot(context.Context) (*Snapshot, error) {
	publicConfig, err := config.LoadPublic(s.configPath)
	if err != nil {
		return nil, err
	}
	effectiveConfig, err := config.LoadWithSecretPath(s.configPath, s.secretPath)
	if err != nil {
		return nil, err
	}

	return &Snapshot{
		Service:  publicConfig.Service,
		HTTP:     publicConfig.HTTP,
		Database: publicConfig.Database,
		Sync:     publicConfig.Sync,
		Admin: config.AdminConfig{
			Enabled:  publicConfig.Admin.Enabled,
			BasePath: publicConfig.Admin.BasePath,
		},
		SecretStatus: SecretStatus{
			AccessTokenSecretConfigured: strings.TrimSpace(effectiveConfig.Auth.AccessTokenSecret) != "",
			DatabaseDSNConfigured:       strings.TrimSpace(effectiveConfig.Database.DSN) != "",
			AdminTokenConfigured:        strings.TrimSpace(effectiveConfig.Admin.AuthToken) != "",
		},
		RestartRequired: false,
	}, nil
}

func (s *service) Update(ctx context.Context, input UpdateInput) (*Snapshot, error) {
	currentConfig, err := config.LoadPublic(s.configPath)
	if err != nil {
		return nil, err
	}

	currentConfig.Service = input.Service
	currentConfig.HTTP = input.HTTP
	currentConfig.Database.Driver = input.Database.Driver
	currentConfig.Sync = input.Sync
	currentConfig.Admin.Enabled = input.Admin.Enabled
	currentConfig.Admin.BasePath = input.Admin.BasePath
	currentConfig.ApplyDefaults()

	if err := config.SavePublic(s.configPath, currentConfig); err != nil {
		return nil, err
	}

	snapshot, err := s.GetSnapshot(ctx)
	if err != nil {
		return nil, err
	}
	snapshot.RestartRequired = true
	return snapshot, nil
}
