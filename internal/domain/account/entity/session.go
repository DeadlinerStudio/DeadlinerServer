package entity

type Session struct {
	ID               int64
	SessionUID       string
	AccountID        int64
	DeviceUID        string
	RefreshTokenHash string
	ExpiresAt        string
}
