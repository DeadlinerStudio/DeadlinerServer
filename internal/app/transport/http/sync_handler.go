package http

import (
	"context"
	"errors"
	"strconv"
	"strings"

	appsync "github.com/aritxonly/deadlinerserver/internal/app/service/sync"
	domainsync "github.com/aritxonly/deadlinerserver/internal/domain/sync"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

type pushChangesRequest struct {
	DeviceUID  string                `json:"device_uid"`
	BaseCursor string                `json:"base_cursor"`
	Mutations  []domainsync.Mutation `json:"mutations"`
}

type pullChangesResponse struct {
	DeadlineChanges []domainsync.DeadlineChange `json:"deadline_changes"`
	HabitChanges    []domainsync.HabitChange    `json:"habit_changes"`
	NextCursor      string                      `json:"next_cursor"`
	HasMore         bool                        `json:"has_more"`
}

type pushChangesResponse struct {
	Results         []domainsync.MutationResult `json:"results"`
	DeadlineChanges []domainsync.DeadlineChange `json:"deadline_changes"`
	HabitChanges    []domainsync.HabitChange    `json:"habit_changes"`
	NextCursor      string                      `json:"next_cursor"`
}

func (h *Handler) pullChanges(ctx context.Context, c *app.RequestContext) {
	if h.syncService == nil {
		writeError(c, consts.StatusInternalServerError, errors.New("sync service is not configured"))
		return
	}

	limit, err := parseOptionalInt32(c.Query("limit"))
	if err != nil {
		writeBadRequest(c, err)
		return
	}

	includeDeleted, err := parseOptionalBool(c.Query("include_deleted"))
	if err != nil {
		writeBadRequest(c, err)
		return
	}

	result, err := h.syncService.PullChanges(
		withRequestAuth(ctx, c),
		appsync.PullChangesInput{
			DeviceUID:      strings.TrimSpace(c.Query("device_uid")),
			Cursor:         strings.TrimSpace(c.Query("cursor")),
			Limit:          limit,
			IncludeDeleted: includeDeleted,
		},
	)
	if err != nil {
		writeSyncError(c, err)
		return
	}

	c.JSON(consts.StatusOK, pullChangesResponse{
		DeadlineChanges: result.DeadlineChanges,
		HabitChanges:    result.HabitChanges,
		NextCursor:      result.NextCursor,
		HasMore:         result.HasMore,
	})
}

func (h *Handler) pushChanges(ctx context.Context, c *app.RequestContext) {
	if h.syncService == nil {
		writeError(c, consts.StatusInternalServerError, errors.New("sync service is not configured"))
		return
	}

	var req pushChangesRequest
	if err := c.BindJSON(&req); err != nil {
		writeBadRequest(c, err)
		return
	}

	result, err := h.syncService.PushChanges(
		withRequestAuth(ctx, c),
		appsync.PushChangesInput{
			DeviceUID:  req.DeviceUID,
			BaseCursor: req.BaseCursor,
			Mutations:  req.Mutations,
		},
	)
	if err != nil {
		writeSyncError(c, err)
		return
	}

	c.JSON(consts.StatusOK, pushChangesResponse{
		Results:         result.Results,
		DeadlineChanges: result.DeadlineChanges,
		HabitChanges:    result.HabitChanges,
		NextCursor:      result.NextCursor,
	})
}

func parseOptionalInt32(value string) (int32, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return 0, nil
	}

	parsed, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		return 0, err
	}
	return int32(parsed), nil
}

func parseOptionalBool(value string) (bool, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return false, nil
	}

	return strconv.ParseBool(value)
}
