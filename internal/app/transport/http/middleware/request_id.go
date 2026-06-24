package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
)

const (
	requestIDHeader = "X-Request-Id"
	requestIDKey    = "deadliner.request_id"
)

func RequestID() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		id := normalizeRequestID(string(c.GetHeader(requestIDHeader)))
		if id == "" {
			id = newRequestID()
		}
		SetRequestID(c, id)
		c.Header(requestIDHeader, id)
		c.Next(ctx)
	}
}

func SetRequestID(c *app.RequestContext, id string) {
	if c == nil {
		return
	}
	c.Set(requestIDKey, strings.TrimSpace(id))
}

func RequestIDFromContext(c *app.RequestContext) string {
	if c == nil {
		return ""
	}
	return strings.TrimSpace(c.GetString(requestIDKey))
}

func newRequestID() string {
	var bytes [8]byte
	if _, err := rand.Read(bytes[:]); err != nil {
		return "deadliner-unknown"
	}
	return hex.EncodeToString(bytes[:])
}

func normalizeRequestID(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	if len(value) > 64 {
		value = value[:64]
	}
	for _, char := range value {
		switch {
		case char >= 'a' && char <= 'z':
		case char >= 'A' && char <= 'Z':
		case char >= '0' && char <= '9':
		case char == '-', char == '_', char == '.':
		default:
			return ""
		}
	}
	return value
}
