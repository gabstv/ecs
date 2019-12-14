package ecs

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// World is the "world" that entities, components and systems live and interact.
type World struct {
	// lock for nextEntity, entities, entityIndexMap
	lock           sync.RWMutex
	nextEntity     uint64
	nextComponent  uint64
	nextView       uint64
	entities       map[Entity]flag
	components     map[flag]*Component
	componentNames map[string]*Component
	systems        []*System
	systemNames    map[string]*System
	views          []*View
	globals        *dict
	ctxbuilder     ContextBuilderFn
}

// NewWorld creates a world and initializes the internal storage (necessary).
func NewWorld() *World {
	w := &World{}
	w.entities = make(map[Entity]flag)
	w.components = make(map[flag]*Component)
	w.componentNames = make(map[string]*Component)
	w.views = make([]*View, 0)
	w.systems = make([]*System, 0)
	w.systemNames = make(map[string]*System)
	w.globals = newdict()
	w.ctxbuilder = DefaultContextBuilder
	return w
}

// NewWorldWithCtx creates a world and sets the context buolder
func NewWorldWithCtx(b ContextBuilderFn) *World {
	w := NewWorld()
	w.ctxbuilder = b
	return w
}

// NewEntity creates and entity and adds it to the current world.
func (w *World) NewEntity() Entity {
	nextid := atomic.AddUint64(&w.nextEntity, 1)
	entity := Entity(nextid - 1)
	w.lock.Lock()
	defer w.lock.Unlock()
	w.entities[entity] = newflag(0, 0, 0, 0)
	return entity
}

// NewEntities creates a "n" amount of entities in the current world.
func (w *World) NewEntities(n int) []Entity {
	w.lock.Lock()
	defer w.lock.Unlock()
	nn := w.nextEntity
	entities := make([]Entity, n)
	for k := range entities {
		nn++
		e := Entity(nn)
		w.entities[e] = newflag(0, 0, 0, 0)
		entities[k] = e
	}
	w.nextEntity = nn
	return entities
}

// ContainsEntity tests if an entity is present in the world.
func (w *World) ContainsEntity(entity Entity) bool {
	w.lock.RLock()
	defer w.lock.RUnlock()
	_, ok := w.entities[entity]
	return ok
}

// AddComponentToEntity adds a component to an entity.
//
// All the existing views are updated if the entity matches the requirements.
func (w *World) AddComponentToEntity(entity Entity, component *Component, data interface{}) error {
	w.lock.RLock()
	if _, ok := w.entities[entity]; !ok {
		w.lock.RUnlock()
		return fmt.Errorf("world does not contain entity")
	}
	if _, ok := w.components[component.flag]; !ok {
		w.lock.RUnlock()
		return fmt.Errorf("world does not contain component")
	}
	w.lock.RUnlock()
	// entity and component are valid
	component.lock.RLock()
	if component.validatedata != nil {
		if !component.validatedata(data) {
			component.lock.RUnlock()
			return fmt.Errorf("invalid component data [ValidateComponentData]")
		}
	}
	component.lock.RUnlock()
	component.lock.Lock()
	component.data[entity] = data
	w.lock.Lock()
	w.entities[entity] = w.entities[entity].or(component.flag)
	w.lock.Unlock()
	component.lock.Unlock()
	//
	// get fingerprint of newly modified instance
	clist := make([]*Component, 0)
	w.lock.RLock()
	for _, v := range w.components {
		v.lock.RLock()
		if _, ok := v.data[entity]; ok {
			clist = append(clist, v)
		}
		v.lock.RUnlock()
	}
	fingerprint := getfingerprint(clist...)
	for _, view := range w.views {
		if containsfingerprint(fingerprint, view.fingerprint) {
			view.upsert(entity)
		} else {
			index := -1
			view.lock.RLock()
			if i2, ok := view.matchmap[entity]; ok {
				index = i2
			}
			view.lock.RUnlock()
			if index != -1 {
				view.remove(entity)
			}
		}
	}
	w.lock.RUnlock()
	//
	return nil
}

// RemoveComponentFromEntity removes a component from an entity.
//
// All the existing views are updated if the entity matches the requirements.
func (w *World) RemoveComponentFromEntity(entity Entity, component *Component) error {
	w.lock.RLock()
	if _, ok := w.entities[entity]; !ok {
		w.lock.RUnlock()
		return fmt.Errorf("world does not contain entity")
	}
	if _, ok := w.components[component.flag]; !ok {
		w.lock.RUnlock()
		return fmt.Errorf("world does not contain component")
	}
	w.lock.RUnlock()
	// entity and component are valid
	var destructor ComponentDestructor
	component.lock.RLock()
	if component.destructor != nil {
		destructor = component.destructor
	}
	ldata := component.data[entity]
	compflag := component.flag
	component.lock.RUnlock()
	// remove from views
	w.lock.RLock()
	for _, view := range w.views {
		view.lock.RLock()
		ookk := view.fingerprint.contains(compflag)
		view.lock.RUnlock()
		if ookk {
			view.remove(entity)
		}
	}
	// remove from data
	w.lock.RUnlock()
	w.lock.Lock()
	w.entities[entity] = w.entities[entity].xor(compflag)
	w.lock.Unlock()
	component.lock.Lock()
	delete(component.data, entity)
	component.lock.Unlock()
	// call destructor if exists
	if destructor != nil {
		destructor(w, entity, ldata)
	}

	return nil
}

// Query will return all entities (and components) that contain the
// combination of entities.
func (w *World) Query(components ...*Component) []QueryMatch {
	flag0 := newflag(0, 0, 0, 0)
	for _, comp := range components {
		flag0 = flag0.or(comp.flag)
	}
	results := make([]QueryMatch, 0)
	w.lock.RLock()
	defer w.lock.RUnlock()
	for entity, eflag := range w.entities {
		if eflag.contains(flag0) {
			mmap := make(map[*Component]interface{})
			for _, comp := range components {
				comp.lock.RLock()
				mmap[comp] = comp.data[entity]
				comp.lock.RUnlock()
			}
			match := QueryMatch{
				Entity:     entity,
				Components: mmap,
			}
			results = append(results, match)
		}
	}
	return results
}

// Run will iterate through all systems loop function (sorted by priority).
func (w *World) Run(delta float64) (taken time.Duration) {
	t0 := time.Now()
	w.lock.RLock()
	allsystems := w.systems
	w.lock.RUnlock()
	rctx := context.Background()
	for _, system := range allsystems {
		system.runfn(w.ctxbuilder(rctx, delta, system, w))
	}
	return time.Now().Sub(t0)
}

// RunWithTag will iterate through all systems that contain the given tag (sorted by priority).
func (w *World) RunWithTag(tag string, delta float64) (taken time.Duration) {
	t0 := time.Now()
	w.lock.RLock()
	allsystems := w.systems
	w.lock.RUnlock()
	rctx := context.Background()
	for _, system := range allsystems {
		if !system.ContainsTag(tag) {
			continue
		}
		system.runfn(w.ctxbuilder(rctx, delta, system, w))
	}
	return time.Now().Sub(t0)
}

// RunWithoutTag will iterate through all systems that don't contain the given tag (sorted by priority).
func (w *World) RunWithoutTag(tag string, delta float64) (taken time.Duration) {
	t0 := time.Now()
	w.lock.RLock()
	allsystems := w.systems
	w.lock.RUnlock()
	rctx := context.Background()
	for _, system := range allsystems {
		if system.ContainsTag(tag) {
			continue
		}
		system.runfn(w.ctxbuilder(rctx, delta, system, w))
	}
	return time.Now().Sub(t0)
}

// Get a global variable
func (w *World) Get(key string) interface{} {
	return w.globals.Get(key)
}

// Set a global variable
func (w *World) Set(key string, val interface{}) {
	w.globals.Set(key, val)
}

type Worlder interface {
	NewEntity() Entity
	NewEntities(n int) []Entity
	ContainsEntity(entity Entity) bool
	AddComponentToEntity(entity Entity, component *Component, data interface{}) error
	RemoveComponentFromEntity(entity Entity, component *Component) error
	Query(components ...*Component) []QueryMatch
	NewComponent(input NewComponentInput) (*Component, error)
	Component(name string) *Component
}

type Dicter interface {
	Get(key string) interface{}
	Set(key string, val interface{})
}

type WorldDicter interface {
	Worlder
	Dicter
}
