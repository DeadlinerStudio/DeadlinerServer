package admin

import (
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
)

func extractAdminToken(c *app.RequestContext) string {
	if c == nil {
		return ""
	}
	if value := strings.TrimSpace(string(c.GetHeader("Authorization"))); value != "" {
		if strings.HasPrefix(strings.ToLower(value), "bearer ") {
			return strings.TrimSpace(value[7:])
		}
		return value
	}
	if value := strings.TrimSpace(string(c.GetHeader("X-Admin-Token"))); value != "" {
		return value
	}
	return ""
}
