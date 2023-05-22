package ecs

import (
	"fmt"
	"reflect"
	"sort"
	"sync"

	"github.com/Pilatuz/bigz/uint256"
	"golang.org/x/exp/slices"
)

type World interface {
	// Exec executes code immediately in the current world. This should be avoided, but it's useful
	// for game engines  to setup things before a world.Step() call.
	// The function passed to Exec will be executed as a startup system. Any call to LocalResource
	// will panic.
	Exec(func(*Context))
	// ShallowCopy will return a new world that shares the same entities, components and resources, but
	// not the same systems. This is useful to separate logic from rendering, for example.
	ShallowCopy() World
	Step()

	addComponentStorage(cs worldComponentStorage)
	addSystem(sys worldSystem) (SystemID, error)
	addStartupSystem(sys System)
	getContext() *Context
	getComponentStorage(Component) worldComponentStorage
	getFatEntity(Entity) *fatEntity
	// getQuery may return nil
	getQuery(TypeTape) any
	newComponentMask() U256
	newEntity() Entity
	removeEntity(*Context, Entity)
	setResource(TypeMapKey, any)
	getResource(TypeMapKey) any
	getEvents() map[TypeMapKey]any
	getComponentAddedEvents() map[TypeMapKey]any
	getComponentRemovedEvents() map[TypeMapKey]any
	gc()
}

type worldImpl struct {
	lastEntityID      uint64
	lastSystemID      uint64
	lastComponentMask U256
	entities          []fatEntity
	events            map[TypeMapKey]any
	components        []worldComponentStorage
	componentsIndex   map[ComponentUUID]int // TypeHash here represents a single type
	componentsAdded   map[TypeMapKey]any
	componentsRemoved map[TypeMapKey]any
	systems           []worldSystem
	startupSystemsBuf []System
	startupSystems    Container[System]
	queries           map[TypeTape]any
	resources         map[TypeMapKey]any

	lastContext      *Context
	lastContextMutex sync.Mutex

	entitiesNeedSorting bool
}

func (w *worldImpl) getComponentAddedEvents() map[TypeMapKey]any {
	return w.componentsAdded
}
func (w *worldImpl) getComponentRemovedEvents() map[TypeMapKey]any {
	return w.componentsRemoved
}

func (w *worldImpl) ShallowCopy() World {
	return &worldShallowCopy{
		parent:            w,
		systems:           make([]worldSystem, 0, 32),
		startupSystemsBuf: make([]System, 0, 16),
		events:            make(map[TypeMapKey]any),
	}
}

func (w *worldImpl) Step() {
	for _, v := range w.events {
		v.(genericEventStorage).step()
	}
	for _, v := range w.componentsAdded {
		v.(genericEventStorage).step()
	}
	for _, v := range w.componentsRemoved {
		v.(genericEventStorage).step()
	}
	ctx := w.getContext()
	ctx.currentSystem = nil
	ctx.isStartupSystem = true
	ctx.currentSystemIndex = -1000
	w.startupSystemsBuf = w.startupSystems.GetAllAndDeleteAll(w.startupSystemsBuf)
	for _, v := range w.startupSystemsBuf {
		v(ctx)
		ctx.run()
		ctx.currentSystemIndex++
	}
	ctx.currentSystemIndex = 0
	w.startupSystemsBuf = w.startupSystemsBuf[:0]
	ctx.isStartupSystem = false
	for i, v := range w.systems {
		ctx.currentSystem = &w.systems[i]
		v.Value(ctx)
		ctx.run()
		ctx.currentSystemIndex++
	}
	ctx.currentSystem = nil
	w.commit()
}

func (w *worldImpl) Exec(fn func(*Context)) {
	commands := w.getContext()
	commands.currentSystem = nil
	commands.isStartupSystem = true
	fn(commands)
	commands.run()
}

func NewWorld() World {
	return &worldImpl{
		components:        make([]worldComponentStorage, 0, 256),
		componentsIndex:   make(map[ComponentUUID]int),
		componentsAdded:   make(map[TypeMapKey]any),
		componentsRemoved: make(map[TypeMapKey]any),
		events:            make(map[TypeMapKey]any),
		queries:           make(map[TypeTape]any),
		resources:         make(map[TypeMapKey]any),
		systems:           make([]worldSystem, 0, 1024),
		startupSystemsBuf: make([]System, 0, 256),
	}
}

// worldImpl implements World interface

func (w *worldImpl) addComponentStorage(cs worldComponentStorage) {
	assert(w.componentsIndex != nil, "w.componentsIndex is nil")
	assert(w.components != nil, "w.components is nil")

	iv := reflect.Zero(cs.ComponentType()).Interface()
	uuid := iv.(Component).ComponentUUID()
	// make sure that the component type is not already registered
	_, exists := w.componentsIndex[uuid]
	assert(!exists, "_BUG_ - component type already registered")
	l := len(w.components)
	w.components = append(w.components, cs)
	w.componentsIndex[uuid] = l
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

// this is thread safe
func (w *worldImpl) addStartupSystem(sys System) {
	w.startupSystems.Insert(sys)
}

func (w *worldImpl) getContext() *Context {
	w.lastContextMutex.Lock()
	defer w.lastContextMutex.Unlock()
	if w.lastContext != nil {
		return w.lastContext
	}
	w.lastContext = &Context{
		world:       w,
		eventRWLock: &sync.Mutex{},
	}
	return w.lastContext
}

func (w *worldImpl) getComponentStorage(zv Component) worldComponentStorage {
	if zv == nil {
		panic("zv is nil")
	}
	if i, ok := w.componentsIndex[zv.ComponentUUID()]; ok {
		// assert(w.components[i].ComponentType() == t, "_BUG_ - component type mismatch")
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

func (w *worldImpl) getEvents() map[TypeMapKey]any {
	return w.events
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

func (w *worldImpl) removeEntity(ctx *Context, e Entity) {
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
			if d := v.removeEntity(ctx, e); d != nil {
				v.fireComponentRemovedEvent(ctx, e, d)
			}
		}
		w.entities[index].ComponentMap = uint256.Zero()
	}
}

func (w *worldImpl) gc() {
	for _, c := range w.components {
		c.gc()
	}
	ecopy := make([]fatEntity, 0, cap(w.entities))
	for _, v := range w.entities {
		if v.IsRemoved {
			continue
		}
		ecopy = append(ecopy, v)
	}
	w.entities = ecopy
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

var getOrCreateComponentStorageLock sync.Mutex

func getOrCreateComponentStorage[T Component](w World, capacity int) *componentStorage[T] {
	getOrCreateComponentStorageLock.Lock()
	defer getOrCreateComponentStorageLock.Unlock()
	if capacity == 0 {
		capacity = 16
	}
	var zt T
	ct := w.getComponentStorage(zt)
	if ct != nil {
		return ct.(*componentStorage[T])
	}
	// create new

	tct := &componentStorage[T]{
		zeroType:  reflect.TypeOf(zt),
		zeroValue: zt,
		mask:      w.newComponentMask(),
		items:     make([]componentStore[T], 0, capacity),
	}
	w.addComponentStorage(tct)
	return tct
}

func (w *worldImpl) commit() {
	if w.entitiesNeedSorting {
		sort.Slice(w.entities, func(i, j int) bool {
			return w.entities[i].Entity < w.entities[j].Entity
		})
		w.entitiesNeedSorting = false
	}
}

/// // //

type worldShallowCopy struct {
	parent *worldImpl

	events            map[TypeMapKey]any
	lastSystemID      uint64
	systems           []worldSystem
	startupSystems    Container[System]
	startupSystemsBuf []System
}

func (w *worldShallowCopy) getComponentAddedEvents() map[TypeMapKey]any {
	return w.parent.getComponentAddedEvents()
}
func (w *worldShallowCopy) getComponentRemovedEvents() map[TypeMapKey]any {
	return w.parent.getComponentRemovedEvents()
}

func (w *worldShallowCopy) Exec(fn func(*Context)) {
	commands := w.getContext()
	commands.currentSystem = nil
	commands.isStartupSystem = true
	fn(commands)
	commands.run()
}

func (w *worldShallowCopy) ShallowCopy() World {
	return &worldShallowCopy{
		parent:            w.parent,
		systems:           make([]worldSystem, 0, 32),
		startupSystemsBuf: make([]System, 0, 16),
	}
}

func (w *worldShallowCopy) Step() {
	for _, v := range w.events {
		v.(genericEventStorage).step()
	}
	pw := w.parent
	commands := w.parent.getContext()
	commands.world = w
	defer func() {
		commands.world = pw
	}()
	w.startupSystemsBuf = w.startupSystems.GetAllAndDeleteAll(w.startupSystemsBuf)
	for _, v := range w.startupSystemsBuf {
		v(commands)
		commands.run()
	}
	w.startupSystemsBuf = w.startupSystemsBuf[:0]
	for _, v := range w.systems {
		v.Value(commands)
		commands.run()
	}
	w.parent.commit()
}

func (w *worldShallowCopy) addComponentStorage(cs worldComponentStorage) {
	w.parent.addComponentStorage(cs)
}

func (w *worldShallowCopy) addSystem(sys worldSystem) (SystemID, error) {
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

// this is thread safe :)
func (w *worldShallowCopy) addStartupSystem(sys System) {
	w.startupSystems.Insert(sys)
}

func (w *worldShallowCopy) getContext() *Context {
	return w.parent.getContext()
}

func (w *worldShallowCopy) getComponentStorage(zv Component) worldComponentStorage {
	return w.parent.getComponentStorage(zv)
}

func (w *worldShallowCopy) getEvents() map[TypeMapKey]any {
	return w.events
}

func (w *worldShallowCopy) getFatEntity(e Entity) *fatEntity {
	return w.parent.getFatEntity(e)
}

func (w *worldShallowCopy) getQuery(tt TypeTape) any {
	return w.parent.getQuery(tt)
}
func (w *worldShallowCopy) newComponentMask() U256 {
	return w.parent.newComponentMask()
}
func (w *worldShallowCopy) newEntity() Entity {
	return w.parent.newEntity()
}
func (w *worldShallowCopy) removeEntity(ctx *Context, e Entity) {
	w.parent.removeEntity(ctx, e)
}
func (w *worldShallowCopy) setResource(k TypeMapKey, r any) {
	w.parent.setResource(k, r)
}
func (w *worldShallowCopy) getResource(k TypeMapKey) any {
	return w.parent.getResource(k)
}

func (w *worldShallowCopy) gc() {
	w.parent.gc()
}
