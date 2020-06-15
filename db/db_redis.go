// Package db Contains all methods used by the other functions
package db

import (
	"time"

	redis "github.com/gomodule/redigo/redis"

	url "github.com/antimatter96/awter-go/db/url"

	// This exposes mysql connector
	_ "github.com/go-sql-driver/mysql"
)

// InitRedis is used to initialize the redis connections
func InitRedis(redisAddress string) *redis.Pool {
	pool := newPool(redisAddress)
	return pool
}

// NewURLInterfaceRedis returns a URLService interface, using redis as backend
func NewURLInterfaceRedis(pool *redis.Pool) (url.Service, error) {
	conn := pool.Get()
	defer conn.Close()

	_, errPingRedis := conn.Do("PING")
	if errPingRedis != nil {
		return nil, errPingRedis
	}

	urlService := url.UrlsRedis{Pool: pool}
	err := urlService.Init()
	if err != nil {
		pool.Close()
		return nil, err
	}
	return &urlService, nil
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
