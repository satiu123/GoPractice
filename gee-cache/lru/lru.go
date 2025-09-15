package lru

import "container/list"

type Cache struct {
	maxBytes int64      //允许使用的最大内存
	nbytes   int64      //当前已使用内存
	ll       *list.List //双向链表
	cache    map[string]*list.Element

	OnEvicted func(key string, value Value)
}

type entry struct {
	key   string
	value Value
}

func (e *entry) Len() int64 {
	return int64(len(e.key) + e.value.Len())
}

type Value interface {
	Len() int
}

func NewCache(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     map[string]*list.Element{},
		OnEvicted: onEvicted,
	}
}
func (c *Cache) Len() int {
	return c.ll.Len()
}
func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv, ok := ele.Value.(*entry)
		return kv.value, ok
	}
	return nil, false
}

func (c *Cache) RemoveOldest() {
	ele := c.ll.Back()
	if ele != nil {
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)
		c.nbytes -= int64(kv.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

func (c *Cache) Add(key string, value Value) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.nbytes += kv.Len()
		kv.value = value
	} else {
		kv := &entry{key, value}
		ele = c.ll.PushFront(kv)
		c.cache[key] = ele
		c.nbytes += kv.Len()
	}
	if c.maxBytes != 0 && c.nbytes > c.maxBytes {
		c.RemoveOldest()
	}
}
