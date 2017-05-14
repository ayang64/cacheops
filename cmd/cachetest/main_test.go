package main

import (
	"cacheops/localmemcache"
	"cacheops/tieredmemcache"
	"github.com/bradfitz/gomemcache/memcache"
	"testing"
)

func TestWork(t *testing.T) {
	tc := localmemcache.New("/tmp/memcached.sock")
	Work(tc)
}

func BenchmarkMemcacheUnix(b *testing.B) {
	mc := memcache.New("/tmp/memcached.sock")
	for n := 0; n < b.N; n++ {
		Work(mc)
	}
}

func BenchmarkMemcacheLocalhost(b *testing.B) {
	mc := memcache.New("127.0.0.1:11211")
	for n := 0; n < b.N; n++ {
		Work(mc)
	}
}

func BenchmarkLocal(b *testing.B) {
	lc := localmemcache.New("/tmp/memcached.sock")
	for n := 0; n < b.N; n++ {
		Work(lc)
	}
}

func BenchmarkTiered(b *testing.B) {
	tc := tieredmemcache.New("/tmp/memcached.sock")
	for n := 0; n < b.N; n++ {
		Work(tc)
	}
}
