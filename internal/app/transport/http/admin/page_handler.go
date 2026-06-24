package admin

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

func Page(_ context.Context, c *app.RequestContext) {
	c.Data(consts.StatusOK, "text/html; charset=utf-8", []byte(configPageHTML))
}
