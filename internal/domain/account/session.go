package account

type Session struct {
	ID               int64
	SessionUID       string
	AccountID        int64
	RefreshTokenHash string
	ExpiresAt        string
}

