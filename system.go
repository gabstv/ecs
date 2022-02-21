package ecs

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

type System[T ComponentType] struct {
	*systemCore
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
		systemCore: newSystemCore(priority, world),
	}
	sys.view = NewView[T](world, sys.onAdded, sys.onRemoved)
	id := world.addSystem(sys)
	sys.id = id
	return sys
}

type System2[T1 ComponentType, T2 ComponentType] struct {
	*systemCore
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
		systemCore: newSystemCore(priority, world),
	}
	sys.view = NewView2[T1, T2](world, sys.onAdded, sys.onRemoved)
	id := world.addSystem(sys)
	sys.id = id
	return sys
}

type System3[T1 ComponentType, T2 ComponentType, T3 ComponentType] struct {
	*systemCore
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
		systemCore: newSystemCore(priority, world),
	}
	sys.view = NewView3[T1, T2, T3](world, sys.onAdded, sys.onRemoved)
	id := world.addSystem(sys)
	sys.id = id
	return sys
}

type System4[T1, T2, T3, T4 ComponentType] struct {
	*systemCore
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
		systemCore: newSystemCore(priority, world),
	}
	sys.view = NewView4[T1, T2, T3, T4](world, sys.onAdded, sys.onRemoved)
	id := world.addSystem(sys)
	sys.id = id
	return sys
}
