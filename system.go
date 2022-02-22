package ecs

import (
	"runtime"
	"sync"
)

type ISystem interface {
	ID() int
	Execute()
	Priority() int
	Flag() int
}

type systemCore struct {
	priority      int
	world         *World
	EntityAdded   func(e Entity)
	EntityRemoved func(e Entity)
	id            int
	flag          int
}

func (s *systemCore) onAdded(e Entity) {
	if s.EntityAdded != nil {
		s.EntityAdded(e)
	}
}

func (s *systemCore) onRemoved(e Entity) {
	if s.EntityRemoved != nil {
		s.EntityRemoved(e)
	}
}

func (s *systemCore) ID() int {
	return s.id
}

func (s *systemCore) Priority() int {
	return s.priority
}

func (s *systemCore) Flag() int {
	return s.flag
}

func (s *systemCore) SetFlag(flag int) {
	s.flag = flag
}

func newSystemCore(priority int, world *World) *systemCore {
	return &systemCore{
		priority: priority,
		world:    world,
	}
}

type ifaceSystemDataProvider struct {
	data interface{}
}

func (s *ifaceSystemDataProvider) Data() interface{} {
	return s.data
}

func (s *ifaceSystemDataProvider) SetData(data interface{}) {

}

func newIfaceSystemDataProvider() *ifaceSystemDataProvider {
	return &ifaceSystemDataProvider{}
}

type System[T ComponentType] struct {
	*systemCore
	*ifaceSystemDataProvider
	view *View[T]
	Run  func(view *View[T])
}

func (s *System[T]) Execute() {
	if s.Run == nil {
		return
	}
	s.Run(s.view)
}

// WarmStart runs the OnEntityAdded callback for all entities in the world.
func (s *System[T]) WarmStart() {
	s.view.Each(func(e Entity, _ *T) {
		s.EntityAdded(e)
	})
}

func NewSystem[T ComponentType](priority int, world *World) *System[T] {
	sys := &System[T]{
		systemCore:              newSystemCore(priority, world),
		ifaceSystemDataProvider: newIfaceSystemDataProvider(),
	}
	sys.view = NewView[T](world, sys.onAdded, sys.onRemoved)
	id := world.addSystem(sys)
	sys.id = id
	return sys
}

type System2[T1 ComponentType, T2 ComponentType] struct {
	*systemCore
	*ifaceSystemDataProvider
	view *View2[T1, T2]
	Run  func(view *View2[T1, T2])
}

func (s *System2[T1, T2]) Execute() {
	if s.Run == nil {
		return
	}
	s.Run(s.view)
}

// WarmStart runs the OnEntityAdded callback for all entities in the world.
func (s *System2[T1, T2]) WarmStart() {
	s.view.Each(func(e Entity, _ *T1, _ *T2) {
		s.EntityAdded(e)
	})
}

func NewSystem2[T1 ComponentType, T2 ComponentType](priority int, world *World) *System2[T1, T2] {
	sys := &System2[T1, T2]{
		systemCore:              newSystemCore(priority, world),
		ifaceSystemDataProvider: newIfaceSystemDataProvider(),
	}
	sys.view = NewView2[T1, T2](world, sys.onAdded, sys.onRemoved)
	id := world.addSystem(sys)
	sys.id = id
	return sys
}

type System3[T1 ComponentType, T2 ComponentType, T3 ComponentType] struct {
	*systemCore
	*ifaceSystemDataProvider
	view *View3[T1, T2, T3]
	Run  func(view *View3[T1, T2, T3])
}

func (s *System3[T1, T2, T3]) Execute() {
	if s.Run == nil {
		return
	}
	s.Run(s.view)
}

// WarmStart runs the OnEntityAdded callback for all entities in the world.
func (s *System3[T1, T2, T3]) WarmStart() {
	s.view.Each(func(e Entity, _ *T1, _ *T2, _ *T3) {
		s.EntityAdded(e)
	})
}

func NewSystem3[T1 ComponentType, T2 ComponentType, T3 ComponentType](priority int, world *World) *System3[T1, T2, T3] {
	sys := &System3[T1, T2, T3]{
		systemCore:              newSystemCore(priority, world),
		ifaceSystemDataProvider: newIfaceSystemDataProvider(),
	}
	sys.view = NewView3[T1, T2, T3](world, sys.onAdded, sys.onRemoved)
	id := world.addSystem(sys)
	sys.id = id
	return sys
}

type System4[T1, T2, T3, T4 ComponentType] struct {
	*systemCore
	*ifaceSystemDataProvider
	view *View4[T1, T2, T3, T4]
	Run  func(view *View4[T1, T2, T3, T4])
}

func (s *System4[T1, T2, T3, T4]) Execute() {
	if s.Run == nil {
		return
	}
	s.Run(s.view)
}

// WarmStart runs the OnEntityAdded callback for all entities in the world.
func (s *System4[T1, T2, T3, T4]) WarmStart() {
	s.view.Each(func(e Entity, _ *T1, _ *T2, _ *T3, _ *T4) {
		s.EntityAdded(e)
	})
}

func NewSystem4[T1, T2, T3, T4 ComponentType](priority int, world *World) *System4[T1, T2, T3, T4] {
	sys := &System4[T1, T2, T3, T4]{
		systemCore:              newSystemCore(priority, world),
		ifaceSystemDataProvider: newIfaceSystemDataProvider(),
	}
	sys.view = NewView4[T1, T2, T3, T4](world, sys.onAdded, sys.onRemoved)
	id := world.addSystem(sys)
	sys.id = id
	return sys
}

// GlobalSystemInfo is the arg of RegisterGlobalSystem
type GlobalSystemInfo[T ComponentType] struct {
	ExecPriority         int
	ExecFlag             int
	ExecBuilder          func(s *System[T]) func(view *View[T])
	EntityAddedBuilder   func(s *System[T]) func(e Entity)
	EntityRemovedBuilder func(s *System[T]) func(e Entity)
	InitialDataBuilder   func() interface{}
	Finalizer            func(s *System[T])
}

// GlobalSystem2Info is the arg of RegisterGlobalSystem2
type GlobalSystem2Info[T1, T2 ComponentType] struct {
	ExecPriority         int
	ExecFlag             int
	ExecBuilder          func(s *System2[T1, T2]) func(view *View2[T1, T2])
	EntityAddedBuilder   func(s *System2[T1, T2]) func(e Entity)
	EntityRemovedBuilder func(s *System2[T1, T2]) func(e Entity)
	InitialDataBuilder   func() interface{}
	Finalizer            func(s *System2[T1, T2])
}

// GlobalSystem3Info is the arg of RegisterGlobalSystem3
type GlobalSystem3Info[T1, T2, T3 ComponentType] struct {
	ExecPriority         int
	ExecFlag             int
	ExecBuilder          func(s *System3[T1, T2, T3]) func(view *View3[T1, T2, T3])
	EntityAddedBuilder   func(s *System3[T1, T2, T3]) func(e Entity)
	EntityRemovedBuilder func(s *System3[T1, T2, T3]) func(e Entity)
	InitialDataBuilder   func() interface{}
	Finalizer            func(s *System3[T1, T2, T3])
}

// GlobalSystem4Info is the arg of RegisterGlobalSystem4
type GlobalSystem4Info[T1, T2, T3, T4 ComponentType] struct {
	ExecPriority         int
	ExecFlag             int
	ExecBuilder          func(s *System4[T1, T2, T3, T4]) func(view *View4[T1, T2, T3, T4])
	EntityAddedBuilder   func(s *System4[T1, T2, T3, T4]) func(e Entity)
	EntityRemovedBuilder func(s *System4[T1, T2, T3, T4]) func(e Entity)
	InitialDataBuilder   func() interface{}
	Finalizer            func(s *System4[T1, T2, T3, T4])
}

var (
	globalSystems = struct {
		lock       sync.Mutex
		sysFactory []func(w *World)
	}{
		sysFactory: make([]func(w *World), 0),
	}
)

// RegisterGlobalSystem registers a system to be included on every new world.
// All worlds initiated with NewWorld() will have this system included.
func RegisterGlobalSystem[T ComponentType](info GlobalSystemInfo[T]) {
	globalSystems.lock.Lock()
	defer globalSystems.lock.Unlock()
	globalSystems.sysFactory = append(globalSystems.sysFactory, func(w *World) {
		sys := NewSystem[T](info.ExecPriority, w)
		sys.Run = info.ExecBuilder(sys)
		sys.EntityAdded = info.EntityAddedBuilder(sys)
		sys.EntityRemoved = info.EntityRemovedBuilder(sys)
		sys.SetFlag(info.ExecFlag)
		sys.SetData(info.InitialDataBuilder())
		if info.Finalizer != nil {
			runtime.SetFinalizer(sys, func(s *System[T]) {
				runtime.SetFinalizer(s, nil)
				info.Finalizer(s)
			})
		}
	})
}

// RegisterGlobalSystem2 registers a system to be included on every new world.
// All worlds initiated with NewWorld() will have this system included.
func RegisterGlobalSystem2[T1, T2 ComponentType](info GlobalSystem2Info[T1, T2]) {
	globalSystems.lock.Lock()
	defer globalSystems.lock.Unlock()
	globalSystems.sysFactory = append(globalSystems.sysFactory, func(w *World) {
		sys := NewSystem2[T1, T2](info.ExecPriority, w)
		sys.Run = info.ExecBuilder(sys)
		sys.EntityAdded = info.EntityAddedBuilder(sys)
		sys.EntityRemoved = info.EntityRemovedBuilder(sys)
		sys.SetFlag(info.ExecFlag)
		sys.SetData(info.InitialDataBuilder())
		if info.Finalizer != nil {
			runtime.SetFinalizer(sys, func(s *System2[T1, T2]) {
				runtime.SetFinalizer(s, nil)
				info.Finalizer(s)
			})
		}
	})
}

// RegisterGlobalSystem3 registers a system to be included on every new world.
// All worlds initiated with NewWorld() will have this system included.
func RegisterGlobalSystem3[T1, T2, T3 ComponentType](info GlobalSystem3Info[T1, T2, T3]) {
	globalSystems.lock.Lock()
	defer globalSystems.lock.Unlock()
	globalSystems.sysFactory = append(globalSystems.sysFactory, func(w *World) {
		sys := NewSystem3[T1, T2, T3](info.ExecPriority, w)
		sys.Run = info.ExecBuilder(sys)
		sys.EntityAdded = info.EntityAddedBuilder(sys)
		sys.EntityRemoved = info.EntityRemovedBuilder(sys)
		sys.SetFlag(info.ExecFlag)
		sys.SetData(info.InitialDataBuilder())
		if info.Finalizer != nil {
			runtime.SetFinalizer(sys, func(s *System3[T1, T2, T3]) {
				runtime.SetFinalizer(s, nil)
				info.Finalizer(s)
			})
		}
	})
}

// RegisterGlobalSystem4 registers a system to be included on every new world.
// All worlds initiated with NewWorld() will have this system included.
func RegisterGlobalSystem4[T1, T2, T3, T4 ComponentType](info GlobalSystem4Info[T1, T2, T3, T4]) {
	globalSystems.lock.Lock()
	defer globalSystems.lock.Unlock()
	globalSystems.sysFactory = append(globalSystems.sysFactory, func(w *World) {
		sys := NewSystem4[T1, T2, T3, T4](info.ExecPriority, w)
		sys.Run = info.ExecBuilder(sys)
		sys.EntityAdded = info.EntityAddedBuilder(sys)
		sys.EntityRemoved = info.EntityRemovedBuilder(sys)
		sys.SetFlag(info.ExecFlag)
		sys.SetData(info.InitialDataBuilder())
		if info.Finalizer != nil {
			runtime.SetFinalizer(sys, func(s *System4[T1, T2, T3, T4]) {
				runtime.SetFinalizer(s, nil)
				info.Finalizer(s)
			})
		}
	})
}
