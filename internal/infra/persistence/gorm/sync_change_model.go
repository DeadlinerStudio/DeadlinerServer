package gorm

import "time"

type SyncChangeModel struct {
	ChangeID    int64     `gorm:"primaryKey;autoIncrement"`
	AccountID   int64     `gorm:"uniqueIndex:uk_sync_changes_account_device_mutation;index:idx_sync_changes_account_change,priority:1;not null"`
	DeviceUID   string    `gorm:"size:64;uniqueIndex:uk_sync_changes_account_device_mutation;not null"`
	MutationID  string    `gorm:"size:128;uniqueIndex:uk_sync_changes_account_device_mutation;not null"`
	EntityKind  string    `gorm:"size:32;not null"`
	EntityUID   string    `gorm:"size:128;not null"`
	Action      string    `gorm:"size:32;not null"`
	Payload     []byte    `gorm:"type:json"`
	CommittedAt time.Time `gorm:"not null"`
}

func (SyncChangeModel) TableName() string { return "sync_changes" }
