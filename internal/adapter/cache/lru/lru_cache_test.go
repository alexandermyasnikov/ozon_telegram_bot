package lrucache //nolint:testpackage

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func assertLRU(t *testing.T, lru *LRUCache, itemsExp []LRUItem) {
	t.Helper()

	size := len(itemsExp)

	assert.Equal(t, size, lru.evictList.Len())
	assert.Equal(t, size, len(lru.items))

	for _, itemExp := range itemsExp {
		itemElemAct, ok := lru.items[itemExp.key]
		assert.Equal(t, true, ok)

		itemAct, ok := itemElemAct.Value.(*LRUItem)
		assert.Equal(t, true, ok)
		assert.Equal(t, itemExp, *itemAct)
	}

	for e, i := lru.evictList.Front(), 0; e != nil; e = e.Next() {
		itemAct, ok := e.Value.(*LRUItem)
		assert.Equal(t, true, ok)
		assert.Equal(t, itemsExp[i], *itemAct)
		i++
	}
}

func Date(offset int) time.Time {
	date, _ := time.Parse(time.RFC1123, "Mon, 30 Oct 2022 16:01:10 MSK")

	return date.Add(time.Duration(offset) * time.Second)
}

func TestLRU_Init(t *testing.T) {
	t.Parallel()

	lru := NewLRUCache(&sync.RWMutex{}, 5)

	assert.NotNil(t, lru.mu)
	assert.NotNil(t, lru.evictList)
	assert.NotNil(t, lru.items)
	assert.Equal(t, 0, lru.evictList.Len())
	assert.Equal(t, 0, len(lru.items))
	assert.Equal(t, 5, lru.capacity)
}

func TestLRU_Get_nonExistentElement_none(t *testing.T) {
	t.Parallel()

	lru := NewLRUCache(&sync.RWMutex{}, 5)

	val, ok := lru.Get(Date(0), "key1")

	assert.Equal(t, false, ok)
	assert.Equal(t, nil, val)
}

func TestLRU_Delete_nonExistentElement_none(t *testing.T) {
	t.Parallel()

	lru := NewLRUCache(&sync.RWMutex{}, 5)

	ok := lru.Delete(Date(0), "key1")

	assert.Equal(t, false, ok)
}

func TestLRU_Add_nonExistentElement_addElement(t *testing.T) {
	t.Parallel()

	lru := NewLRUCache(&sync.RWMutex{}, 5)

	lru.Add(Date(1), "key1", 1001, 101)

	assertLRU(t, lru, []LRUItem{
		{key: "key1", val: 1001, ts: Date(1), ttl: 101},
	})
}

func TestLRU_Add_manyNonExistentElements_addElements(t *testing.T) {
	t.Parallel()

	lru := NewLRUCache(&sync.RWMutex{}, 5)

	lru.Add(Date(1), "key1", 1001, 101)
	lru.Add(Date(2), "key2", 1002, 102)
	lru.Add(Date(3), "key3", 1003, 103)

	assertLRU(t, lru, []LRUItem{
		{key: "key1", val: 1001, ts: Date(1), ttl: 101},
		{key: "key2", val: 1002, ts: Date(2), ttl: 102},
		{key: "key3", val: 1003, ts: Date(3), ttl: 103},
	})
}

func TestLRU_Add_existentElements_overwrite(t *testing.T) {
	t.Parallel()

	lru := NewLRUCache(&sync.RWMutex{}, 5)
	lru.Add(Date(1), "key1", 1001, 101)
	lru.Add(Date(2), "key2", 1002, 102)
	lru.Add(Date(3), "key3", 1003, 103)

	lru.Add(Date(4), "key3", 2003, 203)
	lru.Add(Date(5), "key1", 2001, 201)

	assertLRU(t, lru, []LRUItem{
		{key: "key2", val: 1002, ts: Date(2), ttl: 102},
		{key: "key3", val: 2003, ts: Date(4), ttl: 203},
		{key: "key1", val: 2001, ts: Date(5), ttl: 201},
	})
}

func TestLRU_Add_nonExistentElements_deleteOldElement(t *testing.T) {
	t.Parallel()

	lru := NewLRUCache(&sync.RWMutex{}, 5)
	lru.Add(Date(1), "key1", 1001, 101)
	lru.Add(Date(2), "key2", 1002, 102)
	lru.Add(Date(3), "key3", 1003, 103)
	lru.Add(Date(4), "key4", 1004, 104)
	lru.Add(Date(5), "key5", 1005, 105)

	lru.Add(Date(6), "key6", 1006, 106)

	assertLRU(t, lru,
		[]LRUItem{
			{key: "key2", val: 1002, ts: Date(2), ttl: 102},
			{key: "key3", val: 1003, ts: Date(3), ttl: 103},
			{key: "key4", val: 1004, ts: Date(4), ttl: 104},
			{key: "key5", val: 1005, ts: Date(5), ttl: 105},
			{key: "key6", val: 1006, ts: Date(6), ttl: 106},
		},
	)
}

func TestLRU_Size_nonExistentElements_addSize(t *testing.T) {
	t.Parallel()

	lru := NewLRUCache(&sync.RWMutex{}, 5)
	lru.Add(Date(1), "key1", 1001, 101)

	lru.Add(Date(2), "key2", 1002, 102)

	assert.Equal(t, 2, lru.Size())
}

func TestLRU_Size_existentElements_sameSize(t *testing.T) {
	t.Parallel()

	lru := NewLRUCache(&sync.RWMutex{}, 5)
	lru.Add(Date(1), "key1", 1001, 101)

	lru.Add(Date(2), "key1", 1002, 102)

	assert.Equal(t, 1, lru.Size())
}

func TestLRU_Size_nonExistentElement_sizeEqCap(t *testing.T) {
	t.Parallel()

	lru := NewLRUCache(&sync.RWMutex{}, 5)
	lru.Add(Date(1), "key1", 1001, 101)
	lru.Add(Date(2), "key2", 1002, 102)
	lru.Add(Date(3), "key3", 1003, 103)
	lru.Add(Date(4), "key4", 1004, 104)
	lru.Add(Date(5), "key5", 1005, 105)

	lru.Add(Date(6), "key6", 1006, 106)

	assert.Equal(t, 5, lru.Size())
}

func TestLRU_Get_existentElement_updateTime(t *testing.T) {
	t.Parallel()

	lru := NewLRUCache(&sync.RWMutex{}, 5)
	lru.Add(Date(1), "key1", 1001, 101)

	val, ok := lru.Get(Date(2), "key1")

	assert.Equal(t, true, ok)
	assert.Equal(t, 1001, val)
	assertLRU(t, lru, []LRUItem{
		{key: "key1", val: 1001, ts: Date(2), ttl: 101},
	})
}

func TestLRU_Get_removedElement_None(t *testing.T) {
	t.Parallel()

	lru := NewLRUCache(&sync.RWMutex{}, 5)
	lru.Add(Date(1), "key1", 1001, 101)
	lru.Add(Date(2), "key2", 1002, 102)
	lru.Add(Date(3), "key3", 1003, 103)
	lru.Add(Date(4), "key4", 1004, 104)
	lru.Add(Date(5), "key5", 1005, 105)
	lru.Add(Date(6), "key6", 1006, 106)

	val, ok := lru.Get(Date(7), "key1")

	assert.Equal(t, false, ok)
	assert.Equal(t, nil, val)
}

func TestLRU_Delete_existentElement_deleteElement(t *testing.T) {
	t.Parallel()

	lru := NewLRUCache(&sync.RWMutex{}, 5)
	lru.Add(Date(1), "key1", 1001, 101)
	lru.Add(Date(2), "key2", 1002, 102)
	lru.Add(Date(3), "key3", 1003, 103)

	lru.Delete(Date(4), "key2")

	assertLRU(t, lru, []LRUItem{
		{key: "key1", val: 1001, ts: Date(1), ttl: 101},
		{key: "key3", val: 1003, ts: Date(3), ttl: 103},
	})
}

func TestLRU_Delete_existentElements_deleteAllElements(t *testing.T) {
	t.Parallel()

	lru := NewLRUCache(&sync.RWMutex{}, 5)
	lru.Add(Date(1), "key1", 1001, 101)
	lru.Add(Date(2), "key2", 1002, 102)
	lru.Add(Date(3), "key3", 1003, 103)

	lru.Delete(Date(4), "key1")
	lru.Delete(Date(5), "key2")
	lru.Delete(Date(6), "key3")

	assertLRU(t, lru, []LRUItem{})
}

func TestLRU_Get_nonExistentElements_dontDeleteByTTL(t *testing.T) {
	t.Parallel()

	lru := NewLRUCache(&sync.RWMutex{}, 5)
	lru.Add(Date(1), "key1", 1001, 101)

	lru.Get(Date(102), "key2")

	assertLRU(t, lru, []LRUItem{
		{key: "key1", val: 1001, ts: Date(1), ttl: 101},
	})
}

func TestLRU_Get_nonExistentElements_deleteByTTL(t *testing.T) {
	t.Parallel()

	lru := NewLRUCache(&sync.RWMutex{}, 5)
	lru.Add(Date(1), "key1", 1001, 101)

	lru.Get(Date(103), "key2")

	assertLRU(t, lru, []LRUItem{})
}

func TestLRU_Get_nonExistentElements_deleteAllByTTL(t *testing.T) {
	t.Parallel()

	lru := NewLRUCache(&sync.RWMutex{}, 5)
	lru.Add(Date(1), "key1", 1001, 101)
	lru.Add(Date(2), "key2", 1002, 102)
	lru.Add(Date(3), "key3", 1003, 103)
	lru.Add(Date(4), "key4", 1004, 104)

	lru.Get(Date(200), "key10")

	assertLRU(t, lru, []LRUItem{})
}

func TestLRU_mutex_locked(t *testing.T) {
	t.Parallel()

	mu := &sync.RWMutex{}
	lru := NewLRUCache(mu, 5)

	end := make(chan bool)

	go func() {
		time.Sleep(100 * time.Millisecond)
		end <- true
	}()

	mu.RLock()

	go func() {
		lru.Add(Date(1), "key1", 1001, 101)
		assert.FailNow(t, "Add executed")
	}()

	go func() {
		lru.Get(Date(2), "key1")
		assert.FailNow(t, "Get executed")
	}()

	go func() {
		lru.Delete(Date(3), "key1")
		assert.FailNow(t, "Delete executed")
	}()

	<-end
}
