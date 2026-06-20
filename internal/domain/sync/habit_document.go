package sync

type HabitConfig struct {
	Name           string        `json:"name"`
	Description    string        `json:"description"`
	Color          int32         `json:"color"`
	IconKey        string        `json:"icon_key"`
	Period         HabitPeriod   `json:"period"`
	TimesPerPeriod int32         `json:"times_per_period"`
	GoalType       HabitGoalType `json:"goal_type"`
	TotalTarget    int32         `json:"total_target"`
	CreatedAt      string        `json:"created_at"`
	UpdatedAt      string        `json:"updated_at"`
	Status         HabitStatus   `json:"status"`
	SortOrder      int32         `json:"sort_order"`
	AlarmTime      string        `json:"alarm_time"`
}

type HabitRecord struct {
	Date      string            `json:"date"`
	Count     int32             `json:"count"`
	Status    HabitRecordStatus `json:"status"`
	CreatedAt string            `json:"created_at"`
}

type HabitDocument struct {
	DDLUID  string        `json:"ddl_uid"`
	Habit   HabitConfig   `json:"habit"`
	Records []HabitRecord `json:"records"`
}

