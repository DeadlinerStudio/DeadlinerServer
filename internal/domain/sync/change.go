package sync

type ServerVersion struct {
	ChangeID    int64  `json:"change_id"`
	CommittedAt string `json:"committed_at"`
}

type DeadlineChange struct {
	EntityUID     string           `json:"entity_uid"`
	Deleted       bool             `json:"deleted"`
	ServerVersion ServerVersion    `json:"server_version"`
	Document      DeadlineDocument `json:"doc"`
}

type HabitChange struct {
	EntityUID     string        `json:"entity_uid"`
	Deleted       bool          `json:"deleted"`
	ServerVersion ServerVersion `json:"server_version"`
	Document      HabitDocument `json:"doc"`
}

