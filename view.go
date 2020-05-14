package ecs

import (
	"sync"
	"sync/atomic"
)

// QueryMatch represents an entity and a subset of components
// that belongs to this entity.
type QueryMatch struct {
	Entity     Entity
	Components map[*Component]interface{}
}

// View is the result set of a stored query on the data (entity + components).
type View struct {
	lock        sync.RWMutex
	id          uint64
	world       *World
	matches     []QueryMatch
	matchmap    map[Entity]int
	includemask flag
	excludemask flag
}

func getfingerprint(components ...*Component) flag {
	strslc := newflag(0, 0, 0, 0)
	for _, v := range components {
		strslc = strslc.or(v.flag)
	}
	return strslc
}

func containsfingerprint(bigger, smaller flag) bool {
	return bigger.contains(smaller)
}

// NewView creates a view on the current world based on the combination of components.
func (w *World) NewView(components ...*Component) *View {
	view := &View{
		id:          atomic.AddUint64(&w.nextView, 1) - 1,
		world:       w,
		matchmap:    make(map[Entity]int),
		matches:     w.Query(components...),
		includemask: getfingerprint(components...),
		excludemask: newflag(0, 0, 0, 0),
	}
	for k, v := range view.matches {
		view.matchmap[v.Entity] = k
	}
	w.lock.Lock()
	w.views = append(w.views, view)
	w.lock.Unlock()
	return view
}

// NewMaskView creates a view on the current world based on the combination of components.
func (w *World) NewMaskView(excludemask []*Component, includemask []*Component) *View {
	if excludemask == nil {
		excludemask = make([]*Component, 0)
	}
	if includemask == nil {
		includemask = make([]*Component, 0)
	}
	view := &View{
		id:          atomic.AddUint64(&w.nextView, 1) - 1,
		world:       w,
		matchmap:    make(map[Entity]int),
		matches:     w.QueryMask(excludemask, includemask),
		includemask: getfingerprint(includemask...),
		excludemask: getfingerprint(excludemask...),
	}
	for k, v := range view.matches {
		view.matchmap[v.Entity] = k
	}
	w.lock.Lock()
	w.views = append(w.views, view)
	w.lock.Unlock()
	return view
}

// Matches returns all the matches of a view.
func (v *View) Matches() []QueryMatch {
	v.lock.RLock()
	defer v.lock.RUnlock()
	return v.matches
}

// World return the world the view belongs to.
func (v *View) World() *World {
	v.lock.RLock()
	defer v.lock.RUnlock()
	return v.world
}

func (v *View) upsert(entity Entity) {
	// world is already read locked!
	v.lock.Lock()
	defer v.lock.Unlock()
	if oldindex, ok := v.matchmap[entity]; ok {
		item := v.matches[oldindex]
		// reapply data?
		// world is already read locked
		for _, comp := range v.world.components {
			comp.lock.RLock()
			if cdata, ok := comp.data[entity]; ok {
				if _, ok2 := item.Components[comp]; ok2 {
					item.Components[comp] = cdata
				}
			}
			comp.lock.RUnlock()
		}
		return
	}
	// INSERT
	nextindex := len(v.matches)
	nmap := make(map[*Component]interface{})
	newmatch := QueryMatch{
		Entity: entity,
	}
	// world is already read locked!
	eflag := v.world.entities[entity]
	for cflag, comp := range v.world.components {
		if eflag.contains(cflag) {
			comp.lock.RLock()
			nmap[comp] = comp.data[entity]
			comp.lock.RUnlock()
		}
	}
	newmatch.Components = nmap
	v.matches = append(v.matches, newmatch)
	v.matchmap[entity] = nextindex
}

func (v *View) remove(entity Entity) {
	v.lock.Lock()
	defer v.lock.Unlock()
	kkey, exists := v.matchmap[entity]
	if !exists {
		return
	}
	delete(v.matchmap, entity)
	v.matches = append(v.matches[:kkey], v.matches[kkey+1:]...)
	// reassign indexes
	for i := kkey; i < len(v.matches); i++ {
		v.matchmap[v.matches[i].Entity] = i
	}
}

func (v *View) matchesEntitySignature(f flag) bool {
	//v.lock.RLock()
	//defer v.lock.RUnlock()
	if !f.contains(v.includemask) {
		return false
	}
	if !f.and(v.excludemask).iszero() {
		return false
	}
	return true
}
