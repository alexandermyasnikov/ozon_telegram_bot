package lrucache

import (
	"container/list"
	"sync"
	"time"
)

type LRUItem struct {
	key string
	val interface{}
	ts  time.Time
	ttl int
}

type LRUCache struct {
	mu        *sync.RWMutex
	evictList *list.List
	items     map[string]*list.Element
	capacity  int
}

func NewLRUCache(mu *sync.RWMutex, capacity int) *LRUCache {
	return &LRUCache{
		mu:        mu,
		evictList: list.New(),
		items:     make(map[string]*list.Element, capacity+1),
		capacity:  capacity,
	}
}

func (lru *LRUCache) Add(now time.Time, key string, val interface{}, ttl int) {
	lru.mu.Lock()
	defer lru.mu.Unlock()

	itemElem, ok := lru.items[key]
	if !ok {
		item := &LRUItem{
			key: key,
			val: val,
			ts:  now,
			ttl: ttl,
		}

		lru.items[key] = lru.evictList.PushBack(item)
	} else {
		item, _ := itemElem.Value.(*LRUItem)

		item.val = val
		item.ts = now
		item.ttl = ttl

		lru.evictList.MoveToBack(itemElem)
	}

	lru.evict(now)
}

func (lru *LRUCache) Get(now time.Time, key string) (interface{}, bool) {
	lru.mu.Lock()
	defer lru.mu.Unlock()

	lru.evict(now)

	itemElem, ok := lru.items[key]
	if !ok {
		return nil, false
	}

	item, _ := itemElem.Value.(*LRUItem)

	item.ts = now

	lru.evictList.MoveToBack(itemElem)

	return item.val, true
}

func (lru *LRUCache) Delete(now time.Time, key string) bool {
	lru.mu.Lock()
	defer lru.mu.Unlock()

	lru.evict(now)

	itemElem, ok := lru.items[key]
	if !ok {
		return false
	}

	lru.evictList.Remove(itemElem)
	delete(lru.items, key)

	return true
}

func (lru *LRUCache) Size() int {
	lru.mu.RLock()
	defer lru.mu.RUnlock()

	return len(lru.items)
}

func (lru *LRUCache) evict(now time.Time) {
	for lru.evictList.Len() > 0 {
		itemElem := lru.evictList.Front()
		item, _ := itemElem.Value.(*LRUItem)

		tsExpired := item.ts.Add(time.Duration(item.ttl) * time.Second)
		if lru.evictList.Len() <= lru.capacity && tsExpired.Unix() >= now.Unix() {
			break
		}

		lru.evictList.Remove(itemElem)
		delete(lru.items, item.key)
	}
}
