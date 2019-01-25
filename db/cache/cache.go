package cache

// CacheService is the
type Service interface {
	//Init() error

	Get(string, string) (interface{}, error)
	Set(string, string, interface{}) error
}
