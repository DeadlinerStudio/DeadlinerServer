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
	Service  ServiceConfig  `json:"service"`
	HTTP     HTTPConfig     `json:"http"`
	Auth     AuthConfig     `json:"auth"`
	Database DatabaseConfig `json:"database"`
	Sync     SyncConfig     `json:"sync"`
	Admin    AdminConfig    `json:"admin"`
}

type ServiceConfig struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

type HTTPConfig struct {
	Address                string `json:"address"`
	ReadTimeoutSeconds     int    `json:"readTimeoutSeconds"`
	WriteTimeoutSeconds    int    `json:"writeTimeoutSeconds"`
	IdleTimeoutSeconds     int    `json:"idleTimeoutSeconds"`
	MaxRequestBodyBytes    int    `json:"maxRequestBodyBytes"`
	RateLimitPerMinute     int    `json:"rateLimitPerMinute"`
	RateLimitBurst         int    `json:"rateLimitBurst"`
	AuthRateLimitPerMinute int    `json:"authRateLimitPerMinute"`
	AuthRateLimitBurst     int    `json:"authRateLimitBurst"`
	SyncRateLimitPerMinute int    `json:"syncRateLimitPerMinute"`
	SyncRateLimitBurst     int    `json:"syncRateLimitBurst"`
}

type DatabaseConfig struct {
	Driver string `json:"driver"`
	DSN    string `json:"dsn,omitempty"`
}

type AuthConfig struct {
	AccessTokenSecret     string `json:"accessTokenSecret,omitempty"`
	AccessTokenTTLMinutes int32  `json:"accessTokenTTLMinutes"`
	RefreshTokenTTLHours  int32  `json:"refreshTokenTTLHours"`
	PasswordHashCost      int    `json:"passwordHashCost"`
	RandomTokenBytes      int    `json:"randomTokenBytes"`
}

type SyncConfig struct {
	DefaultPullLimit int32 `json:"defaultPullLimit"`
	MaxPullLimit     int32 `json:"maxPullLimit"`
}

type AdminConfig struct {
	Enabled   bool   `json:"enabled"`
	BasePath  string `json:"basePath"`
	AuthToken string `json:"-"`
}

type secretConfig struct {
	Auth     secretAuthConfig     `json:"auth"`
	Database secretDatabaseConfig `json:"database"`
	Admin    secretAdminConfig    `json:"admin"`
}

type secretAuthConfig struct {
	AccessTokenSecret string `json:"accessTokenSecret"`
}

type secretDatabaseConfig struct {
	DSN string `json:"dsn"`
}

type secretAdminConfig struct {
	AuthToken string `json:"authToken"`
}

func Default() Config {
	return Config{
		Service: ServiceConfig{
			Name:    "deadliner",
			Address: ":8888",
		},
		HTTP: HTTPConfig{
			Address:                ":8080",
			ReadTimeoutSeconds:     15,
			WriteTimeoutSeconds:    15,
			IdleTimeoutSeconds:     60,
			MaxRequestBodyBytes:    1 << 20,
			RateLimitPerMinute:     240,
			RateLimitBurst:         60,
			AuthRateLimitPerMinute: 30,
			AuthRateLimitBurst:     10,
			SyncRateLimitPerMinute: 240,
			SyncRateLimitBurst:     60,
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
		Admin: AdminConfig{
			Enabled:  false,
			BasePath: "/admin",
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
	if c.HTTP.MaxRequestBodyBytes == 0 {
		c.HTTP.MaxRequestBodyBytes = defaults.HTTP.MaxRequestBodyBytes
	}
	if c.HTTP.RateLimitPerMinute == 0 {
		c.HTTP.RateLimitPerMinute = defaults.HTTP.RateLimitPerMinute
	}
	if c.HTTP.RateLimitBurst == 0 {
		c.HTTP.RateLimitBurst = defaults.HTTP.RateLimitBurst
	}
	if c.HTTP.AuthRateLimitPerMinute == 0 {
		c.HTTP.AuthRateLimitPerMinute = defaults.HTTP.AuthRateLimitPerMinute
	}
	if c.HTTP.AuthRateLimitBurst == 0 {
		c.HTTP.AuthRateLimitBurst = defaults.HTTP.AuthRateLimitBurst
	}
	if c.HTTP.SyncRateLimitPerMinute == 0 {
		c.HTTP.SyncRateLimitPerMinute = defaults.HTTP.SyncRateLimitPerMinute
	}
	if c.HTTP.SyncRateLimitBurst == 0 {
		c.HTTP.SyncRateLimitBurst = defaults.HTTP.SyncRateLimitBurst
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
	if c.Admin.BasePath == "" {
		c.Admin.BasePath = defaults.Admin.BasePath
	}
}

func Load(path string) (Config, error) {
	return LoadWithSecretPath(path, ResolveSecretPath())
}

func LoadPublic(path string) (Config, error) {
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
	if err := cfg.ValidatePublic(); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func LoadWithSecretPath(path string, secretPath string) (Config, error) {
	return load(path, secretPath)
}

func load(path string, secretPath string) (Config, error) {
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
	if secret.Admin.AuthToken != "" {
		cfg.Admin.AuthToken = secret.Admin.AuthToken
	}

	return nil
}

func SavePublic(path string, cfg Config) error {
	if path == "" {
		path = DefaultPath
	}

	cfg.ApplyDefaults()
	cfg.Auth.AccessTokenSecret = ""
	cfg.Database.DSN = ""
	cfg.Admin.AuthToken = ""
	if err := cfg.ValidatePublic(); err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("encode public config %s: %w", path, err)
	}
	data = append(data, '\n')

	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("write public config %s: %w", path, err)
	}
	return nil
}

func (c Config) Validate() error {
	if err := c.ValidatePublic(); err != nil {
		return err
	}

	var missing []string

	if strings.TrimSpace(c.Auth.AccessTokenSecret) == "" {
		missing = append(missing, "auth.accessTokenSecret")
	}
	if strings.TrimSpace(c.Database.DSN) == "" {
		missing = append(missing, "database.dsn")
	}
	if c.Admin.Enabled && strings.TrimSpace(c.Admin.AuthToken) == "" {
		missing = append(missing, "admin.authToken")
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

func (c Config) ValidatePublic() error {
	if c.HTTP.ReadTimeoutSeconds <= 0 || c.HTTP.WriteTimeoutSeconds <= 0 || c.HTTP.IdleTimeoutSeconds <= 0 {
		return errors.New("invalid config: http timeouts must be positive")
	}
	if c.HTTP.MaxRequestBodyBytes <= 0 {
		return errors.New("invalid config: http.maxRequestBodyBytes must be positive")
	}
	if c.HTTP.RateLimitPerMinute <= 0 || c.HTTP.RateLimitBurst <= 0 {
		return errors.New("invalid config: http rate limits must be positive")
	}
	if c.HTTP.AuthRateLimitPerMinute <= 0 || c.HTTP.AuthRateLimitBurst <= 0 {
		return errors.New("invalid config: auth rate limits must be positive")
	}
	if c.HTTP.SyncRateLimitPerMinute <= 0 || c.HTTP.SyncRateLimitBurst <= 0 {
		return errors.New("invalid config: sync rate limits must be positive")
	}
	if c.Sync.DefaultPullLimit <= 0 || c.Sync.MaxPullLimit <= 0 {
		return errors.New("invalid config: sync pull limits must be positive")
	}
	if c.Sync.DefaultPullLimit > c.Sync.MaxPullLimit {
		return errors.New("invalid config: sync.defaultPullLimit must not exceed sync.maxPullLimit")
	}
	if c.Admin.Enabled && !strings.HasPrefix(c.Admin.BasePath, "/") {
		return errors.New("invalid config: admin.basePath must start with '/'")
	}
	return nil
}
