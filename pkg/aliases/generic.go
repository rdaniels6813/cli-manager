package aliases

type Generator interface {
	Generate() string
	Install() error
}

func NewGenerator() Generator {
	return nil
}
