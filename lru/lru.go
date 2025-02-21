//无并发控制的LRU缓存
package lru

import (
	"container/list"
	"GoCache/debug"
)

//使用双向链表更适合实现LRU算法
type Cache struct {
	mu sync.Mutex
	cache map[string]*list.Element
	ll *list.List
	maxBytes int64
	nbytes int64

	OnEvicted func(key string, value Value)//Optional 某条记录被移除时的回调函数，可以为nil
}

type entry struct {
	key string
	value Value
}

type Value interface{
	Len() int
}

func New(maxBytes int64, onEvicted func(key string, value Value)) *Cache {
	return &Cache{
		maxBytes: maxBytes,
		ll: list.New(),
		cache: make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

func (c *Cache) Get(key string) (Value, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if ele, ok := c.cache[key]; ok {//如果key存在，则将该条记录移动到队首
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)//类型断言
		debug.Dprintf("Get key=%s, value=%v", key, kv.value)
		return kv.value, true
	}
	return nil, false
}

func (c *Cache) Add(key string, value Value) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.nbytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
		debug.Dprintf("Add key=%s, value=%v", key, value)
	} else {
		ele := c.ll.PushFront(&entry{key, value})
		c.cache[key] = ele
		c.nbytes += int64(len(key)) + int64(value.Len())
		debug.Dprintf("Add key=%s, value=%v", key, value)
	}
	//如果达到了最大的内存限制，则移除最老的记录
	for c.maxBytes != 0 && c.maxBytes < c.nbytes {
		c.RemoveOldest()
	}
}

//缓存删除
func (c *Cache) RemoveOldest() {

	ele := c.ll.Back()
	if ele != nil {
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)
		debug.Dprintf("RemoveOldest key=%s, value=%v", kv.key, kv.value)
		c.nbytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

// Len the number of cache entries
func (c *Cache) Len() int {
	return c.ll.Len()
}