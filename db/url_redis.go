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

func (u *urlsRedis) CreateNoPassword(short, long string) error {
	return nil
}

func (u *urlsRedis) CreatePassword(short, long, password string) error {

	conn := u.pool.Get()
	defer conn.Close()

	_, errPingRedis := conn.Do("PING")
	if errPingRedis != nil {
		return errPingRedis
	}

	conn.Send("MULTI")
	conn.Send("HSET", short, "long", long)
	conn.Send("HSET", short, "pwd", password)

	_, errRedis := conn.Do("EXEC")
	if errRedis != nil {
		return errRedis
	}
	return nil
}

func (u *urlsRedis) GetLong(short string) (bool, string, string, error) {

	conn := u.pool.Get()
	defer conn.Close()

	_, errPingRedis := conn.Do("PING")
	if errPingRedis != nil {
		return false, "", "", errPingRedis
	}

	conn.Send("MULTI")
	conn.Send("HGET", short, "long")
	conn.Send("HGET", short, "pwd")

	resRedis, errGetRedis := redis.Strings(conn.Do("EXEC"))
	if errGetRedis != nil {
		//fmt.Println("REDIS GET error", errGetRedis)
		return false, "", "", errGetRedis
	}

	return true, resRedis[0], resRedis[1], nil
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

func (u *urlsRedis) GetLong2(short string) (bool, map[string]string, error) {

	conn := u.pool.Get()
	defer conn.Close()

	_, errPingRedis := conn.Do("PING")
	if errPingRedis != nil {
		return false, nil, errPingRedis
	}

	resRedis, errGetRedis := redis.StringMap(conn.Do("HGETALL" + short))
	if errGetRedis != nil {
		//fmt.Println("REDIS GET error", errGetRedis)
		return false, nil, errGetRedis
	}

	return true, resRedis, nil
}

func (u *urlsRedis) CreateNew(short, nonce, salt, encrypted, passwordHash string) error {
	conn := u.pool.Get()
	defer conn.Close()

	_, errPingRedis := conn.Do("PING")
	if errPingRedis != nil {
		return errPingRedis
	}

	short += "_new"
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
