package url

// URLService is the
type URLService interface {
	//Init() error

	Create(string, string, string, string, string) error
	GetLong(string) (map[string]string, error)
	PresentShort(string) (bool, error)
}
