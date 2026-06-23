package http

import (
	"context"
	"strings"

	appauth "github.com/aritxonly/deadlinerserver/internal/app/auth"
	"github.com/cloudwego/hertz/pkg/app"
)

func withRequestAuth(ctx context.Context, reqCtx *app.RequestContext) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	if reqCtx == nil {
		return ctx
	}

	if value := strings.TrimSpace(string(reqCtx.GetHeader("Authorization"))); value != "" {
		ctx = appauth.WithAuthorization(ctx, value)
	}

	for _, key := range []string{"X-Deadliner-Access-Token", "Deadliner-Access-Token"} {
		if value := strings.TrimSpace(string(reqCtx.GetHeader(key))); value != "" {
			ctx = appauth.WithAccessToken(ctx, value)
			break
		}
	}

	return ctx
}
