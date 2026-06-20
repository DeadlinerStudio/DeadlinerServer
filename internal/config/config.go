package config

import (
	"encoding/json"
	"fmt"
	"os"
)

const DefaultPath = "conf/config.json"

type Config struct {
	Service  ServiceConfig
	Database DatabaseConfig
	Sync     SyncConfig
}

type ServiceConfig struct {
	Name    string
	Address string
}

type DatabaseConfig struct {
	Driver string
	DSN    string
}

type SyncConfig struct {
	DefaultPullLimit int32
	MaxPullLimit     int32
}

func Default() Config {
	return Config{
		Service: ServiceConfig{
			Name:    "deadliner",
			Address: ":8888",
		},
		Database: DatabaseConfig{
			Driver: "mysql",
			DSN:    "deadliner:deadliner@tcp(127.0.0.1:3306)/deadliner?charset=utf8mb4&parseTime=True&loc=Local",
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
	if c.Database.Driver == "" {
		c.Database.Driver = defaults.Database.Driver
	}
	if c.Database.DSN == "" {
		c.Database.DSN = defaults.Database.DSN
	}
	if c.Sync.DefaultPullLimit == 0 {
		c.Sync.DefaultPullLimit = defaults.Sync.DefaultPullLimit
	}
	if c.Sync.MaxPullLimit == 0 {
		c.Sync.MaxPullLimit = defaults.Sync.MaxPullLimit
	}
}

func Load(path string) (Config, error) {
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
	return cfg, nil
}
