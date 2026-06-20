package repo

import (
	"github.com/aritxonly/deadlinerserver/internal/domain/sync"
	syncrepo "github.com/aritxonly/deadlinerserver/internal/infra/repo/sync"
	"gorm.io/gorm"
)

func NewDeadlineRepo(db *gorm.DB) sync.DeadlineRepository {
	return syncrepo.NewDeadlineRepo(db)
}
