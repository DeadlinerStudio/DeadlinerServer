package state

import documentpkg "github.com/aritxonly/deadlinerserver/internal/domain/sync/document"

type ServerVersion struct {
	ChangeID    int64  `json:"change_id"`
	CommittedAt string `json:"committed_at"`
}

type DeadlineChange struct {
	EntityUID     string                       `json:"entity_uid"`
	Deleted       bool                         `json:"deleted"`
	ServerVersion ServerVersion                `json:"server_version"`
	Document      documentpkg.DeadlineDocument `json:"doc"`
}

type HabitChange struct {
	EntityUID     string                    `json:"entity_uid"`
	Deleted       bool                      `json:"deleted"`
	ServerVersion ServerVersion             `json:"server_version"`
	Document      documentpkg.HabitDocument `json:"doc"`
}
