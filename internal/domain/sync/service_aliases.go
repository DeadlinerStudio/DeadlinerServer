package sync

import (
	"github.com/aritxonly/deadlinerserver/internal/domain/account"
	portpkg "github.com/aritxonly/deadlinerserver/internal/domain/sync/port"
	servicepkg "github.com/aritxonly/deadlinerserver/internal/domain/sync/service"
)

type Service = servicepkg.Service

func NewService(
	accountRepo account.Repository,
	deadlineRepo portpkg.DeadlineRepository,
	mutationReceiptRepo portpkg.MutationReceiptRepository,
	syncChangeRepo portpkg.SyncChangeRepository,
) Service {
	return servicepkg.NewService(accountRepo, deadlineRepo, mutationReceiptRepo, syncChangeRepo)
}
