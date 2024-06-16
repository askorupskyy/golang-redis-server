package cache

import (
	"fmt"
	"log"
	"sync"
	"time"
)

type CacheType struct {
	sync.Map
}

type SetCacheValueArgs struct {
	Expiry int64
}

type CacheValue struct {
	Value  any
	Expiry time.Time
}

var Cache CacheType

func (c *CacheType) Get(key string) (CacheValue, bool) {
	val, exists := c.Map.Load(key)
	if !exists {
		return CacheValue{
			Value:  nil,
			Expiry: time.Time{},
		}, false
	}
	return val.(CacheValue), exists
}

func (c *CacheType) Set(key string, val any, args SetCacheValueArgs) {
	c.Map.Store(key, CacheValue{
		Value:  val,
		Expiry: time.Now().Add(time.Millisecond * time.Duration(args.Expiry)),
	})

	if args.Expiry > 0 {
		go func() {
			time.Sleep(time.Millisecond * time.Duration(args.Expiry))
			deleted := c.Del(key)
			log.Println(deleted)
		}()
	}
}

func (c *CacheType) Del(key string) bool {
	_, exists := c.Get(key)
	if !exists {
		return false
	}

	c.Map.Delete(key)
	return true
}

func (c *CacheType) Flush() {
	c.Map.Range(func(key, value interface{}) bool {
		return c.Del(fmt.Sprintf(key.(string)))
	})
}
