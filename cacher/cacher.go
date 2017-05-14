package cacher

import (
	"github.com/bradfitz/gomemcache/memcache"
)

type Cacher interface {
	Set(*memcache.Item) error
	Get(string) (*memcache.Item, error)
	GetMulti([]string) (map[string]*memcache.Item, error)
}
