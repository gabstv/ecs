package ecs

import (
	"errors"
	"sort"
	"sync"
)

type world struct {
	l          sync.RWMutex
	lentity    Entity
	lflag      uint8
	components map[string]Component
	systems    []System
	syscache   map[string]System

	entities []EntityFlag
	key      [4]byte

	ei    int64
	evts  map[int64]*EventListener
	evtfs map[EventType]map[int64]*EventListener

	flaggroupsm sync.RWMutex
	flaggroups  map[string]Flag

	locker Locker
}

func (w *world) RegisterComponent(c Component) {
	w.l.Lock()
	defer w.l.Unlock()
	if _, ok := w.components[c.UUID()]; ok {
		panic("component " + c.Name() + " already registered (" + c.UUID() + ")")
	}
	w.components[c.UUID()] = c
	c.Setup(w, NewFlag(w.lflag), w.key)
	w.lflag++
}

func (w *world) IsRegistered(id string) bool {
	w.l.RLock()
	defer w.l.RUnlock()
	_, ok := w.components[id]
	return ok
}

func (w *world) entityindex(e Entity) int {
	i := sort.Search(len(w.entities), func(i int) bool { return w.entities[i].Entity >= e })
	if i < len(w.entities) && w.entities[i].Entity == e {
		return i
	}
	return -1
}

func (w *world) CFlag(e Entity) Flag {
	w.l.RLock()
	defer w.l.RUnlock()
	i := w.entityindex(e)
	if i == -1 {
		return NewFlagRaw(0, 0, 0, 0)
	}
	return w.entities[i].Flag
}

func (w *world) C(id string) Component {
	// All componets should already be loaded at this point,
	// so no locking is done
	//
	// w.l.Lock()
	// defer w.l.Unlock()
	return w.components[id]
}

func (w *world) S(id string) System {
	w.l.RLock()
	defer w.l.RUnlock()
	return w.syscache[id]
}

func (w *world) CAdded(e Entity, c Component, key [4]byte) {
	if w.key[0] != key[0] || w.key[1] != key[1] || w.key[2] != key[2] || w.key[3] != key[3] {
		panic("CAdded forbidden")
	}
	i := w.entityindex(e)
	w.entities[i].Flag = w.entities[i].Flag.Or(c.Flag())
	for _, sys := range w.systems {
		sys.ComponentAdded(e, w.entities[i].Flag)
	}
}

func (w *world) CRemoved(e Entity, c Component, key [4]byte) {
	if w.key[0] != key[0] || w.key[1] != key[1] || w.key[2] != key[2] || w.key[3] != key[3] {
		panic("CRemoved forbidden")
	}
	i := w.entityindex(e)
	w.entities[i].Flag = w.entities[i].Flag.Xor(c.Flag())
	for _, sys := range w.systems {
		sys.ComponentRemoved(e, w.entities[i].Flag)
	}
}

func (w *world) CWillResize(c Component, key [4]byte) {
	if w.key[0] != key[0] || w.key[1] != key[1] || w.key[2] != key[2] || w.key[3] != key[3] {
		panic("CWillResize forbidden")
	}
	for _, sys := range w.systems {
		sys.ComponentWillResize(c.Flag())
	}
}

func (w *world) CResized(c Component, key [4]byte) {
	if w.key[0] != key[0] || w.key[1] != key[1] || w.key[2] != key[2] || w.key[3] != key[3] {
		panic("CResized forbidden")
	}
	for _, sys := range w.systems {
		sys.ComponentResized(c.Flag())
	}
}

func (w *world) NewEntity() Entity {
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

func (w *world) RemoveEntity(e Entity) bool {
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
func (w *world) AddSystem(s System) error {
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

func (w *world) RemoveSystem(s System) {
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

func (w *world) SetFlagGroup(name string, f Flag) {
	w.flaggroupsm.Lock()
	defer w.flaggroupsm.Unlock()
	if w.flaggroups == nil {
		w.flaggroups = make(map[string]Flag)
	}
	w.flaggroups[name] = f
}

func (w *world) FlagGroup(name string) Flag {
	w.flaggroupsm.RLock()
	defer w.flaggroupsm.RUnlock()
	if w.flaggroups == nil {
		return Flag{}
	}
	return w.flaggroups[name]
}

func (w *world) LGet(name string) interface{} {
	return w.locker.Item(name)
}

func (w *world) LSet(name string, value interface{}) {
	w.locker.SetItem(name, value)
}

func (w *world) Init() {
	w.components = make(map[string]Component)
	w.systems = make([]System, 0)
	w.syscache = make(map[string]System)
	w.entities = make([]EntityFlag, 0)
	w.key = [4]byte{10, 227, 227, 9}
	w.evts = make(map[int64]*EventListener)
	w.evtfs = make(map[EventType]map[int64]*EventListener)
	for _, t := range allevents {
		w.evtfs[t] = make(map[int64]*EventListener)
	}
}

func (w *world) EachSystem(fn func(s System) bool) {
	w.l.RLock()
	clone := make([]System, 0, len(w.systems))
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

func (w *world) Dispatch(e Event) {
	w.l.RLock()
	m := w.evtfs[e.Type]
	evs := make([]EventFn, 0, len(m))
	for _, v := range m {
		evs = append(evs, v.Fn)
	}
	w.l.RUnlock()
	for _, v := range evs {
		v(e)
	}
}

func (w *world) Listen(mask EventType, fn EventFn) int64 {
	w.l.Lock()
	defer w.l.Unlock()
	w.ei++
	id := w.ei
	x := &EventListener{
		ID:   id,
		Fn:   fn,
		Mask: mask,
	}
	w.evts[id] = x
	for _, t := range allevents {
		if mask&t == t {
			w.evtfs[t][id] = x
		}
	}
	return id
}

func (w *world) RemoveListener(id int64) {
	w.l.Lock()
	defer w.l.Unlock()
	l, ok := w.evts[id]
	if !ok {
		return
	}
	for _, t := range allevents {
		if l.Mask&t == t {
			delete(w.evtfs[t], l.ID)
		}
	}
	delete(w.evts, l.ID)
}

func NewWorld() World {
	w := &world{}
	w.Init()
	return w
}
