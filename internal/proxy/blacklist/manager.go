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

func (blm *BlackListManager) UpdateCache(key string) {
	cacheValue, _ := blm.Cache.Get(key)
	floatValue := convertCacheValueToFloat64(cacheValue)
	exp := math.Min(floatValue+1, 6)
	newExpirationTime := time.Second * time.Duration(math.Pow(2, exp+1))
	blm.Cache.SetWithExpire(key, exp, newExpirationTime)
	log.Printf("next request from source %s will be blocked for %d seconds", key, newExpirationTime/2e9)
}

func (blm *BlackListManager) isIntervalSufficient(key string) bool {
	var isSufficient = true
	lastAccessTime, present := blm.LastAccess[key]
	interval := time.Now().Sub(lastAccessTime).Seconds()
	if present {
		isSufficient = interval >= blm.MinimalRequestInterval.Seconds()
		if !isSufficient {
			log.Printf("INSUFFICIENT INTERVAL: %f\n", interval)
		}
	}

	return isSufficient
}

func (blm *BlackListManager) ShouldRequestBeBlocked(key string) bool {
	val, _ := blm.Cache.Get(key)

	return !blm.isIntervalSufficient(key) || convertCacheValueToFloat64(val) > 0
}

func (blm *BlackListManager) PurgeCache() {
	blm.Cache.Purge()
}

func (blm *BlackListManager) BlockProcessingTraffic(key string) {
	cacheValue, _ := blm.Cache.Get(key)
	floatValue := convertCacheValueToFloat64(cacheValue)
	time.Sleep(getSufficientSleepTime(floatValue))
}

func getSufficientSleepTime(exp float64) time.Duration {
	finalExp := math.Max(exp, 4)
	return time.Duration(math.Pow(2, finalExp)) * time.Second
}

func convertCacheValueToFloat64(value interface{}) float64 {
	convertedValue, ok := value.(float64)
	if ok {
		return convertedValue
	}

	return 0
}
