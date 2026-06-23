package repo

import (
	"github.com/aritxonly/deadlinerserver/internal/domain/sync"
	syncrepo "github.com/aritxonly/deadlinerserver/internal/infra/repo/sync"
	"gorm.io/gorm"
)

func NewHabitRepo(db *gorm.DB) sync.HabitRepository {
	return syncrepo.NewHabitRepo(db)
}
