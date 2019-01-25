package cache

import (
	redis "github.com/gomodule/redigo/redis"
)

type Redis struct {
	Pool *redis.Pool
}

// Init creates all the prepared statements
func (u *Redis) Init() error {
	return nil
}

func (u *Redis) Get(key, field string) (interface{}, error) {
	conn := u.Pool.Get()
	defer conn.Close()

	_, errPingRedis := conn.Do("PING")
	if errPingRedis != nil {
		return nil, errPingRedis
	}

	value, errorGetRedis := redis.Int(conn.Do("HGET", key, field))
	if errorGetRedis != nil {
		return nil, errorGetRedis
	}
	return value, nil
}

func (u *Redis) Set(key, field string, values interface{}) error {

	conn := u.Pool.Get()
	defer conn.Close()

	_, errPingRedis := conn.Do("PING")
	if errPingRedis != nil {
		return errPingRedis
	}

	_, errGetRedis := conn.Do("HSET", key, field, values)
	if errGetRedis != nil {
		//fmt.Println("REDIS GET error", errGetRedis)
		return errGetRedis
	}

	return nil
}
