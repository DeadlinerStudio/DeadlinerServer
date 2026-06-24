package admin

import (
	"context"
	"strings"

	appadmin "github.com/aritxonly/deadlinerserver/internal/app/adminconfig"
	httpmiddleware "github.com/aritxonly/deadlinerserver/internal/app/transport/http/middleware"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

type APIHandler struct {
	service   appadmin.Service
	authToken string
}

func NewAPIHandler(service appadmin.Service, authToken string) *APIHandler {
	return &APIHandler{
		service:   service,
		authToken: strings.TrimSpace(authToken),
	}
}

func (h *APIHandler) GetConfig(ctx context.Context, c *app.RequestContext) {
	if !h.authorize(c) {
		return
	}
	snapshot, err := h.service.GetSnapshot(ctx)
	if err != nil {
		writeError(c, consts.StatusInternalServerError, "internal server error")
		return
	}
	c.JSON(consts.StatusOK, snapshot)
}

func (h *APIHandler) UpdateConfig(ctx context.Context, c *app.RequestContext) {
	if !h.authorize(c) {
		return
	}

	var input appadmin.UpdateInput
	if err := c.BindJSON(&input); err != nil {
		writeError(c, consts.StatusBadRequest, err.Error())
		return
	}

	snapshot, err := h.service.Update(ctx, input)
	if err != nil {
		writeError(c, consts.StatusBadRequest, err.Error())
		return
	}
	c.JSON(consts.StatusOK, snapshot)
}

func (h *APIHandler) authorize(c *app.RequestContext) bool {
	if strings.TrimSpace(h.authToken) == "" {
		writeError(c, consts.StatusServiceUnavailable, "admin backend is not configured")
		return false
	}
	if extractAdminToken(c) != h.authToken {
		c.Header("WWW-Authenticate", "Bearer")
		writeError(c, consts.StatusUnauthorized, "unauthorized")
		return false
	}
	return true
}

func writeError(c *app.RequestContext, statusCode int, message string) {
	if c == nil {
		return
	}
	c.AbortWithStatusJSON(statusCode, map[string]string{
		"error":      message,
		"request_id": httpmiddleware.RequestIDFromContext(c),
	})
}
