package ConcurrencyCache

import (
	"sync"
	"GoCache/debug"
	"GoCache/lru"
)
//抽象一个只读的数据结构来表示缓存值
type ByteView struct {
	bytes []byte
}

func (v ByteView) Len() int {
	return len(v.bytes)
}

//由于slice是浅拷贝，为了防止修改原始值，需要复制一段传出去
func (v ByteView) ByteSlice() []byte{
	return cloneBytes(v)
}

func cloneBytes(b []byte) []byte{
	c:=make([]byte,len(b))
	copy(c,b)
	return c
} 
func (v ByteView) String() string{
	return string(v.bytes)
}

