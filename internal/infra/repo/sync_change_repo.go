package repo

import (
	"github.com/aritxonly/deadlinerserver/internal/domain/sync"
	syncrepo "github.com/aritxonly/deadlinerserver/internal/infra/repo/sync"
	"gorm.io/gorm"
)

func NewSyncChangeRepo(db *gorm.DB) sync.SyncChangeRepository {
	return syncrepo.NewSyncChangeRepo(db)
}
