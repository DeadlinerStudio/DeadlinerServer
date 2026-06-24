package http

import (
	"errors"
	"log"

	appauth "github.com/aritxonly/deadlinerserver/internal/app/auth"
	httpmiddleware "github.com/aritxonly/deadlinerserver/internal/app/transport/http/middleware"
	domainaccount "github.com/aritxonly/deadlinerserver/internal/domain/account"
	"github.com/aritxonly/deadlinerserver/internal/infra/provider"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

type errorResponse struct {
	Error     string `json:"error"`
	RequestID string `json:"request_id,omitempty"`
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

	requestID := httpmiddleware.RequestIDFromContext(c)
	message := "internal server error"
	if statusCode >= 500 {
		if err != nil {
			log.Printf(
				"HTTP_ERR rid=%s status=%d path=%s err=%q",
				requestID,
				statusCode,
				string(c.Path()),
				err.Error(),
			)
		}
	} else if err != nil {
		message = err.Error()
	}
	c.AbortWithStatusJSON(statusCode, errorResponse{
		Error:     message,
		RequestID: requestID,
	})
}
