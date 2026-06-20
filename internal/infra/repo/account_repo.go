package repo

import (
	domainAccount "github.com/aritxonly/deadlinerserver/internal/domain/account"
	accountrepo "github.com/aritxonly/deadlinerserver/internal/infra/repo/account"
	"gorm.io/gorm"
)

func NewAccountRepo(db *gorm.DB) domainAccount.Repository {
	return accountrepo.NewAccountRepo(db)
}
