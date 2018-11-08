// Package db Contains all methods used by the other functions
package db

import (
	"time"

	"../constants"

	redis "github.com/gomodule/redigo/redis"

	// This exposes mysql connector
	_ "github.com/go-sql-driver/mysql"
)

var pool *redis.Pool

func InitRedis() {
	x := constants.Value("redisAddress").(string)
	pool = newPool(x)
}

// The main db object

func NewURLInterfaceRedis() URLService {
	urlService := urlsRedis{pool: pool}
	err := urlService.Init()
	if err != nil {
		panic(err.Error())
	}
	return &urlService
}

func checkStatusRedis() bool {
	conn := pool.Get()
	defer conn.Close()

	_, errPingRedis := conn.Do("PING")
	if errPingRedis != nil {
		return false
	}
	return true
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