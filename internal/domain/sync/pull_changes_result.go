package sync

type PullChangesResult struct {
	DeadlineChanges []DeadlineChange
	HabitChanges    []HabitChange
	NextCursor      string
	HasMore         bool
}

