package tieredmemcache

import (
	"cacheops/localmemcache"
	"github.com/bradfitz/gomemcache/memcache"
	"time"
)

/*
 * TieredCache is a type that implements a subset of memcache's
 * functionality and conforms to our Cacher interface.
 */
type TieredCache struct {
	local  *localmemcache.LocalCache
	remote *memcache.Client
}

func (tc *TieredCache) Set(i *memcache.Item) error {
	/* update remote cache asyncronously. */
	tc.remote.Set(i)

	tc.local.Set(i)
	return nil
}

func (tc *TieredCache) Get(s string) (*memcache.Item, error) {
	if rc, err := tc.local.Get(s); err == nil {
		/* if our local version is out of date, try fetching a remote version. */
		if now := int32(time.Now().Unix()); now-rc.Expiration > 0 {
			/* local version is 'fresh' -- lets return it. */
			return rc, err
		}
	}

	/* if we're here, either the local cached version is stale or non-existent */
	if rc, err := tc.remote.Get(s); err == nil {
		tc.local.Set(rc)
		return rc, err
	}
	return nil, memcache.ErrCacheMiss
}

func (tc *TieredCache) GetMulti(keys []string) (map[string]*memcache.Item, error) {
	if rc, err := tc.local.GetMulti(keys); err != nil {
		return rc, nil
	} else {
		return nil, err
	}

	if rc, err := tc.remote.GetMulti(keys); err != nil {
		for _, i := range rc {
			tc.local.Set(i)
		}
		return rc, nil
	} else {
		return nil, err
	}

	return nil, nil
}

func New(s string) *TieredCache {
	var rc TieredCache
	rc.remote = memcache.New(s)
	rc.local = localmemcache.New(s)
	return &rc
}
