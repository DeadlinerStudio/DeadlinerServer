package kitex

import (
	appaccount "github.com/aritxonly/deadlinerserver/internal/app/service/account"
	appsync "github.com/aritxonly/deadlinerserver/internal/app/service/sync"
)

type Handler struct {
	accountService appaccount.Service
	syncService    appsync.Service
}

func NewHandler(accountService appaccount.Service, syncService appsync.Service) *Handler {
	return &Handler{
		accountService: accountService,
		syncService:    syncService,
	}
}
