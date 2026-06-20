package port

type TokenGenerator interface {
	Generate() (string, error)
}
