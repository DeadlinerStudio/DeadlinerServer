package gorm

import "time"

type DeadlineItemModel struct {
	ID                 int64   `gorm:"primaryKey;autoIncrement"`
	AccountID          int64   `gorm:"uniqueIndex:uk_deadline_items_account_uid;not null"`
	UID                string  `gorm:"size:128;uniqueIndex:uk_deadline_items_account_uid;not null"`
	LegacyID           int64   `gorm:"not null"`
	Name               string  `gorm:"size:255;not null"`
	StartTime          string  `gorm:"size:64;not null"`
	EndTime            string  `gorm:"size:64;not null"`
	State              string  `gorm:"size:64;not null"`
	CompleteTime       string  `gorm:"size:64;not null"`
	Note               string  `gorm:"type:text;not null"`
	IsStared           bool    `gorm:"not null"`
	Type               string  `gorm:"size:32;not null"`
	HabitCount         int32   `gorm:"not null"`
	HabitTotalCount    int32   `gorm:"not null"`
	CalendarEvent      int64   `gorm:"not null"`
	BusinessTimestamp  string  `gorm:"size:64;not null"`
	SubTasks           []byte  `gorm:"type:json;not null"`
	Deleted            bool    `gorm:"not null"`
	ClientVerTS        *string `gorm:"size:64"`
	ClientVerCtr       *int32
	ClientVerDev       *string   `gorm:"size:64"`
	ServerChangeID     int64     `gorm:"index:idx_deadline_items_account_change,priority:2;not null"`
	CommittedAt        time.Time `gorm:"not null"`
	UpdatedByDeviceUID *string   `gorm:"size:64"`
}

func (DeadlineItemModel) TableName() string { return "deadline_items" }
