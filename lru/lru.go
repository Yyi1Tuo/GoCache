package lru

import (
	"sync"
	"fmt"
	"container/list"
)

type entry struct {
	key string
	value Value
}

type Value interface{
	Len() int
}

