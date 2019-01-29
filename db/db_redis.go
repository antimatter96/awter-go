// Package db Contains all methods used by the other functions
package db

import (
	"time"

	"github.com/antimatter96/awter-go/constants"

	redis "github.com/gomodule/redigo/redis"

	cache "github.com/antimatter96/awter-go/db/cache"
	url "github.com/antimatter96/awter-go/db/url"

	// This exposes mysql connector
	_ "github.com/go-sql-driver/mysql"
)

var pool *redis.Pool

// InitRedis is used to initialize the redis connections
func InitRedis() {
	x := constants.Value("redisAddress").(string)
	pool = newPool(x)
}

// NewURLInterfaceRedis returns a URLService interface, using redis as backend
func NewURLInterfaceRedis() url.Service {
	urlService := url.UrlsRedis{Pool: pool}
	err := urlService.Init()
	if err != nil {
		panic(err.Error())
	}
	return &urlService
}

func NewCacheInterfaceRedis() cache.Service {
	cacheService := cache.Redis{Pool: pool}
	err := cacheService.Init()
	if err != nil {
		panic(err.Error())
	}
	return &cacheService
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
