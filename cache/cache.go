//

package cache

import (
	"github.com/antimatter96/awter-go/db"
	"github.com/antimatter96/awter-go/db/cache"
)

var cacheService cache.Service

func Init(store string) {
	switch store {
	case "redis":
		cacheService = db.NewCacheInterfaceRedis()
	case "mysql":
		//cacheService = db.NewCacheInterfaceMySQL()
	}

}

func GetSessionValue(sessionId, key string) (interface{}, error) {
	val, err := cacheService.Get(sessionId, key)
	return val, err
}

func SetSessionValue(sessionId, key, value string) error {
	return cacheService.Set(sessionId, key, value)
}
