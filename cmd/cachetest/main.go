package main

import (
	"cacheops/cacher"
	"fmt"
	"github.com/bradfitz/gomemcache/memcache"
	"time"
)

func Work(c cacher.Cacher) {
	/* set some values */
	c.Set(&memcache.Item{Key: "key_one", Value: []byte("michael"), Expiration: int32(time.Now().Unix()) + 15})
	c.Set(&memcache.Item{Key: "key_two", Value: []byte("programming"), Expiration: int32(time.Now().Unix()) + 15})
	c.Set(&memcache.Item{Key: "foo", Value: []byte("bar"), Expiration: int32(time.Now().Unix()) + 15})

	/* get a single value */
	_, err := c.Get("key_one")

	/* fmt.Printf("c.Get(\"key_one\") = %s\n", val.Value) */

	/* get multiple values */
	_, err = c.GetMulti([]string{"key_one", "key_two", "foo"})

	if err != nil {
		/* fmt.Println(err) */
		return
	}

	/* fmt.Printf("Result of c.GetMulti([]string{\"key_one\", ...}):\n") */
	/* it's important to note here that `range` iterates in a
	 * *random* order */
}

func main() {
	fmt.Printf("Real memcache...\n")
	mc := memcache.New("127.0.0.1:11211")
	Work(mc)
}
