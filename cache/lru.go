package cache

import (
    "container/list"
    "sync"
)

type entry struct {
    key   string
    value string
}

type LRUCache struct {
    mu        sync.Mutex
    cap       int
    items     map[string]*list.Element
    evictList *list.List
}

func NewLRU(capacity int) *LRUCache {
    return &LRUCache{
        cap:       capacity,
        items:     make(map[string]*list.Element),
        evictList: list.New(),
    }
}

func (c *LRUCache) Get(key string) (string, bool) {
    c.mu.Lock()
    defer c.mu.Unlock()

    if el, found := c.items[key]; found {
        c.evictList.MoveToFront(el)
        return el.Value.(*entry).value, true
    }
    return "", false
}

func (c *LRUCache) Set(key, value string) {
    c.mu.Lock()
    defer c.mu.Unlock()

    if el, found := c.items[key]; found {
        el.Value.(*entry).value = value
        c.evictList.MoveToFront(el)
        return
    }

    el := c.evictList.PushFront(&entry{key, value})
    c.items[key] = el

    if c.evictList.Len() > c.cap {
        c.removeOldest()
    }
}

func (c *LRUCache) Clear() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.items = make(map[string]*list.Element)
    c.evictList.Init()
}

func (c *LRUCache) removeOldest() {
    el := c.evictList.Back()
    if el != nil {
        c.evictList.Remove(el)
        kv := el.Value.(*entry)
        delete(c.items, kv.key)
    }
}
