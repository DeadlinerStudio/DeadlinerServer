package sync

import documentpkg "github.com/aritxonly/deadlinerserver/internal/domain/sync/document"

type DeadlineDocument = documentpkg.DeadlineDocument
type DeadlineState = documentpkg.DeadlineState
type DeadlineType = documentpkg.DeadlineType
type HabitConfig = documentpkg.HabitConfig
type HabitDocument = documentpkg.HabitDocument
type HabitGoalType = documentpkg.HabitGoalType
type HabitPeriod = documentpkg.HabitPeriod
type HabitRecord = documentpkg.HabitRecord
type HabitRecordStatus = documentpkg.HabitRecordStatus
type HabitStatus = documentpkg.HabitStatus
type SubTask = documentpkg.SubTask

const (
	DeadlineStateActive            = documentpkg.DeadlineStateActive
	DeadlineStateCompleted         = documentpkg.DeadlineStateCompleted
	DeadlineStateArchived          = documentpkg.DeadlineStateArchived
	DeadlineStateAbandoned         = documentpkg.DeadlineStateAbandoned
	DeadlineStateAbandonedArchived = documentpkg.DeadlineStateAbandonedArchived
	DeadlineTypeTask               = documentpkg.DeadlineTypeTask
	DeadlineTypeHabit              = documentpkg.DeadlineTypeHabit
	HabitPeriodDaily               = documentpkg.HabitPeriodDaily
	HabitPeriodWeekly              = documentpkg.HabitPeriodWeekly
	HabitPeriodMonthly             = documentpkg.HabitPeriodMonthly
	HabitPeriodOnce                = documentpkg.HabitPeriodOnce
	HabitPeriodEbbinghaus          = documentpkg.HabitPeriodEbbinghaus
	HabitGoalTypePerPeriod         = documentpkg.HabitGoalTypePerPeriod
	HabitGoalTypeTotal             = documentpkg.HabitGoalTypeTotal
	HabitStatusActive              = documentpkg.HabitStatusActive
	HabitStatusArchived            = documentpkg.HabitStatusArchived
	HabitRecordStatusCompleted     = documentpkg.HabitRecordStatusCompleted
	HabitRecordStatusSkipped       = documentpkg.HabitRecordStatusSkipped
	HabitRecordStatusFailed        = documentpkg.HabitRecordStatusFailed
)
