package sync

import statepkg "github.com/aritxonly/deadlinerserver/internal/domain/sync/state"

type ServerVersion = statepkg.ServerVersion
type DeadlineChange = statepkg.DeadlineChange
type HabitChange = statepkg.HabitChange
type MutationReceipt = statepkg.MutationReceipt
type SyncChange = statepkg.SyncChange
type LogicalVersion = statepkg.LogicalVersion

var CompareLogicalVersion = statepkg.CompareLogicalVersion
