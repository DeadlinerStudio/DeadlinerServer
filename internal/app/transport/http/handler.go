package http

import (
	"context"

	appauth "github.com/aritxonly/deadlinerserver/internal/app/auth"
	appaccount "github.com/aritxonly/deadlinerserver/internal/app/service/account"
	appsync "github.com/aritxonly/deadlinerserver/internal/app/service/sync"
	httpmiddleware "github.com/aritxonly/deadlinerserver/internal/app/transport/http/middleware"
	"github.com/aritxonly/deadlinerserver/internal/config"
	"github.com/cloudwego/hertz/pkg/app"
	hertzserver "github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

type Handler struct {
	accountService    appaccount.Service
	syncService       appsync.Service
	accessTokenParser appauth.AccessTokenParser
	httpConfig        config.HTTPConfig
}

func NewHandler(
	accountService appaccount.Service,
	syncService appsync.Service,
	accessTokenParser appauth.AccessTokenParser,
	httpConfig config.HTTPConfig,
) *Handler {
	return &Handler{
		accountService:    accountService,
		syncService:       syncService,
		accessTokenParser: accessTokenParser,
		httpConfig:        httpConfig,
	}
}

func (h *Handler) RegisterRoutes(server *hertzserver.Hertz) {
	server.GET("/healthz", func(ctx context.Context, c *app.RequestContext) {
		c.JSON(consts.StatusOK, utils.H{"status": "ok"})
	})

	v1 := server.Group("/v1")
	v1.Use(httpmiddleware.RequireJSONMethods("POST"))
	auth := v1.Group("/auth")
	auth.Use(httpmiddleware.LimitByClientIP("auth", h.httpConfig.AuthRateLimitPerMinute, h.httpConfig.AuthRateLimitBurst))
	auth.POST("/register", h.register)
	auth.POST("/login", h.login)
	auth.POST("/refresh", h.refreshSession)

	sync := v1.Group("/sync")
	sync.Use(
		httpmiddleware.RequireAccessToken(h.accessTokenParser),
		httpmiddleware.LimitByClientIP("sync", h.httpConfig.SyncRateLimitPerMinute, h.httpConfig.SyncRateLimitBurst),
	)
	sync.GET("/pull", h.pullChanges)
	sync.POST("/push", h.pushChanges)
}
