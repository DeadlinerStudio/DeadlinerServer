package middleware

import (
	"context"
	"strings"

	appauth "github.com/aritxonly/deadlinerserver/internal/app/auth"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

const accountUIDKey = "deadliner.account_uid"

func RequireAccessToken(parser appauth.AccessTokenParser) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		if parser == nil {
			writeJSONError(c, consts.StatusInternalServerError, "internal server error")
			return
		}

		token := extractAccessToken(c)
		if token == "" {
			c.Header("WWW-Authenticate", "Bearer")
			writeJSONError(c, consts.StatusUnauthorized, "unauthorized")
			return
		}

		claims, err := parser.Parse(token)
		if err != nil || claims == nil || strings.TrimSpace(claims.AccountUID) == "" {
			c.Header("WWW-Authenticate", "Bearer")
			writeJSONError(c, consts.StatusUnauthorized, "unauthorized")
			return
		}

		SetAccountUID(c, claims.AccountUID)
		c.Next(ctx)
	}
}

func SetAccountUID(c *app.RequestContext, accountUID string) {
	if c == nil {
		return
	}
	c.Set(accountUIDKey, strings.TrimSpace(accountUID))
}

func AccountUID(c *app.RequestContext) string {
	if c == nil {
		return ""
	}
	return strings.TrimSpace(c.GetString(accountUIDKey))
}

func extractAccessToken(c *app.RequestContext) string {
	if c == nil {
		return ""
	}

	if value := strings.TrimSpace(string(c.GetHeader("Authorization"))); value != "" {
		return stripBearerPrefix(value)
	}
	for _, key := range []string{"X-Deadliner-Access-Token", "Deadliner-Access-Token"} {
		if value := strings.TrimSpace(string(c.GetHeader(key))); value != "" {
			return value
		}
	}
	return ""
}

func stripBearerPrefix(value string) string {
	value = strings.TrimSpace(value)
	if strings.HasPrefix(strings.ToLower(value), "bearer ") {
		return strings.TrimSpace(value[7:])
	}
	return value
}
