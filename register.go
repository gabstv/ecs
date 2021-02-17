package ecs

import (
	"sync"
)

var defaultLock sync.Mutex
var defaultComponents = make(map[string]RegisterComponentFn, 0)
var defaultSystems = make(map[string]RegisterSystemFn, 0)

func RegisterComponent(fn RegisterComponentFn) {
	defaultLock.Lock()
	defer defaultLock.Unlock()
	c := fn()
	if _, ok := defaultComponents[c.UUID()]; ok {
		panic("components with duplicate UUIDs: " + c.UUID())
	}
	if len(defaultComponents) >= MaxFlagCapacity {
		panic("component max flag capacity reached")
	}
	defaultComponents[c.UUID()] = fn
}

func RegisterSystem(fn RegisterSystemFn) {
	defaultLock.Lock()
	defer defaultLock.Unlock()
	c := fn()
	if _, ok := defaultSystems[c.UUID()]; ok {
		panic("systems with duplicate UUIDs: " + c.UUID())
	}
	defaultSystems[c.UUID()] = fn
}

func RegisterWorldDefaults(w World) {
	defaultLock.Lock()
	defer defaultLock.Unlock()
	for uuid, fn := range defaultComponents {
		if w.IsRegistered(uuid) {
			println("WARNING: Component " + uuid + " already registered!")
			continue
		}
		w.RegisterComponent(fn())
	}
	for uuid, fn := range defaultSystems {
		if w.IsRegistered(uuid) {
			println("WARNING: Component " + uuid + " already registered!")
			continue
		}
		_ = w.AddSystem(fn())
	}
}
