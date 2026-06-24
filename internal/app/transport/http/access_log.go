package http

import (
	"context"
	"log"
	"time"

	httpmiddleware "github.com/aritxonly/deadlinerserver/internal/app/transport/http/middleware"
	"github.com/aritxonly/deadlinerserver/internal/utils/logutil"
	"github.com/cloudwego/hertz/pkg/app"
)

func AccessLog() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		startedAt := time.Now()
		method := string(c.Method())
		path := string(c.Path())
		clientIP := c.ClientIP()
		userAgent := string(c.UserAgent())
		reqBytes := len(c.Request.Body())

		c.Next(ctx)

		route := c.FullPath()
		if route == "" {
			route = path
		}

		routeSuffix := ""
		if route != "" && route != path {
			routeSuffix = " route=" + route
		}
		accountSuffix := ""
		if accountUID := httpmiddleware.AccountUID(c); accountUID != "" {
			accountSuffix = " account=" + accountUID
		}

		log.Printf(
			"HTTP rid=%s %d %s %s%s%s dur=%s ip=%s bytes=%d/%d ua=%q",
			httpmiddleware.RequestIDFromContext(c),
			c.Response.StatusCode(),
			method,
			path,
			routeSuffix,
			accountSuffix,
			logutil.Duration(time.Since(startedAt)),
			clientIP,
			reqBytes,
			len(c.Response.Body()),
			logutil.NormalizeUserAgent(userAgent),
		)
	}
}
