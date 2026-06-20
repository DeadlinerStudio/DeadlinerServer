package syncrepo

import (
	"context"
	"errors"

	domainSync "github.com/aritxonly/deadlinerserver/internal/domain/sync"
	persistencegorm "github.com/aritxonly/deadlinerserver/internal/infra/persistence/gorm"
	"gorm.io/gorm"
)

type mutationReceiptRepo struct {
	db *gorm.DB
}

func NewMutationReceiptRepo(db *gorm.DB) domainSync.MutationReceiptRepository {
	return &mutationReceiptRepo{db: db}
}

func (r *mutationReceiptRepo) Find(
	ctx context.Context,
	accountID int64,
	deviceUID, mutationID string,
) (*domainSync.MutationReceipt, error) {
	var model persistencegorm.MutationReceiptModel
	err := r.db.WithContext(ctx).
		Where("account_id = ? AND device_uid = ? AND mutation_id = ?", accountID, deviceUID, mutationID).
		Take(&model).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return toMutationReceipt(model), nil
}

func (r *mutationReceiptRepo) Save(ctx context.Context, receipt *domainSync.MutationReceipt) error {
	if receipt == nil {
		return nil
	}

	model := toMutationReceiptModel(receipt)

	var existing persistencegorm.MutationReceiptModel
	err := r.db.WithContext(ctx).
		Where("account_id = ? AND device_uid = ? AND mutation_id = ?", receipt.AccountID, receipt.DeviceUID, receipt.MutationID).
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

func toMutationReceiptModel(receipt *domainSync.MutationReceipt) persistencegorm.MutationReceiptModel {
	model := persistencegorm.MutationReceiptModel{
		AccountID:     receipt.AccountID,
		DeviceUID:     receipt.DeviceUID,
		MutationID:    receipt.MutationID,
		EntityKind:    receipt.EntityKind,
		EntityUID:     receipt.EntityUID,
		Status:        receipt.Status,
		Replayed:      receipt.Replayed,
		ResultPayload: receipt.ResultPayload,
	}

	if receipt.ResultChangeID > 0 {
		model.ResultChangeID = &receipt.ResultChangeID
	}

	return model
}

func toMutationReceipt(model persistencegorm.MutationReceiptModel) *domainSync.MutationReceipt {
	receipt := &domainSync.MutationReceipt{
		AccountID:     model.AccountID,
		DeviceUID:     model.DeviceUID,
		MutationID:    model.MutationID,
		EntityKind:    model.EntityKind,
		EntityUID:     model.EntityUID,
		Status:        model.Status,
		Replayed:      model.Replayed,
		ResultPayload: model.ResultPayload,
	}

	if model.ResultChangeID != nil {
		receipt.ResultChangeID = *model.ResultChangeID
	}

	return receipt
}
