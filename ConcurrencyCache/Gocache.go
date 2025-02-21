//单机并发缓存

package	ConcurrencyCache

import (
	"sync"
	"GoCache/debug"
	"GoCache/lru"
)

//实例化lru，封装get add方法并加上互斥锁mu
type cache struct{
	mu sync.Mutex
	lru *lru.Cache
	cacheBytes int64
}

func (c *cache) add(key string,value Byteview){
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil{
		c.lru:=lru.New(c.cacheBytes,nil)
		//延迟初始化(Lazy Initialization)，
		//一个对象的延迟初始化意味着该对象的创建将会延迟至第一次使用该对象时。
		//主要用于提高性能，并减少程序内存要求。
	}
	c.lru.Add(key,value)
}

func (c *cache) get(key string) (value ByteView, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		return
	}

	if v, ok := c.lru.Get(key); ok {
		return v.(ByteView), ok
	}

	return
}