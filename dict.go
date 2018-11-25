package ecs

import (
	"sync"
)

type dict struct {
	lock sync.RWMutex
	m    map[string]interface{}
}

func (d *dict) Get(key string) interface{} {
	d.lock.RLock()
	defer d.lock.RUnlock()
	return d.m[key]
}

func (d *dict) Set(key string, val interface{}) {
	d.lock.Lock()
	defer d.lock.Unlock()
	d.m[key] = val
}

func newdict() *dict {
	return &dict{
		m: make(map[string]interface{}),
	}
}
