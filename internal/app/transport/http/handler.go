package http

import (
	"context"

	appaccount "github.com/aritxonly/deadlinerserver/internal/app/service/account"
	appsync "github.com/aritxonly/deadlinerserver/internal/app/service/sync"
	"github.com/cloudwego/hertz/pkg/app"
	hertzserver "github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
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

func (h *Handler) RegisterRoutes(server *hertzserver.Hertz) {
	server.GET("/healthz", func(ctx context.Context, c *app.RequestContext) {
		c.JSON(consts.StatusOK, utils.H{"status": "ok"})
	})

	v1 := server.Group("/v1")
	auth := v1.Group("/auth")
	auth.POST("/register", h.register)
	auth.POST("/login", h.login)
	auth.POST("/refresh", h.refreshSession)

	sync := v1.Group("/sync")
	sync.GET("/pull", h.pullChanges)
	sync.POST("/push", h.pushChanges)
}
