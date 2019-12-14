package ecs

import (
	"sort"
	"sync"
)

// SystemExec is the loop function of a system.
//
// dt (delta time): the time taken between the last World.Run
//
// view: the combination of entity + component(s) data
type SystemExec func(ctx Context)

// SystemMiddleware // TODO: godoc
type SystemMiddleware func(next SystemExec) SystemExec

// System is the brain of an ECS.
// The system performs global actions on every Entity that possesses
// a Component (or a combination of components) of the same aspect as
// that System.
type System struct {
	name     string
	priority int
	view     *View
	runfn    SystemExec
	tags     map[string]bool
	taglock  sync.RWMutex
	dict     *dict
	world    *World
}

type sortedSystems []*System

func (a sortedSystems) Len() int {
	return len(a)
}
func (a sortedSystems) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a sortedSystems) Less(i, j int) bool {
	return a[i].priority > a[j].priority
}

// NewSystem constructs a new system.
//
// The priority isused to sort the execution between systems.
//
// fn is the loop function of a system.
//
// comps is the combination of components that a system will iterate on.
func (w *World) NewSystem(name string, priority int, fn SystemExec, comps ...*Component) *System {
	sys := &System{
		name:     name,
		priority: priority,
		runfn:    fn,
		tags:     make(map[string]bool),
		dict:     newdict(),
		world:    w,
	}
	sys.view = w.NewView(comps...)
	w.lock.Lock()
	w.systems = append(w.systems, sys)
	if name != "" {
		if _, ok := w.systemNames[name]; ok {
			panic("system " + name + " already exists on world")
		}
		w.systemNames[name] = sys
	}
	sort.Sort(sortedSystems(w.systems))
	w.lock.Unlock()
	return sys
}

// AddTag adds a tag to the system. It is used to filter execution with
// World.RunWithTag and World.RunWithoutTag
func (sys *System) AddTag(tag string) *System {
	sys.taglock.Lock()
	sys.tags[tag] = true
	sys.taglock.Unlock()
	return sys
}

// RemoveTag removes a tag from the system. It is used to filter execution with
// World.RunWithTag and World.RunWithoutTag
func (sys *System) RemoveTag(tag string) *System {
	sys.taglock.Lock()
	delete(sys.tags, tag)
	sys.taglock.Unlock()
	return sys
}

// ContainsTag checks if the system contains the tag. It is used to filter execution with
// World.RunWithTag and World.RunWithoutTag
func (sys *System) ContainsTag(tag string) bool {
	sys.taglock.RLock()
	defer sys.taglock.RUnlock()
	return sys.tags[tag]
}

// Get a variable from the system dictionary
func (sys *System) Get(key string) interface{} {
	return sys.dict.Get(key)
}

// Set a variable to the system dictionary
func (sys *System) Set(key string, val interface{}) {
	sys.dict.Set(key, val)
}

// World returns the world the system belongs to.
func (sys *System) World() *World {
	return sys.world
}

// View returns the current view of the system
func (sys *System) View() *View {
	return sys.view
}

func SysWrapFn(fn SystemExec, mid ...SystemMiddleware) SystemExec {
	return func(ctx Context) {
		for _, m := range mid {
			lfn := fn
			fn = m(fn)
			if fn == nil {
				lfn(ctx)
				return
			}
		}
		if fn == nil {
			return
		}
		fn(ctx)
	}
}

// System returns a registered system by name
func (w *World) System(name string) *System {
	w.lock.RLock()
	defer w.lock.RUnlock()
	return w.systemNames[name]
}
