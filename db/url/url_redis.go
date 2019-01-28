package url

import (
	redis "github.com/gomodule/redigo/redis"
)

// UrlsRedis is
type UrlsRedis struct {
	Pool *redis.Pool
}

// Init is an empty method,
func (u *UrlsRedis) Init() error {
	return nil
}

// Present checks the presence of given short URL
func (u *UrlsRedis) Present(short string) (bool, error) {
	conn := u.Pool.Get()
	defer conn.Close()

	_, errPingRedis := conn.Do("PING")
	if errPingRedis != nil {
		return false, errPingRedis
	}

	present, errorGetRedis := redis.Int(conn.Do("HLEN", short))
	if errorGetRedis != nil {
		return false, errorGetRedis
	}
	if present == 0 {
		return false, nil
	}
	return true, nil
}

// GetLong returns all the info related to the given short URL
func (u *UrlsRedis) GetLong(short string) (map[string]string, error) {

	conn := u.Pool.Get()
	defer conn.Close()

	_, errPingRedis := conn.Do("PING")
	if errPingRedis != nil {
		return nil, errPingRedis
	}

	resRedis, errGetRedis := redis.StringMap(conn.Do("HGETALL", short))
	if errGetRedis != nil {
		//fmt.Println("REDIS GET error", errGetRedis)
		return nil, errGetRedis
	}

	return resRedis, nil
}

// Create is used to create an entry in the datastore
func (u *UrlsRedis) Create(short, nonce, salt, encrypted, passwordHash string) error {
	conn := u.Pool.Get()
	defer conn.Close()

	_, errPingRedis := conn.Do("PING")
	if errPingRedis != nil {
		return errPingRedis
	}

	conn.Send("MULTI")
	conn.Send("HSET", short, "encrypted", encrypted)
	conn.Send("HSET", short, "salt", salt)
	conn.Send("HSET", short, "nonce", nonce)
	conn.Send("HSET", short, "passwordHash", passwordHash)

	_, errRedis := conn.Do("EXEC")
	if errRedis != nil {
		return errRedis
	}
	return nil
}
