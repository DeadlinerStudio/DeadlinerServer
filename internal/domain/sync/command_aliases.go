package sync

import commandpkg "github.com/aritxonly/deadlinerserver/internal/domain/sync/command"

type Mutation = commandpkg.Mutation
type DeadlinePatch = commandpkg.DeadlinePatch
type HabitPatch = commandpkg.HabitPatch
type MutationResult = commandpkg.MutationResult
type PullChangesCommand = commandpkg.PullChangesCommand
type PullChangesResult = commandpkg.PullChangesResult
type PushChangesCommand = commandpkg.PushChangesCommand
type PushChangesResult = commandpkg.PushChangesResult

const (
	MutationStatusApplied  = commandpkg.MutationStatusApplied
	MutationStatusReplayed = commandpkg.MutationStatusReplayed
	MutationStatusConflict = commandpkg.MutationStatusConflict
	MutationStatusRejected = commandpkg.MutationStatusRejected
)
