package gorm

import "time"

type HabitDocModel struct {
	ID                 int64   `gorm:"primaryKey;autoIncrement"`
	AccountID          int64   `gorm:"uniqueIndex:uk_habit_docs_account_ddl_uid;index:idx_habit_docs_account_change,priority:1;not null"`
	DDLUID             string  `gorm:"size:128;uniqueIndex:uk_habit_docs_account_ddl_uid;not null"`
	Payload            []byte  `gorm:"type:json"`
	Deleted            bool    `gorm:"not null"`
	ClientVerTS        *string `gorm:"size:64"`
	ClientVerCtr       *int32
	ClientVerDev       *string   `gorm:"size:64"`
	ServerChangeID     int64     `gorm:"index:idx_habit_docs_account_change,priority:2;not null"`
	CommittedAt        time.Time `gorm:"not null"`
	UpdatedByDeviceUID *string   `gorm:"size:64"`
}

func (HabitDocModel) TableName() string { return "habit_docs" }
