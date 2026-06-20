package bootstrap

import (
	"fmt"
	"net"

	"github.com/aritxonly/deadlinerserver/internal/config"
	deadlinerv1 "github.com/aritxonly/deadlinerserver/kitex_gen/deadliner/v1/deadlinerservice"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/pkg/transmeta"
	"github.com/cloudwego/kitex/server"
)

func NewKitexServer(cfg config.Config) (server.Server, error) {
	handler, err := NewKitexHandler(cfg)
	if err != nil {
		return nil, err
	}

	addr, err := net.ResolveTCPAddr("tcp", cfg.Service.Address)
	if err != nil {
		return nil, fmt.Errorf("resolve service address %s: %w", cfg.Service.Address, err)
	}

	return deadlinerv1.NewServer(
		handler,
		server.WithServiceAddr(addr),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
			ServiceName: cfg.Service.Name,
		}),
		server.WithMetaHandler(transmeta.MetainfoServerHandler),
	), nil
}
