package sync

type Mutation struct {
	MutationID    string         `json:"mutation_id"`
	DeviceUID     string         `json:"device_uid"`
	EntityUID     string         `json:"entity_uid"`
	ClientVersion LogicalVersion `json:"client_version"`
	BaseChangeID  int64          `json:"base_change_id"`
	Deadline      *DeadlinePatch `json:"deadline,omitempty"`
	Habit         *HabitPatch    `json:"habit,omitempty"`
}

type DeadlinePatch struct {
	Deleted  bool             `json:"deleted"`
	Document DeadlineDocument `json:"doc"`
}

type HabitPatch struct {
	Deleted  bool          `json:"deleted"`
	Document HabitDocument `json:"doc"`
}

type MutationResult struct {
	MutationID      string        `json:"mutation_id"`
	EntityUID       string        `json:"entity_uid"`
	Accepted        bool          `json:"accepted"`
	RejectionReason string        `json:"rejection_reason"`
	ServerVersion   ServerVersion `json:"server_version"`
	Replayed        bool          `json:"replayed"`
	Status          string        `json:"status"`
}

