package url

// Service is the
type Service interface {
	Init() error

	Create(string, string, string, string, string) error
	GetLong(string) (map[string]string, error)
	Present(string) (bool, error)
}
