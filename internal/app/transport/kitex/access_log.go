package kitex

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/aritxonly/deadlinerserver/internal/utils/logutil"
	"github.com/cloudwego/kitex/pkg/endpoint"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
)

func AccessLogMiddleware() endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, req, resp interface{}) error {
			startedAt := time.Now()
			err := next(ctx, req, resp)

			serviceName, packageName, methodName, seqID, callerAddr, calleeAddr := extractRPCFields(rpcinfo.GetRPCInfo(ctx))
			outcome := "OK"
			if err != nil {
				outcome = "ERR"
			}
			log.Printf(
				"KITEX %s %s.%s dur=%s seq=%d from=%s to=%s err=%q",
				outcome,
				fullServiceName(packageName, serviceName),
				methodName,
				logutil.Duration(time.Since(startedAt)),
				seqID,
				callerAddr,
				calleeAddr,
				logutil.NormalizeErr(err),
			)

			return err
		}
	}
}

func extractRPCFields(ri rpcinfo.RPCInfo) (string, string, string, int32, string, string) {
	if ri == nil {
		return "", "", "", 0, "", ""
	}

	invocation := ri.Invocation()
	packageName := ""
	serviceName := ""
	methodName := ""
	var seqID int32
	if invocation != nil {
		packageName = invocation.PackageName()
		serviceName = invocation.ServiceName()
		methodName = invocation.MethodName()
		seqID = invocation.SeqID()
	}

	return serviceName, packageName, methodName, seqID, endpointAddr(ri.From()), endpointAddr(ri.To())
}

func endpointAddr(info rpcinfo.EndpointInfo) string {
	if info == nil {
		return ""
	}
	return formatAddr(info.Address())
}

func formatAddr(addr net.Addr) string {
	if addr == nil {
		return ""
	}
	return addr.String()
}

func fullServiceName(packageName string, serviceName string) string {
	if packageName == "" {
		return serviceName
	}
	if serviceName == "" {
		return packageName
	}
	return packageName + "." + serviceName
}
