package adminconfig

import "github.com/aritxonly/deadlinerserver/internal/config"

type SecretStatus struct {
	AccessTokenSecretConfigured bool `json:"accessTokenSecretConfigured"`
	DatabaseDSNConfigured       bool `json:"databaseDSNConfigured"`
	AdminTokenConfigured        bool `json:"adminTokenConfigured"`
}

type Snapshot struct {
	Service         config.ServiceConfig  `json:"service"`
	HTTP            config.HTTPConfig     `json:"http"`
	Database        config.DatabaseConfig `json:"database"`
	Sync            config.SyncConfig     `json:"sync"`
	Admin           config.AdminConfig    `json:"admin"`
	SecretStatus    SecretStatus          `json:"secretStatus"`
	RestartRequired bool                  `json:"restartRequired"`
}

type UpdateInput struct {
	Service  config.ServiceConfig  `json:"service"`
	HTTP     config.HTTPConfig     `json:"http"`
	Database config.DatabaseConfig `json:"database"`
	Sync     config.SyncConfig     `json:"sync"`
	Admin    config.AdminConfig    `json:"admin"`
}
