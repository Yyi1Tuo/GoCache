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

type GetterFunc func(key string) ([]byte,error)//接口型函数

func (f GetterFunc) Get(key string) ([]byte,error){
	return f(key)
}

//定义一个 Group 结构体,一个 Group 可以认为是一个缓存的命名空间
type Group struct{
	name string
	getter Getter
	mainCache cache
}

var (
	groupMu sync.RWMutex //使用读写锁提高效率
	groups = make(map[string]*Group)
)

// NewGroup 创建一个新的 Group 实例
func NewGroup(name string, cacheBytes int64, getter Getter) *Group{
	if getter == nil {
		panic("nil Getter")
	}
	groupMu.Lock()
	defer groupMu.Unlock()
	g := &Group{
		name: name,
		getter: getter,
		mainCache: cache{cacheBytes: cacheBytes},
	}
	groups[name] = g
	return g
}
func GetGroup(name string) *Group{
	groupMu.RLock()//读锁
	g := groups[name]
	groupMu.RUnlock()
	return g
}

//实现group的get方法
func (g *Group) Get(key string) (ByteView,error){
	if key == "" {
		return ByteView{},fmt.Errorf("key is required")
	}
	//如果缓存命中则返回缓存值
	if v,ok := g.mainCache.get(get);ok{
		DPrintf("[Gocache] hit")
		return v,nil
	}
	//如果缓存未命中，则调用load方法
	return g.load(key)
}

func (g *Group) load(key string) (value ByteView,err error){
	//这里需要分情况讨论从哪里load数据
	//目前先考虑单机模式
	return g.getLocally(key)
}

func (g *Group) getLocally(key string) (ByteView,error){
	bytes,err := g.getter.Get(key)
	if err != nil {
		return ByteView{},err
	}
	value := ByteView{b: bytes}
	g.populateCache(key,value)
	return value,nil
}

//将数据添加到缓存
func (g *Group) populateCache(key string,value ByteView){
	g.mainCache.add(key,value)
}



