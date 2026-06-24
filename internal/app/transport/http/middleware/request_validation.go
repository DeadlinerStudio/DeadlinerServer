package middleware

import (
	"context"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

type errorResponse struct {
	Error     string `json:"error"`
	RequestID string `json:"request_id,omitempty"`
}

func RequireJSONMethods(methods ...string) app.HandlerFunc {
	allowed := make(map[string]struct{}, len(methods))
	for _, method := range methods {
		allowed[strings.ToUpper(strings.TrimSpace(method))] = struct{}{}
	}

	return func(ctx context.Context, c *app.RequestContext) {
		if _, ok := allowed[strings.ToUpper(string(c.Method()))]; ok {
			contentType := strings.ToLower(strings.TrimSpace(string(c.GetHeader("Content-Type"))))
			if contentType == "" || !strings.HasPrefix(contentType, "application/json") {
				writeJSONError(c, consts.StatusUnsupportedMediaType, "content type must be application/json")
				return
			}
		}
		c.Next(ctx)
	}
}

func EnforceMaxBodyBytes(maxBytes int) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		if maxBytes > 0 && len(c.Request.Body()) > maxBytes {
			writeJSONError(c, consts.StatusRequestEntityTooLarge, "request body is too large")
			return
		}
		c.Next(ctx)
	}
}

func writeJSONError(c *app.RequestContext, statusCode int, message string) {
	if c == nil {
		return
	}
	c.AbortWithStatusJSON(statusCode, errorResponse{
		Error:     message,
		RequestID: RequestIDFromContext(c),
	})
}
