package accountrepo

import (
	"context"
	"errors"

	domainAccount "github.com/aritxonly/deadlinerserver/internal/domain/account"
	persistencegorm "github.com/aritxonly/deadlinerserver/internal/infra/persistence/gorm"
	"gorm.io/gorm"
)

type accountRepo struct {
	db *gorm.DB
}

func NewAccountRepo(db *gorm.DB) domainAccount.Repository {
	return &accountRepo{db: db}
}

func (r *accountRepo) FindAccountByEmail(ctx context.Context, email string) (*domainAccount.Account, error) {
	var model persistencegorm.AccountModel
	err := r.db.WithContext(ctx).
		Where("email = ?", email).
		Take(&model).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return toAccount(model), nil
}

func (r *accountRepo) FindAccountByUID(ctx context.Context, uid string) (*domainAccount.Account, error) {
	var model persistencegorm.AccountModel
	err := r.db.WithContext(ctx).
		Where("account_uid = ?", uid).
		Take(&model).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return toAccount(model), nil
}

func (r *accountRepo) FindAccountByID(ctx context.Context, id int64) (*domainAccount.Account, error) {
	var model persistencegorm.AccountModel
	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		Take(&model).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return toAccount(model), nil
}

func (r *accountRepo) FindSessionByRefreshTokenHash(
	ctx context.Context,
	refreshTokenHash string,
) (*domainAccount.Session, error) {
	var model persistencegorm.SessionModel
	err := r.db.WithContext(ctx).
		Where("refresh_token_hash = ?", refreshTokenHash).
		Take(&model).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return toSession(model), nil
}

func (r *accountRepo) SaveAccount(ctx context.Context, acc *domainAccount.Account) error {
	model := toAccountModel(acc)

	var existing persistencegorm.AccountModel
	err := r.db.WithContext(ctx).
		Where("account_uid = ?", acc.AccountUID).
		Take(&existing).Error
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
			return err
		}
		acc.ID = model.ID
		return nil
	case err != nil:
		return err
	default:
		model.ID = existing.ID
		if err := r.db.WithContext(ctx).Save(&model).Error; err != nil {
			return err
		}
		acc.ID = model.ID
		return nil
	}
}

func (r *accountRepo) SaveDevice(ctx context.Context, device *domainAccount.Device) error {
	model := toDeviceModel(device)

	var existing persistencegorm.DeviceModel
	err := r.db.WithContext(ctx).
		Where("device_uid = ?", device.DeviceUID).
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

func (r *accountRepo) SaveSession(ctx context.Context, session *domainAccount.Session) error {
	model, err := toSessionModel(session)
	if err != nil {
		return err
	}

	var existing persistencegorm.SessionModel
	err = r.db.WithContext(ctx).
		Where("session_uid = ?", session.SessionUID).
		Take(&existing).Error
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
			return err
		}
		session.ID = model.ID
		return nil
	case err != nil:
		return err
	default:
		model.ID = existing.ID
		if err := r.db.WithContext(ctx).Save(&model).Error; err != nil {
			return err
		}
		session.ID = model.ID
		return nil
	}
}
