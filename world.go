package ecs

import (
	"errors"
	"sort"
	"sync"
)

type World struct {
	l          sync.RWMutex
	lentity    Entity
	lflag      uint8
	components map[string]BaseComponent
	systems    []BaseSystem
	syscache   map[string]BaseSystem

	entities []EntityFlag
	key      [4]byte
}

func (w *World) RegisterComponent(c BaseComponent) {
	w.l.Lock()
	defer w.l.Unlock()
	if _, ok := w.components[c.UUID()]; ok {
		panic("component " + c.Name() + " already registered (" + c.UUID() + ")")
	}
	w.components[c.UUID()] = c
	c.Setup(w, NewFlag(w.lflag), w.key)
	w.lflag++
}

func (w *World) IsRegistered(id string) bool {
	w.l.RLock()
	defer w.l.RUnlock()
	_, ok := w.components[id]
	return ok
}

func (w *World) entityindex(e Entity) int {
	i := sort.Search(len(w.entities), func(i int) bool { return w.entities[i].Entity >= e })
	if i < len(w.entities) && w.entities[i].Entity == e {
		return i
	}
	return -1
}

func (w *World) CFlag(e Entity) Flag {
	w.l.RLock()
	defer w.l.RUnlock()
	i := w.entityindex(e)
	if i == -1 {
		return NewFlagRaw(0, 0, 0, 0)
	}
	return w.entities[i].Flag
}

func (w *World) C(id string) BaseComponent {
	// All componets should already be loaded at this point,
	// so no locking is done
	//
	// w.l.Lock()
	// defer w.l.Unlock()
	return w.components[id]
}

func (w *World) S(id string) BaseSystem {
	w.l.RLock()
	defer w.l.RUnlock()
	return w.syscache[id]
}

func (w *World) CAdded(e Entity, c BaseComponent, key [4]byte) {
	if w.key[0] != key[0] || w.key[1] != key[1] || w.key[2] != key[2] || w.key[3] != key[3] {
		panic("CAdded forbidden")
	}
	i := w.entityindex(e)
	w.entities[i].Flag = w.entities[i].Flag.Or(c.Flag())
	for _, sys := range w.systems {
		sys.ComponentAdded(e, w.entities[i].Flag)
	}
}

func (w *World) CRemoved(e Entity, c BaseComponent, key [4]byte) {
	if w.key[0] != key[0] || w.key[1] != key[1] || w.key[2] != key[2] || w.key[3] != key[3] {
		panic("CRemoved forbidden")
	}
	i := w.entityindex(e)
	w.entities[i].Flag = w.entities[i].Flag.Xor(c.Flag())
	for _, sys := range w.systems {
		sys.ComponentRemoved(e, w.entities[i].Flag)
	}
}

func (w *World) CResized(c BaseComponent, key [4]byte) {
	if w.key[0] != key[0] || w.key[1] != key[1] || w.key[2] != key[2] || w.key[3] != key[3] {
		panic("CResized forbidden")
	}
	for _, sys := range w.systems {
		sys.ComponentResized(c.Flag())
	}
}

func (w *World) NewEntity() Entity {
	w.l.Lock()
	w.lentity++
	e := w.lentity
	w.entities = append(w.entities, EntityFlag{
		Entity: e,
		Flag:   NewFlagRaw(0, 0, 0, 0),
	})
	w.l.Unlock()
	return e
}

func (w *World) RemoveEntity(e Entity) bool {
	w.l.RLock()
	i := w.entityindex(e)
	if i == -1 {
		w.l.RUnlock()
		return false
	}
	f := w.entities[i].Flag
	w.l.RUnlock()
	for _, comp := range w.components {
		if f.Contains(comp.Flag()) {
			//TODO: optimize by ignoring CRemoved from this entity
			comp.Remove(e)
		}
	}
	w.l.Lock()
	w.entities = w.entities[:i+copy(w.entities[i:], w.entities[i+1:])]
	w.l.Unlock()
	return true
}

// AddSystem returns an error if the system was already added
func (w *World) AddSystem(s BaseSystem) error {
	w.l.Lock()
	defer w.l.Unlock()
	if _, ok := w.syscache[s.UUID()]; ok {
		return errors.New("system already added (UUID: " + s.UUID() + " Name: " + s.Name() + ")")
	}
	w.syscache[s.UUID()] = s
	w.systems = append(w.systems, s)
	if len(w.systems) > 1 {
		sort.SliceStable(w.systems, func(i, j int) bool {
			return w.systems[i].Priority() > w.systems[j].Priority()
		})
	}
	s.Setup(w)
	return nil
}

func (w *World) RemoveSystem(s BaseSystem) {
	w.l.Lock()
	defer w.l.Unlock()
	if _, ok := w.syscache[s.UUID()]; !ok {
		return
	}
	delete(w.syscache, s.UUID())
	i := -1
	for k, v := range w.systems {
		if v.UUID() == s.UUID() {
			i = k
			break
		}
	}
	if i != -1 {
		w.systems = w.systems[:i+copy(w.systems[i:], w.systems[i+1:])]
		if len(w.systems) > 1 {
			sort.SliceStable(w.systems, func(i, j int) bool {
				return w.systems[i].Priority() > w.systems[j].Priority()
			})
		}
	}
}

func (w *World) Init() {
	w.components = make(map[string]BaseComponent)
	w.systems = make([]BaseSystem, 0)
	w.syscache = make(map[string]BaseSystem)
	w.entities = make([]EntityFlag, 0)
	w.key = [4]byte{10, 227, 227, 9}
}

func (w *World) EachSystem(fn func(s BaseSystem) bool) {
	w.l.RLock()
	clone := make([]BaseSystem, 0, len(w.systems))
	for _, v := range w.systems {
		if v.Enabled() {
			clone = append(clone, v)
		}
	}
	w.l.RUnlock()
	for _, s := range clone {
		if !fn(s) {
			return
		}
	}
}

func NewWorld() BaseWorld {
	w := &World{}
	w.Init()
	return w
}
