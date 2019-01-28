package url

// Service is the
type Service interface {
	Init() error

	Create(ShortURL) error
	GetLong(string) (*ShortURL, error)
	Present(string) (bool, error)
}

// ShortURL is the basic structure of datastore entry
type ShortURL struct {
	Short         string
	EncryptedLong string
	Salt          string
	Nonce         string
	PasswordHash  string
}
