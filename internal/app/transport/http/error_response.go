package http

import (
	"errors"

	appauth "github.com/aritxonly/deadlinerserver/internal/app/auth"
	domainaccount "github.com/aritxonly/deadlinerserver/internal/domain/account"
	"github.com/aritxonly/deadlinerserver/internal/infra/provider"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

type errorResponse struct {
	Error string `json:"error"`
}

func writeBadRequest(c *app.RequestContext, err error) {
	writeError(c, consts.StatusBadRequest, err)
}

func writeAccountError(c *app.RequestContext, err error) {
	switch {
	case errors.Is(err, domainaccount.ErrAccountAlreadyExists):
		writeError(c, consts.StatusConflict, err)
	case errors.Is(err, domainaccount.ErrDeviceMismatch):
		writeError(c, consts.StatusConflict, err)
	case errors.Is(err, domainaccount.ErrInvalidCredentials),
		errors.Is(err, domainaccount.ErrInvalidRefreshToken),
		errors.Is(err, domainaccount.ErrExpiredRefreshToken),
		errors.Is(err, appauth.ErrUnauthenticated),
		errors.Is(err, provider.ErrInvalidAccessToken):
		writeError(c, consts.StatusUnauthorized, err)
	default:
		writeError(c, consts.StatusInternalServerError, err)
	}
}

func writeSyncError(c *app.RequestContext, err error) {
	switch {
	case errors.Is(err, appauth.ErrUnauthenticated),
		errors.Is(err, provider.ErrInvalidAccessToken):
		writeError(c, consts.StatusUnauthorized, err)
	default:
		writeError(c, consts.StatusInternalServerError, err)
	}
}

func writeError(c *app.RequestContext, statusCode int, err error) {
	if c == nil {
		return
	}
	message := "internal server error"
	if err != nil {
		message = err.Error()
	}
	c.AbortWithStatusJSON(statusCode, errorResponse{Error: message})
}
