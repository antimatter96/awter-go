package db

import (
	redis "github.com/gomodule/redigo/redis"
)

// import (
// 	"database/sql"

// 	// M
// 	_ "github.com/go-sql-driver/mysql"
// )

type urlsRedis struct {
	pool *redis.Pool
}

// Init creates all the prepared statements
func (u *urlsRedis) Init() error {
	return nil
}

func (u *urlsRedis) PresentShort(short string) (bool, error) {
	conn := u.pool.Get()
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

func (u *urlsRedis) GetLong(short string) (map[string]string, error) {

	conn := u.pool.Get()
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

func (u *urlsRedis) Create(short, nonce, salt, encrypted, passwordHash string) error {
	conn := u.pool.Get()
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
