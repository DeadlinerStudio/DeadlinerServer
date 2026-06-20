package repo

import (
	"github.com/aritxonly/deadlinerserver/internal/domain/sync"
	syncrepo "github.com/aritxonly/deadlinerserver/internal/infra/repo/sync"
	"gorm.io/gorm"
)

func NewMutationReceiptRepo(db *gorm.DB) sync.MutationReceiptRepository {
	return syncrepo.NewMutationReceiptRepo(db)
}
