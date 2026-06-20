package sync

type DeadlineState string

const (
	DeadlineStateActive            DeadlineState = "active"
	DeadlineStateCompleted         DeadlineState = "completed"
	DeadlineStateArchived          DeadlineState = "archived"
	DeadlineStateAbandoned         DeadlineState = "abandoned"
	DeadlineStateAbandonedArchived DeadlineState = "abandonedArchived"
)

