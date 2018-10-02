//

package cache

import (
	"time"

	"../constants"
	redis "github.com/gomodule/redigo/redis"
)

var pool *redis.Pool

// Init
func Init() {
	x := constants.Value("redisAddress").(string)
	pool = newPool(x)
}

// newPool generates a common pool from which we can access connections
func newPool(addr string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     10,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", addr)
		},
	}
}

//
func GetSessionValue(sessionId, key string) (string, error) {
	conn := pool.Get()
	defer conn.Close()

	_, errPingRedis := conn.Do("PING") //REDIS IS DOWN
	if errPingRedis != nil {
		//fmt.Println(errPingRedis)
		return "-", errPingRedis
	}

	resRedis, errGetRedis := redis.String(conn.Do("HGET", sessionId, key))
	if errGetRedis != nil {
		return "-", errGetRedis
	}

	return resRedis, nil
}

func SetSessionValue(sessionId, key, value string) error {
	conn := pool.Get()
	defer conn.Close()

	_, errPingRedis := conn.Do("PING") //REDIS IS DOWN
	if errPingRedis != nil {
		return errPingRedis
	}

	_, errGetRedis := conn.Do("HSET", sessionId, key, value)
	if errGetRedis != nil {
		return errGetRedis
	}

	return nil
}
