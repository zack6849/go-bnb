package query_cache

import (
	"fmt"
	"reflect"
	"sync"

	"golang.org/x/sync/singleflight"
)

type CacheStore struct {
	mu    sync.RWMutex
	wg    singleflight.Group
	cache map[string]any
}

func (cacheStore *CacheStore) GetItemByFieldValue[T any](
	field string,
	value any,
	resolve func() (T, error),
	cache bool,
) (T, error) {
	var zero T
	typeKey := reflect.TypeFor[T]().String()
	cacheKey := fmt.Sprintf("%s:%s:%s", typeKey, field, value)
	val, err, _ := cacheStore.wg.Do(cacheKey, func() (any, error) {
		if cache {
			cacheStore.mu.RLock()
			v, hit := cacheStore.cache[cacheKey]
			//always release the read lock, whether hit or miss
			cacheStore.mu.RUnlock()
			if hit {
				typed, ok := v.(T)
				if !ok {
					//cache hit but returned the wrong type
					return nil, fmt.Errorf("mismatched cache type for key %+v", cacheKey)
				}
				//cache hit, same type, bingo!
				return typed, nil
			}
		}
		//resolve the new object to put into the cache
		obj, err := resolve()
		if err != nil {
			//if there's an error, bail, don't cache it.
			return nil, err
		}
		if cache {
			//write lock the cache, we're about to add to it
			cacheStore.mu.Lock()
			//ensure we unlock it
			defer cacheStore.mu.Unlock()
			cacheStore.cache[cacheKey] = obj
		}
		return obj, nil
	})

	if err != nil {
		return zero, err
	}
	return val.(T), nil
}

func NewCache() *CacheStore {
	return &CacheStore{
		cache: make(map[string]any),
	}
}
