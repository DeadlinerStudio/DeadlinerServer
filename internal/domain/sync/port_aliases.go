package sync

import portpkg "github.com/aritxonly/deadlinerserver/internal/domain/sync/port"

type DeadlineRepository = portpkg.DeadlineRepository
type HabitRepository = portpkg.HabitRepository
type MutationReceiptRepository = portpkg.MutationReceiptRepository
type SyncChangeRepository = portpkg.SyncChangeRepository
type SaveDeadlineParams = portpkg.SaveDeadlineParams
type SaveHabitParams = portpkg.SaveHabitParams
type AppendSyncChangeParams = portpkg.AppendSyncChangeParams
