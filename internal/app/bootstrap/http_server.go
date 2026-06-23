package bootstrap

import (
	"errors"
	"time"

	transporthttp "github.com/aritxonly/deadlinerserver/internal/app/transport/http"
	"github.com/aritxonly/deadlinerserver/internal/config"
	hertzserver "github.com/cloudwego/hertz/pkg/app/server"
)

func NewHertzServer(cfg config.Config, runtime *Runtime) (*hertzserver.Hertz, error) {
	if runtime == nil {
		return nil, errors.New("runtime is required")
	}

	server := hertzserver.Default(
		hertzserver.WithHostPorts(cfg.HTTP.Address),
		hertzserver.WithReadTimeout(time.Duration(cfg.HTTP.ReadTimeoutSeconds)*time.Second),
		hertzserver.WithWriteTimeout(time.Duration(cfg.HTTP.WriteTimeoutSeconds)*time.Second),
		hertzserver.WithIdleTimeout(time.Duration(cfg.HTTP.IdleTimeoutSeconds)*time.Second),
	)
	server.Use(transporthttp.AccessLog())

	handler := transporthttp.NewHandler(runtime.AccountService, runtime.NewHTTPSyncService())
	handler.RegisterRoutes(server)
	return server, nil
}
