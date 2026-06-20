package sync

import portpkg "github.com/aritxonly/deadlinerserver/internal/domain/sync/port"

type DeadlineRepository = portpkg.DeadlineRepository
type MutationReceiptRepository = portpkg.MutationReceiptRepository
type SyncChangeRepository = portpkg.SyncChangeRepository
type SaveDeadlineParams = portpkg.SaveDeadlineParams
type AppendSyncChangeParams = portpkg.AppendSyncChangeParams
