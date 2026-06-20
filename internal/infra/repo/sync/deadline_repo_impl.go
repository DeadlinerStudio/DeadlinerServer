package syncrepo

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	domainSync "github.com/aritxonly/deadlinerserver/internal/domain/sync"
	persistencegorm "github.com/aritxonly/deadlinerserver/internal/infra/persistence/gorm"
	"gorm.io/gorm"
)

type deadlineRepo struct {
	db *gorm.DB
}

func NewDeadlineRepo(db *gorm.DB) domainSync.DeadlineRepository {
	return &deadlineRepo{db: db}
}

func (r *deadlineRepo) FindByUID(ctx context.Context, accountID int64, uid string) (*domainSync.DeadlineChange, error) {
	var model persistencegorm.DeadlineItemModel
	err := r.db.WithContext(ctx).
		Where("account_id = ? AND uid = ?", accountID, uid).
		Take(&model).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return toDeadlineChange(model)
}

func (r *deadlineRepo) Save(ctx context.Context, params domainSync.SaveDeadlineParams) error {
	model, err := toDeadlineModel(params)
	if err != nil {
		return err
	}

	var existing persistencegorm.DeadlineItemModel
	err = r.db.WithContext(ctx).
		Where("account_id = ? AND uid = ?", params.AccountID, params.Document.UID).
		Take(&existing).Error
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return r.db.WithContext(ctx).Create(&model).Error
	case err != nil:
		return err
	default:
		model.ID = existing.ID
		return r.db.WithContext(ctx).Save(&model).Error
	}
}

func (r *deadlineRepo) ListAfterChangeID(
	ctx context.Context,
	accountID int64,
	afterChangeID int64,
	limit int,
	includeDeleted bool,
) ([]domainSync.DeadlineChange, error) {
	query := r.db.WithContext(ctx).
		Where("account_id = ? AND server_change_id > ?", accountID, afterChangeID)
	if !includeDeleted {
		query = query.Where("deleted = ?", false)
	}
	if limit > 0 {
		query = query.Limit(limit)
	}

	var models []persistencegorm.DeadlineItemModel
	if err := query.Order("server_change_id ASC").Find(&models).Error; err != nil {
		return nil, err
	}

	changes := make([]domainSync.DeadlineChange, 0, len(models))
	for _, model := range models {
		change, err := toDeadlineChange(model)
		if err != nil {
			return nil, err
		}
		changes = append(changes, *change)
	}

	return changes, nil
}

func toDeadlineModel(params domainSync.SaveDeadlineParams) (persistencegorm.DeadlineItemModel, error) {
	subTasks, err := json.Marshal(params.Document.SubTasks)
	if err != nil {
		return persistencegorm.DeadlineItemModel{}, err
	}

	committedAt, err := time.Parse(time.RFC3339, params.ServerVersion.CommittedAt)
	if err != nil {
		return persistencegorm.DeadlineItemModel{}, err
	}

	model := persistencegorm.DeadlineItemModel{
		AccountID:         params.AccountID,
		UID:               params.Document.UID,
		LegacyID:          params.Document.LegacyID,
		Name:              params.Document.Name,
		StartTime:         params.Document.StartTime,
		EndTime:           params.Document.EndTime,
		State:             string(params.Document.State),
		CompleteTime:      params.Document.CompleteTime,
		Note:              params.Document.Note,
		IsStared:          params.Document.IsStared,
		Type:              string(params.Document.Type),
		HabitCount:        params.Document.HabitCount,
		HabitTotalCount:   params.Document.HabitTotalCount,
		CalendarEvent:     params.Document.CalendarEvent,
		BusinessTimestamp: params.Document.Timestamp,
		SubTasks:          subTasks,
		Deleted:           params.Deleted,
		ServerChangeID:    params.ServerVersion.ChangeID,
		CommittedAt:       committedAt,
	}

	if params.ClientVersion != nil {
		model.ClientVerTS = &params.ClientVersion.TS
		model.ClientVerCtr = &params.ClientVersion.Ctr
		model.ClientVerDev = &params.ClientVersion.Dev
	}
	if params.UpdatedByDeviceUID != "" {
		model.UpdatedByDeviceUID = &params.UpdatedByDeviceUID
	}

	return model, nil
}

func toDeadlineChange(model persistencegorm.DeadlineItemModel) (*domainSync.DeadlineChange, error) {
	var subTasks []domainSync.SubTask
	if len(model.SubTasks) > 0 {
		if err := json.Unmarshal(model.SubTasks, &subTasks); err != nil {
			return nil, err
		}
	}

	return &domainSync.DeadlineChange{
		EntityUID: model.UID,
		Deleted:   model.Deleted,
		ServerVersion: domainSync.ServerVersion{
			ChangeID:    model.ServerChangeID,
			CommittedAt: model.CommittedAt.UTC().Format(time.RFC3339),
		},
		Document: domainSync.DeadlineDocument{
			UID:             model.UID,
			LegacyID:        model.LegacyID,
			Name:            model.Name,
			StartTime:       model.StartTime,
			EndTime:         model.EndTime,
			State:           domainSync.DeadlineState(model.State),
			CompleteTime:    model.CompleteTime,
			Note:            model.Note,
			IsStared:        model.IsStared,
			Type:            domainSync.DeadlineType(model.Type),
			HabitCount:      model.HabitCount,
			HabitTotalCount: model.HabitTotalCount,
			CalendarEvent:   model.CalendarEvent,
			Timestamp:       model.BusinessTimestamp,
			SubTasks:        subTasks,
		},
	}, nil
}
