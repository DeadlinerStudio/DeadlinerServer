package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadAppliesDefaults(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")
	secretPath := filepath.Join(dir, "secret.json")

	content := []byte(`{
  "service": {
    "name": "deadliner-dev"
  },
  "database": {
    "driver": "mysql"
  }
}`)

	if err := os.WriteFile(path, content, 0o644); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}
	if err := os.WriteFile(secretPath, []byte(`{
  "auth": {
    "accessTokenSecret": "test-secret"
  },
  "database": {
    "dsn": "deadliner:test@tcp(127.0.0.1:3306)/deadliner?charset=utf8mb4&parseTime=True&loc=Local"
  }
}`), 0o644); err != nil {
		t.Fatalf("WriteFile secret failed: %v", err)
	}

	cfg, err := LoadWithSecretPath(path, secretPath)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if cfg.Service.Name != "deadliner-dev" {
		t.Fatalf("unexpected service name: %s", cfg.Service.Name)
	}
	if cfg.Service.Address == "" {
		t.Fatal("expected default service address")
	}
	if cfg.HTTP.Address == "" {
		t.Fatal("expected default http address")
	}
	if cfg.HTTP.ReadTimeoutSeconds == 0 || cfg.HTTP.WriteTimeoutSeconds == 0 || cfg.HTTP.IdleTimeoutSeconds == 0 {
		t.Fatal("expected http timeout defaults")
	}
	if cfg.HTTP.MaxRequestBodyBytes == 0 {
		t.Fatal("expected max request body default")
	}
	if cfg.HTTP.RateLimitPerMinute == 0 || cfg.HTTP.RateLimitBurst == 0 {
		t.Fatal("expected http rate limit defaults")
	}
	if cfg.HTTP.AuthRateLimitPerMinute == 0 || cfg.HTTP.AuthRateLimitBurst == 0 {
		t.Fatal("expected auth rate limit defaults")
	}
	if cfg.HTTP.SyncRateLimitPerMinute == 0 || cfg.HTTP.SyncRateLimitBurst == 0 {
		t.Fatal("expected sync rate limit defaults")
	}
	if cfg.Auth.AccessTokenSecret != "test-secret" {
		t.Fatalf("unexpected access token secret: %s", cfg.Auth.AccessTokenSecret)
	}
	if cfg.Auth.AccessTokenTTLMinutes == 0 || cfg.Auth.RefreshTokenTTLHours == 0 {
		t.Fatal("expected auth ttl defaults")
	}
	if cfg.Database.DSN == "" {
		t.Fatal("expected merged database dsn")
	}
	if cfg.Sync.DefaultPullLimit == 0 || cfg.Sync.MaxPullLimit == 0 {
		t.Fatal("expected sync defaults")
	}
	if cfg.Admin.BasePath == "" {
		t.Fatal("expected admin base path default")
	}
}

func TestLoadRejectsMissingSensitiveConfig(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")

	if err := os.WriteFile(path, []byte(`{
  "service": {
    "name": "deadliner-dev"
  }
}`), 0o644); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}

	_, err := LoadWithSecretPath(path, filepath.Join(dir, "missing-secret.json"))
	if err == nil {
		t.Fatal("expected missing sensitive config error")
	}
}

func TestSavePublicOmitsSensitiveFields(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")

	cfg := Default()
	cfg.Auth.AccessTokenSecret = "secret-value"
	cfg.Database.DSN = "dsn-value"
	cfg.Admin.Enabled = true
	cfg.Admin.AuthToken = "admin-secret"

	if err := SavePublic(path, cfg); err != nil {
		t.Fatalf("SavePublic failed: %v", err)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile failed: %v", err)
	}

	text := string(content)
	for _, secret := range []string{"secret-value", "dsn-value", "admin-secret"} {
		if strings.Contains(text, secret) {
			t.Fatalf("public config leaked secret %q", secret)
		}
	}
}
