package accountrepo

import (
	"time"

	domainAccount "github.com/aritxonly/deadlinerserver/internal/domain/account"
	persistencegorm "github.com/aritxonly/deadlinerserver/internal/infra/persistence/gorm"
)

func toAccount(model persistencegorm.AccountModel) *domainAccount.Account {
	return &domainAccount.Account{
		ID:           model.ID,
		AccountUID:   model.AccountUID,
		Email:        model.Email,
		PasswordHash: model.PasswordHash,
		DisplayName:  model.DisplayName,
	}
}

func toAccountModel(acc *domainAccount.Account) persistencegorm.AccountModel {
	return persistencegorm.AccountModel{
		ID:           acc.ID,
		AccountUID:   acc.AccountUID,
		Email:        acc.Email,
		PasswordHash: acc.PasswordHash,
		DisplayName:  acc.DisplayName,
	}
}

func toDeviceModel(device *domainAccount.Device) persistencegorm.DeviceModel {
	return persistencegorm.DeviceModel{
		ID:         device.ID,
		DeviceUID:  device.DeviceUID,
		AccountID:  device.AccountID,
		Platform:   device.Platform,
		DeviceName: device.DeviceName,
	}
}

func toSessionModel(session *domainAccount.Session) (persistencegorm.SessionModel, error) {
	expiresAt, err := time.Parse(time.RFC3339, session.ExpiresAt)
	if err != nil {
		return persistencegorm.SessionModel{}, err
	}

	return persistencegorm.SessionModel{
		ID:               session.ID,
		SessionUID:       session.SessionUID,
		AccountID:        session.AccountID,
		DeviceUID:        session.DeviceUID,
		RefreshTokenHash: session.RefreshTokenHash,
		ExpiresAt:        expiresAt.UTC(),
	}, nil
}

func toSession(model persistencegorm.SessionModel) *domainAccount.Session {
	return &domainAccount.Session{
		ID:               model.ID,
		SessionUID:       model.SessionUID,
		AccountID:        model.AccountID,
		DeviceUID:        model.DeviceUID,
		RefreshTokenHash: model.RefreshTokenHash,
		ExpiresAt:        model.ExpiresAt.UTC().Format(time.RFC3339),
	}
}
