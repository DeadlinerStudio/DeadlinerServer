package entity

type Account struct {
	ID           int64
	AccountUID   string
	Email        string
	PasswordHash string
	DisplayName  string
}
