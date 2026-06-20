package sync

type PushChangesResult struct {
	Results         []MutationResult
	DeadlineChanges []DeadlineChange
	HabitChanges    []HabitChange
	NextCursor      string
}

