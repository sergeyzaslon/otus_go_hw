package hw04lrucache

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	newCacheItem := cacheItem{key: key, value: value}
	if mapValue, ok := c.items[key]; ok { // елемент есть всписке
		mapValue.Value = newCacheItem
		c.queue.MoveToFront(mapValue)
		c.items[key] = c.queue.Front()
		return true
	}
	// елемента нет в списке
	if c.queue.Len() == c.capacity { // достигли емкости, удаляем последний елемент
		last := c.queue.Back()
		if last != nil {
			c.queue.Remove(last)
			delete(c.items, last.Value.(cacheItem).key)
		}
	}
	li := c.queue.PushFront(newCacheItem)
	c.items[key] = li
	return false
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	if mapValue, ok := c.items[key]; ok {
		c.queue.MoveToFront(mapValue)
		return mapValue.Value.(cacheItem).value, true
	}
	return nil, false
}

func (c *lruCache) Clear() {
	c.queue.Clear()
	c.items = make(map[Key]*ListItem, c.capacity)
}

type cacheItem struct {
	key   Key
	value interface{}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
