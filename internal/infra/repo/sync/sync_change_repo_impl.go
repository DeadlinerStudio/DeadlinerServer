package syncrepo

import (
	"context"
	"time"

	domainSync "github.com/aritxonly/deadlinerserver/internal/domain/sync"
	persistencegorm "github.com/aritxonly/deadlinerserver/internal/infra/persistence/gorm"
	"gorm.io/gorm"
)

type syncChangeRepo struct {
	db *gorm.DB
}

func NewSyncChangeRepo(db *gorm.DB) domainSync.SyncChangeRepository {
	return &syncChangeRepo{db: db}
}

func (r *syncChangeRepo) Append(
	ctx context.Context,
	params domainSync.AppendSyncChangeParams,
) (*domainSync.SyncChange, error) {
	model, err := toSyncChangeModel(params)
	if err != nil {
		return nil, err
	}

	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		return nil, err
	}

	return toSyncChange(model), nil
}

func (r *syncChangeRepo) ListAfterChangeID(
	ctx context.Context,
	accountID, afterChangeID int64,
	limit int,
) ([]domainSync.SyncChange, error) {
	query := r.db.WithContext(ctx).
		Where("account_id = ? AND change_id > ?", accountID, afterChangeID).
		Order("change_id ASC")
	if limit > 0 {
		query = query.Limit(limit)
	}

	var models []persistencegorm.SyncChangeModel
	if err := query.Find(&models).Error; err != nil {
		return nil, err
	}

	changes := make([]domainSync.SyncChange, 0, len(models))
	for _, model := range models {
		changes = append(changes, *toSyncChange(model))
	}

	return changes, nil
}

func toSyncChangeModel(params domainSync.AppendSyncChangeParams) (persistencegorm.SyncChangeModel, error) {
	committedAt := time.Now().UTC()
	if params.CommittedAt != "" {
		parsed, err := time.Parse(time.RFC3339, params.CommittedAt)
		if err != nil {
			return persistencegorm.SyncChangeModel{}, err
		}
		committedAt = parsed.UTC()
	}

	return persistencegorm.SyncChangeModel{
		AccountID:   params.AccountID,
		DeviceUID:   params.DeviceUID,
		MutationID:  params.MutationID,
		EntityKind:  params.EntityKind,
		EntityUID:   params.EntityUID,
		Action:      params.Action,
		Payload:     params.Payload,
		CommittedAt: committedAt,
	}, nil
}

func toSyncChange(model persistencegorm.SyncChangeModel) *domainSync.SyncChange {
	return &domainSync.SyncChange{
		ChangeID:    model.ChangeID,
		AccountID:   model.AccountID,
		DeviceUID:   model.DeviceUID,
		MutationID:  model.MutationID,
		EntityKind:  model.EntityKind,
		EntityUID:   model.EntityUID,
		Action:      model.Action,
		Payload:     model.Payload,
		CommittedAt: model.CommittedAt.UTC().Format(time.RFC3339),
	}
}
