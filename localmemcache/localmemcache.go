package localmemcache

import (
	"github.com/bradfitz/gomemcache/memcache"
	"sync"
	"time"
)

/*
 * LocalCache is a type that implements a subset of memcache's
 * functionality and conforms to our Cacher interface.
 */
type LocalCache struct {
	Data  map[string]*memcache.Item
	mutex *sync.RWMutex
}

func (lc *LocalCache) Set(i *memcache.Item) error {
	lc.mutex.Lock()
	defer lc.mutex.Unlock()

	lc.Data[i.Key] = i
	return nil
}

func (lc *LocalCache) Get(s string) (*memcache.Item, error) {
	/* lock our map for reading. */
	lc.mutex.RLock()
	defer lc.mutex.RUnlock()

	if v, ok := lc.Data[s]; ok == true {
		return v, nil
	}
	return nil, memcache.ErrCacheMiss
}

func (lc *LocalCache) GetMulti(keys []string) (map[string]*memcache.Item, error) {
	/* lock our map for reading. */
	lc.mutex.RLock()
	defer lc.mutex.RUnlock()

	rc := make(map[string]*memcache.Item)
	for _, cur := range keys {
		if v, ok := lc.Data[cur]; ok == true {
			rc[cur] = v
		}
	}
	return rc, nil
}

func (lc *LocalCache) AsyncExpire() {

	for {
		time.Sleep(1 * time.Second)
		now := int32(time.Now().Unix())
		lc.mutex.Lock()
		for key, v := range lc.Data {
			if now-v.Expiration < 0 {
				/* entry is stale. it is safe to delete from lc.Data while ranging over it*/
				delete(lc.Data, key)
			}
		}
		lc.mutex.Unlock()
	}
}

func New(s string) *LocalCache {
	/* since this is purely for testing, we ignore the string
	 * passed to this routine as it is mainly for package level
	 * API compatibility.  */
	var rc LocalCache
	rc.Data = make(map[string]*memcache.Item)
	rc.mutex = &sync.RWMutex{}

	go rc.AsyncExpire()

	return &rc
}
