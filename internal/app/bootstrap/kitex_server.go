package bootstrap

import (
	"errors"
	"fmt"
	"net"

	transportkitex "github.com/aritxonly/deadlinerserver/internal/app/transport/kitex"
	"github.com/aritxonly/deadlinerserver/internal/config"
	deadlinerv1 "github.com/aritxonly/deadlinerserver/kitex_gen/deadliner/v1/deadlinerservice"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/pkg/transmeta"
	"github.com/cloudwego/kitex/server"
)

func NewKitexServer(cfg config.Config, runtime *Runtime) (server.Server, error) {
	if runtime == nil {
		return nil, errors.New("runtime is required")
	}

	addr, err := net.ResolveTCPAddr("tcp", cfg.Service.Address)
	if err != nil {
		return nil, fmt.Errorf("resolve service address %s: %w", cfg.Service.Address, err)
	}

	return deadlinerv1.NewServer(
		transportkitex.NewHandler(runtime.AccountService, runtime.NewKitexSyncService()),
		server.WithServiceAddr(addr),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
			ServiceName: cfg.Service.Name,
		}),
		server.WithMiddleware(transportkitex.AccessLogMiddleware()),
		server.WithMetaHandler(transmeta.MetainfoServerHandler),
	), nil
}
