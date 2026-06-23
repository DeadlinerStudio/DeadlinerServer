package http

import (
	"context"
	"errors"

	appaccount "github.com/aritxonly/deadlinerserver/internal/app/service/account"
	domainaccount "github.com/aritxonly/deadlinerserver/internal/domain/account"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

type registerRequest struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	DisplayName string `json:"display_name"`
	DeviceUID   string `json:"device_uid"`
	DeviceName  string `json:"device_name"`
	Platform    string `json:"platform"`
}

type loginRequest struct {
	Email      string `json:"email"`
	Password   string `json:"password"`
	DeviceUID  string `json:"device_uid"`
	DeviceName string `json:"device_name"`
	Platform   string `json:"platform"`
}

type refreshSessionRequest struct {
	RefreshToken string `json:"refresh_token"`
	DeviceUID    string `json:"device_uid"`
}

type sessionEnvelope struct {
	Session *sessionPayload `json:"session"`
}

type sessionPayload struct {
	AccountUID   string `json:"account_uid"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    string `json:"expires_at"`
}

func (h *Handler) register(ctx context.Context, c *app.RequestContext) {
	if h.accountService == nil {
		writeError(c, consts.StatusInternalServerError, errors.New("account service is not configured"))
		return
	}

	var req registerRequest
	if err := c.BindJSON(&req); err != nil {
		writeBadRequest(c, err)
		return
	}

	result, err := h.accountService.Register(ctx, appaccount.RegisterInput{
		Email:       req.Email,
		Password:    req.Password,
		DisplayName: req.DisplayName,
		DeviceUID:   req.DeviceUID,
		DeviceName:  req.DeviceName,
		Platform:    req.Platform,
	})
	if err != nil {
		writeAccountError(c, err)
		return
	}

	c.JSON(consts.StatusOK, toSessionEnvelope(result))
}

func (h *Handler) login(ctx context.Context, c *app.RequestContext) {
	if h.accountService == nil {
		writeError(c, consts.StatusInternalServerError, errors.New("account service is not configured"))
		return
	}

	var req loginRequest
	if err := c.BindJSON(&req); err != nil {
		writeBadRequest(c, err)
		return
	}

	result, err := h.accountService.Login(ctx, appaccount.LoginInput{
		Email:      req.Email,
		Password:   req.Password,
		DeviceUID:  req.DeviceUID,
		DeviceName: req.DeviceName,
		Platform:   req.Platform,
	})
	if err != nil {
		writeAccountError(c, err)
		return
	}

	c.JSON(consts.StatusOK, toSessionEnvelope(result))
}

func (h *Handler) refreshSession(ctx context.Context, c *app.RequestContext) {
	if h.accountService == nil {
		writeError(c, consts.StatusInternalServerError, errors.New("account service is not configured"))
		return
	}

	var req refreshSessionRequest
	if err := c.BindJSON(&req); err != nil {
		writeBadRequest(c, err)
		return
	}

	result, err := h.accountService.RefreshSession(ctx, appaccount.RefreshSessionInput{
		RefreshToken: req.RefreshToken,
		DeviceUID:    req.DeviceUID,
	})
	if err != nil {
		writeAccountError(c, err)
		return
	}

	c.JSON(consts.StatusOK, toSessionEnvelope(result))
}

func toSessionEnvelope(bundle *domainaccount.SessionBundle) sessionEnvelope {
	if bundle == nil {
		return sessionEnvelope{}
	}

	return sessionEnvelope{
		Session: &sessionPayload{
			AccountUID:   bundle.AccountUID,
			AccessToken:  bundle.AccessToken,
			RefreshToken: bundle.RefreshToken,
			ExpiresAt:    bundle.ExpiresAt,
		},
	}
}
