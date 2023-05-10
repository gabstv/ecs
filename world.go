package ecs

import (
	"fmt"
	"reflect"
	"sort"

	"github.com/Pilatuz/bigz/uint256"
	"golang.org/x/exp/slices"
)

type World interface {
	Step()

	addComponentStorage(cs worldComponentStorage)
	addSystem(sys worldSystem) (SystemID, error)
	addStartupSystem(sys System)
	getCommands() *Commands
	getComponentStorage(reflect.Type) worldComponentStorage
	getFatEntity(Entity) *fatEntity
	// getQuery may return nil
	getQuery(TypeTape) any
	newComponentMask() U256
	newEntity() Entity
	removeEntity(Entity)
	setResource(TypeMapKey, any)
	getResource(TypeMapKey) any
}

type worldImpl struct {
	lastEntityID      uint64
	lastSystemID      uint64
	lastComponentMask U256
	entities          []fatEntity
	components        []worldComponentStorage
	componentsIndex   map[TypeMapKey]int // TypeHash here represents a single type
	systems           []worldSystem
	startupSystems    []System
	queries           map[TypeTape]any
	resources         map[TypeMapKey]any

	lastCommands *Commands

	entitiesNeedSorting bool
}

func (w *worldImpl) Step() {
	commands := w.getCommands()
	for _, v := range w.startupSystems {
		v(commands)
		commands.run()
		w.clearCommands()
	}
	w.startupSystems = w.startupSystems[:0]
	for _, v := range w.systems {
		v.Value(commands)
		commands.run()
		w.clearCommands()
	}
	w.commit()
}

func NewWorld() World {
	return &worldImpl{
		components:      make([]worldComponentStorage, 0, 256),
		componentsIndex: make(map[TypeMapKey]int),
		queries:         make(map[TypeTape]any),
		resources:       make(map[TypeMapKey]any),
		systems:         make([]worldSystem, 0, 1024),
		startupSystems:  make([]System, 0, 256),
	}
}

// worldImpl implements World interface

func (w *worldImpl) addComponentStorage(cs worldComponentStorage) {
	assert(w.componentsIndex != nil, "w.componentsIndex is nil")
	assert(w.components != nil, "w.components is nil")

	th := typeMapKeyOf(cs.ComponentType())
	// make sure that the component type is not already registered
	_, exists := w.componentsIndex[th]
	assert(!exists, "_BUG_ - component type already registered")
	l := len(w.components)
	w.components = append(w.components, cs)
	w.componentsIndex[th] = l
}

func (w *worldImpl) addSystem(sys worldSystem) (SystemID, error) {
	w.lastSystemID++
	sys.ID = SystemID(w.lastSystemID)
	w.systems = append(w.systems, sys)
	if len(w.systems) == 1 {
		return sys.ID, nil
	}
	if w.systems[len(w.systems)-1].SortPriority >= w.systems[len(w.systems)-2].SortPriority {
		return sys.ID, nil
	}
	// sort
	sort.Slice(w.systems, func(i, j int) bool {
		if w.systems[i].SortPriority == w.systems[j].SortPriority {
			return w.systems[i].ID < w.systems[j].ID
		}
		return w.systems[i].SortPriority < w.systems[j].SortPriority
	})
	return sys.ID, nil
}

func (w *worldImpl) addStartupSystem(sys System) {
	w.startupSystems = append(w.startupSystems, sys)
}

func (w *worldImpl) getCommands() *Commands {
	if w.lastCommands != nil {
		return w.lastCommands
	}
	w.lastCommands = &Commands{
		world: w,
	}
	return w.lastCommands
}

func (w *worldImpl) getComponentStorage(t reflect.Type) worldComponentStorage {
	th := typeMapKeyOf(t)
	if i, ok := w.componentsIndex[th]; ok {
		assert(w.components[i].ComponentType() == t, "_BUG_ - component type mismatch")
		return w.components[i]
	}
	return nil
}

func (w *worldImpl) getFatEntity(e Entity) *fatEntity {
	if w.entities[len(w.entities)-1].Entity == e {
		return &w.entities[len(w.entities)-1]
	}
	if index, ok := slices.BinarySearchFunc(w.entities, e, func(fe fatEntity, target Entity) int {
		if fe.Entity < target {
			return -1
		}
		if fe.Entity > target {
			return 1
		}
		return 0
	}); ok {
		return &w.entities[index]
	}
	return nil
}

func (w *worldImpl) getQuery(tt TypeTape) any {
	return w.queries[tt]
}

func (w *worldImpl) newComponentMask() U256 {
	if w.lastComponentMask.IsZero() {
		w.lastComponentMask = uint256.One()
	} else {
		w.lastComponentMask = w.lastComponentMask.Lsh(1)
	}
	return w.lastComponentMask
}

func (w *worldImpl) newEntity() Entity {
	w.lastEntityID++
	if len(w.entities) > 0 {
		if w.entities[len(w.entities)-1].Entity > Entity(w.lastEntityID) {
			w.entitiesNeedSorting = true
		}
	}
	w.entities = append(w.entities, fatEntity{
		Entity:       Entity(w.lastEntityID),
		ComponentMap: uint256.Zero(),
	})
	return Entity(w.lastEntityID)
}

func (w *worldImpl) removeEntity(e Entity) {
	if index, ok := slices.BinarySearchFunc(w.entities, e, func(fe fatEntity, target Entity) int {
		if fe.Entity < target {
			return -1
		}
		if fe.Entity > target {
			return 1
		}
		return 0
	}); ok {
		w.entities[index].IsRemoved = true
		// remove all components
		for i, v := range w.components {
			bm := uint256.One().Lsh(uint(i))
			if w.entities[index].ComponentMap.And(bm).IsZero() {
				continue
			}
			v.removeEntity(e)
		}
		w.entities[index].ComponentMap = uint256.Zero()
	}
}

func (w *worldImpl) getResource(k TypeMapKey) any {
	if r, ok := w.resources[k]; ok {
		return r
	}
	return nil
}

func (w *worldImpl) setResource(k TypeMapKey, r any) {
	if _, ok := w.resources[k]; ok {
		kv := reflect.Value(k)
		panic(fmt.Sprintf("resource already set %v", kv.Type().String()))
	}
	w.resources[k] = r
	if dr, ok := r.(WorldIniter); ok {
		dr.Init(w)
	}
}

func getOrCreateComponentStorage[T Component](w World) *componentStorage[T] {
	var zt T
	ct := w.getComponentStorage(reflect.TypeOf(zt))
	if ct != nil {
		return ct.(*componentStorage[T])
	}
	// create new

	tct := &componentStorage[T]{
		zeroType:  reflect.TypeOf(zt),
		zeroValue: zt,
		mask:      w.newComponentMask(),
		Items:     make([]componentStore[T], 0, 16),
	}
	w.addComponentStorage(tct)
	return tct
}

func (w *worldImpl) clearCommands() {
	if w.lastCommands == nil {
		return
	}
	w.lastCommands.list = w.lastCommands.list[:0]
}

func (w *worldImpl) commit() {
	if w.entitiesNeedSorting {
		sort.Slice(w.entities, func(i, j int) bool {
			return w.entities[i].Entity < w.entities[j].Entity
		})
		w.entitiesNeedSorting = false
	}
}
