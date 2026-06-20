package document

type HabitPeriod string

const (
	HabitPeriodDaily      HabitPeriod = "DAILY"
	HabitPeriodWeekly     HabitPeriod = "WEEKLY"
	HabitPeriodMonthly    HabitPeriod = "MONTHLY"
	HabitPeriodOnce       HabitPeriod = "ONCE"
	HabitPeriodEbbinghaus HabitPeriod = "EBBINGHAUS"
)

type HabitGoalType string

const (
	HabitGoalTypePerPeriod HabitGoalType = "PER_PERIOD"
	HabitGoalTypeTotal     HabitGoalType = "TOTAL"
)

type HabitStatus string

const (
	HabitStatusActive   HabitStatus = "ACTIVE"
	HabitStatusArchived HabitStatus = "ARCHIVED"
)

type HabitRecordStatus string

const (
	HabitRecordStatusCompleted HabitRecordStatus = "COMPLETED"
	HabitRecordStatusSkipped   HabitRecordStatus = "SKIPPED"
	HabitRecordStatusFailed    HabitRecordStatus = "FAILED"
)
