package gorm

import "time"

type AccountModel struct {
	ID           int64     `gorm:"primaryKey;autoIncrement"`
	AccountUID   string    `gorm:"size:64;uniqueIndex;not null"`
	Email        string    `gorm:"size:255;uniqueIndex;not null"`
	PasswordHash string    `gorm:"size:255;not null"`
	DisplayName  string    `gorm:"size:255;not null"`
	CreatedAt    time.Time `gorm:"not null"`
	UpdatedAt    time.Time `gorm:"not null"`
}

func (AccountModel) TableName() string { return "accounts" }

type DeviceModel struct {
	ID         int64  `gorm:"primaryKey;autoIncrement"`
	DeviceUID  string `gorm:"size:64;uniqueIndex;not null"`
	AccountID  int64  `gorm:"index;not null"`
	Platform   string `gorm:"size:64;not null"`
	DeviceName string `gorm:"size:255;not null"`
	LastSeenAt *time.Time
	CreatedAt  time.Time `gorm:"not null"`
	UpdatedAt  time.Time `gorm:"not null"`
}

func (DeviceModel) TableName() string { return "devices" }

type SessionModel struct {
	ID               int64     `gorm:"primaryKey;autoIncrement"`
	SessionUID       string    `gorm:"size:64;uniqueIndex;not null"`
	AccountID        int64     `gorm:"index;not null"`
	DeviceUID        string    `gorm:"size:64;index;not null"`
	RefreshTokenHash string    `gorm:"size:255;not null"`
	ExpiresAt        time.Time `gorm:"not null"`
	RevokedAt        *time.Time
	CreatedAt        time.Time `gorm:"not null"`
}

func (SessionModel) TableName() string { return "sessions" }

type MutationReceiptModel struct {
	ID             int64  `gorm:"primaryKey;autoIncrement"`
	AccountID      int64  `gorm:"uniqueIndex:uk_receipt;not null"`
	DeviceUID      string `gorm:"size:64;uniqueIndex:uk_receipt;not null"`
	MutationID     string `gorm:"size:128;uniqueIndex:uk_receipt;not null"`
	EntityKind     string `gorm:"size:32;not null"`
	EntityUID      string `gorm:"size:128;not null"`
	Status         string `gorm:"size:32;not null"`
	Replayed       bool   `gorm:"not null"`
	ResultChangeID *int64
	ResultPayload  []byte    `gorm:"type:json"`
	CreatedAt      time.Time `gorm:"not null"`
	UpdatedAt      time.Time `gorm:"not null"`
}

func (MutationReceiptModel) TableName() string { return "mutation_receipts" }
