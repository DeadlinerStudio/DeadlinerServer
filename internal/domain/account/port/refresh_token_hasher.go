package port

type RefreshTokenHasher interface {
	Hash(token string) string
}
