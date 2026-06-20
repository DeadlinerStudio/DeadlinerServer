package sync

type DeadlineDocument struct {
	UID             string        `json:"uid"`
	LegacyID        int64         `json:"legacy_id"`
	Name            string        `json:"name"`
	StartTime       string        `json:"start_time"`
	EndTime         string        `json:"end_time"`
	State           DeadlineState `json:"state"`
	CompleteTime    string        `json:"complete_time"`
	Note            string        `json:"note"`
	IsStared        bool          `json:"is_stared"`
	Type            DeadlineType  `json:"type"`
	HabitCount      int32         `json:"habit_count"`
	HabitTotalCount int32         `json:"habit_total_count"`
	CalendarEvent   int64         `json:"calendar_event"`
	Timestamp       string        `json:"timestamp"`
	SubTasks        []SubTask     `json:"sub_tasks"`
}

