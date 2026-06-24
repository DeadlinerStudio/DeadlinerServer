package adminconfig

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/aritxonly/deadlinerserver/internal/config"
)

func TestUpdatePersistsPublicConfig(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.json")
	secretPath := filepath.Join(dir, "secret.json")

	initial := config.Default()
	initial.Admin.Enabled = true
	if err := config.SavePublic(configPath, initial); err != nil {
		t.Fatalf("SavePublic initial failed: %v", err)
	}
	if err := os.WriteFile(secretPath, []byte(`{
  "auth": { "accessTokenSecret": "test-secret" },
  "database": { "dsn": "deadliner:test@tcp(127.0.0.1:3306)/deadliner?charset=utf8mb4&parseTime=True&loc=Local" },
  "admin": { "authToken": "admin-secret" }
}`), 0o644); err != nil {
		t.Fatalf("write secret failed: %v", err)
	}

	service := NewService(configPath, secretPath)
	snapshot, err := service.Update(context.Background(), UpdateInput{
		Service: config.ServiceConfig{
			Name:    "deadliner",
			Address: ":9999",
		},
		HTTP: config.HTTPConfig{
			Address:                ":18080",
			ReadTimeoutSeconds:     20,
			WriteTimeoutSeconds:    20,
			IdleTimeoutSeconds:     90,
			MaxRequestBodyBytes:    2048,
			RateLimitPerMinute:     120,
			RateLimitBurst:         20,
			AuthRateLimitPerMinute: 15,
			AuthRateLimitBurst:     5,
			SyncRateLimitPerMinute: 180,
			SyncRateLimitBurst:     30,
		},
		Database: config.DatabaseConfig{
			Driver: "mysql",
		},
		Sync: config.SyncConfig{
			DefaultPullLimit: 80,
			MaxPullLimit:     300,
		},
		Admin: config.AdminConfig{
			Enabled:  true,
			BasePath: "/ops",
		},
	})
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}
	if !snapshot.RestartRequired {
		t.Fatal("expected restartRequired=true")
	}
	if snapshot.HTTP.Address != ":18080" {
		t.Fatalf("unexpected http address: %s", snapshot.HTTP.Address)
	}

	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("ReadFile failed: %v", err)
	}
	text := string(content)
	if strings.Contains(text, "test-secret") || strings.Contains(text, "admin-secret") {
		t.Fatalf("public config leaked secrets: %s", text)
	}
}
