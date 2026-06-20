package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadAppliesDefaults(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")

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

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if cfg.Service.Name != "deadliner-dev" {
		t.Fatalf("unexpected service name: %s", cfg.Service.Name)
	}
	if cfg.Service.Address == "" {
		t.Fatal("expected default service address")
	}
	if cfg.Database.DSN == "" {
		t.Fatal("expected default database dsn")
	}
	if cfg.Sync.DefaultPullLimit == 0 || cfg.Sync.MaxPullLimit == 0 {
		t.Fatal("expected sync defaults")
	}
}
