package blacklist

import (
	"github.com/bluele/gcache"
	"log"
	"math"
	"time"
)

type BlackListManager struct {
	Cache                  gcache.Cache
	LastAccess             map[string]time.Time
	MinimalRequestInterval time.Duration
}

func NewBlackListManager(interval time.Duration) *BlackListManager {
	return &BlackListManager{
		Cache:                  newBlackListCache(),
		LastAccess:             make(map[string]time.Time),
		MinimalRequestInterval: interval,
	}
}

func newBlackListCache() gcache.Cache {
	return gcache.New(100).LFU().LoaderFunc(func(key interface{}) (interface{}, error) {
		return 0, nil
	}).EvictedFunc(func(key, val interface{}) {
		strKey, ok := key.(string)
		if ok {
			log.Printf("evicted entry from cache - %s can connect again without delay", strKey)
		}

	}).Build()
}

func (blm *BlackListManager) UpdateLastAccess(key string) {
	blm.LastAccess[key] = time.Now()
}

func (blm *BlackListManager) UpdateCache(key string, value interface{}) {
	intValue, ok := value.(int)
	if !ok {
		log.Printf("failed to load cache value under '%s' key - value is not an integer", key)
		return
	}
	exp := math.Min(float64(intValue)+1, 6)
	newBlockTime := time.Second * time.Duration(math.Pow(2, exp+1))
	blm.Cache.SetWithExpire(key, exp, newBlockTime)
	log.Printf("next request from source %s will be blocked for %d seconds", key, newBlockTime/1e9)
}

func (blm *BlackListManager) isIntervalSufficient(key string) bool {
	lastAccessTime, present := blm.LastAccess[key]
	if present {
		return time.Now().Sub(lastAccessTime).Seconds() < blm.MinimalRequestInterval.Seconds()
	}
	return true
}

func (blm *BlackListManager) ShouldRequestBeBlocked(key string) bool {
	val, _ := blm.Cache.Get(key)
	return !blm.isIntervalSufficient(key) && val.(int) > 0
}

func (blm *BlackListManager) BlockRequestFromSource(key string) {
	cacheValue, _ := blm.Cache.Get(key)
	intValue := convertCacheValueToInt(cacheValue)
	if intValue > 0 {
		timeOfSleeping := math.Pow(2, float64(intValue))
		time.Sleep(time.Duration(timeOfSleeping) * time.Second)
	}
}

func convertCacheValueToInt(value interface{}) int {
	intValue, ok := value.(int)
	if ok {
		return intValue
	}

	return 0
}
