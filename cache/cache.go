//

package cache

import (
	"github.com/antimatter96/awter-go/db"
	"github.com/antimatter96/awter-go/db/cache"
)

var cacheService cache.Service

func Init() {
	cacheService = db.NewCacheInterfaceRedis()
}

func GetSessionValue(sessionId, key string) (interface{}, error) {
	val, err := cacheService.Get(sessionId, key)
	return val, err
}

func SetSessionValue(sessionId, key, value string) error {
	return cacheService.Set(sessionId, key, value)
}
