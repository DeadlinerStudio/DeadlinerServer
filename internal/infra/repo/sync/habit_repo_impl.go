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

type habitRepo struct {
	db *gorm.DB
}

func NewHabitRepo(db *gorm.DB) domainSync.HabitRepository {
	return &habitRepo{db: db}
}

func (r *habitRepo) FindByDDLUID(ctx context.Context, accountID int64, ddlUID string) (*domainSync.HabitChange, error) {
	var model persistencegorm.HabitDocModel
	err := r.db.WithContext(ctx).
		Where("account_id = ? AND ddl_uid = ?", accountID, ddlUID).
		Take(&model).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return toHabitChange(model)
}

func (r *habitRepo) Save(ctx context.Context, params domainSync.SaveHabitParams) error {
	model, err := toHabitModel(params)
	if err != nil {
		return err
	}

	var existing persistencegorm.HabitDocModel
	err = r.db.WithContext(ctx).
		Where("account_id = ? AND ddl_uid = ?", params.AccountID, params.Document.DDLUID).
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

func (r *habitRepo) ListAfterChangeID(
	ctx context.Context,
	accountID int64,
	afterChangeID int64,
	limit int,
	includeDeleted bool,
) ([]domainSync.HabitChange, error) {
	query := r.db.WithContext(ctx).
		Where("account_id = ? AND server_change_id > ?", accountID, afterChangeID)
	if !includeDeleted {
		query = query.Where("deleted = ?", false)
	}
	if limit > 0 {
		query = query.Limit(limit)
	}

	var models []persistencegorm.HabitDocModel
	if err := query.Order("server_change_id ASC").Find(&models).Error; err != nil {
		return nil, err
	}

	changes := make([]domainSync.HabitChange, 0, len(models))
	for _, model := range models {
		change, err := toHabitChange(model)
		if err != nil {
			return nil, err
		}
		changes = append(changes, *change)
	}

	return changes, nil
}

func toHabitModel(params domainSync.SaveHabitParams) (persistencegorm.HabitDocModel, error) {
	payload, err := json.Marshal(params.Document)
	if err != nil {
		return persistencegorm.HabitDocModel{}, err
	}

	committedAt, err := time.Parse(time.RFC3339, params.ServerVersion.CommittedAt)
	if err != nil {
		return persistencegorm.HabitDocModel{}, err
	}

	model := persistencegorm.HabitDocModel{
		AccountID:      params.AccountID,
		DDLUID:         params.Document.DDLUID,
		Payload:        payload,
		Deleted:        params.Deleted,
		ServerChangeID: params.ServerVersion.ChangeID,
		CommittedAt:    committedAt,
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

func toHabitChange(model persistencegorm.HabitDocModel) (*domainSync.HabitChange, error) {
	var document domainSync.HabitDocument
	if len(model.Payload) > 0 {
		if err := json.Unmarshal(model.Payload, &document); err != nil {
			return nil, err
		}
	}
	if document.DDLUID == "" {
		document.DDLUID = model.DDLUID
	}

	return &domainSync.HabitChange{
		EntityUID: model.DDLUID,
		Deleted:   model.Deleted,
		ServerVersion: domainSync.ServerVersion{
			ChangeID:    model.ServerChangeID,
			CommittedAt: model.CommittedAt.UTC().Format(time.RFC3339),
		},
		Document: document,
	}, nil
}
