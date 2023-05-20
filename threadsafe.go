package ecs

import (
	"sync"

	"golang.org/x/exp/slices"
)

type Container[T any] struct {
	items []T
	lock  sync.RWMutex
}

func (c *Container[T]) Len() int {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return len(c.items)
}

func (c *Container[T]) Insert(item T) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.items = append(c.items, item)
}

func (c *Container[T]) DeleteAll() {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.items = c.items[:0]
}

func (c *Container[T]) DeleteAt(index int) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.items = slices.Delete(c.items, index, index+1)
}

func (c *Container[T]) InsertAt(index int, item T) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.items = slices.Insert(c.items, index, item)
}

func (c *Container[T]) GetAll(buf []T) []T {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return append(buf, c.items...)
}

func (c *Container[T]) GetAllAndDeleteAll(buf []T) []T {
	c.lock.RLock()
	defer c.DeleteAll()
	defer c.lock.RUnlock()
	return append(buf, c.items...)
}
