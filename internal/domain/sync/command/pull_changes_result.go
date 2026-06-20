package command

import statepkg "github.com/aritxonly/deadlinerserver/internal/domain/sync/state"

type PullChangesResult struct {
	DeadlineChanges []statepkg.DeadlineChange
	HabitChanges    []statepkg.HabitChange
	NextCursor      string
	HasMore         bool
}
