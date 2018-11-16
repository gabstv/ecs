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
type SystemExec func(dt float64, view *View)

// System is the brain of an ECS.
// The system performs global actions on every Entity that possesses
// a Component (or a combination of components) of the same aspect as
// that System.
type System struct {
	priority int
	view     *View
	runfn    SystemExec
	tags     map[string]bool
	taglock  sync.RWMutex
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
func (w *World) NewSystem(priority int, fn SystemExec, comps ...*Component) *System {
	sys := &System{
		priority: priority,
		runfn:    fn,
		tags:     make(map[string]bool),
	}
	sys.view = w.NewView(comps...)
	w.lock.Lock()
	w.systems = append(w.systems, sys)
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
