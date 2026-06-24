package bootstrap

import (
	"errors"
	"time"

	transporthttp "github.com/aritxonly/deadlinerserver/internal/app/transport/http"
	httpmiddleware "github.com/aritxonly/deadlinerserver/internal/app/transport/http/middleware"
	"github.com/aritxonly/deadlinerserver/internal/config"
	hertzserver "github.com/cloudwego/hertz/pkg/app/server"
)

func NewHertzServer(cfg config.Config, runtime *Runtime) (*hertzserver.Hertz, error) {
	if runtime == nil {
		return nil, errors.New("runtime is required")
	}

	server := hertzserver.New(
		hertzserver.WithHostPorts(cfg.HTTP.Address),
		hertzserver.WithReadTimeout(time.Duration(cfg.HTTP.ReadTimeoutSeconds)*time.Second),
		hertzserver.WithWriteTimeout(time.Duration(cfg.HTTP.WriteTimeoutSeconds)*time.Second),
		hertzserver.WithIdleTimeout(time.Duration(cfg.HTTP.IdleTimeoutSeconds)*time.Second),
	)
	server.Use(
		httpmiddleware.RequestID(),
		httpmiddleware.Recovery(),
		httpmiddleware.SecurityHeaders(),
		httpmiddleware.EnforceMaxBodyBytes(cfg.HTTP.MaxRequestBodyBytes),
		httpmiddleware.LimitByClientIP("http", cfg.HTTP.RateLimitPerMinute, cfg.HTTP.RateLimitBurst),
		transporthttp.AccessLog(),
	)

	handler := transporthttp.NewHandler(
		runtime.AccountService,
		runtime.NewHTTPSyncService(),
		runtime.AdminConfigService,
		runtime.AccessTokenCodec,
		cfg.HTTP,
		runtime.AdminRuntimeConfig,
	)
	handler.RegisterRoutes(server)
	return server, nil
}
