package ecs

import "sync"

type Locker struct {
	m sync.RWMutex
	v map[string]interface{}
}

func (l *Locker) Item(name string) interface{} {
	l.m.RLock()
	defer l.m.RUnlock()
	if l.v == nil {
		return nil
	}
	return l.v[name]
}

func (l *Locker) SetItem(name string, value interface{}) {
	l.m.Lock()
	defer l.m.Unlock()
	if l.v == nil {
		l.v = make(map[string]interface{})
	}
	l.v[name] = value
}
