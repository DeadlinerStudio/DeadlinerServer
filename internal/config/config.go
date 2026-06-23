package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
)

const DefaultPath = "conf/config.json"
const DefaultSecretPath = "conf/secret.json"
const SecretPathEnv = "DEADLINER_SECRET_CONFIG"

type Config struct {
	Service  ServiceConfig
	HTTP     HTTPConfig
	Auth     AuthConfig
	Database DatabaseConfig
	Sync     SyncConfig
}

type ServiceConfig struct {
	Name    string
	Address string
}

type HTTPConfig struct {
	Address             string
	ReadTimeoutSeconds  int
	WriteTimeoutSeconds int
	IdleTimeoutSeconds  int
}

type DatabaseConfig struct {
	Driver string
	DSN    string
}

type AuthConfig struct {
	AccessTokenSecret     string
	AccessTokenTTLMinutes int32
	RefreshTokenTTLHours  int32
	PasswordHashCost      int
	RandomTokenBytes      int
}

type SyncConfig struct {
	DefaultPullLimit int32
	MaxPullLimit     int32
}

type secretConfig struct {
	Auth     secretAuthConfig     `json:"auth"`
	Database secretDatabaseConfig `json:"database"`
}

type secretAuthConfig struct {
	AccessTokenSecret string `json:"accessTokenSecret"`
}

type secretDatabaseConfig struct {
	DSN string `json:"dsn"`
}

func Default() Config {
	return Config{
		Service: ServiceConfig{
			Name:    "deadliner",
			Address: ":8888",
		},
		HTTP: HTTPConfig{
			Address:             ":8080",
			ReadTimeoutSeconds:  15,
			WriteTimeoutSeconds: 15,
			IdleTimeoutSeconds:  60,
		},
		Auth: AuthConfig{
			AccessTokenTTLMinutes: 60 * 24,
			RefreshTokenTTLHours:  24 * 30,
			PasswordHashCost:      12,
			RandomTokenBytes:      32,
		},
		Database: DatabaseConfig{
			Driver: "mysql",
		},
		Sync: SyncConfig{
			DefaultPullLimit: 100,
			MaxPullLimit:     500,
		},
	}
}

func (c *Config) ApplyDefaults() {
	defaults := Default()

	if c.Service.Name == "" {
		c.Service.Name = defaults.Service.Name
	}
	if c.Service.Address == "" {
		c.Service.Address = defaults.Service.Address
	}
	if c.HTTP.Address == "" {
		c.HTTP.Address = defaults.HTTP.Address
	}
	if c.HTTP.ReadTimeoutSeconds == 0 {
		c.HTTP.ReadTimeoutSeconds = defaults.HTTP.ReadTimeoutSeconds
	}
	if c.HTTP.WriteTimeoutSeconds == 0 {
		c.HTTP.WriteTimeoutSeconds = defaults.HTTP.WriteTimeoutSeconds
	}
	if c.HTTP.IdleTimeoutSeconds == 0 {
		c.HTTP.IdleTimeoutSeconds = defaults.HTTP.IdleTimeoutSeconds
	}
	if c.Auth.AccessTokenTTLMinutes == 0 {
		c.Auth.AccessTokenTTLMinutes = defaults.Auth.AccessTokenTTLMinutes
	}
	if c.Auth.RefreshTokenTTLHours == 0 {
		c.Auth.RefreshTokenTTLHours = defaults.Auth.RefreshTokenTTLHours
	}
	if c.Auth.PasswordHashCost == 0 {
		c.Auth.PasswordHashCost = defaults.Auth.PasswordHashCost
	}
	if c.Auth.RandomTokenBytes == 0 {
		c.Auth.RandomTokenBytes = defaults.Auth.RandomTokenBytes
	}
	if c.Database.Driver == "" {
		c.Database.Driver = defaults.Database.Driver
	}
	if c.Sync.DefaultPullLimit == 0 {
		c.Sync.DefaultPullLimit = defaults.Sync.DefaultPullLimit
	}
	if c.Sync.MaxPullLimit == 0 {
		c.Sync.MaxPullLimit = defaults.Sync.MaxPullLimit
	}
}

func Load(path string) (Config, error) {
	return LoadWithSecretPath(path, ResolveSecretPath())
}

func LoadWithSecretPath(path string, secretPath string) (Config, error) {
	if path == "" {
		path = DefaultPath
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("read config %s: %w", path, err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("decode config %s: %w", path, err)
	}

	cfg.ApplyDefaults()
	if err := mergeSecretConfig(&cfg, secretPath); err != nil {
		return Config{}, err
	}
	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func ResolveSecretPath() string {
	if path := strings.TrimSpace(os.Getenv(SecretPathEnv)); path != "" {
		return path
	}
	return DefaultSecretPath
}

func mergeSecretConfig(cfg *Config, path string) error {
	if cfg == nil {
		return nil
	}
	path = strings.TrimSpace(path)
	if path == "" {
		return nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return fmt.Errorf("read secret config %s: %w", path, err)
	}

	var secret secretConfig
	if err := json.Unmarshal(data, &secret); err != nil {
		return fmt.Errorf("decode secret config %s: %w", path, err)
	}

	if secret.Auth.AccessTokenSecret != "" {
		cfg.Auth.AccessTokenSecret = secret.Auth.AccessTokenSecret
	}
	if secret.Database.DSN != "" {
		cfg.Database.DSN = secret.Database.DSN
	}

	return nil
}

func (c Config) Validate() error {
	var missing []string

	if strings.TrimSpace(c.Auth.AccessTokenSecret) == "" {
		missing = append(missing, "auth.accessTokenSecret")
	}
	if strings.TrimSpace(c.Database.DSN) == "" {
		missing = append(missing, "database.dsn")
	}

	if len(missing) > 0 {
		return fmt.Errorf(
			"invalid config: missing sensitive settings %s; put them in %s or set %s",
			strings.Join(missing, ", "),
			DefaultSecretPath,
			SecretPathEnv,
		)
	}

	return nil
}
