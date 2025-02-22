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

var once sync.Once 

func (c *cache) add(key string,value Byteview){
	c.mu.Lock()
	defer c.mu.Unlock()
	
	once.Do(func(){
		c.lru = lru.New(c.cacheBytes,nil)
	})
	//延迟初始化(Lazy Initialization)，
	//一个对象的延迟初始化意味着该对象的创建将会延迟至第一次使用该对象时。
	//主要用于提高性能，并减少程序内存要求。
	c.lru.Add(key,value)
}

func (c *cache) get(key string) (value ByteView, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		//这里需要添加一个回调函数，以从数据库中读取数据
		return
	}

	if v, ok := c.lru.Get(key); ok {
		return v.(ByteView), ok
	}

	return
}

//由于数据源种类可能很多，定义一个接口用来读取数据源到缓存
type Getter interface{
	Get(key string) ([]byte,error)
}
