package command

import statepkg "github.com/aritxonly/deadlinerserver/internal/domain/sync/state"

type PushChangesResult struct {
	Results         []MutationResult
	DeadlineChanges []statepkg.DeadlineChange
	HabitChanges    []statepkg.HabitChange
	NextCursor      string
}
